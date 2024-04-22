package bptree

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
