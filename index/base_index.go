package index

import "Relatdb/common"

type BaseIndex struct {
	FLag int
}

func (self *BaseIndex) IsPrimary() bool {
	return self.FLag&common.PRIMARY_KEY_FLAG != 0
}
func (self *BaseIndex) IsUnique() bool {
	return self.FLag&common.UNIQUE_KEY_FLAG != 0
}

func (self *BaseIndex) RecyclePageNo(pageNo uint) {

}
