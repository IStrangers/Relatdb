package parser

type UseStatement struct {
	_Statement_

	UseIndex uint64
	Database *Identifier
}

func (self *UseStatement) StartIndex() uint64 {
	return self.UseIndex
}

func (self *UseStatement) EndIndex() uint64 {
	return self.Database.EndIndex()
}
