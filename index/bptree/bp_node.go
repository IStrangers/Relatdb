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

func (self *BPNode) isLeafSplit(key *meta.IndexEntry) bool {
	return self.Page.remainFreeSpace() < index.GetItemLength(key)
}

func (self *BPNode) innerInsert(key *meta.IndexEntry) {
	insertIndex := len(self.Entries)
	if insertIndex == 0 {
		self.Entries = append(self.Entries, key)
		return
	}
	for i, entry := range self.Entries {
		if key.Compare(entry) == 0 || key.Compare(entry) < 0 {
			insertIndex = i
			break
		}
	}
	self.Entries = append(self.Entries[:insertIndex+1], self.Entries[insertIndex:]...)
	self.Entries[insertIndex] = key
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
			if key.Compare(entry) != -1 {
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
		return errors.New("duplicated Key error")
	}
	if self.isLeaf {
		if !self.isLeafSplit(key) {
			self.innerInsert(key)
			return nil
		}
		left := CreateBPNode(self.OwnerTree, false, true)
		right := CreateBPNode(self.OwnerTree, false, true)
		if self.Prev != nil {
			self.Prev.Next = left
			left.Prev = self.Prev
		} else {
			bpTree.setHead(left)
		}
		if self.Next != nil {
			self.Next.Prev = right
			right.Next = self.Next
		}
		left.Next = right
		right.Prev = left
		self.Prev = nil
		self.Next = nil

		self.innerInsert(key)
		leftSize := len(self.Entries) / 2
		rightSize := len(self.Entries) - leftSize
		for i := range leftSize {
			left.Entries = append(left.Entries, self.Entries[i])
		}
		for i := range rightSize {
			right.Entries = append(right.Entries, self.Entries[leftSize+i])
		}

		return nil
	}
	if key.Compare(self.Entries[0]) < 0 {
		return self.Children[0].Insert(key, bpTree, isUnique)
	} else if key.Compare(self.Entries[len(self.Entries)-1]) >= 0 {
		return self.Children[len(self.Children)-1].Insert(key, bpTree, isUnique)
	}
	for i, entry := range self.Entries {
		if key.Compare(entry) != -1 {
			continue
		}
		return self.Children[i].Insert(key, bpTree, isUnique)
	}
	return nil
}

func (self *BPNode) Remove(key *meta.IndexEntry, bpTree *BPTree) {

}
