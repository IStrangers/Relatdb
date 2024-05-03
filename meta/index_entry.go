package meta

type IndexEntry interface {
	GetValues() []Value
	GetDesc() *IndexDesc
	GetLength() uint
	GetCompareEntry() IndexEntry
	GetDeleteCompareEntry() IndexEntry
	CompareEntry(IndexEntry) int8
	CompareDeleteEntry(IndexEntry) int8
}

type BaseIndexEntry struct {
	Values     []Value
	Desc       *IndexDesc
	IndexEntry IndexEntry
	IsAllNull  bool
}

func NewIndexEntry(values []Value, desc *IndexDesc) *BaseIndexEntry {
	entry := &BaseIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *BaseIndexEntry) innerCompare(entry IndexEntry) int8 {
	selfValues := self.Values
	selfValuesLength := len(selfValues)
	entryValues := entry.GetValues()
	entryValuesLength := len(entryValues)
	minLength := min(selfValuesLength, entryValuesLength)
	for i := 0; i < minLength; i++ {
		selfValue := selfValues[i]
		entryValue := entryValues[i]
		if selfValue == nil && entryValue == nil {
			continue
		}
		if selfValue == nil {
			return -1
		}
		if entryValue == nil {
			return 1
		}
		if comp := selfValue.Compare(entryValue); comp != 0 {
			return comp
		}
	}
	if selfValuesLength < entryValuesLength {
		return -1
	}
	if selfValuesLength > entryValuesLength {
		return 1
	}
	return 0
}

func (self *BaseIndexEntry) GetValues() []Value {
	return self.Values
}

func (self *BaseIndexEntry) GetDesc() *IndexDesc {
	return self.Desc
}

func (self *BaseIndexEntry) GetLength() uint {
	length := uint(0)
	for _, value := range self.Values {
		length += value.GetLength()
	}
	return length
}

func (self *BaseIndexEntry) GetCompareEntry() IndexEntry {
	if self.CompareEntry == nil {
		self.IndexEntry = NewNotLeafIndexEntry(self.Values[:], self.Desc)
	}
	return self.IndexEntry
}

func (self *BaseIndexEntry) GetDeleteCompareEntry() IndexEntry {
	return self
}

func (self *BaseIndexEntry) CompareEntry(compareEntry IndexEntry) int8 {
	return self.innerCompare(compareEntry.GetCompareEntry())
}

func (self *BaseIndexEntry) CompareDeleteEntry(compareEntry IndexEntry) int8 {
	return self.innerCompare(compareEntry.GetDeleteCompareEntry())
}

type NotLeafIndexEntry struct {
	BaseIndexEntry
}

func NewNotLeafIndexEntry(values []Value, desc *IndexDesc) *NotLeafIndexEntry {
	entry := &NotLeafIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *NotLeafIndexEntry) GetCompareEntry() IndexEntry {
	return self
}

type ClusterIndexEntry struct {
	BaseIndexEntry
}

func NewClusterIndexEntry(values []Value, desc *IndexDesc) *ClusterIndexEntry {
	entry := &ClusterIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *ClusterIndexEntry) GetCompareEntry() IndexEntry {
	if self.CompareEntry == nil {
		primaryAttr := self.Desc.PrimaryFiled
		rowId := self.Values[primaryAttr.Index]
		desc := NewIndexDesc([]*Field{primaryAttr})
		self.IndexEntry = NewNotLeafIndexEntry([]Value{rowId}, desc)
	}
	return self.IndexEntry
}

func (self *ClusterIndexEntry) GetDeleteCompareEntry() IndexEntry {
	return self.GetCompareEntry()
}
