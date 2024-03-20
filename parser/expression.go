package parser

import (
	"strconv"
)

func (self *Parser) parseExpression() Expression {
	left := self.parseAssignExpression()

	return left
}

func (self *Parser) parseAssignExpression() Expression {
	left := self.parseConditionalExpression()

	if self.token == ASSIGN {
		operator := self.token
		self.expectToken(operator)
		return &AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    self.parseAssignExpression(),
		}
	}

	return left
}

func (self *Parser) parseConditionalExpression() Expression {
	left := self.parseLogicalOrExpression()

	return left
}

func (self *Parser) parseLogicalOrExpression() Expression {
	left := self.parseLogicalAndExpression()

	for {
		switch self.token {
		case OR:
			left = &BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseLogicalAndExpression(),
			}
		default:
			return left
		}
	}
}

func (self *Parser) parseLogicalAndExpression() Expression {
	left := self.parseEqualityExpression()

	for {
		switch self.token {
		case AND:
			left = &BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseEqualityExpression(),
			}
		default:
			return left
		}
	}
}

func (self *Parser) parseEqualityExpression() Expression {
	left := self.parseRelationalExpression()

	for {
		if self.token == EQUAL || self.token == NOT_EQUAL || (self.scope.inWhere && self.token == ASSIGN) {
			left = &BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseRelationalExpression(),
			}
		} else {
			return left
		}
	}
}

func (self *Parser) parseRelationalExpression() Expression {
	left := self.parseAdditiveExpression()

	for {
		switch self.token {
		case LESS, LESS_OR_EQUAL, GREATER, GREATER_OR_EQUAL:
			left = &BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseAdditiveExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseAdditiveExpression() Expression {
	left := parser.parseMultiplicativeExpression()

	for {
		switch parser.token {
		case ADDITION, SUBTRACT:
			left = &BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseMultiplicativeExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseMultiplicativeExpression() Expression {
	left := parser.parseUnaryExpression()

	for {
		switch parser.token {
		case MULTIPLY, DIVIDE, REMAINDER:
			left = &BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseUnaryExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseUnaryExpression() Expression {

	tkn := parser.token
	switch tkn {
	case NOT, ADDITION, SUBTRACT:
		unaryExpression := &UnaryExpression{
			Index:    parser.expect(tkn),
			Operator: tkn,
			Operand:  parser.parseUnaryExpression(),
		}
		return unaryExpression
	}

	left := parser.parseLeftHandSideExpressionAllowCall([]Token{})

	return left
}

func (parser *Parser) parseLeftHandSideExpressionAllowCall(stopTokens []Token) (left Expression) {
	isStopToken := func(token Token) bool {
		for _, stopToken := range stopTokens {
			if token == stopToken {
				return true
			}
		}
		return false
	}

	left = parser.parsePrimaryExpression()

	for !isStopToken(parser.token) {
		switch parser.token {
		case LEFT_PARENTHESIS:
			left = parser.parseCallExpression(left)
			continue
		}
		break
	}

	return left
}

func (self *Parser) parsePrimaryExpression() Expression {
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

func (self *Parser) parseNumberLiteral() Expression {
	defer self.expect(NUMBER)
	return &NumberLiteral{
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

func (self *Parser) parseStringLiteral() Expression {
	defer self.expectToken(STRING)
	return &StringLiteral{
		Index:   self.index,
		Literal: self.literal,
		Value:   self.value,
	}
}

func (self *Parser) parseBooleanLiteral() Expression {
	defer self.expectToken(BOOLEAN)
	return &BooleanLiteral{
		Index: self.index,
		Value: self.value == "true",
	}
}

func (self *Parser) parseNullLiteral() Expression {
	defer self.expectToken(NULL)
	return &NullLiteral{
		Index: self.index,
	}
}

func (self *Parser) parseIdentifier() *Identifier {
	defer self.expectToken(IDENTIFIER)
	return &Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseKeyWordIdentifier(keyword Token) *Identifier {
	defer self.expectToken(keyword)
	return &Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseStringLiteralOrIdentifier() Expression {
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

func (self *Parser) parseTableName() *TableName {
	tableName := &TableName{
		Name: self.parseStringLiteralOrIdentifier(),
	}
	if self.expectEqualsToken(DOT) {
		tableName.Schema = tableName.Name
		tableName.Name = self.parseStringLiteralOrIdentifier()
	}
	return tableName
}

func (self *Parser) parseTableNames() (names []*TableName) {
	for {
		names = append(names, self.parseTableName())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseColumnName() *ColumnName {
	columnName := &ColumnName{
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

func (self *Parser) parseColumnNames() (names []*ColumnName) {
	for {
		names = append(names, self.parseColumnName())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseWhereExpression() Expression {
	self.scope.inWhere = true
	expr := self.parseExpression()
	self.scope.inWhere = false
	return expr
}

func (self *Parser) parseCallExpression(left Expression) Expression {
	leftParenthesis, arguments, rightParenthesis := self.parseArguments()
	return &CallExpression{
		Callee:           left,
		LeftParenthesis:  leftParenthesis,
		Arguments:        arguments,
		RightParenthesis: rightParenthesis,
	}
}

func (self *Parser) parseArguments() (leftParenthesis uint64, arguments []Expression, rightParenthesis uint64) {
	leftParenthesis = self.expect(LEFT_PARENTHESIS)
	for self.token != RIGHT_PARENTHESIS {
		arguments = append(arguments, self.parseExpression())
		if self.token != COMMA {
			break
		}
		self.expect(COMMA)
	}
	rightParenthesis = self.expect(RIGHT_PARENTHESIS)
	return
}
