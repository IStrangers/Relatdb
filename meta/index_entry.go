package meta

type CompareEntry interface {
	GetCompareEntry() CompareEntry
	GetDeleteCompareEntry() CompareEntry
	Compare(CompareEntry) int8
}

type IndexEntry struct {
	Values       []Value
	Desc         *IndexDesc
	CompareEntry CompareEntry
	IsAllNull    bool
}

func CreateIndexEntry(values []Value, desc *IndexDesc) *IndexEntry {
	entry := &IndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *IndexEntry) Compare(compareEntry CompareEntry) int8 {
	return self.innerCompare(compareEntry.GetCompareEntry())
}

func (self *IndexEntry) innerCompare(entry CompareEntry) int8 {
	selfValues := self.Values
	selfValuesLength := len(selfValues)
	entryValues := entry.(*IndexEntry).Values
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

func (self *IndexEntry) GetLength() uint {
	length := uint(0)
	for _, value := range self.Values {
		length += value.GetLength()
	}
	return length
}

func (self *IndexEntry) GetCompareEntry() CompareEntry {
	if self.CompareEntry == nil {
		self.CompareEntry = CreateNotLeafIndexEntry(self.Values[:], self.Desc)
	}
	return self.CompareEntry
}

func (self *IndexEntry) GetDeleteCompareEntry() CompareEntry {
	return self
}

type NotLeafIndexEntry struct {
	IndexEntry
}

func CreateNotLeafIndexEntry(values []Value, desc *IndexDesc) *NotLeafIndexEntry {
	entry := &NotLeafIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *NotLeafIndexEntry) GetCompareEntry() CompareEntry {
	return self
}

type ClusterIndexEntry struct {
	IndexEntry
}

func CreateClusterIndexEntry(values []Value, desc *IndexDesc) *ClusterIndexEntry {
	entry := &ClusterIndexEntry{}
	entry.Values = values
	entry.Desc = desc
	return entry
}

func (self *ClusterIndexEntry) GetCompareEntry() CompareEntry {
	if self.CompareEntry == nil {
		primaryAttr := self.Desc.PrimaryAttr
		rowId := self.Values[primaryAttr.Index]
		desc := CreateIndexDesc([]*Attribute{primaryAttr})
		self.CompareEntry = CreateNotLeafIndexEntry([]Value{rowId}, desc)
	}
	return self.CompareEntry
}

func (self *ClusterIndexEntry) GetDeleteCompareEntry() CompareEntry {
	return self.GetCompareEntry()
}
