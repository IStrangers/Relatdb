package store

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

func NewPageHeader(size int) *PageHeader {
	magicWordLength := len([]byte(MAGIC_WORD)) + 1
	lowerOffset := magicWordLength + 4 + 4 + 4 + 4
	upperOffset := size - DEFAULT_SPECIAL_POINT_LENGTH
	return &PageHeader{
		LowerOffset:  lowerOffset,
		UpperOffset:  upperOffset,
		Special:      upperOffset,
		HeaderLength: lowerOffset,
	}
}

type Page struct {
	Header *PageHeader
	Buffer *Buffer
	Length int
	Dirty  bool
}

func NewPage() *Page {
	return NewPageBySize(DEFAULT_PAGE_SIZE)
}

func NewPageByBuffer(buffer *Buffer) *Page {
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

func NewPageBySize(size int) *Page {
	pageHeader := NewPageHeader(size)
	page := &Page{
		Header: pageHeader,
		Buffer: NewBufferBySize(size),
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

func (self *Page) readItemPointer() *ItemPointer {
	return NewItemPointer(self.Buffer.ReadInt(), self.Buffer.ReadInt())
}

func (self *Page) readItemData(itemPointer *ItemPointer) *ItemData {
	data := self.Buffer.ReadBytesByOffset(itemPointer.Offset, itemPointer.TupleLength)
	return NewItemData(data, len(data), itemPointer.Offset)
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

func (self *Page) writeItem(item *Item) {

}
