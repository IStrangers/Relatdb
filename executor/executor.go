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
	GetStore() store.Store
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

func (self *Executor) evalExpressionOrDefaultValue(expr ast.Expression, defaultValue any) meta.Value {
	if expr == nil {
		return meta.ToValue(defaultValue)
	}
	return self.evalExpression(expr)
}

func (self *Executor) evalExpression(expr ast.Expression) meta.Value {
	switch expr := expr.(type) {
	case *ast.TableName:
		return self.evalExpression(expr.Name)
	case *ast.Identifier:
		return meta.ToValue(expr.Name)
	case *ast.StringLiteral:
		return meta.ToValue(expr.Value)
	case *ast.NumberLiteral:
		return meta.ToValue(expr.Value)
	case *ast.BooleanLiteral:
		if expr.Value {
			return meta.ToValue(1)
		}
		return meta.ToValue(0)
	default:
		panic(fmt.Errorf("unsupported expression type: %T", expr))
	}
}

func (self *Executor) Execute() RecordSet {
	switch stmt := self.stmt.(type) {
	case *ast.ShowStatement:
		return self.executeShowStatement(stmt)
	case *ast.CreateTableStatement:
		return self.executeCreateTableStatement(stmt)
	case *ast.SelectStatement:
		return self.executeSelectStatement(stmt)
	default:
		panic(fmt.Errorf("unsupported statement type: %T", stmt))
	}
}

func (self *Executor) executeShowStatement(stmt *ast.ShowStatement) RecordSet {
	keyWord := self.evalExpressionOrDefaultValue(stmt.KeyWord, "")
	rows := make([][]meta.Value, 0)
	switch stmt.Type {
	case ast.ShowEngines:
	case ast.ShowDatabases:
		rows = append(rows, []meta.Value{meta.StringValue("default")})
	case ast.ShowTables:
		rows = append(rows, []meta.Value{meta.StringValue("test")})
	case ast.ShowColumns:
	case ast.ShowVariables:
	}
	return NewRecordSet(0, 0, []meta.Value{keyWord}, rows)
}

func (self *Executor) executeCreateTableStatement(stmt *ast.CreateTableStatement) RecordSet {
	fieldLength := len(stmt.ColumnDefinitions)
	fields := make([]*meta.Field, fieldLength)
	var primaryFiled *meta.Field
	fieldMap := make(map[string]uint, fieldLength)
	var clusterIndex meta.Index
	var secondaryIndexes []meta.Index
	for i, definition := range stmt.ColumnDefinitions {
		field := meta.NewField(
			uint(i), self.evalExpression(definition.Name).ToString(), definition.Type, definition.Flag,
			meta.ToValue(self.evalExpressionOrDefaultValue(definition.DefaultValue, nil)),
			self.evalExpressionOrDefaultValue(definition.Comment, "").ToString(),
		)
		if field.Flag&common.PRIMARY_KEY_FLAG != 0 {
			primaryFiled = field
			clusterIndex = bptree.NewBPTree(field.Name, []*meta.Field{primaryFiled}, field.Flag)
		}
		fields[i] = field
		fieldMap[field.Name] = field.Index
	}
	table := meta.NewTable(self.evalExpression(stmt.Name).ToString(), fields, primaryFiled, fieldMap, clusterIndex, secondaryIndexes)
	store := self.ctx.GetStore()
	store.CreateTable(table)
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeSelectStatement(stmt *ast.SelectStatement) RecordSet {
	columns := make([]meta.Value, len(stmt.Fields))
	rows := make([][]meta.Value, 0)
	for i, field := range stmt.Fields {
		columns[i] = self.evalExpressionOrDefaultValue(field.AsName, self.evalExpression(field.Expr))
	}
	return NewRecordSet(0, 0, columns, rows)
}
