package ast

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

type SelectStatement struct {
	_DMLStatement_

	SelectIndex uint64
	Fields      []*SelectField
	From        ResultSet
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
