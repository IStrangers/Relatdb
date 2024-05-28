package context

import (
	"Relatdb/store"
)

type Connection interface {
	GetDatabase() string
	SetDatabase(database string)
}

type Session interface {
	GetVariable(name string) string
	SetVariable(name string, value string)
}

type ExecuteContext interface {
	GetConnection() Connection
	GetSession() Session
	GetStore() store.Store
}
