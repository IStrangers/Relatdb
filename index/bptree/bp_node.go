package bptree

import (
	"Relatdb/index"
	"Relatdb/meta"
)

type BPNode struct {
	OwnerTree *BPTree
	IsRoot    bool
	isLeaf    bool
	Parent    *BPNode
	prev      *BPNode
	next      *BPNode
	Entries   []*meta.IndexEntry
	children  []*BPNode
	Page      *BPPage
}

func CreateBPNode(ownerTree *BPTree, isRoot bool, isLeaf bool) *BPNode {
	bpNode := &BPNode{
		OwnerTree: ownerTree,
		IsRoot:    isRoot,
		isLeaf:    isLeaf,
		Entries:   []*meta.IndexEntry{},
	}
	if !isLeaf {
		bpNode.children = make([]*BPNode, 3)
	}
	return bpNode
}

func (self *BPNode) Get(key *meta.IndexEntry, compareType index.CompareType) *BPPosition {
	if self.isLeaf {
		if compareType == index.COMPARE_EQUAL {
			for i, entry := range self.Entries {
				if key.Compare(entry) != 0 {
					continue
				}
				return CreateBPPosition(nil, uint(i), self)
			}
		} else if compareType == index.COMPARE_LOW {
			return CreateBPPosition(nil, 0, self)
		} else {
			return CreateBPPosition(nil, uint(len(self.Entries)-1), self)
		}
	}
	if firstEntry := self.Entries[0]; key.Compare(firstEntry) < 0 {
		return self.children[0].Get(key, compareType)
	} else if lastEntry := self.Entries[len(self.Entries)-1]; key.Compare(lastEntry) >= 0 {
		return self.children[len(self.children)-1].Get(key, compareType)
	} else {
		for i, entry := range self.Entries {
			if key.Compare(entry) > -1 {
				continue
			}
			return self.children[i].Get(key, compareType)
		}
	}
	return nil
}
