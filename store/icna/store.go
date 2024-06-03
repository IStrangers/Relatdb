package icna

import (
	"Relatdb/common"
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
	path        string
	databaseMap map[string]*meta.DataBase
}

func NewIcnaStore(options *Options) *IcnaStore {
	store := &IcnaStore{
		path:        options.Path,
		databaseMap: make(map[string]*meta.DataBase),
	}
	_ = os.MkdirAll(store.path, os.ModePerm)
	return store
}

func (self *IcnaStore) Init() {
	self.InitDatabases()
	self.InitTables()
}

func (self *IcnaStore) InitDatabases() {
	self.databaseMap["default"] = meta.NewDataBase("default")
}

func (self *IcnaStore) InitTables() {
	metaDir, _ := os.ReadDir(self.path)
	for _, entry := range metaDir {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, META_SUFFIX) {
			continue
		}
		table := self.readTable(utils.ConcatFilePaths(self.path, fileName))
		database := self.databaseMap[table.DatabaseName]
		database.TableMap[table.Name] = table
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
	databaseName := entries[1].GetValues()[0].ToString()
	tableName := entries[2].GetValues()[0].ToString()
	var fields []*meta.Field
	var primaryFiled *meta.Field
	fieldMap := make(map[string]*meta.Field, metaQuantity-3)
	for i := 3; i <= metaQuantity; i++ {
		values := entries[i].GetValues()
		field := meta.NewFieldByValues(values)
		if field.Flag&common.PRIMARY_KEY_FLAG != 0 {
			primaryFiled = field
		}
		fields = append(fields, field)
		fieldMap[field.Name] = field
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
	table := meta.NewTable(databaseName, tableName, fields, primaryFiled, fieldMap, clusterIndex, secondaryIndexes)
	table.MetaPath = path
	table.DataPath = strings.ReplaceAll(path, META_SUFFIX, DATA_SUFFIX)
	return table
}

func (self *IcnaStore) writeTable(table *meta.Table) {
	pageStore := store.NewPageStore(table.MetaPath)

	//写入字段
	var fields []*store.Item
	for _, field := range table.Fields {
		fields = append(fields, store.IndexEntryToItem(meta.NewIndexEntry(meta.FieldToValues(field), nil)))
	}
	metaQuantity := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.IntValue(len(fields) + 2)}, nil))
	databaseName := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.StringValue(table.DatabaseName)}, nil))
	tableName := store.IndexEntryToItem(meta.NewIndexEntry([]meta.Value{meta.StringValue(table.Name)}, nil))
	page := store.NewPage()
	page.WriteItem(metaQuantity)
	page.WriteItem(databaseName)
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

func (self *IcnaStore) CreateDatabase(database *meta.DataBase) {
	self.databaseMap[database.Name] = database
}

func (self *IcnaStore) DropDatabase(databaseName string) {
	database := self.databaseMap[databaseName]
	for _, table := range database.TableMap {
		self.DropTable(database.Name, table.Name)
	}
	delete(self.databaseMap, database.Name)
}

func (self *IcnaStore) GetDatabase(databaseName string) *meta.DataBase {
	if databaseName == "" {
		panic("no database selected")
	}
	database := self.databaseMap[databaseName]
	if database == nil {
		panic("database not exists: " + databaseName)
	}
	return database
}

func (self *IcnaStore) CreateTable(table *meta.Table) {
	database := self.GetDatabase(table.DatabaseName)
	if database.GetTable(table.Name) != nil {
		panic("table already exists: " + table.Name)
	}
	if table.ClusterIndex == nil {
		panic("cluster index is required: " + table.Name)
	}
	table.MetaPath = utils.ConcatFilePaths(self.path, table.Name+META_SUFFIX)
	table.DataPath = utils.ConcatFilePaths(self.path, table.Name+DATA_SUFFIX)
	self.writeTable(table)

	database.TableMap[table.Name] = table
}

func (self *IcnaStore) DropTable(databaseName string, tableName string) {
	database := self.GetDatabase(databaseName)
	table := database.GetTable(tableName)
	if table == nil {
		panic("table not exists: " + tableName)
	}
	os.Remove(table.DataPath)
	os.Remove(table.MetaPath)
	delete(database.TableMap, tableName)
}

func (self *IcnaStore) GetTable(databaseName string, tableName string) *meta.Table {
	database := self.GetDatabase(databaseName)
	table := database.GetTable(tableName)
	if table == nil {
		panic("table not exists: " + tableName)
	}
	return table
}

func (self *IcnaStore) ExistTable(databaseName string, tableName string) bool {
	database := self.GetDatabase(databaseName)
	table := database.GetTable(tableName)
	return table != nil
}

func (self *IcnaStore) Insert(databaseName string, tableName string, columns []string, rows [][]meta.Value) {
	table := self.GetTable(databaseName, tableName)
	columnMap := make(map[string]int, len(columns))
	for i, column := range columns {
		columnMap[column] = i
	}
	hasColumn := len(columnMap) > 0
	for _, values := range rows {
		fullValues := make([]meta.Value, len(table.Fields))
		if hasColumn {
			for i, field := range table.Fields {
				if index, ok := columnMap[field.Name]; ok {
					fullValues[i] = values[index]
				}
			}
		} else {
			for i := 0; i < min(len(table.Fields), len(values)); i++ {
				fullValues[i] = values[i]
			}
		}
		desc := meta.NewIndexDescByAllArgs(table.Fields, table.PrimaryFiled, table.FieldMap)
		entry := meta.NewClusterIndexEntry(fullValues, desc)
		table.Insert(entry)
	}
}
