package meta

type IndexEntry struct {
	Values       []Value
	Desc         *IndexDesc
	CompareEntry *IndexEntry
	IsAllNull    bool
}

func (self *IndexEntry) Compare(indexEntry *IndexEntry) int8 {
	return 0
}

func (self *IndexEntry) GetLength() uint {
	length := uint(0)
	for _, value := range self.Values {
		length += value.GetLength()
	}
	return length
}

func (self *IndexEntry) GetCompareEntry() *IndexEntry {
	if self.CompareEntry == nil {

	}
	return self.CompareEntry
}
