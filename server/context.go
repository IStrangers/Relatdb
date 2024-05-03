package server

import (
	"Relatdb/executor"
	"Relatdb/parser/ast"
)

type Context struct {
	conn    *Connection
	session *Session
}

func NewContext(conn *Connection) *Context {
	return &Context{
		conn:    conn,
		session: NewSession(),
	}
}

func (self *Context) executeStmt(stmt ast.Statement) {
	executor := executor.NewExecutor(stmt)
	recordSet, err := executor.Execute()
	if err != nil {
		return
	}
	println(recordSet)
}
