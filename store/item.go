package store

import (
	"Relatdb/common"
	"Relatdb/meta"
)

func GetItemLength(indexEntry meta.IndexEntry) uint {
	return ITEM_POINTER_LENGTH + indexEntry.GetLength()
}

const ITEM_POINTER_LENGTH = 8
const ITEM_INT_LENGTH = ITEM_POINTER_LENGTH + 5

type ItemPointer struct {
	Offset      int
	TupleLength int
}

func NewItemPointer(offset int, tupleLength int) *ItemPointer {
	return &ItemPointer{
		Offset:      offset,
		TupleLength: tupleLength,
	}
}

type ItemData struct {
	Data   []byte
	Length int
}

func NewItemData(data []byte, length int) *ItemData {
	return &ItemData{
		Data:   data,
		Length: length,
	}
}

type Item struct {
	Pointer *ItemPointer
	Data    *ItemData
}

func NewItem(pointer *ItemPointer, data *ItemData) *Item {
	return &Item{
		Pointer: pointer,
		Data:    data,
	}
}

func IndexToItems(index meta.Index) []*Item {
	itemSize := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(1 + len(index.GetFields()) + 1)}, nil))
	itemName := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.StringValue(index.GetName())}, nil))
	itemFlag := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(index.GetFlag())}, nil))
	items := []*Item{itemSize, itemName, itemFlag}
	for _, field := range index.GetFields() {
		items = append(items, IndexEntryToItem(meta.NewIndexEntry(meta.FieldToValues(field), nil)))
	}
	return items
}

func ItemToIndexEntry(item *Item) meta.IndexEntry {
	buffer := common.NewBuffer(item.Data.Data)
	var values []meta.Value
	for buffer.Remaining() > 0 {
		var value meta.Value
		fieldType := meta.ValueType(buffer.ReadByte())
		switch fieldType {
		case meta.StringValueType:
			length := buffer.ReadInt()
			value = meta.StringValue(buffer.ReadBytes(uint(length)))
		case meta.Int64ValueType:
			value = meta.Int64Value(buffer.ReadInt64())
		case meta.IntValueType:
			value = meta.IntValue(buffer.ReadInt())
		case meta.NullValueType:
			value = meta.CONST_NULL_VALUE
		}
		values = append(values, value)
	}
	return meta.NewIndexEntry(values, nil)
}

func IndexEntryToItem(entry meta.IndexEntry) *Item {
	var data []byte
	for _, value := range entry.GetValues() {
		data = append(data, value.ToBytes()...)
	}
	itemPointer := NewItemPointer(-1, len(data))
	itemData := NewItemData(data, itemPointer.TupleLength)
	return NewItem(itemPointer, itemData)
}
