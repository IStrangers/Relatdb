package store

import (
	"Relatdb/meta"
	"Relatdb/utils"
	"errors"
	"os"
	"strings"
)

const (
	META_SUFFIX = ".meta"
	DATA_SUFFIX = ".data"
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
		table, err := self.readTable(utils.ConcatFilePaths(self.path, fileName))
		if err != nil {
			panic(err)
		}
		self.tableMap[table.Name] = table
	}
}

func (self *Store) readTable(path string) (*meta.Table, error) {
	pageStore := NewPageStore(path)
	page := pageStore.readPage(0)
	items := page.readItems()
	var entries []meta.IndexEntry
	for _, item := range items {
		entries = append(entries, ItemToIndexEntry(item))
	}
	return &meta.Table{}, nil
}

func (self *Store) writeTable(table *meta.Table) error {
	pageStore := NewPageStore(table.MetaPath)

	var items []*Item
	for _, field := range table.Fields {
		items = append(items, IndexEntryToItem(meta.NewIndexEntry(meta.FieldToValues(field), nil)))
	}
	itemName := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.StringValue(table.Name)}, nil))
	itemSize := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(len(items) + 1)}, nil))
	page := NewPage()
	page.writeItem(itemSize)
	page.writeItem(itemName)
	page.writeItem(items...)

	indexNum := IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(len(table.SecondaryIndexes) + 1)}, nil))
	page.writeItem(indexNum)

	pageStore.writePage(page, 0)
	return nil
}

func (self *Store) CreateTable(table *meta.Table) error {
	if self.tableMap[table.Name] != nil {
		return errors.New("table already exists: " + table.Name)
	}
	table.MetaPath = utils.ConcatFilePaths(self.path, table.Name+META_SUFFIX)
	table.DataPath = utils.ConcatFilePaths(self.path, table.Name+DATA_SUFFIX)
	err := self.writeTable(table)
	if err != nil {
		return err
	}
	self.tableMap[table.Name] = table
	return nil
}
