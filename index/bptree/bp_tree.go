package bptree

import (
	"Relatdb/meta"
)

type BPTree struct {
	meta.BaseIndex
	Root *BPNode
	Head *BPNode
}

func NewBPTree(name string, fields []*meta.Field, flag uint) *BPTree {
	bpTree := &BPTree{}
	bpTree.Name = name
	bpTree.Fields = fields
	bpTree.FLag = flag
	bpTree.Root = NewBPNode(bpTree, true, true)
	bpTree.Head = bpTree.Root
	return bpTree
}

func (self *BPTree) Insert(entry meta.IndexEntry) {
	//self.Root.Insert(entry, self, self.IsUnique())
}
