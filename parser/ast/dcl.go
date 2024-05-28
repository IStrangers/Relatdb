package ast

type UseStatement struct {
	_Statement_

	UseIndex uint64
	Database Expression
}

func (self *UseStatement) StartIndex() uint64 {
	return self.UseIndex
}

func (self *UseStatement) EndIndex() uint64 {
	return self.Database.EndIndex()
}

type SetVariableStatement struct {
	_Statement_

	SetIndex uint64
	Name     *Identifier
	Value    Expression
}

func (self *SetVariableStatement) StartIndex() uint64 {
	return self.SetIndex
}

func (self *SetVariableStatement) EndIndex() uint64 {
	return self.Value.EndIndex()
}
