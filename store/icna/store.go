package icna

import (
	"Relatdb/index/bptree"
	"Relatdb/meta"
	"Relatdb/store"
	"Relatdb/utils"
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

type IcnaStore struct {
	path     string
	tableMap map[string]*meta.Table
}

func NewIcnaStore(options *Options) *IcnaStore {
	store := &IcnaStore{
		path:     options.Path,
		tableMap: make(map[string]*meta.Table),
	}
	_ = os.MkdirAll(store.path, os.ModePerm)
	return store
}

func (self *IcnaStore) Init() {
	metaDir, _ := os.ReadDir(self.path)
	for _, entry := range metaDir {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, META_SUFFIX) {
			continue
		}
		table := self.readTable(utils.ConcatFilePaths(self.path, fileName))
		self.tableMap[table.Name] = table
	}
}

func (self *IcnaStore) readTable(path string) *meta.Table {
	pageStore := store.NewPageStore(path)
	page := pageStore.ReadPage(0)
	items := page.ReadItems()
	var entries []meta.IndexEntry
	for _, item := range items {
		entries = append(entries, store.ItemToIndexEntry(item))
	}
	metaQuantity := entries[0].GetValues()[0].ToInt()
	tableName := entries[1].GetValues()[0].ToString()
	var fields []*meta.Field
	for i := 2; i <= metaQuantity; i++ {
		values := entries[i].GetValues()
		fields = append(fields, meta.NewFieldByValues(values))
	}
	var clusterIndex meta.Index
	var secondaryIndexes []meta.Index
	indexQuantity := entries[metaQuantity+1].GetValues()[0].ToInt()
	indexStartOffset := metaQuantity + 2
	for i := range indexQuantity {
		indexMetaSize := entries[indexStartOffset].GetValues()[0].ToInt()
		indexName := entries[indexStartOffset+1].GetValues()[0].ToString()
		indexFlag := entries[indexStartOffset+2].GetValues()[0].ToInt()
		var indexFields []*meta.Field
		for j := indexStartOffset + 3; j < indexMetaSize+indexStartOffset+1; j++ {
			values := entries[j].GetValues()
			indexFields = append(indexFields, meta.NewFieldByValues(values))
		}
		index := bptree.NewBPTree(indexName, indexFields, uint(indexFlag))
		if i == 0 {
			clusterIndex = index
		} else {
			secondaryIndexes = append(secondaryIndexes, index)
		}
		indexStartOffset += indexMetaSize + 1
	}
	return &meta.Table{
		MetaPath:         path,
		DataPath:         strings.ReplaceAll(path, META_SUFFIX, DATA_SUFFIX),
		Name:             tableName,
		Fields:           fields,
		ClusterIndex:     clusterIndex,
		SecondaryIndexes: secondaryIndexes,
	}

}

func (self *IcnaStore) writeTable(table *meta.Table) {
	pageStore := store.NewPageStore(table.MetaPath)

	//写入字段
	var fields []*store.Item
	for _, field := range table.Fields {
		fields = append(fields, store.IndexEntryToItem(meta.NewIndexEntry(meta.FieldToValues(field), nil)))
	}
	metaQuantity := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(len(fields) + 1)}, nil))
	tableName := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.StringValue(table.Name)}, nil))
	page := store.NewPage()
	page.WriteItem(metaQuantity)
	page.WriteItem(tableName)
	page.WriteItem(fields...)
	//写入索引
	indexQuantity := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(len(table.SecondaryIndexes) + 1)}, nil))
	page.WriteItem(indexQuantity)
	page.WriteItem(store.IndexToItems(table.ClusterIndex)...)
	for _, secondaryIndex := range table.SecondaryIndexes {
		page.WriteItem(store.IndexToItems(secondaryIndex)...)
	}

	pageStore.WritePage(page, 0)
}

func (self *IcnaStore) CreateTable(table *meta.Table) {
	if self.tableMap[table.Name] != nil {
		panic("table already exists: " + table.Name)
	}
	if table.ClusterIndex == nil {
		panic("cluster index is required: " + table.Name)
	}
	table.MetaPath = utils.ConcatFilePaths(self.path, table.Name+META_SUFFIX)
	table.DataPath = utils.ConcatFilePaths(self.path, table.Name+DATA_SUFFIX)
	self.writeTable(table)

	self.tableMap[table.Name] = table
}
