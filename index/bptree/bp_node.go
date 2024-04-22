package bptree

import (
	"Relatdb/index"
	"Relatdb/meta"
	"Relatdb/store"
	"errors"
)

type BPNode struct {
	OwnerTree *BPTree
	IsRoot    bool
	isLeaf    bool
	Parent    *BPNode
	Prev      *BPNode
	Next      *BPNode
	Entries   []*meta.IndexEntry
	Children  []*BPNode
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
		bpNode.Children = make([]*BPNode, 3)
	}
	return bpNode
}

func (self *BPNode) getBorrowKeyLength(key *meta.IndexEntry) uint {
	itemLength := index.GetItemLength(key)
	if !self.isLeaf {
		itemLength += store.ITEM_INT_LENGTH
	}
	return itemLength
}

func (self *BPNode) innerCheckExist(key *meta.IndexEntry) bool {
	for _, entry := range self.Entries {
		if key.GetCompareEntry().Compare(entry) == 0 {
			return true
		}
	}
	return false
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
		return self.Children[0].Get(key, compareType)
	} else if lastEntry := self.Entries[len(self.Entries)-1]; key.Compare(lastEntry) >= 0 {
		return self.Children[len(self.Children)-1].Get(key, compareType)
	} else {
		for i, entry := range self.Entries {
			if key.Compare(entry) > -1 {
				continue
			}
			return self.Children[i].Get(key, compareType)
		}
	}
	return nil
}

func (self *BPNode) Insert(key *meta.IndexEntry, bpTree *BPTree, isUnique bool) error {
	if self.getBorrowKeyLength(key) > self.Page.getInitFreeSpace()/3 {
		return errors.New("entry size must <= Max/3")
	} else if isUnique && self.innerCheckExist(key) {
		return errors.New("Duplicated Key error")
	}

	return nil
}

func (self *BPNode) Remove(key *meta.IndexEntry, bpTree *BPTree) {

}
