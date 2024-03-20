package parser

import (
	"Relatdb/parser/ast"
	"strconv"
)

func (self *Parser) parseExpression() ast.Expression {
	switch self.token {
	case NUMBER:
		return self.parseNumberLiteral()
	case STRING:
		return self.parseStringLiteral()
	case BOOLEAN:
		return self.parseBooleanLiteral()
	case NULL:
		return self.parseNullLiteral()
	case IDENTIFIER:
		return self.parseIdentifier()
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseNumberLiteral() ast.Expression {
	defer self.expect(NUMBER)
	return &ast.NumberLiteral{
		Index:   self.index,
		Literal: self.literal,
		Value:   self.parseNumberLiteralValue(self.value),
	}
}

func (self *Parser) parseNumberLiteralValue(literal string) any {
	var value any = 0
	updateValue := func(v any, err error) bool {
		if err != nil {
			return false
		}
		value = v
		return true
	}
	intValue, err := strconv.ParseInt(literal, 0, 64)
	if updateValue(intValue, err) {
		return value
	}
	floatValue, err := strconv.ParseFloat(literal, 64)
	if updateValue(floatValue, err) {
		return value
	}
	return value
}

func (self *Parser) parseStringLiteral() ast.Expression {
	defer self.expectToken(STRING)
	return &ast.StringLiteral{
		Index:   self.index,
		Literal: self.literal,
		Value:   self.value,
	}
}

func (self *Parser) parseBooleanLiteral() ast.Expression {
	defer self.expectToken(BOOLEAN)
	return &ast.BooleanLiteral{
		Index: self.index,
		Value: self.value == "true",
	}
}

func (self *Parser) parseNullLiteral() ast.Expression {
	defer self.expectToken(NULL)
	return &ast.NullLiteral{
		Index: self.index,
	}
}

func (self *Parser) parseIdentifier() *ast.Identifier {
	defer self.expectToken(IDENTIFIER)
	return &ast.Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseKeyWordIdentifier(keyword Token) *ast.Identifier {
	defer self.expectToken(keyword)
	return &ast.Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseStringLiteralOrIdentifier() ast.Expression {
	switch self.token {
	case STRING:
		return self.parseStringLiteral()
	case IDENTIFIER:
		return self.parseIdentifier()
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseTableName() *ast.TableName {
	tableName := &ast.TableName{
		Name: self.parseStringLiteralOrIdentifier(),
	}
	if self.expectEqualsToken(DOT) {
		tableName.Schema = tableName.Name
		tableName.Name = self.parseStringLiteralOrIdentifier()
	}
	return tableName
}

func (self *Parser) parseTableNames() (names []*ast.TableName) {
	for {
		names = append(names, self.parseTableName())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseColumnName() *ast.ColumnName {
	columnName := &ast.ColumnName{
		Name: self.parseStringLiteralOrIdentifier(),
	}
	if self.expectEqualsToken(DOT) {
		columnName.Table = columnName.Name
		columnName.Name = self.parseStringLiteralOrIdentifier()
	}
	if self.expectEqualsToken(DOT) {
		columnName.Schema = columnName.Table
		columnName.Table = columnName.Name
		columnName.Name = self.parseStringLiteralOrIdentifier()
	}
	return columnName
}

func (self *Parser) parseColumnNames() (names []*ast.ColumnName) {
	for {
		names = append(names, self.parseColumnName())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}
