package bptree

import "Relatdb/meta"

type BPPosition struct {
	SearchEntry *meta.IndexEntry
	Position    uint
	Node        *BPNode
}

func NewBPPosition(searchEntry *meta.IndexEntry, position uint, node *BPNode) *BPPosition {
	return &BPPosition{
		SearchEntry: searchEntry,
		Position:    position,
		Node:        node,
	}
}
