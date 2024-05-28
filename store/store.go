package store

import "Relatdb/meta"

type Store interface {
	Init()
	CreateTable(table *meta.Table)
	DropTable(databaseName string, tableName string)
}
