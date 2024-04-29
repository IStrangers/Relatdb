package bptree

import "Relatdb/index"

type BPTree struct {
	index.BaseIndex
	Root *BPNode
	Head *BPNode
}
