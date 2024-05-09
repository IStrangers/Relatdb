package store

import (
	"Relatdb/common"
	"os"
)

type PageStore struct {
	path string
	file *os.File
}

func NewPageStore(path string) *PageStore {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	return &PageStore{
		path: path,
		file: file,
	}
}

func (self *PageStore) ReadPage(pageIndex int) *Page {
	readPos := int64(pageIndex * DEFAULT_PAGE_SIZE)
	buf := make([]byte, DEFAULT_PAGE_SIZE)
	self.file.Seek(readPos, 0)
	self.file.Read(buf)
	buffer := common.NewBuffer(buf)
	return NewPageByBuffer(buffer)
}

func (self *PageStore) WritePage(page *Page, pageIndex int) {
	writePos := int64(pageIndex * DEFAULT_PAGE_SIZE)
	self.file.Seek(writePos, 0)
	self.file.Write(page.Buffer.Data)
}
