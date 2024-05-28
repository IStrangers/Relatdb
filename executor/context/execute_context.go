package context

import "Relatdb/store"

type Connection interface {
	GetDatabase() string
}

type ExecuteContext interface {
	GetConnection() Connection
	GetStore() store.Store
}
