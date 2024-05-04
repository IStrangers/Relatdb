package executor

import (
	"Relatdb/meta"
	"Relatdb/parser/ast"
	"Relatdb/store"
	"fmt"
)

type ExecuteContext interface {
	GetStore() *store.Store
}

type Executor struct {
	ctx  ExecuteContext
	stmt ast.Statement
}

func NewExecutor(ctx ExecuteContext, stmt ast.Statement) *Executor {
	return &Executor{
		ctx:  ctx,
		stmt: stmt,
	}
}

func (self *Executor) evalExpression(expr ast.Expression) any {
	switch expr := expr.(type) {
	case *ast.TableName:
		return self.evalExpression(expr.Name)
	case *ast.Identifier:
		return expr.Name
	case *ast.StringLiteral:
		return expr.Value
	default:
		panic(fmt.Errorf("unsupported expression type: %T", expr))
	}
}

func (self *Executor) Execute() (RecordSet, error) {
	switch stmt := self.stmt.(type) {
	case *ast.CreateTableStatement:
		return self.executeCreateTableStatement(stmt)
	}
	return nil, nil
}

func (self *Executor) executeCreateTableStatement(stmt *ast.CreateTableStatement) (RecordSet, error) {
	fields := make([]*meta.Field, len(stmt.ColumnDefinitions))
	for i, definition := range stmt.ColumnDefinitions {
		fields = append(fields, meta.NewField(
			uint(i), self.evalExpression(definition.Name).(string), definition.Type,
			self.evalExpression(definition.Comment).(string), definition.Flag),
		)
	}
	table := meta.NewTable(self.evalExpression(stmt.Name).(string), fields)
	store := self.ctx.GetStore()
	err := store.CreateTable(table)
	return nil, err
}
