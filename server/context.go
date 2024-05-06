package server

import (
	"Relatdb/executor"
	"Relatdb/parser/ast"
	"Relatdb/store"
)

type Context struct {
	conn    *Connection
	session *Session
}

func (self *Context) GetStore() *store.Store {
	return self.conn.server.store
}

func NewContext(conn *Connection) *Context {
	return &Context{
		conn:    conn,
		session: NewSession(),
	}
}

func (self *Context) executeStmt(stmt ast.Statement) {
	executor := executor.NewExecutor(self, stmt)
	recordSet := executor.Execute()
	println(recordSet)
}
