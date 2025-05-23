/*
Copyright 2022 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package decimal

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strings"

	"github.com/vedadiyan/sqlparser/pkg/mysql/fastparse"
)

var errOverflow = errors.New("overflow")

func parseDecimal64(s []byte) (Decimal, error) {
	const cutoff = math.MaxUint64/10 + 1
	var n uint64
	var dot = -1

	for i, c := range s {
		var d byte
		switch {
		case c == '.':
			if dot > -1 {
				return Decimal{}, fmt.Errorf("too many .s")
			}
			dot = i
			continue
		case '0' <= c && c <= '9':
			d = c - '0'
		default:
			return Decimal{}, fmt.Errorf("unexpected character %q", c)
		}

		if n >= cutoff {
			// n*base overflows
			return Decimal{}, errOverflow
		}
		n *= 10
		n1 := n + uint64(d)
		if n1 < n {
			return Decimal{}, errOverflow
		}
		n = n1
	}

	var exp int32
	if dot != -1 {
		exp = -int32(len(s) - dot - 1)
	}
	return Decimal{
		value: new(big.Int).SetUint64(n),
		exp:   exp,
	}, nil
}

// SizeAndScaleFromString gets the size and scale for the decimal value without needing to parse it.
func SizeAndScaleFromString(s string) (int32, int32) {
	switch s[0] {
	case '+', '-':
		s = s[1:]
	}
	totalLen := len(s)
	idx := strings.Index(s, ".")
	if idx == -1 {
		return int32(totalLen), 0
	}
	return int32(totalLen - 1), int32(totalLen - 1 - idx)
}

func NewFromMySQL(s []byte) (Decimal, error) {
	var original = s
	var neg bool

	if len(s) > 0 {
		switch s[0] {
		case '+':
			s = s[1:]
		case '-':
			neg = true
			s = s[1:]
		}
	}

	if len(s) == 0 {
		return Decimal{}, fmt.Errorf("can't convert %q to decimal: too short", original)
	}

	if len(s) <= 18 {
		dec, err := parseDecimal64(s)
		if err == nil {
			if neg {
				dec.value.Neg(dec.value)
			}
			return dec, nil
		}
		if err != errOverflow {
			return Decimal{}, fmt.Errorf("can't convert %s to decimal: %v", original, err)
		}
	}

	var fractional, integral []byte
	if pIndex := bytes.IndexByte(s, '.'); pIndex >= 0 {
		if bytes.IndexByte(s[pIndex+1:], '.') != -1 {
			return Decimal{}, fmt.Errorf("can't convert %s to decimal: too many .s", original)
		}
		if pIndex+1 < len(s) {
			integral = s[:pIndex]
			fractional = s[pIndex+1:]
		} else {
			integral = s[:pIndex]
		}
	} else {
		integral = s
	}

	// Check if the size of this bigint would fit in the limits
	// that MySQL has by default. To do that, we must convert the
	// length of our integral and fractional part to "mysql digits"
	myintg := myBigDigits(int32(len(integral)))
	myfrac := myBigDigits(int32(len(fractional)))
	if myintg > MyMaxBigDigits {
		return largestForm(MyMaxPrecision, 0, neg), nil
	}
	if myintg+myfrac > MyMaxBigDigits {
		fractional = fractional[:int((MyMaxBigDigits-myintg)*9)]
	}
	value, err := parseLargeDecimal(integral, fractional)
	if err != nil {
		return Decimal{}, err
	}
	if neg {
		value.Neg(value)
	}
	return Decimal{value: value, exp: -int32(len(fractional))}, nil
}

const ExponentLimit = 1024

// NewFromString returns a new Decimal from a string representation.
// Trailing zeroes are not trimmed.
// In case of an error, we still return the parsed value up to that
// point.
//
// Example:
//
//	d, err := NewFromString("-123.45")
//	d2, err := NewFromString(".0001")
//	d3, err := NewFromString("1.47000")
func NewFromString(s string) (d Decimal, err error) {
	maxLen := len(s)
	if maxLen > math.MaxInt32 {
		maxLen = math.MaxInt32
	}

	dotPos := -1
	expPos := -1
	i := 0
	var num bool
	var exp int64

	for i < maxLen {
		if !isSpace(s[i]) {
			break
		}
		i++
	}
next:
	for i < maxLen {
		switch {
		case s[i] == '-':
			// Negative sign is allowed at the start and at the start
			// of the exponent.
			if i != 0 && expPos == -1 && i != expPos+1 {
				break next
			}
		case s[i] >= '0' && s[i] <= '9':
			num = true
		case s[i] == '.':
			if dotPos == -1 && expPos == -1 {
				dotPos = i
			} else {
				break next
			}
		case s[i] == 'e' || s[i] == 'E':
			if expPos == -1 {
				expPos = i
				num = false
			} else {
				break next
			}
		default:
			break next
		}
		i++
	}

	// If we have a small total string or until the first dot,
	// we can fast parse it as an integer.
	var si string
	switch {
	case dotPos == -1 && expPos == -1:
		si = s[:i]
	case expPos == -1:
		si = s[:dotPos] + s[dotPos+1:i]
		exp -= int64(i - dotPos - 1)
	case dotPos == -1:
		si = s[:expPos]
	default:
		si = s[:dotPos] + s[dotPos+1:expPos]
		exp -= int64(expPos - dotPos - 1)
	}

	if len(si) <= 18 {
		var v int64
		v, err = fastparse.ParseInt64(si, 10)
		d.value = big.NewInt(v)
	} else {
		d.value = new(big.Int)
		d.value.SetString(si, 10)
	}

	var expOverflow bool
	if expPos != -1 {
		e, _ := fastparse.ParseInt64(s[expPos+1:i], 10)
		switch {
		case e > ExponentLimit:
			e = ExponentLimit
			expOverflow = true
		case e < -ExponentLimit:
			e = -ExponentLimit
			expOverflow = true
		}
		exp += e
	}

	d.exp = int32(exp)

	for i < maxLen {
		if !isSpace(s[i]) {
			break
		}
		i++
	}

	if !num || i < maxLen || expOverflow {
		err = fmt.Errorf("invalid decimal string: %q", s)
	}
	return d, err
}

func mulWW(x, y big.Word) (z1, z0 big.Word) {
	zz1, zz0 := bits.Mul(uint(x), uint(y))
	return big.Word(zz1), big.Word(zz0)
}

func mulAddWWW(x, y, c big.Word) (z1, z0 big.Word) {
	z1, zz0 := mulWW(x, y)
	if z0 = zz0 + c; z0 < zz0 {
		z1++
	}
	return z1, z0
}

func mulAddVWW(z, x []big.Word, y, r big.Word) (c big.Word) {
	c = r
	for i := range z {
		c, z[i] = mulAddWWW(x[i], y, c)
	}
	return c
}

func mulAddWW(z, x []big.Word, y, r big.Word) []big.Word {
	m := len(x)
	z = z[:m+1]
	z[m] = mulAddVWW(z[0:m], x, y, r)
	return z
}

func pow(x big.Word, n int) (p big.Word) {
	// n == sum of bi * 2**i, for 0 <= i < imax, and bi is 0 or 1
	// thus x**n == product of x**(2**i) for all i where bi == 1
	// (Russian Peasant Method for exponentiation)
	p = 1
	for n > 0 {
		if n&1 != 0 {
			p *= x
		}
		x *= x
		n >>= 1
	}
	return
}

func parseLargeDecimal(integral, fractional []byte) (*big.Int, error) {
	var (
		di = big.Word(0) // 0 <= di < b1**i < bn
		i  = 0           // 0 <= i < n
		// s is the largest possible size for a MySQL decimal; anything
		// that doesn't fit in s words won't make it to this func
		z = make([]big.Word, 0, s)
	)

	parseChunk := func(partial []byte) error {
		for _, ch := range partial {
			var d1 big.Word
			switch {
			case ch == '.':
				continue
			case '0' <= ch && ch <= '9':
				d1 = big.Word(ch - '0')
			default:
				return fmt.Errorf("unexpected character %q", ch)
			}

			// collect d1 in di
			di = di*b1 + d1
			i++

			// if di is "full", add it to the result
			if i == n {
				z = mulAddWW(z, z, bn, di)
				di = 0
				i = 0
			}
		}
		return nil
	}

	if err := parseChunk(integral); err != nil {
		return nil, err
	}
	if err := parseChunk(fractional); err != nil {
		return nil, err
	}
	if i > 0 {
		z = mulAddWW(z, z, pow(b1, i), di)
	}
	return new(big.Int).SetBits(z), nil
}

func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
