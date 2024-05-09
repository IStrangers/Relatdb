package meta

import "Relatdb/common"

type CompareType = uint

const (
	_ CompareType = iota
	COMPARE_EQUAL
	COMPARE_LOW
	COMPARE_UP
)

type Index interface {
	GetName() string
	GetFields() []*Field
	GetFlag() uint
	IsPrimary() bool
	IsUnique() bool
}

type BaseIndex struct {
	Name   string
	Fields []*Field
	FLag   uint
}

func (self *BaseIndex) GetName() string {
	return self.Name
}

func (self *BaseIndex) GetFields() []*Field {
	return self.Fields
}

func (self *BaseIndex) GetFlag() uint {
	return self.FLag
}

func (self *BaseIndex) IsPrimary() bool {
	return self.FLag&common.PRIMARY_KEY_FLAG != 0
}

func (self *BaseIndex) IsUnique() bool {
	return self.FLag&common.UNIQUE_KEY_FLAG != 0
}

func (self *BaseIndex) RecyclePageNo(pageNo uint) {

}
