package store

import (
	"Relatdb/meta"
	"Relatdb/utils"
	"os"
)

type Options struct {
	Path string
}

type Store struct {
	metaPath string
	dataPath string
	tableMap map[string]*meta.Table
}

func CreateStore(options *Options) *Store {
	return &Store{
		metaPath: utils.ConcatFilePaths(options.Path, "meta"),
		dataPath: utils.ConcatFilePaths(options.Path, "data"),
	}
}

func (self *Store) Init() {
	metaDir, err := os.ReadDir(self.metaPath)
	if err != nil {
		panic(err)
	}
	for _, entry := range metaDir {
		table := self.ReadTable(utils.ConcatFilePaths(self.metaPath, entry.Name()))
		self.tableMap[table.Name] = table
	}
}

func (self *Store) ReadTable(path string) *meta.Table {
	return &meta.Table{}
}
