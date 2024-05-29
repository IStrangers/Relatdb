package executor

import (
	"Relatdb/common"
	"Relatdb/executor/context"
	"Relatdb/index/bptree"
	"Relatdb/meta"
	"Relatdb/parser/ast"
	"fmt"
)

type Executor struct {
	ctx  context.ExecuteContext
	stmt ast.Statement
}

func NewExecutor(ctx context.ExecuteContext, stmt ast.Statement) *Executor {
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
		return meta.StringValue(expr.Name)
	case *ast.StringLiteral:
		return meta.StringValue(expr.Value)
	case *ast.NumberLiteral:
		return meta.ToValue(expr.Value)
	case *ast.BooleanLiteral:
		if expr.Value {
			return meta.IntValue(1)
		}
		return meta.IntValue(0)
	case *ast.VariableName:
		variableName := self.evalExpression(expr.Name).ToString()
		return meta.StringValue(self.ctx.GetSession().GetVariable(variableName))
	case *ast.VariableRef:
		variableName := self.evalExpression(expr.Name).ToString()
		return meta.StringValue(self.ctx.GetSession().GetVariable(variableName))
	default:
		panic(fmt.Errorf("unsupported expression type: %T", expr))
	}
}

func (self *Executor) Execute() RecordSet {
	switch stmt := self.stmt.(type) {
	case *ast.CreateDatabaseStatement:
		return self.executeCreateDatabaseStatement(stmt)
	case *ast.DropDatabaseStatement:
		return self.executeDropDatabaseStatement(stmt)
	case *ast.UseStatement:
		return self.executeUseStatement(stmt)
	case *ast.ShowStatement:
		return self.executeShowStatement(stmt)
	case *ast.SetVariableStatement:
		return self.executeSetVariableStatement(stmt)
	case *ast.CreateTableStatement:
		return self.executeCreateTableStatement(stmt)
	case *ast.DropTableStatement:
		return self.executeDropTableStatement(stmt)
	case *ast.SelectStatement:
		return self.executeSelectStatement(stmt)
	default:
		panic(fmt.Errorf("unsupported statement type: %T", stmt))
	}
}

func (self *Executor) executeCreateDatabaseStatement(stmt *ast.CreateDatabaseStatement) RecordSet {
	store := self.ctx.GetStore()
	databaseName := self.evalExpression(stmt.Name).ToString()
	database := meta.NewDataBase(databaseName)
	store.CreateDatabase(database)
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeDropDatabaseStatement(stmt *ast.DropDatabaseStatement) RecordSet {
	store := self.ctx.GetStore()
	databaseName := self.evalExpression(stmt.Name).ToString()
	store.DropDatabase(databaseName)
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeUseStatement(stmt *ast.UseStatement) RecordSet {
	connection := self.ctx.GetConnection()
	connection.SetDatabase(self.evalExpression(stmt.Database).ToString())
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeShowStatement(stmt *ast.ShowStatement) RecordSet {
	columns := make([]meta.Value, 0)
	rows := make([][]meta.Value, 0)
	switch stmt.Type {
	case ast.ShowEngines:
		columns = []meta.Value{
			meta.StringValue("Engine"), meta.StringValue("Support"), meta.StringValue("Comment"),
			meta.StringValue("Transactions"), meta.StringValue("XA"), meta.StringValue("Savepoints"),
		}
	case ast.ShowDatabases:
		columns = []meta.Value{meta.StringValue("Database")}
		rows = append(rows, []meta.Value{meta.StringValue("default")})
	case ast.ShowTables:
		columns = []meta.Value{meta.StringValue("Tables_in_")}
		rows = append(rows, []meta.Value{meta.StringValue("test")})
	case ast.ShowColumns:
		columns = []meta.Value{
			meta.StringValue("Field"), meta.StringValue("Type"), meta.StringValue("Null"),
			meta.StringValue("Key"), meta.StringValue("Default"), meta.StringValue("Extra"),
		}
	case ast.ShowVariables:
		columns = []meta.Value{meta.StringValue("Variable_name"), meta.StringValue("Value")}
	case ast.ShowStatus:
		columns = []meta.Value{meta.StringValue("Variable_name"), meta.StringValue("Value")}
	}
	return NewRecordSet(0, 0, columns, rows)
}

func (self *Executor) executeSetVariableStatement(stmt *ast.SetVariableStatement) RecordSet {
	session := self.ctx.GetSession()
	name := self.evalExpression(stmt.Name).ToString()
	value := self.evalExpression(stmt.Value).ToString()
	session.SetVariable(name, value)
	return NewRecordSet(0, 0, nil, nil)
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
	connection := self.ctx.GetConnection()
	table := meta.NewTable(connection.GetDatabase(), self.evalExpression(stmt.Name).ToString(), fields, primaryFiled, fieldMap, clusterIndex, secondaryIndexes)
	store := self.ctx.GetStore()
	store.CreateTable(table)
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeDropTableStatement(stmt *ast.DropTableStatement) RecordSet {
	connection := self.ctx.GetConnection()
	store := self.ctx.GetStore()
	for _, name := range stmt.Names {
		store.DropTable(connection.GetDatabase(), self.evalExpression(name.Name).ToString())
	}
	return NewRecordSet(0, 0, nil, nil)
}

func (self *Executor) executeSelectStatement(stmt *ast.SelectStatement) RecordSet {
	columns := make([]meta.Value, len(stmt.Fields))
	rows := make([][]meta.Value, 0)
	row := make([]meta.Value, len(stmt.Fields))
	for i, field := range stmt.Fields {
		columns[i] = self.evalExpressionOrDefaultValue(field.AsName, self.evalExpression(field.Expr))
		row[i] = columns[i]
	}
	rows = append(rows, row)
	return NewRecordSet(0, 0, columns, rows)
}
