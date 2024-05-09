package bptree

import (
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

func NewBPNode(ownerTree *BPTree, isRoot bool, isLeaf bool) *BPNode {
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

func (self *BPNode) addEntries(key ...meta.IndexEntry) {
	self.Entries = append(self.Entries, key...)
}

func (self *BPNode) addEntriesByIndex(index int, key ...meta.IndexEntry) {
	self.Entries = slices.Insert(self.Entries, index, key...)
}

func (self *BPNode) findDeleteEntriesIndex(key meta.IndexEntry) int {
	for i, entry := range self.Entries {
		if entry.GetDeleteCompareEntry().CompareDeleteEntry(key) == 0 {
			return i
		}
	}
	return -1
}

func (self *BPNode) removeEntriesByIndex(index int) meta.IndexEntry {
	key := self.Entries[index]
	self.Entries = slices.Delete(self.Entries, index, index)
	return key
}

func (self *BPNode) setEntriesByIndex(index int, key meta.IndexEntry) {
	self.Entries[index] = key
}

func (self *BPNode) addChildren(node ...*BPNode) {
	self.Children = append(self.Children, node...)
}

func (self *BPNode) addChildrenByIndex(index int, node ...*BPNode) {
	self.Children = slices.Insert(self.Children, index, node...)
}

func (self *BPNode) findChildrenIndex(node *BPNode) int {
	return slices.Index(self.Children, node)
}

func (self *BPNode) removeChildren(node *BPNode) int {
	index := self.Parent.findChildrenIndex(node)
	if index < 0 {
		return -1
	}
	self.removeChildrenByIndex(index)
	return index
}

func (self *BPNode) removeChildrenByIndex(index int) *BPNode {
	child := self.Children[index]
	self.Children = slices.Delete(self.Children, int(index), int(index))
	return child
}

func (self *BPNode) getBorrowKeyLength(key meta.IndexEntry) uint {
	itemLength := store.GetItemLength(key)
	if !self.isLeaf {
		itemLength += store.ITEM_INT_LENGTH
	}
	return itemLength
}

func (self *BPNode) internalCheckExist(key meta.IndexEntry) bool {
	for _, entry := range self.Entries {
		if key.GetCompareEntry().CompareEntry(entry) == 0 {
			return true
		}
	}
	return false
}

// 叶子节点是否需要分裂
func (self *BPNode) isLeafSplit(key meta.IndexEntry) bool {
	return self.Page.remainFreeSpace() < store.GetItemLength(key)
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
		if key.CompareEntry(entry) == 0 || key.CompareEntry(entry) < 0 {
			insertIndex = i
			break
		}
	}
	self.addEntriesByIndex(insertIndex, key)
}

// 回收
func (self *BPNode) recycle() {
	self.OwnerTree.RecyclePageNo(self.Page.PageNo)
	self.OwnerTree = nil
	self.Parent = nil
	self.Prev = nil
	self.Next = nil
	self.Entries = nil
	self.Children = nil
	self.Page = nil
}

func (self *BPNode) handlingParent(bpTree *BPTree, left *BPNode, right *BPNode) {
	//根节点
	if self.IsRoot {
		self.IsRoot = false
		//创建新的根节点
		root := NewBPNode(self.OwnerTree, true, false)
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
		self.Parent.addChildrenByIndex(index, left)
		self.Parent.addChildrenByIndex(index+1, right)
		//将当前节点中间的key添加到父节点
		self.Parent.internalInsert(right.Entries[0].GetCompareEntry())
		//父节点进行分裂
		self.Parent.internalSplit(bpTree)
	}
	//回收
	self.recycle()
}

// 内部节点分裂
func (self *BPNode) internalSplit(bpTree *BPTree) {
	if !self.isInternalSplit() {
		return
	}
	left := NewBPNode(self.OwnerTree, false, false)
	right := NewBPNode(self.OwnerTree, false, false)

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

func (self *BPNode) internalRemove(key meta.IndexEntry) bool {
	index := self.findDeleteEntriesIndex(key)
	if index < 0 {
		return false
	}
	self.removeEntriesByIndex(index)
	return true
}

// 是否是平衡的叶子节点
func (self *BPNode) isBalancedLeaf(key meta.IndexEntry) bool {
	return self.Page.getContentSize()-store.GetItemLength(key) > self.Page.getInitFreeSpace()/2
}

// 上一个叶子节点是否可借用
func (self *BPNode) prevLeafCanBorrow() bool {
	if self.Prev == nil || len(self.Prev.Entries) < 2 || self.Prev.Parent != self.Parent {
		return false
	}
	return self.isBalancedLeaf(self.Prev.Entries[len(self.Prev.Entries)-1]) ||
		len(self.Entries) == 1 && len(self.Prev.Entries) >= 2
}

// 下一个叶子节点是否可借用
func (self *BPNode) nextLeafCanBorrow() bool {
	if self.Next == nil || len(self.Next.Entries) < 2 || self.Next.Parent != self.Parent {
		return false
	}
	return self.isBalancedLeaf(self.Next.Entries[len(self.Next.Entries)-1]) ||
		len(self.Entries) == 1 && len(self.Next.Entries) >= 2
}

// 叶子节点是否可合并
func (self *BPNode) leafCanMerge(leaf *BPNode) bool {
	return leaf != nil && leaf.Parent == self.Parent &&
		leaf.Page.getContentSize() < leaf.Page.getInitFreeSpace()/2 &&
		leaf.Page.getContentSize() <= self.Page.remainFreeSpace()
}

// 合并Prev节点
func (self *BPNode) mergePrevNode(prev *BPNode) {
	if !prev.isLeaf {
		index := self.Parent.findChildrenIndex(self)
		self.addEntriesByIndex(0, self.Parent.Entries[index])
		for _, child := range prev.Children {
			child.Parent = self
		}
		self.addChildrenByIndex(0, prev.Children...)
	}
	self.addEntriesByIndex(0, prev.Entries...)
}

// 合并Next节点
func (self *BPNode) mergeNextNode(next *BPNode) {
	if !next.isLeaf {
		index := self.Parent.findChildrenIndex(next)
		self.addEntries(next.Parent.Entries[index])
	}
	self.addEntries(next.Entries...)
	if len(next.Children) != 0 {
		for _, child := range next.Children {
			child.Parent = self
		}
		self.addChildren(next.Children...)
	}
}

// 上一个内部节点是否可借用
func (self *BPNode) prevInternalCanBorrow(prev *BPNode) bool {
	return prev != nil && len(prev.Entries) >= 2 &&
		self.internalCanBorrow(prev, prev.Entries[len(prev.Entries)-1])
}

// 下一个内部节点是否可借用
func (self *BPNode) nextInternalCanBorrow(next *BPNode) bool {
	return next != nil && len(next.Entries) >= 2 &&
		self.internalCanBorrow(next, next.Entries[0])
}

// 内部节点是否可借用
func (self *BPNode) internalCanBorrow(node *BPNode, key meta.IndexEntry) bool {
	if node == nil {
		return false
	}
	if len(self.Entries) == 0 && len(node.Entries) >= 2 {
		return true
	}
	borrowKeyLength := self.getBorrowKeyLength(key)
	return node.Parent == self.Parent && len(node.Entries) >= 2 &&
		node.Page.getContentSize()-borrowKeyLength > node.Page.getInitFreeSpace()/2 &&
		borrowKeyLength <= self.Page.remainFreeSpace()
}

func (self *BPNode) prevCanMerge(prev *BPNode) bool {
	if prev == nil {
		return false
	}
	if prev.Parent == self.Parent {
		adjustSize := uint(0)
		if !prev.isLeaf {
			selfIndex := self.Parent.findChildrenIndex(self)
			downKey := self.Parent.Entries[selfIndex]
			adjustSize = store.GetItemLength(downKey)
		}
		return prev.Page.getContentSize()+adjustSize <= self.Page.remainFreeSpace()
	}
	return false
}

func (self *BPNode) nextCanMerge(next *BPNode) bool {
	if next == nil {
		return false
	}
	if next.Parent == self.Parent {
		adjustSize := uint(0)
		if !next.isLeaf {
			nextIndex := self.Parent.findChildrenIndex(next)
			downKey := self.Parent.Entries[nextIndex]
			adjustSize = store.GetItemLength(downKey)
		}
		return next.Page.getContentSize()+adjustSize <= self.Page.remainFreeSpace()
	}
	return false
}

// 内部节点合并
func (self *BPNode) internalMerge(bpTree *BPTree) {
	if len(self.Children) < 2 || self.Page.getContentSize() < self.Page.getInitFreeSpace()/2 {
		if self.IsRoot && len(self.Children) < 2 {
			//与子节点合并，子节点成为根节点
			child := self.Children[0]
			bpTree.Root = child
			self.recycle()
		} else {
			selfIndex := self.Parent.findChildrenIndex(self)
			prevIndex := selfIndex - 1
			nextIndex := selfIndex + 1
			var prev, next *BPNode
			if prevIndex >= 0 {
				prev = self.Children[prevIndex]
			}
			if nextIndex < len(self.Children) {
				next = self.Children[nextIndex]
			}
			defer self.Parent.internalMerge(bpTree)
			// 上一个内部节点是否可借用
			if self.prevInternalCanBorrow(prev) {
				// 下放key
				downKey := self.Parent.Entries[selfIndex]
				self.addEntriesByIndex(0, downKey)
				// Prev key上提
				prevLastKeyIndex := len(prev.Entries) - 1
				self.Parent.setEntriesByIndex(selfIndex, prev.Entries[prevLastKeyIndex])
				prev.removeEntriesByIndex(prevLastKeyIndex)
				// 子节点也借用
				borrowChild := prev.removeChildrenByIndex(len(prev.Children) - 1)
				borrowChild.Parent = self
				self.addChildrenByIndex(0, borrowChild)
				return
			}
			// 下一个内部节点是否可借用
			if self.nextInternalCanBorrow(next) {
				// 下放key
				downKey := self.Parent.Entries[nextIndex]
				self.addEntries(downKey)
				// Next key上提
				self.Parent.setEntriesByIndex(nextIndex, next.Entries[0])
				next.removeEntriesByIndex(0)
				// 子节点也借用
				borrowChild := next.removeChildrenByIndex(0)
				borrowChild.Parent = self
				self.addChildren(borrowChild)
				return
			}
			if self.prevCanMerge(prev) {
				self.mergePrevNode(prev)
				self.Parent.removeEntriesByIndex(selfIndex)
				self.Parent.removeChildren(prev)
				prev.recycle()
				return
			}
			if self.nextCanMerge(next) {
				self.mergeNextNode(next)
				self.Parent.removeEntriesByIndex(nextIndex)
				self.Parent.removeChildren(next)
				next.recycle()
				return
			}
		}
		return
	}
	if self.Page.getContentSize() > self.Page.getInitFreeSpace() {
		/*
			因为在更新的时候,由于key值大小不定,可能导致虽然删除了关键字,但是由于
			更新了新的长的key,导致比删除之前的size还要大,所以就有可能导致分裂
			即changeKeySize - deleteKeySize > 0的某些情况下会导致分裂
		*/
		self.internalSplit(bpTree)
		return
	}
}

// 获取
func (self *BPNode) Get(key meta.IndexEntry, compareType meta.CompareType) *BPPosition {
	//叶子节点
	if self.isLeaf {
		if compareType == meta.COMPARE_EQUAL {
			//查找相等的Entries
			for i, entry := range self.Entries {
				if key.CompareEntry(entry) != 0 {
					continue
				}
				return NewBPPosition(nil, uint(i), self)
			}
		} else if compareType == meta.COMPARE_LOW {
			return NewBPPosition(nil, 0, self)
		} else {
			return NewBPPosition(nil, uint(len(self.Entries)-1), self)
		}
	}
	//非叶子节点
	if firstEntry := self.Entries[0]; key.CompareEntry(firstEntry) < 0 { //小于第一个Entries
		return self.Children[0].Get(key, compareType)
	} else if lastEntry := self.Entries[len(self.Entries)-1]; key.CompareEntry(lastEntry) >= 0 { //大于等于最后一个Entries
		return self.Children[len(self.Children)-1].Get(key, compareType)
	}
	//查找大于等于的Entries
	for i, entry := range self.Entries {
		if key.CompareEntry(entry) != -1 {
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
		left := NewBPNode(self.OwnerTree, false, true)  //左叶子节点
		right := NewBPNode(self.OwnerTree, false, true) //右叶子节点
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
	if key.CompareEntry(self.Entries[0]) < 0 { //小于第一个Entries
		return self.Children[0].Insert(key, bpTree, isUnique)
	} else if key.CompareEntry(self.Entries[len(self.Entries)-1]) >= 0 { //大于等于最后一个Entries
		return self.Children[len(self.Children)-1].Insert(key, bpTree, isUnique)
	}
	//查找大于等于的Entries
	for i, entry := range self.Entries {
		if key.CompareEntry(entry) != -1 {
			continue
		}
		return self.Children[i].Insert(key, bpTree, isUnique)
	}
	return nil
}

// 删除
func (self *BPNode) Remove(key meta.IndexEntry, bpTree *BPTree) bool {
	//叶子节点
	if self.isLeaf {
		//不包含key直接返回
		if self.findDeleteEntriesIndex(key) == -1 {
			return false
		}
		removeOk := self.internalRemove(key)
		//叶子节点并且根节点，表明只有一个节点
		if self.IsRoot {
			return removeOk
		}
		defer self.internalMerge(bpTree)
		//不平衡的叶子节点，不能直接删除
		if !self.isBalancedLeaf(key) {
			//上一个叶子节点是否可借用
			if self.prevLeafCanBorrow() {
				//借用Prev最后一个key添加到当前Entries最前面
				key := self.Prev.removeEntriesByIndex(len(self.Prev.Entries) - 1)
				self.addEntriesByIndex(0, key)
				//更新key到父节点Entries
				index := self.Parent.findChildrenIndex(self)
				self.Parent.setEntriesByIndex(index, key)
				return removeOk
			}
			//下一个叶子节点是否可借用
			if self.nextLeafCanBorrow() {
				//借用Next第一个key添加到当前Entries最后面
				key := self.Next.removeEntriesByIndex(0)
				self.addEntries(key)
				//将Next的第一个key上提到父节点Entries
				index := self.Parent.findChildrenIndex(self.Next)
				self.Parent.setEntriesByIndex(index, self.Next.Entries[0])
				return removeOk
			}
			//Prev是否可合并
			if self.leafCanMerge(self.Prev) {
				//合并Prev节点
				self.mergePrevNode(self.Prev)
				//父节点删除Prev节点
				prevIndex := self.Parent.removeChildren(self.Prev)
				//父节点删除当前节点的key
				self.Parent.removeEntriesByIndex(prevIndex + 1)
				//更新节点指向
				if self.Prev.Prev != nil {
					self.Prev = self.Prev.Prev
					self.Prev.Next = self
				} else {
					self.Prev = nil
					bpTree.Head = self
				}
				//回收Prev节点
				self.Prev.recycle()
				return removeOk
			}
			//Next是否可合并
			if self.leafCanMerge(self.Next) {
				//合并Next节点
				self.mergeNextNode(self.Next)
				//父节点删除Next节点
				nextIndex := self.Parent.removeChildren(self.Next)
				//父节点删除Next节点的key
				self.Parent.removeEntriesByIndex(nextIndex)
				//更新节点指向
				if self.Next.Next != nil {
					self.Next = self.Next.Next
					self.Next.Prev = self
				} else {
					self.Next = nil
				}
				//回收Next节点
				self.Next.recycle()
				return removeOk
			}
		}
		return removeOk
	}
	//非叶子节点
	if key.CompareEntry(self.Entries[0]) < 0 { //小于第一个Entries
		return self.Children[0].Remove(key, bpTree)
	} else if key.CompareEntry(self.Entries[len(self.Entries)-1]) >= 0 { //大于等于最后一个Entries
		return self.Children[len(self.Children)-1].Remove(key, bpTree)
	}
	//查找大于等于的Entries
	for i, entry := range self.Entries {
		if key.CompareEntry(entry) != -1 {
			continue
		}
		return self.Children[i].Remove(key, bpTree)
	}
	return false
}
