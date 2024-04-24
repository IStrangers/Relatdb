package bptree

import (
	"Relatdb/index"
	"Relatdb/store"
)

type BPPage struct {
	Node              *BPNode
	NodeInitFreeSpace uint
	LeafInitFreeSpace uint
}

func (self *BPPage) getInitFreeSpace() uint {
	if self.Node.isLeaf {
		return self.LeafInitFreeSpace
	}
	return self.NodeInitFreeSpace
}

func (self *BPPage) remainFreeSpace() uint {
	return self.getInitFreeSpace() - self.getContentSize()
}

func (self *BPPage) getContentSize() uint {
	size := uint(0)
	for _, entry := range self.Node.Entries {
		size += index.GetItemLength(entry)
	}
	if !self.Node.isLeaf {
		size += uint(len(self.Node.Children) * store.ITEM_INT_LENGTH)
	}
	return size
}
