package executor

import "Relatdb/parser/ast"

type Executor struct {
	stmt ast.Statement
}

func NewExecutor(stmt ast.Statement) *Executor {
	return &Executor{
		stmt: stmt,
	}
}

func (self *Executor) Execute() (RecordSet, error) {
	return nil, nil
}
