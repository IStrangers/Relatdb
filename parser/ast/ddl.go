package ast

type DDLStatement interface {
	Statement
	ddlStatement()
}

type _DDLStatement_ struct {
	_Statement_
}

func (self *_DDLStatement_) ddlStatement() {
}

type ColumnDefinition struct {
	_Statement_

	Name         Expression
	Type         byte       //字段类型
	Flag         uint       //字段标记: NotNull, Unsigned, PriKey
	Length       int        //字段长度
	Decimal      int        //小数位数
	DefaultValue Expression //默认值
	Comment      Expression //注释
}

func (self *ColumnDefinition) StartIndex() uint64 {
	return self.Name.StartIndex()
}

func (self *ColumnDefinition) EndIndex() uint64 {
	if self.Comment != nil {
		return self.Comment.EndIndex()
	}
	if self.DefaultValue != nil {
		return self.DefaultValue.EndIndex()
	}
	return self.Name.EndIndex()
}

type CreateDatabaseStatement struct {
	_DDLStatement_

	CreateIndex uint64
	IfNotExists bool
	Name        Expression
}

func (self *CreateDatabaseStatement) StartIndex() uint64 {
	return self.CreateIndex
}

func (self *CreateDatabaseStatement) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type DropDatabaseStatement struct {
	_DDLStatement_

	DropIndex uint64
	IfExists  bool
	Name      Expression
}

func (self *DropDatabaseStatement) StartIndex() uint64 {
	return self.DropIndex
}

func (self *DropDatabaseStatement) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type CreateTableStatement struct {
	_DDLStatement_

	CreateIndex       uint64
	IfNotExists       bool
	Name              *TableName
	ColumnDefinitions []*ColumnDefinition
	RightParenthesis  uint64
}

func (self *CreateTableStatement) StartIndex() uint64 {
	return self.CreateIndex
}

func (self *CreateTableStatement) EndIndex() uint64 {
	return self.RightParenthesis
}

type DropTableStatement struct {
	_DDLStatement_

	DropIndex uint64
	IfExists  bool
	Names     []*TableName
}

func (self *DropTableStatement) StartIndex() uint64 {
	return self.DropIndex
}

func (self *DropTableStatement) EndIndex() uint64 {
	return self.Names[len(self.Names)-1].EndIndex()
}

type IndexType int

const (
	IndexTypeNone IndexType = iota
	IndexTypeUnique
	IndexTypeSpatial
	IndexTypeFullText
)

type CreateIndexStatement struct {
	_DDLStatement_

	CreateIndex uint64
	IfNotExists bool
	Name        Expression
	TableName   *TableName
	ColumnNames []*ColumnName
	Type        IndexType
}

func (self *CreateIndexStatement) StartIndex() uint64 {
	return self.CreateIndex
}

func (self *CreateIndexStatement) EndIndex() uint64 {
	return self.Name.EndIndex()
}

type DropIndexStatement struct {
	_DDLStatement_

	DropIndex uint64
	IfExists  bool
	Name      Expression
	TableName *TableName
}

func (self *DropIndexStatement) StartIndex() uint64 {
	return self.DropIndex
}

func (self *DropIndexStatement) EndIndex() uint64 {
	return self.TableName.EndIndex()
}
