package ast

type DMLStatement interface {
	Statement
	dmlStatement()
}

type _DMLStatement_ struct {
	_Statement_
}

func (self *_DMLStatement_) dmlStatement() {
}

type TableName struct {
	_Expression_

	Schema Expression
	Name   Expression
}

func (self *TableName) StartIndex() uint64 {
	if self.Schema != nil {
		return self.Schema.StartIndex()
	}
	return self.Name.StartIndex()
}

func (self *TableName) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type TableSource struct {
	_ResultSet_

	TableName *TableName
	AsName    Expression
}

func (self *TableSource) StartIndex() uint64 {
	return self.TableName.StartIndex()
}

func (self *TableSource) EndIndex() uint64 {
	if self.AsName != nil {
		return self.AsName.EndIndex()
	}
	return self.TableName.EndIndex()
}

type ColumnName struct {
	_Expression_

	Schema Expression
	Table  Expression
	Name   Expression
}

func (self *ColumnName) StartIndex() uint64 {
	if self.Schema != nil {
		return self.Schema.StartIndex()
	}
	if self.Table != nil {
		return self.Table.StartIndex()
	}
	return self.Name.StartIndex()
}

func (self *ColumnName) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type VariableName struct {
	_Expression_

	AtIndex uint64
	Name    *Identifier
}

func (self *VariableName) StartIndex() uint64 {
	return self.AtIndex
}

func (self *VariableName) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type VariableRef struct {
	_Expression_

	AtIndex uint64
	Name    *Identifier
}

func (self *VariableRef) StartIndex() uint64 {
	return self.AtIndex
}

func (self *VariableRef) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type JoinType int

const (
	CrossJoin JoinType = iota + 1
	InnerJoin
	LeftJoin
	RightJoin
)

type Join struct {
	_Statement_
	_ResultSet_

	Left     ResultSet
	Right    ResultSet
	JoinType JoinType
	On       *OnCondition
}

func (self *Join) StartIndex() uint64 {
	return self.Left.StartIndex()
}

func (self *Join) EndIndex() uint64 {
	if self.On != nil {
		return self.On.EndIndex()
	}
	if self.Right != nil {
		return self.Right.EndIndex()
	}
	if self.Left != nil {
		return self.Left.EndIndex()
	}
	return self.Left.EndIndex()
}

type ShowStatementType int

const (
	_ = iota
	ShowEngines
	ShowDatabases
	ShowTables
	ShowColumns
	ShowVariables
)

type ShowStatement struct {
	_DMLStatement_

	ShowIndex uint64
	Type      ShowStatementType
	KeyWord   *Identifier
}

func (self *ShowStatement) StartIndex() uint64 {
	return self.ShowIndex
}

func (self *ShowStatement) EndIndex() uint64 {
	switch self.Type {
	case ShowDatabases, ShowTables, ShowVariables:
		return self.KeyWord.EndIndex()
	}
	return self.ShowIndex
}

type InsertStatement struct {
	_DMLStatement_

	InsertIndex uint64
	TableName   *TableName
	ColumnNames []*ColumnName
	Values      [][]Expression
}

func (self *InsertStatement) StartIndex() uint64 {
	return self.InsertIndex
}

func (self *InsertStatement) EndIndex() uint64 {
	lastRow := self.Values[len(self.Values)-1]
	return lastRow[len(lastRow)-1].EndIndex()
}

type DeleteStatement struct {
	_DMLStatement_

	DeleteIndex uint64
	TableName   *TableName
	Where       Expression
	Order       *OrderByClause
	Limit       *Limit
}

func (self *DeleteStatement) StartIndex() uint64 {
	return self.DeleteIndex
}

func (self *DeleteStatement) EndIndex() uint64 {
	if self.Limit != nil {
		return self.Limit.EndIndex()
	}
	if self.Order != nil {
		return self.Order.EndIndex()
	}
	if self.Where != nil {
		return self.Where.EndIndex()
	}
	return self.TableName.EndIndex()
}

type UpdateStatement struct {
	_DMLStatement_

	UpdateIndex       uint64
	TableName         *TableName
	AssignExpressions []Expression
	Where             Expression
	Order             *OrderByClause
	Limit             *Limit
}

func (self *UpdateStatement) StartIndex() uint64 {
	return self.UpdateIndex
}

func (self *UpdateStatement) EndIndex() uint64 {
	if self.Limit != nil {
		return self.Limit.EndIndex()
	}
	if self.Order != nil {
		return self.Order.EndIndex()
	}
	if self.Where != nil {
		return self.Where.EndIndex()
	}
	return self.TableName.EndIndex()
}
