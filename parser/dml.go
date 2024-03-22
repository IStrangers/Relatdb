package parser

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

type SelectField struct {
	_Expression_

	Expr   Expression
	AsName Expression
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
	Desc       bool
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

type SelectStatement struct {
	_DMLStatement_

	SelectIndex uint64
	Fields      []*SelectField
	From        *TableRefsClause
	Where       Expression
	GroupBy     *GroupByClause
	Having      *HavingClause
	Order       *OrderByClause
	Limit       *Limit
}

func (self *SelectStatement) StartIndex() uint64 {
	return self.SelectIndex
}

func (self *SelectStatement) EndIndex() uint64 {
	if self.Limit != nil {
		return self.Limit.EndIndex()
	}
	if self.Order != nil {
		return self.Order.EndIndex()
	}
	if self.Having != nil {
		return self.Having.EndIndex()
	}
	if self.GroupBy != nil {
		return self.GroupBy.EndIndex()
	}
	if self.Where != nil {
		return self.Where.EndIndex()
	}
	return self.From.EndIndex()
}

type SubqueryExpression struct {
	_ResultSet_

	LeftParenthesis  uint64
	Select           *SelectStatement
	RightParenthesis uint64
	AsName           Expression
}

func (self *SubqueryExpression) StartIndex() uint64 {
	return self.LeftParenthesis
}

func (self *SubqueryExpression) EndIndex() uint64 {
	if self.AsName != nil {
		return self.AsName.EndIndex()
	}
	return self.RightParenthesis
}
