package parser

import (
	"Relatdb/parser/ast"
	"Relatdb/parser/token"
	"fmt"
	"strconv"
)

func (self *Parser) parseExpression() ast.Expression {
	left := self.parseAssignExpression()

	return left
}

func (self *Parser) parseAssignExpression() ast.Expression {
	left := self.parseConditionalExpression()

	if self.token == token.ASSIGN {
		operator := self.token
		self.expectToken(operator)
		return &ast.AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    self.parseAssignExpression(),
		}
	}

	return left
}

func (self *Parser) parseConditionalExpression() ast.Expression {
	left := self.parseLogicalOrExpression()

	return left
}

func (self *Parser) parseLogicalOrExpression() ast.Expression {
	left := self.parseLogicalAndExpression()

	for {
		switch self.token {
		case token.OR:
			left = &ast.BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseLogicalAndExpression(),
			}
		default:
			return left
		}
	}
}

func (self *Parser) parseLogicalAndExpression() ast.Expression {
	left := self.parseEqualityExpression()

	for {
		switch self.token {
		case token.AND:
			left = &ast.BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseEqualityExpression(),
			}
		default:
			return left
		}
	}
}

func (self *Parser) parseEqualityExpression() ast.Expression {
	left := self.parseRelationalExpression()

	for {
		if self.token == token.EQUAL || self.token == token.NOT_EQUAL || (self.scope.inWhere && self.token == token.ASSIGN) {
			left = &ast.BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseRelationalExpression(),
			}
		} else {
			return left
		}
	}
}

func (self *Parser) parseRelationalExpression() ast.Expression {
	left := self.parseAdditiveExpression()

	for {
		switch self.token {
		case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUAL:
			left = &ast.BinaryExpression{
				Operator: self.expectToken(self.token),
				Left:     left,
				Right:    self.parseAdditiveExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseAdditiveExpression() ast.Expression {
	left := parser.parseMultiplicativeExpression()

	for {
		switch parser.token {
		case token.ADDITION, token.SUBTRACT:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseMultiplicativeExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseMultiplicativeExpression() ast.Expression {
	left := parser.parseUnaryExpression()

	for {
		switch parser.token {
		case token.MULTIPLY, token.DIVIDE, token.REMAINDER:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseUnaryExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseUnaryExpression() ast.Expression {

	tkn := parser.token
	switch tkn {
	case token.NOT, token.ADDITION, token.SUBTRACT:
		unaryExpression := &ast.UnaryExpression{
			Index:    parser.expect(tkn),
			Operator: tkn,
			Operand:  parser.parseUnaryExpression(),
		}
		return unaryExpression
	}

	left := parser.parseLeftHandSideExpressionAllowCall([]token.Token{})

	return left
}

func (parser *Parser) parseLeftHandSideExpressionAllowCall(stopTokens []token.Token) (left ast.Expression) {
	isStopToken := func(token token.Token) bool {
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
		case token.LEFT_PARENTHESIS:
			if columnName, ok := left.(*ast.ColumnName); ok {
				left = columnName.Name
			}
			left = parser.parseCallExpression(left)
			if parser.scope.inSelectField {
				identifier := left.(*ast.CallExpression).Callee.(*ast.Identifier)
				identifier.Name += "()"
				left = identifier
			}
			continue
		}
		break
	}

	return left
}

func (self *Parser) parsePrimaryExpression() ast.Expression {
	var expr ast.Expression

	switch self.token {
	case token.NUMBER:
		expr = self.parseNumberLiteral()
	case token.STRING:
		expr = self.parseStringLiteral()
	case token.BOOLEAN:
		expr = self.parseBooleanLiteral()
	case token.NULL:
		expr = self.parseNullLiteral()
	case token.IDENTIFIER:
		if self.scope.inSelectField || self.scope.inWhere {
			expr = self.parseColumnName()
		} else {
			expr = self.parseIdentifier()
		}
	case token.LEFT_PARENTHESIS:
		expr = self.parseSubqueryExpression()
	case token.AT_IDENTIFIER:
		atIndex := self.expect(token.AT_IDENTIFIER)
		if self.token == token.AT_IDENTIFIER {
			expr = self.parseVariableRef(atIndex)
		} else {
			expr = self.parseVariableName(atIndex)
		}
	default:
		if _, ok := token.IsKeyword(self.literal); ok {
			expr = self.parseKeyWordIdentifier(self.token)
		} else {
			self.errorUnexpectedToken(self.token)
		}
	}

	return expr
}

func (self *Parser) parseNumberLiteral() *ast.NumberLiteral {
	defer self.expect(token.NUMBER)
	numberLiteral := &ast.NumberLiteral{
		Index:   self.index,
		Literal: self.literal,
	}
	numberLiteral.Value, numberLiteral.IsDecimal = self.parseNumberLiteralValue(self.value)
	return numberLiteral
}

func (self *Parser) parseNumberLiteralValue(literal string) (any, bool) {
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
		return value, false
	}
	floatValue, err := strconv.ParseFloat(literal, 64)
	if updateValue(floatValue, err) {
		return value, true
	}
	return value, false
}

func (self *Parser) parseStringLiteral() *ast.StringLiteral {
	defer self.expectToken(token.STRING)
	return &ast.StringLiteral{
		Index:   self.index,
		Literal: self.literal,
		Value:   self.value,
	}
}

func (self *Parser) parseBooleanLiteral() *ast.BooleanLiteral {
	defer self.expectToken(token.BOOLEAN)
	return &ast.BooleanLiteral{
		Index: self.index,
		Value: self.value == "true",
	}
}

func (self *Parser) parseNullLiteral() *ast.NullLiteral {
	defer self.expectToken(token.NULL)
	return &ast.NullLiteral{
		Index: self.index,
	}
}

func (self *Parser) parseIdentifier() *ast.Identifier {
	defer self.expectToken(token.IDENTIFIER)
	return &ast.Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseStringLiteralOrIdentifier() ast.Expression {
	switch self.token {
	case token.STRING:
		return self.parseStringLiteral()
	default:
		return self.parseIdentifier()
	}
}

func (self *Parser) parseKeyWordIdentifier(keyword token.Token) *ast.Identifier {
	defer self.expectToken(keyword)
	return &ast.Identifier{
		Index: self.index,
		Name:  self.value,
	}
}

func (self *Parser) parseVariableRef(atIndex uint64) *ast.VariableRef {
	self.expectToken(token.AT_IDENTIFIER)
	return &ast.VariableRef{
		AtIndex: atIndex,
		Name:    self.parseIdentifier(),
	}
}

func (self *Parser) parseVariableName(atIndex uint64) *ast.VariableName {
	return &ast.VariableName{
		AtIndex: atIndex,
		Name:    self.parseIdentifier(),
	}
}

func (self *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	leftParenthesis, arguments, rightParenthesis := self.parseArguments()
	return &ast.CallExpression{
		Callee:           left,
		LeftParenthesis:  leftParenthesis,
		Arguments:        arguments,
		RightParenthesis: rightParenthesis,
	}
}

func (self *Parser) parseArguments() (leftParenthesis uint64, arguments []ast.Expression, rightParenthesis uint64) {
	leftParenthesis = self.expect(token.LEFT_PARENTHESIS)
	for self.token != token.RIGHT_PARENTHESIS {
		arguments = append(arguments, self.parseExpression())
		if self.token != token.COMMA {
			break
		}
		self.expect(token.COMMA)
	}
	rightParenthesis = self.expect(token.RIGHT_PARENTHESIS)
	return
}

func (self *Parser) parseTableName() *ast.TableName {
	tableName := &ast.TableName{
		Name: self.parseStringLiteralOrIdentifier(),
	}
	if self.expectEqualsToken(token.DOT) {
		tableName.Schema = tableName.Name
		tableName.Name = self.parseStringLiteralOrIdentifier()
	}
	return tableName
}

func (self *Parser) parseTableNames() (names []*ast.TableName) {
	for {
		names = append(names, self.parseTableName())
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	return
}

func (self *Parser) parseColumnName() *ast.ColumnName {
	columnName := &ast.ColumnName{
		Name: self.parseStringLiteralOrIdentifier(),
	}
	if self.expectEqualsToken(token.DOT) {
		columnName.Table = columnName.Name
		columnName.Name = self.parseStringLiteralOrIdentifier()
	}
	if self.expectEqualsToken(token.DOT) {
		columnName.Schema = columnName.Table
		columnName.Table = columnName.Name
		columnName.Name = self.parseStringLiteralOrIdentifier()
	}
	return columnName
}

func (self *Parser) parseColumnNames() (names []*ast.ColumnName) {
	for {
		names = append(names, self.parseColumnName())
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	return
}

func (self *Parser) parseWhereExpression() ast.Expression {
	self.scope.inWhere = true
	expr := self.parseExpression()
	self.scope.inWhere = false
	return expr
}

func (self *Parser) parsePrimaryResultSet() ast.ResultSet {
	switch self.token {
	case token.LEFT_PARENTHESIS:
		return self.parseSubqueryExpression()
	case token.IDENTIFIER:
		return self.parseTableSource()
	default:
		self.errorUnexpectedMsg(fmt.Sprintf("Unexpected result set: %v", self.token))
		return nil
	}
}

func (self *Parser) parseResultSet() ast.ResultSet {
	left := self.parsePrimaryResultSet()

	for {
		switch self.token {
		case token.COMMA, token.JOIN, token.INNER, token.LEFT, token.RIGHT:
			left = self.parseJoin(left)
		default:
			return left
		}
	}
}

func (self *Parser) parseJoin(left ast.ResultSet) ast.ResultSet {
	join := &ast.Join{
		Left: left,
	}
	switch self.token {
	case token.COMMA:
		self.expectToken(token.COMMA)
		join.JoinType = ast.CrossJoin
	case token.JOIN, token.INNER:
		if self.token == token.INNER {
			self.expectToken(token.INNER)
		}
		self.expectToken(token.JOIN)
		join.JoinType = ast.InnerJoin
	case token.LEFT:
		self.expectToken(token.LEFT)
		self.expectToken(token.JOIN)
		join.JoinType = ast.LeftJoin
	case token.RIGHT:
		self.expectToken(token.RIGHT)
		self.expectToken(token.JOIN)
		join.JoinType = ast.RightJoin
	}
	if join.JoinType != 0 {
		join.Right = self.parsePrimaryResultSet()
		if self.token == token.ON {
			join.On = self.parseOnCondition()
		}
	}
	return join
}

func (self *Parser) parseOnCondition() *ast.OnCondition {
	self.expectToken(token.ON)
	return &ast.OnCondition{
		Expr: self.parseWhereExpression(),
	}
}

func (self *Parser) parseSubqueryExpression() ast.ResultSet {
	subqueryExpression := &ast.SubqueryExpression{
		LeftParenthesis:  self.expect(token.LEFT_PARENTHESIS),
		Select:           self.parseSelectStatement(),
		RightParenthesis: self.expect(token.RIGHT_PARENTHESIS),
	}
	return subqueryExpression
}

func (self *Parser) parseTableSource() ast.ResultSet {
	tableSource := &ast.TableSource{
		TableName: self.parseTableName(),
	}
	if self.expectEqualsToken(token.AS) || self.token == token.IDENTIFIER {
		tableSource.AsName = self.parseIdentifier()
	}
	return tableSource
}
