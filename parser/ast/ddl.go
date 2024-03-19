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
}

func (self *CreateTableStatement) StartIndex() uint64 {
	return self.CreateIndex
}

func (self *CreateTableStatement) EndIndex() uint64 {
	return self.Name.EndIndex()
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
	TableName *TableName
	Name      Expression
}

func (self *DropIndexStatement) StartIndex() uint64 {
	return self.DropIndex
}

func (self *DropIndexStatement) EndIndex() uint64 {
	return self.Name.EndIndex()
}
