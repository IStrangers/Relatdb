package store

import "Relatdb/meta"

func GetItemLength(indexEntry meta.IndexEntry) uint {
	return ITEM_POINTER_LENGTH + indexEntry.GetLength()
}

const ITEM_POINTER_LENGTH = 8

type ItemPointer struct {
	Offset      int
	TupleLength int
}

func CreateItemPointer(offset int, tupleLength int) *ItemPointer {
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

func CreateItemData(data []byte, length int, offset int) *ItemData {
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

func CreateItem(pointer *ItemPointer, data *ItemData) *Item {
	return &Item{
		Pointer: pointer,
		Data:    data,
	}
}
