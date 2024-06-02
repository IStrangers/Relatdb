package store

import "Relatdb/meta"

type Store interface {
	Init()
	CreateDatabase(database *meta.DataBase)
	DropDatabase(databaseName string)
	GetDatabase(databaseName string) *meta.DataBase
	CreateTable(table *meta.Table)
	DropTable(databaseName string, tableName string)
	GetTable(databaseName string, tableName string) *meta.Table
	Insert(databaseName string, tableName string, columns []string, rows [][]meta.Value)
}
