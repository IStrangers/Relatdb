package server

import "Relatdb/parser/ast"

type Context struct {
	session *Session
}

func NewContext() *Context {
	return &Context{
		session: NewSession(),
	}
}

func (self *Context) executeStmt(stmt ast.Statement) {
	recordSet, err := self.session.executeStmt(stmt)
	if err != nil {
		return
	}
	println(recordSet)
}
