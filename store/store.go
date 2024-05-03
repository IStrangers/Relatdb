package store

import (
	"Relatdb/meta"
	"Relatdb/utils"
	"os"
	"strings"
)

const (
	META_SUFFIX = ".meta"
)

type Options struct {
	Path string
}

type Store struct {
	path     string
	tableMap map[string]*meta.Table
}

func NewStore(options *Options) *Store {
	return &Store{
		path: utils.ConcatFilePaths(options.Path),
	}
}

func (self *Store) Init() {
	metaDir, err := os.ReadDir(self.path)
	if err != nil {
		panic(err)
	}
	for _, entry := range metaDir {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, META_SUFFIX) {
			continue
		}
		table, err := self.ReadTable(utils.ConcatFilePaths(self.path, fileName))
		if err != nil {
			panic(err)
		}
		self.tableMap[table.Name] = table
	}
}

func (self *Store) ReadTable(path string) (*meta.Table, error) {
	pageStore, err := NewPageStore(path)
	if err != nil {
		return nil, err
	}
	page := pageStore.ReadPage(0)
	items := page.ReadItems()
	var entries []meta.IndexEntry
	for _, item := range items {
		entries = append(entries, ItemToIndexEntry(item))
	}
	return &meta.Table{}, nil
}
