package executor

import (
	"Relatdb/common"
	"Relatdb/index/bptree"
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

func (self *Executor) Execute() RecordSet {
	switch stmt := self.stmt.(type) {
	case *ast.CreateTableStatement:
		return self.executeCreateTableStatement(stmt)
	}
	return nil
}

func (self *Executor) executeCreateTableStatement(stmt *ast.CreateTableStatement) RecordSet {
	fields := make([]*meta.Field, len(stmt.ColumnDefinitions))
	var primaryFiled *meta.Field
	var fieldMap map[string]uint
	var clusterIndex meta.Index
	var secondaryIndexes []meta.Index
	for i, definition := range stmt.ColumnDefinitions {
		comment := ""
		if definition.Comment != nil {
			comment = self.evalExpression(definition.Comment).(string)
		}
		field := meta.NewField(
			uint(i), self.evalExpression(definition.Name).(string),
			definition.Type, comment, definition.Flag,
		)
		if field.Flag&common.PRIMARY_KEY_FLAG != 0 {
			primaryFiled = field
			clusterIndex = bptree.NewBPTree(field.Name, []*meta.Field{primaryFiled}, field.Flag)
		}
		fields[i] = field
		fieldMap[field.Name] = field.Index
	}
	table := meta.NewTable(self.evalExpression(stmt.Name).(string), fields, primaryFiled, fieldMap, clusterIndex, secondaryIndexes)
	store := self.ctx.GetStore()
	store.CreateTable(table)
	return nil
}
