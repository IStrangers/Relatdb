package bptree

import (
	"Relatdb/store"
)

type BPPage struct {
	*store.Page
	PageNo            uint
	Node              *BPNode
	NodeInitFreeSpace uint
	LeafInitFreeSpace uint
}

func NewBPPage(node *BPNode) *BPPage {
	bpPage := &BPPage{
		Page:   store.NewPage(),
		PageNo: 1,
		Node:   node,
	}
	nodeInitFreeSpace := bpPage.Length - store.DEFAULT_SPECIAL_POINT_LENGTH - store.PAGE_HEADER_SIZE - store.ITEM_INT_LENGTH*7
	leafInitFreeSpace := bpPage.Length - store.DEFAULT_SPECIAL_POINT_LENGTH - store.PAGE_HEADER_SIZE - store.ITEM_INT_LENGTH*6
	bpPage.NodeInitFreeSpace = nodeInitFreeSpace
	bpPage.LeafInitFreeSpace = leafInitFreeSpace
	return bpPage
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
		size += store.GetItemLength(entry)
	}
	if !self.Node.isLeaf {
		size += uint(len(self.Node.Children) * store.ITEM_INT_LENGTH)
	}
	return size
}
