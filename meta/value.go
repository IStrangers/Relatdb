package meta

import (
	"strconv"
	"strings"
)

type Value interface {
	ToString() string
	ToInt() int
	ToInt64() int64
	GetLength() uint
	Compare(Value) int
}

type StringValue string

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

func (self StringValue) GetLength() uint {
	return uint(len(string(self)))
}

func (self StringValue) Compare(value Value) int {
	return strings.Compare(self.ToString(), value.ToString())
}

type Int64Value int64

func (self Int64Value) ToString() string {
	return strconv.FormatInt(self.ToInt64(), 10)
}

func (self Int64Value) ToInt() int {
	return int(self)
}

func (self Int64Value) ToInt64() int64 {
	return int64(self)
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

func (self IntValue) ToString() string {
	return strconv.FormatInt(self.ToInt64(), 10)
}

func (self IntValue) ToInt() int {
	return int(self)
}

func (self IntValue) ToInt64() int64 {
	return int64(self)
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
