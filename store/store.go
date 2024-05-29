package store

import "Relatdb/meta"

type Store interface {
	Init()
	CreateDatabase(database *meta.DataBase)
	DropDatabase(databaseName string)
	CreateTable(table *meta.Table)
	DropTable(databaseName string, tableName string)
}
