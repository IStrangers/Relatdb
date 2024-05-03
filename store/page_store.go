package store

import (
	"os"
)

type PageStore struct {
	path string
	file *os.File
}

func CreatePageStore(path string) (*PageStore, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &PageStore{
		path: path,
		file: file,
	}, nil
}

func (self *PageStore) ReadPage(pageIndex int) *Page {
	readPos := int64(pageIndex * DEFAULT_PAGE_SIZE)
	buf := make([]byte, DEFAULT_PAGE_SIZE)
	self.file.Seek(readPos, 0)
	self.file.Read(buf)
	buffer := CreateBuffer(buf)
	return CreatePage(buffer)
}
