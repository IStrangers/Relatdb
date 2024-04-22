package index

import (
	"Relatdb/meta"
	"Relatdb/store"
)

type CompareType = uint

const (
	_ CompareType = iota
	COMPARE_EQUAL
	COMPARE_LOW
	COMPARE_UP
)

func GetItemLength(indexEntry *meta.IndexEntry) uint {
	return store.ITEM_POINTER_LENGTH + indexEntry.GetLength()
}
