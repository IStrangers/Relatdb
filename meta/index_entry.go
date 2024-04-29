package meta

type IndexEntry interface {
	GetValues() []Value
	GetDesc() *IndexDesc
	GetLength() uint
	GetCompareEntry() IndexEntry
	GetDeleteCompareEntry() IndexEntry
	Compare(IndexEntry) int8
}

type BaseIndexEntry struct {
	Values       []Value
	Desc         *IndexDesc
	CompareEntry IndexEntry
	IsAllNull    bool
}

func CreateIndexEntry(values []Value, desc *IndexDesc) *BaseIndexEntry {
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
		self.CompareEntry = CreateNotLeafIndexEntry(self.Values[:], self.Desc)
	}
	return self.CompareEntry
}

func (self *BaseIndexEntry) GetDeleteCompareEntry() IndexEntry {
	return self
}

func (self *BaseIndexEntry) Compare(compareEntry IndexEntry) int8 {
	return self.innerCompare(compareEntry.GetCompareEntry())
}

type NotLeafIndexEntry struct {
	BaseIndexEntry
}

func CreateNotLeafIndexEntry(values []Value, desc *IndexDesc) *NotLeafIndexEntry {
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

func CreateClusterIndexEntry(values []Value, desc *IndexDesc) *ClusterIndexEntry {
	entry := &ClusterIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *ClusterIndexEntry) GetCompareEntry() IndexEntry {
	if self.CompareEntry == nil {
		primaryAttr := self.Desc.PrimaryAttr
		rowId := self.Values[primaryAttr.Index]
		desc := CreateIndexDesc([]*Attribute{primaryAttr})
		self.CompareEntry = CreateNotLeafIndexEntry([]Value{rowId}, desc)
	}
	return self.CompareEntry
}

func (self *ClusterIndexEntry) GetDeleteCompareEntry() IndexEntry {
	return self.GetCompareEntry()
}
