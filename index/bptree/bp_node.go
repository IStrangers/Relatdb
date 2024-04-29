package bptree

import (
	"Relatdb/index"
	"Relatdb/meta"
	"Relatdb/store"
	"errors"
	"slices"
)

type BPNode struct {
	OwnerTree *BPTree           //所属树
	IsRoot    bool              //是否是根节点
	isLeaf    bool              //是否是叶子节点
	Parent    *BPNode           //父节点
	Prev      *BPNode           //上一个叶子节点
	Next      *BPNode           //下一个叶子节点
	Entries   []meta.IndexEntry //关键字
	Children  []*BPNode         //子节点
	Page      *BPPage           //页
}

func CreateBPNode(ownerTree *BPTree, isRoot bool, isLeaf bool) *BPNode {
	bpNode := &BPNode{
		OwnerTree: ownerTree,
		IsRoot:    isRoot,
		isLeaf:    isLeaf,
		Entries:   make([]meta.IndexEntry, 1),
	}
	if !isLeaf {
		bpNode.Children = make([]*BPNode, 3)
	}
	return bpNode
}

func (self *BPNode) addEntriesByIndex(index uint, key meta.IndexEntry) {
	self.Entries = slices.Insert(self.Entries, int(index), key)
}

func (self *BPNode) addChildren(node *BPNode) {
	self.Children = append(self.Children, node)
}

func (self *BPNode) addChildrenByIndex(index uint, node *BPNode) {
	self.Children = slices.Insert(self.Children, int(index), node)
}

func (self *BPNode) findChildrenIndex(node *BPNode) int {
	return slices.Index(self.Children, node)
}

func (self *BPNode) removeChildren(node *BPNode) int {
	index := self.Parent.findChildrenIndex(node)
	if index < 0 {
		return index
	}
	self.Children = slices.Delete(self.Children, index, index)
	return index
}

func (self *BPNode) getBorrowKeyLength(key meta.IndexEntry) uint {
	itemLength := index.GetItemLength(key)
	if !self.isLeaf {
		itemLength += store.ITEM_INT_LENGTH
	}
	return itemLength
}

func (self *BPNode) internalCheckExist(key meta.IndexEntry) bool {
	for _, entry := range self.Entries {
		if key.GetCompareEntry().Compare(entry) == 0 {
			return true
		}
	}
	return false
}

// 叶子节点是否需要分裂
func (self *BPNode) isLeafSplit(key meta.IndexEntry) bool {
	return self.Page.remainFreeSpace() < index.GetItemLength(key)
}

// 内部节点是否需要分裂
func (self *BPNode) isInternalSplit() bool {
	return self.Page.remainFreeSpace() < 0
}

// 内部节点插入
func (self *BPNode) internalInsert(key meta.IndexEntry) {
	insertIndex := len(self.Entries)
	if insertIndex == 0 {
		self.Entries = append(self.Entries, key)
		return
	}
	//插入在大于等于的Entries前面
	for i, entry := range self.Entries {
		if key.Compare(entry) == 0 || key.Compare(entry) < 0 {
			insertIndex = i
			break
		}
	}
	self.addEntriesByIndex(uint(insertIndex), key)
}

// 回收
func (self *BPNode) recycle() {
	self.Entries = nil
	self.Children = nil
	self.OwnerTree.RecyclePageNo(self.Page.PageNo)
}

func (self *BPNode) handlingParent(bpTree *BPTree, left *BPNode, right *BPNode) {
	//根节点
	if self.IsRoot {
		self.IsRoot = false
		//创建新的根节点
		root := CreateBPNode(self.OwnerTree, true, false)
		//更新节点指向
		bpTree.Root = root
		left.Parent = root
		right.Parent = root
		//将left和right叶子节点添加到root
		root.addChildren(left)
		root.addChildren(right)
		//将当前节点中间的key添加到root
		root.internalInsert(right.Entries[0].GetCompareEntry())
		//root节点进行分裂
		root.internalSplit(bpTree)
	} else {
		//删除当前叶子节点并返回在父节点的下标位置
		index := self.Parent.removeChildren(self)
		//更新父节点指向
		left.Parent = self.Parent
		right.Parent = self.Parent
		//将left和right叶子节点添加到父节点
		self.Parent.addChildrenByIndex(uint(index), left)
		self.Parent.addChildrenByIndex(uint(index+1), right)
		//将当前节点中间的key添加到父节点
		self.Parent.internalInsert(right.Entries[0].GetCompareEntry())
		//父节点进行分裂
		self.Parent.internalSplit(bpTree)
		self.Parent = nil
	}
	//回收
	self.recycle()
}

// 内部节点分裂
func (self *BPNode) internalSplit(bpTree *BPTree) {
	if !self.isInternalSplit() {
		return
	}
	left := CreateBPNode(self.OwnerTree, false, false)
	right := CreateBPNode(self.OwnerTree, false, false)

	leftSize := len(self.Entries) / 2
	rightSize := len(self.Entries) - leftSize
	//将当前节点的key复制到新的left和right节点
	for i := range leftSize {
		left.Entries = append(left.Entries, self.Entries[i])
	}
	for i := range rightSize {
		right.Entries = append(right.Entries, self.Entries[leftSize+i])
	}
	//将当前节点的children复制到新的left和right节点
	for i := range leftSize {
		children := self.Children[i]
		children.Parent = left
		left.addChildren(children)
	}
	for i := range rightSize {
		children := self.Children[leftSize+i]
		children.Parent = right
		right.addChildren(children)
	}
	self.handlingParent(bpTree, left, right)
}

// 获取
func (self *BPNode) Get(key meta.IndexEntry, compareType index.CompareType) *BPPosition {
	//叶子节点
	if self.isLeaf {
		if compareType == index.COMPARE_EQUAL {
			//查找相等的Entries
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
	//非叶子节点
	if firstEntry := self.Entries[0]; key.Compare(firstEntry) < 0 { //小于第一个Entries
		return self.Children[0].Get(key, compareType)
	} else if lastEntry := self.Entries[len(self.Entries)-1]; key.Compare(lastEntry) >= 0 { //大于等于最后一个Entries
		return self.Children[len(self.Children)-1].Get(key, compareType)
	}
	//查找大于等于的Entries
	for i, entry := range self.Entries {
		if key.Compare(entry) != -1 {
			continue
		}
		return self.Children[i].Get(key, compareType)
	}
	return nil
}

// 插入
func (self *BPNode) Insert(key meta.IndexEntry, bpTree *BPTree, isUnique bool) error {
	if self.getBorrowKeyLength(key) > self.Page.getInitFreeSpace()/3 {
		return errors.New("entry size must <= Max/3")
	} else if isUnique && self.internalCheckExist(key) {
		return errors.New("duplicated Key error")
	}
	//叶子节点
	if self.isLeaf {
		//叶子节点不用分裂直接插入key
		if !self.isLeafSplit(key) {
			self.internalInsert(key)
			return nil
		}
		//叶子节点需要分裂，并将当前叶子节点分裂成两个新的叶子节点
		left := CreateBPNode(self.OwnerTree, false, true)  //左叶子节点
		right := CreateBPNode(self.OwnerTree, false, true) //右叶子节点
		//将Prev和Next叶子节点指向新的left和right
		if self.Prev != nil {
			self.Prev.Next = left
			left.Prev = self.Prev
		} else {
			bpTree.Head = left
		}
		if self.Next != nil {
			self.Next.Prev = right
			right.Next = self.Next
		}
		left.Next = right
		right.Prev = left
		self.Prev = nil
		self.Next = nil

		//先插入key
		self.internalInsert(key)
		//将当前叶子节点的key复制到新的left和right节点
		leftSize := len(self.Entries) / 2
		rightSize := len(self.Entries) - leftSize
		for i := range leftSize {
			left.Entries = append(left.Entries, self.Entries[i])
		}
		for i := range rightSize {
			right.Entries = append(right.Entries, self.Entries[leftSize+i])
		}
		self.handlingParent(bpTree, left, right)
		return nil
	}
	//非叶子节点
	if key.Compare(self.Entries[0]) < 0 { //小于第一个Entries
		return self.Children[0].Insert(key, bpTree, isUnique)
	} else if key.Compare(self.Entries[len(self.Entries)-1]) >= 0 { //大于等于最后一个Entries
		return self.Children[len(self.Children)-1].Insert(key, bpTree, isUnique)
	}
	//查找大于等于的Entries
	for i, entry := range self.Entries {
		if key.Compare(entry) != -1 {
			continue
		}
		return self.Children[i].Insert(key, bpTree, isUnique)
	}
	return nil
}

// 删除
func (self *BPNode) Remove(key meta.IndexEntry, bpTree *BPTree) {

}
