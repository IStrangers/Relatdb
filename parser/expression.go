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

func (parser *Parser) parseNumberLiteral() ast.Expression {
	defer parser.expect(NUMBER)
	return &ast.NumberLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.parseNumberLiteralValue(parser.value),
	}
}

func (parser *Parser) parseNumberLiteralValue(literal string) any {
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

func (parser *Parser) parseStringLiteral() ast.Expression {
	defer parser.expectToken(STRING)
	return &ast.StringLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.value,
	}
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
	defer parser.expectToken(BOOLEAN)
	return &ast.BooleanLiteral{
		Index: parser.index,
		Value: parser.value == "true",
	}
}

func (parser *Parser) parseNullLiteral() ast.Expression {
	defer parser.expectToken(NULL)
	return &ast.NullLiteral{
		Index: parser.index,
	}
}

func (self *Parser) parseIdentifier() *ast.Identifier {
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
	for self.token == COMMA {
		names = append(names, self.parseTableName())
	}
	return
}
