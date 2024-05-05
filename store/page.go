package store

import (
	"Relatdb/common"
	"errors"
)

const (
	DEFAULT_PAGE_SIZE            = 4096
	DEFAULT_SPECIAL_POINT_LENGTH = 64
)

const (
	PAGE_HEADER_SIZE    = 24
	MAGIC_WORD          = "StorePage"
	LOWER_POINTER       = 8
	UPPER_POINTER       = 12
	SPECIAL_POINTER     = 16
	TUPLE_COUNT_POINTER = 20
)

type PageHeader struct {
	LowerOffset  int
	UpperOffset  int
	Special      int
	TupleCount   int
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

func (self *Page) remainFreeSpace() int {
	return self.Header.UpperOffset - self.Header.LowerOffset
}

func (self *Page) writePageHeader() {
	self.Buffer.WriteStringWithZero(MAGIC_WORD)
	self.Buffer.WriteInt(self.Header.LowerOffset)
	self.Buffer.WriteInt(self.Header.UpperOffset)
	self.Buffer.WriteInt(self.Header.Special)
	self.Buffer.WriteInt(self.Header.TupleCount)
}

func (self *Page) readItemPointer() *ItemPointer {
	return NewItemPointer(self.Buffer.ReadInt(), self.Buffer.ReadInt())
}

func (self *Page) readItemData(itemPointer *ItemPointer) *ItemData {
	data := self.Buffer.ReadBytesByOffset(itemPointer.Offset, itemPointer.TupleLength)
	return NewItemData(data, len(data))
}

func (self *Page) readItem() *Item {
	itemPointer := self.readItemPointer()
	if itemPointer.TupleLength == -1 {
		return nil
	}
	itemData := self.readItemData(itemPointer)
	return NewItem(itemPointer, itemData)
}

func (self *Page) readItems() (items []*Item) {
	for _ = range self.Header.TupleCount {
		item := self.readItem()
		if item == nil {
			continue
		}
		items = append(items, item)
	}
	return
}

func (self *Page) writeItem(items ...*Item) error {
	for _, item := range items {
		if self.remainFreeSpace() < item.Data.Length+ITEM_POINTER_LENGTH {
			return errors.New("page remaining space insufficient")
		}
		writePos := self.Header.UpperOffset - item.Data.Length
		self.Buffer.WriteBytes(item.Data.Data)
		self.Header.UpperOffset = writePos
		self.Dirty = true
	}
	return nil
}
