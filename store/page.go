package store

import "bytes"

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

func CreatePageHeader(size int) *PageHeader {
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
	Header PageHeader
	Buff   *bytes.Buffer
	Length int
	Dirty  bool
}

func ReadPage(pageIndex int, isIndex bool) *Page {
	return &Page{}
}
