package ast

type Node interface {
	StartIndex() uint64
	EndIndex() uint64
}

type Expression interface {
	Node
	expression()
}

type Statement interface {
	Node
	statement()
}

type _Expression_ struct {
	Expression
}

func (self *_Expression_) expression() {
}

type _Statement_ struct {
	Statement
}

func (self *_Statement_) statement() {
}

type ResultSet interface {
	Expression

	resultSet()
}

type _ResultSet_ struct {
	ResultSet
}

func (self *_ResultSet_) resultSet() {
}

type Identifier struct {
	_Expression_

	Index uint64
	Name  string
}

func (self *Identifier) StartIndex() uint64 {
	return self.Index
}

func (self *Identifier) EndIndex() uint64 {
	return self.Index + uint64(len(self.Name))
}

type Assignment struct {
	_Expression_

	Left  Expression
	Right Expression
}

func (self *Assignment) StartIndex() uint64 {
	return self.Left.StartIndex()
}

func (self *Assignment) EndIndex() uint64 {
	return self.Right.EndIndex()
}

type TableName struct {
	_Expression_
	_ResultSet_

	Schema *Identifier
	Name   *Identifier
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

type ColumnName struct {
	_Expression_

	Schema *Identifier
	Table  *Identifier
	Name   *Identifier
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

type SelectField struct {
	_Expression_

	Expr   Expression
	AsName *Identifier
}

func (self *SelectField) StartIndex() uint64 {
	return self.Expr.StartIndex()
}

func (self *SelectField) EndIndex() uint64 {
	if self.AsName != nil {
		return self.AsName.EndIndex()
	}
	return self.Expr.EndIndex()
}

type TableRefsClause struct {
	_Statement_

	TableRefs *Join
}

func (self *TableRefsClause) StartIndex() uint64 {
	return self.TableRefs.StartIndex()
}

func (self *TableRefsClause) EndIndex() uint64 {
	return self.TableRefs.EndIndex()
}

type JoinType int

const (
	CrossJoin JoinType = iota + 1
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

type OnCondition struct {
	_Expression_

	Expr Expression
}

func (self *OnCondition) StartIndex() uint64 {
	return self.Expr.StartIndex()
}

func (self *OnCondition) EndIndex() uint64 {
	return self.Expr.EndIndex()
}

type GroupByClause struct {
	_Statement_

	GroupByIndex uint64
	Items        []*ColumnName
}

func (self *GroupByClause) StartIndex() uint64 {
	return self.GroupByIndex
}

func (self *GroupByClause) EndIndex() uint64 {
	return self.Items[len(self.Items)-1].EndIndex()
}

type HavingClause struct {
	_Statement_

	HavingIndex uint64
	Expr        Expression
}

func (self *HavingClause) StartIndex() uint64 {
	return self.HavingIndex
}

func (self *HavingClause) EndIndex() uint64 {
	return self.Expr.EndIndex()
}

type OrderByClause struct {
	_Statement_

	OrderByIndex uint64
	Items        []*OrderItem
}

func (self *OrderByClause) StartIndex() uint64 {
	return self.OrderByIndex
}

func (self *OrderByClause) EndIndex() uint64 {
	return self.Items[len(self.Items)-1].EndIndex()
}

type OrderItem struct {
	_Statement_

	ColumnName *ColumnName
	Order      *Identifier
}

func (self *OrderItem) StartIndex() uint64 {
	return self.ColumnName.StartIndex()
}

func (self *OrderItem) EndIndex() uint64 {
	if self.Order != nil {
		return self.Order.EndIndex()
	}
	return self.ColumnName.EndIndex()
}

type Limit struct {
	_Statement_

	LimitIndex uint64
	Offset     Expression
	Count      Expression
}

func (self *Limit) StartIndex() uint64 {
	return self.LimitIndex
}

func (self *Limit) EndIndex() uint64 {
	return self.Count.EndIndex()
}
