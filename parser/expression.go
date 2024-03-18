package parser

import "Relatdb/parser/ast"

func (self *Parser) parseExpression() ast.Expression {
	switch self.token {
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseTableName() *ast.TableName {
	return &ast.TableName{}
}

func (self *Parser) parseTableNames() (names []*ast.TableName) {
	return
}
