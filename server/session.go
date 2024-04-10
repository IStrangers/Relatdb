package server

import (
	"Relatdb/executor"
	"Relatdb/parser/ast"
)

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (s Session) executeStmt(stmt ast.Statement) (executor.RecordSet, error) {
	executor := executor.NewExecutor(stmt)
	return executor.Execute()
}
