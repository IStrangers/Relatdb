package executor

import (
	"Relatdb/parser/ast"
)

type Executor struct {
	stmt ast.Statement
}

func NewExecutor(stmt ast.Statement) *Executor {
	return &Executor{}
}

func (self *Executor) Execute() (RecordSet, error) {
	switch stmt := self.stmt.(type) {
	case *ast.CreateTableStatement:
		return self.executeCreateTableStatement(stmt)
	}
	return nil, nil
}

func (self *Executor) executeCreateTableStatement(stmt *ast.CreateTableStatement) (RecordSet, error) {
	return nil, nil
}
