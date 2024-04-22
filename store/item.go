package store

const ITEM_POINTER_LENGTH = 8

type ItemPointer struct {
	Offset      uint
	TupleLength uint
}

type ItemData struct {
	Data   []byte
	Length uint
	Offset uint
}

const ITEM_INT_LENGTH = ITEM_POINTER_LENGTH + 5

type Item struct {
	Pointer *ItemPointer
	Data    *ItemData
}
