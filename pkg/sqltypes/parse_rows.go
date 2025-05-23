/*
Copyright 2023 The Vitess Authors.

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

package sqltypes

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"text/scanner"

	querypb "github.com/vedadiyan/sqlparser/pkg/query"
)

// ParseRows parses the output generated by fmt.Sprintf("#v", rows), and reifies the original []sqltypes.Row
// NOTE: This is not meant for production use!
func ParseRows(input string) ([]Row, error) {
	type state int
	const (
		stInvalid state = iota
		stInit
		stBeginRow
		stInRow
		stInValue0
		stInValue1
		stInValue2
	)

	var (
		scan   scanner.Scanner
		result []Row
		row    Row
		vtype  int32
		st     = stInit
	)

	scan.Init(strings.NewReader(input))

	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		var next state

		switch st {
		case stInit:
			if tok == '[' {
				next = stBeginRow
			}
		case stBeginRow:
			switch tok {
			case '[':
				next = stInRow
			case ']':
				return result, nil
			}
		case stInRow:
			switch tok {
			case ']':
				result = append(result, row)
				row = nil
				next = stBeginRow
			case scanner.Ident:
				ident := scan.TokenText()

				if ident == "NULL" {
					row = append(row, NULL)
					continue
				}

				var ok bool
				vtype, ok = querypb.Type_value[ident]
				if !ok {
					return nil, fmt.Errorf("unknown SQL type %q at %s", ident, scan.Position)
				}
				next = stInValue0
			}
		case stInValue0:
			if tok == '(' {
				next = stInValue1
			}
		case stInValue1:
			literal := scan.TokenText()
			switch tok {
			case scanner.String:
				var err error
				literal, err = strconv.Unquote(literal)
				if err != nil {
					return nil, fmt.Errorf("failed to parse literal string at %s: %w", scan.Position, err)
				}
				fallthrough
			case scanner.Int, scanner.Float:
				row = append(row, MakeTrusted(Type(vtype), []byte(literal)))
				next = stInValue2
			}
		case stInValue2:
			if tok == ')' {
				next = stInRow
			}
		}
		if next == stInvalid {
			return nil, fmt.Errorf("unexpected token '%s' at %s", scan.TokenText(), scan.Position)
		}
		st = next
	}
	return nil, io.ErrUnexpectedEOF
}

type RowMismatchError struct {
	err       error
	want, got []Row
}

func (e *RowMismatchError) Error() string {
	return fmt.Sprintf("results differ: %v\n\twant: %v\n\tgot:  %v", e.err, e.want, e.got)
}

func RowEqual(want, got Row) bool {
	return slices.EqualFunc(want, got, func(a, b Value) bool {
		return a.Equal(b)
	})
}

func RowsEquals(want, got []Row) error {
	if len(want) != len(got) {
		return &RowMismatchError{
			err:  fmt.Errorf("expected %d rows in result, got %d", len(want), len(got)),
			want: want,
			got:  got,
		}
	}

	var matched = make([]bool, len(want))
	for _, aa := range want {
		var ok bool
		for i, bb := range got {
			if matched[i] {
				continue
			}
			if RowEqual(aa, bb) {
				matched[i] = true
				ok = true
				break
			}
		}
		if !ok {
			return &RowMismatchError{
				err:  fmt.Errorf("row %v is missing from result", aa),
				want: want,
				got:  got,
			}
		}
	}
	for _, m := range matched {
		if !m {
			return fmt.Errorf("not all elements matched")
		}
	}
	return nil
}

func RowsEqualsStr(wantStr string, got []Row) error {
	want, err := ParseRows(wantStr)
	if err != nil {
		return fmt.Errorf("malformed row assertion: %w", err)
	}
	return RowsEquals(want, got)
}
