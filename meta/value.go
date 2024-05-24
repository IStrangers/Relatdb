package meta

import (
	"Relatdb/common"
	"fmt"
	"strconv"
	"strings"
)

type ValueType byte

const (
	_ ValueType = iota
	StringValueType
	Int64ValueType
	IntValueType
	NullValueType
)

var (
	CONST_NULL_VALUE = &NullValue{}
)

type Value interface {
	GetType() ValueType
	ToString() string
	ToInt() int
	ToInt64() int64
	ToBytes() []byte
	ToValueBytes() []byte
	GetLength() uint
	Compare(Value) int
}

type StringValue string

func (v StringValue) GetType() ValueType {
	return StringValueType
}

func (self StringValue) ToString() string {
	return string(self)
}

func (self StringValue) ToInt() int {
	return int(self.ToInt64())
}

func (self StringValue) ToInt64() int64 {
	total := int64(0)
	for _, char := range self {
		total += int64(char)
	}
	return total
}

func (self StringValue) ToBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength())
	buffer.WriteByte(byte(self.GetType()))
	buffer.WriteInt(len(self))
	buffer.WriteString(string(self))
	return buffer.Data
}

func (self StringValue) ToValueBytes() []byte {
	return []byte(self)
}

func (self StringValue) GetLength() uint {
	return 4 + 1 + uint(len(string(self)))
}

func (self StringValue) Compare(value Value) int {
	return strings.Compare(self.ToString(), value.ToString())
}

type Int64Value int64

func (v Int64Value) GetType() ValueType {
	return Int64ValueType
}

func (self Int64Value) ToString() string {
	return strconv.FormatInt(self.ToInt64(), 10)
}

func (self Int64Value) ToInt() int {
	return int(self)
}

func (self Int64Value) ToInt64() int64 {
	return int64(self)
}

func (self Int64Value) ToBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength())
	buffer.WriteByte(byte(self.GetType()))
	buffer.WriteInt64(int64(self))
	return buffer.Data
}

func (self Int64Value) ToValueBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength() - 1)
	buffer.WriteInt64(int64(self))
	return buffer.Data
}

func (self Int64Value) GetLength() uint {
	return 8 + 1
}

func (self Int64Value) Compare(value Value) int {
	if self.ToInt64() < value.ToInt64() {
		return -1
	}
	if self.ToInt64() > value.ToInt64() {
		return 1
	}
	return 0
}

type IntValue int

func (v IntValue) GetType() ValueType {
	return IntValueType
}

func (self IntValue) ToString() string {
	return strconv.FormatInt(self.ToInt64(), 10)
}

func (self IntValue) ToInt() int {
	return int(self)
}

func (self IntValue) ToInt64() int64 {
	return int64(self)
}

func (self IntValue) ToBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength())
	buffer.WriteByte(byte(self.GetType()))
	buffer.WriteInt(int(self))
	return buffer.Data
}

func (self IntValue) ToValueBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength() - 1)
	buffer.WriteInt(int(self))
	return buffer.Data
}

func (self IntValue) GetLength() uint {
	return 4 + 1
}

func (self IntValue) Compare(value Value) int {
	if self.ToInt64() < value.ToInt64() {
		return -1
	}
	if self.ToInt64() > value.ToInt64() {
		return 1
	}
	return 0
}

type NullValue struct{}

func (self NullValue) GetType() ValueType {
	return NullValueType
}

func (self NullValue) ToString() string {
	return ""
}

func (self NullValue) ToInt() int {
	return 0
}

func (self NullValue) ToInt64() int64 {
	return 0
}

func (self NullValue) ToBytes() []byte {
	buffer := common.NewBufferBySize(self.GetLength())
	buffer.WriteByte(byte(self.GetType()))
	return buffer.Data
}

func (self NullValue) ToValueBytes() []byte {
	return []byte{}
}

func (self NullValue) GetLength() uint {
	return 1
}

func (self NullValue) Compare(value Value) int {
	if value.GetType() == NullValueType {
		return 0
	}
	return -1
}

func ToValue(val any) Value {
	if val == nil {
		return CONST_NULL_VALUE
	}
	switch v := val.(type) {
	case string:
		return StringValue(v)
	case int:
		return IntValue(v)
	case uint:
		return IntValue(v)
	case int64:
		return Int64Value(v)
	case uint64:
		return Int64Value(v)
	case Value:
		return v
	default:
		panic(fmt.Sprintf("TOValue:unknown value type: %T", v))
	}
}

func FieldToValues(field *Field) []Value {
	return []Value{
		IntValue(field.Index),
		StringValue(field.Name),
		IntValue(field.Type),
		IntValue(field.Flag),
		field.DefaultValue,
		StringValue(field.Comment),
	}
}
