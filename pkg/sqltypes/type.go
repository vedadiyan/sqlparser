/*
Copyright 2019 The Vitess Authors.

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

	querypb "github.com/vedadiyan/sqlparser/pkg/query"
)

type Type = querypb.Type

// This file provides wrappers and support
// functions for querypb.Type.

// These bit flags can be used to query on the
// common properties of types.
const (
	flagIsIntegral = int(querypb.Flag_ISINTEGRAL)
	flagIsUnsigned = int(querypb.Flag_ISUNSIGNED)
	flagIsFloat    = int(querypb.Flag_ISFLOAT)
	flagIsQuoted   = int(querypb.Flag_ISQUOTED)
	flagIsText     = int(querypb.Flag_ISTEXT)
	flagIsBinary   = int(querypb.Flag_ISBINARY)
)

const (
	TimestampFormat           = "2006-01-02 15:04:05"
	TimestampFormatPrecision3 = "2006-01-02 15:04:05.000"
	TimestampFormatPrecision6 = "2006-01-02 15:04:05.000000"
)

// IsIntegral returns true if querypb.Type is an integral
// (signed/unsigned) that can be represented using
// up to 64 binary bits.
// If you have a Value object, use its member function.
func IsIntegral(t querypb.Type) bool {
	return int(t)&flagIsIntegral == flagIsIntegral
}

// IsSigned returns true if querypb.Type is a signed integral.
// If you have a Value object, use its member function.
func IsSigned(t querypb.Type) bool {
	return int(t)&(flagIsIntegral|flagIsUnsigned) == flagIsIntegral
}

// IsUnsigned returns true if querypb.Type is an unsigned integral.
// Caution: this is not the same as !IsSigned.
// If you have a Value object, use its member function.
func IsUnsigned(t querypb.Type) bool {
	return int(t)&(flagIsIntegral|flagIsUnsigned) == flagIsIntegral|flagIsUnsigned
}

// IsFloat returns true is querypb.Type is a floating point.
// If you have a Value object, use its member function.
func IsFloat(t querypb.Type) bool {
	return int(t)&flagIsFloat == flagIsFloat
}

// IsDecimal returns true is querypb.Type is a decimal.
// If you have a Value object, use its member function.
func IsDecimal(t querypb.Type) bool {
	return t == Decimal
}

// IsQuoted returns true if querypb.Type is a quoted text or binary.
// If you have a Value object, use its member function.
func IsQuoted(t querypb.Type) bool {
	return (int(t)&flagIsQuoted == flagIsQuoted) && t != Bit
}

// IsText returns true if querypb.Type is a text.
// If you have a Value object, use its member function.
func IsText(t querypb.Type) bool {
	return int(t)&flagIsText == flagIsText
}

func IsTextOrBinary(t querypb.Type) bool {
	return int(t)&flagIsText == flagIsText || int(t)&flagIsBinary == flagIsBinary
}

// IsBinary returns true if querypb.Type is a binary.
// If you have a Value object, use its member function.
func IsBinary(t querypb.Type) bool {
	return int(t)&flagIsBinary == flagIsBinary
}

// IsNumber returns true if the type is any type of number.
func IsNumber(t querypb.Type) bool {
	return IsIntegral(t) || IsFloat(t) || t == Decimal
}

// IsDateOrTime returns true if the type represents a date and/or time.
func IsDateOrTime(t querypb.Type) bool {
	return t == Datetime || t == Date || t == Timestamp || t == Time
}

// IsDate returns true if the type has a date component
func IsDate(t querypb.Type) bool {
	return t == Datetime || t == Date || t == Timestamp
}

// IsNull returns true if the type is NULL type
func IsNull(t querypb.Type) bool {
	return t == Null
}

// IsEnum returns true if the type is Enum type
func IsEnum(t querypb.Type) bool {
	return t == Enum
}

// IsSet returns true if the type is Set type
func IsSet(t querypb.Type) bool {
	return t == Set
}

// Vitess data types. These are idiomatically named synonyms for the querypb.Type values.
// Although these constants are interchangeable, they should be treated as different from querypb.Type.
// Use the synonyms only to refer to the type in Value. For proto variables, use the querypb.Type constants instead.
// The following is a complete listing of types that match each classification function in this API:
//
//	IsSigned(): INT8, INT16, INT24, INT32, INT64
//	IsFloat(): FLOAT32, FLOAT64
//	IsUnsigned(): UINT8, UINT16, UINT24, UINT32, UINT64, YEAR
//	IsIntegral(): INT8, UINT8, INT16, UINT16, INT24, UINT24, INT32, UINT32, INT64, UINT64, YEAR
//	IsText(): TEXT, VARCHAR, CHAR, HEXNUM, HEXVAL, BITNUM
//	IsNumber(): INT8, UINT8, INT16, UINT16, INT24, UINT24, INT32, UINT32, INT64, UINT64, FLOAT32, FLOAT64, YEAR, DECIMAL
//	IsQuoted(): TIMESTAMP, DATE, TIME, DATETIME, TEXT, BLOB, VARCHAR, VARBINARY, CHAR, BINARY, ENUM, SET, GEOMETRY, JSON
//	IsBinary(): BLOB, VARBINARY, BINARY
//	IsDate(): TIMESTAMP, DATE, TIME, DATETIME
//	IsNull(): NULL_TYPE
//
// TODO(sougou): provide a categorization function
// that returns enums, which will allow for cleaner
// switch statements for those who want to cover types
// by their category.
const (
	Unknown    = querypb.Type(-1)
	Null       = querypb.Type_NULL_TYPE
	Int8       = querypb.Type_INT8
	Uint8      = querypb.Type_UINT8
	Int16      = querypb.Type_INT16
	Uint16     = querypb.Type_UINT16
	Int24      = querypb.Type_INT24
	Uint24     = querypb.Type_UINT24
	Int32      = querypb.Type_INT32
	Uint32     = querypb.Type_UINT32
	Int64      = querypb.Type_INT64
	Uint64     = querypb.Type_UINT64
	Float32    = querypb.Type_FLOAT32
	Float64    = querypb.Type_FLOAT64
	Timestamp  = querypb.Type_TIMESTAMP
	Date       = querypb.Type_DATE
	Time       = querypb.Type_TIME
	Datetime   = querypb.Type_DATETIME
	Year       = querypb.Type_YEAR
	Decimal    = querypb.Type_DECIMAL
	Text       = querypb.Type_TEXT
	Blob       = querypb.Type_BLOB
	VarChar    = querypb.Type_VARCHAR
	VarBinary  = querypb.Type_VARBINARY
	Char       = querypb.Type_CHAR
	Binary     = querypb.Type_BINARY
	Bit        = querypb.Type_BIT
	Enum       = querypb.Type_ENUM
	Set        = querypb.Type_SET
	Geometry   = querypb.Type_GEOMETRY
	TypeJSON   = querypb.Type_JSON
	Expression = querypb.Type_EXPRESSION
	HexNum     = querypb.Type_HEXNUM
	HexVal     = querypb.Type_HEXVAL
	Tuple      = querypb.Type_TUPLE
	BitNum     = querypb.Type_BITNUM
	Vector     = querypb.Type_VECTOR
)

// bit-shift the mysql flags by two byte so we
// can merge them with the mysql or vitess types.
const (
	mysqlUnsigned = 32
	mysqlBinary   = 128
	mysqlEnum     = 256
	mysqlSet      = 2048
)

// If you add to this map, make sure you add a test case
// in tabletserver/endtoend.
var mysqlToType = map[byte]querypb.Type{
	0:   Decimal,
	1:   Int8,
	2:   Int16,
	3:   Int32,
	4:   Float32,
	5:   Float64,
	6:   Null,
	7:   Timestamp,
	8:   Int64,
	9:   Int24,
	10:  Date,
	11:  Time,
	12:  Datetime,
	13:  Year,
	15:  VarChar,
	16:  Bit,
	17:  Timestamp,
	18:  Datetime,
	19:  Time,
	242: Vector,
	245: TypeJSON,
	246: Decimal,
	247: Enum,
	248: Set,
	249: Text,
	250: Text,
	251: Text,
	252: Text,
	253: VarChar,
	254: Char,
	255: Geometry,
}

// modifyType modifies the vitess type based on the
// mysql flag. The function checks specific flags based
// on the type. This allows us to ignore stray flags
// that MySQL occasionally sets.
func modifyType(typ querypb.Type, flags int64) querypb.Type {
	switch typ {
	case Int8:
		if flags&mysqlUnsigned != 0 {
			return Uint8
		}
	case Int16:
		if flags&mysqlUnsigned != 0 {
			return Uint16
		}
	case Int32:
		if flags&mysqlUnsigned != 0 {
			return Uint32
		}
	case Int64:
		if flags&mysqlUnsigned != 0 {
			return Uint64
		}
	case Int24:
		if flags&mysqlUnsigned != 0 {
			return Uint24
		}
	case Text:
		if flags&mysqlBinary != 0 {
			return Blob
		}
	case VarChar:
		if flags&mysqlBinary != 0 {
			return VarBinary
		}
	case Char:
		if flags&mysqlBinary != 0 {
			return Binary
		}
		if flags&mysqlEnum != 0 {
			return Enum
		}
		if flags&mysqlSet != 0 {
			return Set
		}
	}
	return typ
}

// MySQLToType computes the vitess type from mysql type and flags.
func MySQLToType(mysqlType byte, flags int64) (typ querypb.Type, err error) {
	result, ok := mysqlToType[mysqlType]
	if !ok {
		return 0, fmt.Errorf("unsupported type: %d", mysqlType)
	}
	return modifyType(result, flags), nil
}

// AreTypesEquivalent returns whether two types are equivalent.
func AreTypesEquivalent(mysqlTypeFromBinlog, mysqlTypeFromSchema querypb.Type) bool {
	return (mysqlTypeFromBinlog == mysqlTypeFromSchema) ||
		(mysqlTypeFromBinlog == VarChar && mysqlTypeFromSchema == VarBinary) ||
		// Binlog only has base type. But doesn't have per-column-flags to differentiate
		// various logical types. For Binary, Enum, Set types, binlog only returns Char
		// as data type.
		(mysqlTypeFromBinlog == Char && mysqlTypeFromSchema == Binary) ||
		(mysqlTypeFromBinlog == Char && mysqlTypeFromSchema == Enum) ||
		(mysqlTypeFromBinlog == Char && mysqlTypeFromSchema == Set) ||
		(mysqlTypeFromBinlog == Text && mysqlTypeFromSchema == Blob) ||
		(mysqlTypeFromBinlog == Int8 && mysqlTypeFromSchema == Uint8) ||
		(mysqlTypeFromBinlog == Int16 && mysqlTypeFromSchema == Uint16) ||
		(mysqlTypeFromBinlog == Int24 && mysqlTypeFromSchema == Uint24) ||
		(mysqlTypeFromBinlog == Int32 && mysqlTypeFromSchema == Uint32) ||
		(mysqlTypeFromBinlog == Int64 && mysqlTypeFromSchema == Uint64)
}

// typeToMySQL is the reverse of mysqlToType.
var typeToMySQL = map[querypb.Type]struct {
	typ   byte
	flags int64
}{
	Int8:      {typ: 1},
	Uint8:     {typ: 1, flags: mysqlUnsigned},
	Int16:     {typ: 2},
	Uint16:    {typ: 2, flags: mysqlUnsigned},
	Int32:     {typ: 3},
	Uint32:    {typ: 3, flags: mysqlUnsigned},
	Float32:   {typ: 4},
	Float64:   {typ: 5},
	Null:      {typ: 6, flags: mysqlBinary},
	Timestamp: {typ: 7},
	Int64:     {typ: 8},
	Uint64:    {typ: 8, flags: mysqlUnsigned},
	Int24:     {typ: 9},
	Uint24:    {typ: 9, flags: mysqlUnsigned},
	Date:      {typ: 10, flags: mysqlBinary},
	Time:      {typ: 11, flags: mysqlBinary},
	Datetime:  {typ: 12, flags: mysqlBinary},
	Year:      {typ: 13, flags: mysqlUnsigned},
	Bit:       {typ: 16, flags: mysqlUnsigned},
	Vector:    {typ: 242},
	TypeJSON:  {typ: 245},
	Decimal:   {typ: 246},
	Text:      {typ: 252},
	Blob:      {typ: 252, flags: mysqlBinary},
	BitNum:    {typ: 253, flags: mysqlBinary},
	HexNum:    {typ: 253, flags: mysqlBinary},
	HexVal:    {typ: 253, flags: mysqlBinary},
	VarChar:   {typ: 253},
	VarBinary: {typ: 253, flags: mysqlBinary},
	Char:      {typ: 254},
	Binary:    {typ: 254, flags: mysqlBinary},
	Enum:      {typ: 254, flags: mysqlEnum},
	Set:       {typ: 254, flags: mysqlSet},
	Geometry:  {typ: 255},
}

// TypeToMySQL returns the equivalent mysql type and flag for a vitess type.
func TypeToMySQL(typ querypb.Type) (mysqlType byte, flags int64) {
	val := typeToMySQL[typ]
	return val.typ, val.flags
}
