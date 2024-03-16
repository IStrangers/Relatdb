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
