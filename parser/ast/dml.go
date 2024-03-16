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
	AssignmentColumns []*Assignment
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
