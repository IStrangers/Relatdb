package store

import (
	"Relatdb/meta"
	"Relatdb/utils"
)

func GetItemLength(indexEntry meta.IndexEntry) uint {
	return ITEM_POINTER_LENGTH + indexEntry.GetLength()
}

const ITEM_POINTER_LENGTH = 8

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
	Offset int
}

func NewItemData(data []byte, length int, offset int) *ItemData {
	return &ItemData{
		Data:   data,
		Length: length,
		Offset: offset,
	}
}

const ITEM_INT_LENGTH = ITEM_POINTER_LENGTH + 5

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

func ItemToIndexEntry(item *Item) meta.IndexEntry {
	bytesReader := utils.NewBytesReader(item.Data.Data)
	var values []meta.Value
	for bytesReader.Remaining() > 0 {
		fieldType := bytesReader.ReadByte()
		println(fieldType)
		break
	}
	return meta.NewIndexEntry(values, nil)
}

func IndexEntryToItem(entry meta.IndexEntry) *Item {

	return NewItem(nil, nil)
}
