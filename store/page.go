package store

import (
	"Relatdb/common"
)

/*
1.页存储结构, 写入逻辑：ItemPointer向右写入, ItemData向左写入, 往中间靠拢，中间未写入的就是剩余空间
MagicWord | LowerOffset | UpperOffset | Special | TupleCount | ItemPointer1 | ItemPointer2 | ItemPointer3
										剩余空间...
ItemData1 | ItemData2 | ItemData3 |
*/

const (
	DEFAULT_PAGE_SIZE            = 4096
	DEFAULT_SPECIAL_POINT_LENGTH = 64
)

const (
	PAGE_HEADER_SIZE          = 24
	MAGIC_WORD                = "MagicWord"
	LOWER_WRITE_POINTER       = uint(len(MAGIC_WORD)) + 1
	UPPER_WRITE_POINTER       = LOWER_WRITE_POINTER + 4
	SPECIAL_WRITE_POINTER     = UPPER_WRITE_POINTER + 4
	TUPLE_COUNT_WRITE_POINTER = SPECIAL_WRITE_POINTER + 4
)

type PageHeader struct {
	//剩余空间起始偏移
	LowerOffset int
	//剩余空间末尾偏移
	UpperOffset int
	//下一页储偏指针信息的存移
	Special int
	//Item数量
	TupleCount int
	//页头长度
	HeaderLength int
}

func NewPageHeader(size uint) *PageHeader {
	magicWordLength := len([]byte(MAGIC_WORD)) + 1
	lowerOffset := magicWordLength + 4 + 4 + 4 + 4
	upperOffset := int(size) - DEFAULT_SPECIAL_POINT_LENGTH
	return &PageHeader{
		LowerOffset:  lowerOffset,
		UpperOffset:  upperOffset,
		Special:      upperOffset,
		HeaderLength: lowerOffset,
	}
}

type Page struct {
	Header *PageHeader
	Buffer *common.Buffer
	Length uint
	Dirty  bool
}

func NewPage() *Page {
	return NewPageBySize(DEFAULT_PAGE_SIZE)
}

func NewPageByBuffer(buffer *common.Buffer) *Page {
	page := NewPageBySize(buffer.Length)
	pageHeader := page.Header
	magicWord := buffer.ReadStringWithZero()
	if magicWord != MAGIC_WORD {
	}
	pageHeader.LowerOffset = buffer.ReadInt()
	pageHeader.UpperOffset = buffer.ReadInt()
	pageHeader.Special = buffer.ReadInt()
	pageHeader.TupleCount = buffer.ReadInt()
	page.Buffer = buffer
	return page
}

func NewPageBySize(size uint) *Page {
	pageHeader := NewPageHeader(size)
	page := &Page{
		Header: pageHeader,
		Buffer: common.NewBufferBySize(size),
		Length: size,
		Dirty:  false,
	}
	page.writePageHeader()
	return page
}

func (self *Page) writePageHeader() {
	self.Buffer.WriteStringWithZero(MAGIC_WORD)
	self.Buffer.WriteInt(self.Header.LowerOffset)
	self.Buffer.WriteInt(self.Header.UpperOffset)
	self.Buffer.WriteInt(self.Header.Special)
	self.Buffer.WriteInt(self.Header.TupleCount)
}

// 剩余可用空间
func (self *Page) remainFreeSpace() int {
	return self.Header.UpperOffset - self.Header.LowerOffset
}

// 更新剩余空间起始偏移
func (self *Page) updateHeaderLowerOffset(lowerOffset int) {
	self.Header.LowerOffset = lowerOffset
	self.Buffer.WriteIntByPos(LOWER_WRITE_POINTER, lowerOffset)
}

// 更新剩余空间末尾偏移
func (self *Page) updateHeaderUpperOffset(upperOffset int) {
	self.Header.UpperOffset = upperOffset
	self.Buffer.WriteIntByPos(UPPER_WRITE_POINTER, upperOffset)
}

// 更新Item写入数量
func (self *Page) updateHeaderTupleCount(tupleCount int) {
	self.Header.TupleCount = tupleCount
	self.Buffer.WriteIntByPos(TUPLE_COUNT_WRITE_POINTER, tupleCount)
}

func (self *Page) readItemPointer() *ItemPointer {
	return NewItemPointer(self.Buffer.ReadInt(), self.Buffer.ReadInt())
}

func (self *Page) readItemData(itemPointer *ItemPointer) *ItemData {
	data := self.Buffer.ReadBytesByOffset(itemPointer.Offset, itemPointer.TupleLength)
	return NewItemData(data, itemPointer.TupleLength)
}

func (self *Page) readItem() *Item {
	itemPointer := self.readItemPointer()
	if itemPointer.TupleLength == -1 {
		return nil
	}
	itemData := self.readItemData(itemPointer)
	return NewItem(itemPointer, itemData)
}

func (self *Page) ReadItems() (items []*Item) {
	for _ = range self.Header.TupleCount {
		item := self.readItem()
		if item == nil {
			continue
		}
		items = append(items, item)
	}
	return
}

func (self *Page) WriteItem(items ...*Item) {
	for _, item := range items {
		data := item.Data
		pointer := item.Pointer
		if self.remainFreeSpace() < data.Length+ITEM_POINTER_LENGTH {
			panic("page remaining space insufficient")
		}
		//写入ItemData
		writePos := self.Header.UpperOffset - data.Length
		self.Buffer.WriteBytesByPos(uint(writePos), data.Data)
		self.updateHeaderUpperOffset(writePos)

		//写入ItemPointer
		self.Buffer.WriteInt(writePos)
		self.Buffer.WriteInt(pointer.TupleLength)
		self.updateHeaderLowerOffset(self.Header.LowerOffset + ITEM_POINTER_LENGTH)

		self.updateHeaderTupleCount(self.Header.TupleCount + 1)
		//标记为脏页
		self.Dirty = true
	}
}
