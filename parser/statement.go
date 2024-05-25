package parser

import (
	"Relatdb/common"
	"Relatdb/parser/ast"
	"Relatdb/parser/token"
	"Relatdb/utils"
	"fmt"
)

func (self *Parser) parseStatement() ast.Statement {
	defer self.closeScope()
	self.openScope()

	switch self.token {
	case token.SHOW:
		return self.parseShowStatement()
	case token.USE:
		return self.parseUseStatement()
	case token.SET:
		return self.parseSetVariableStatement()
	case token.CREATE:
		return self.parseCreateStatement()
	case token.DROP:
		return self.parseDropStatement()
	case token.INSERT:
		return self.parseInsertStatement()
	case token.DELETE:
		return self.parseDeleteStatement()
	case token.UPDATE:
		return self.parseUpdateStatement()
	case token.SELECT:
		return self.parseSelectStatement()
	default:
		return self.parseExpressionStatement()
	}
}

func (self *Parser) parseStatements() (statements []ast.Statement) {
	for self.token != token.EOF {
		statements = append(statements, self.parseStatement())
		self.expectEqualsToken(token.SEMICOLON)
	}
	return
}

func (self *Parser) parseShowStatement() ast.Statement {
	showIndex := self.expect(token.SHOW)
	switch self.token {
	case token.DATABASES, token.STATUS:
		showType := map[token.Token]ast.ShowStatementType{
			token.DATABASES: ast.ShowDatabases,
			token.STATUS:    ast.ShowStatus,
		}[self.token]
		return &ast.ShowStatement{
			ShowIndex: showIndex,
			Type:      showType,
			KeyWord:   self.parseKeyWordIdentifier(self.token),
		}
	case token.VARIABLES:
		showStatement := &ast.ShowStatement{
			ShowIndex: showIndex,
			Type:      ast.ShowVariables,
			KeyWord:   self.parseKeyWordIdentifier(self.token),
		}
		if self.expectEqualsToken(token.LIKE) {
			showStatement.Where = self.parseExpression()
		}
		return showStatement
	case token.TABLES:
		showStatement := &ast.ShowStatement{
			ShowIndex: showIndex,
			Type:      ast.ShowTables,
			KeyWord:   self.parseKeyWordIdentifier(self.token),
		}
		self.expectToken(token.FROM)
		showStatement.TableName = self.parseTableName()
		return showStatement
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseUseStatement() ast.Statement {
	useIndex := self.expect(token.USE)
	return &ast.UseStatement{
		UseIndex: useIndex,
		Database: self.parseIdentifier(),
	}
}

func (self *Parser) parseSetVariableStatement() ast.Statement {
	setIndex := self.expect(token.SET)
	return &ast.SetVariableStatement{
		SetIndex: setIndex,
		Name:     self.parseIdentifier(),
		Value:    self.parseExpression(),
	}
}

func (self *Parser) parseCreateStatement() ast.Statement {
	createIndex := self.expect(token.CREATE)
	switch self.token {
	case token.DATABASE:
		return self.parseCreateDatabaseStatement(createIndex)
	case token.TABLE:
		return self.parseCreateTableStatement(createIndex)
	case token.INDEX, token.UNIQUE, token.SPATIAL, token.FULLTEXT:
		return self.parseCreateIndexStatement(createIndex)
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseCreateDatabaseStatement(createIndex uint64) ast.Statement {
	self.expectToken(token.DATABASE)
	return &ast.CreateDatabaseStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.NOT) && self.expectEqualsToken(token.EXISTS),
		Name:        self.parseIdentifier(),
	}
}

func (self *Parser) parseCreateTableStatement(createIndex uint64) ast.Statement {
	self.expectToken(token.TABLE)
	return &ast.CreateTableStatement{
		CreateIndex:       createIndex,
		IfNotExists:       self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.NOT) && self.expectEqualsToken(token.EXISTS),
		Name:              self.parseTableName(),
		ColumnDefinitions: self.parseColumnDefinitions(),
		RightParenthesis:  self.expect(token.RIGHT_PARENTHESIS),
	}
}

func (self *Parser) parseColumnDefinitions() (columnDefinitions []*ast.ColumnDefinition) {
	self.expectToken(token.LEFT_PARENTHESIS)
	for self.token != token.RIGHT_PARENTHESIS && self.token != token.EOF {
		columnDefinitions = append(columnDefinitions, self.parseColumnDefinition())
		self.expectEqualsToken(token.COMMA)
	}
	return
}

func (self *Parser) parseColumnDefinition() *ast.ColumnDefinition {
	columnDefinition := &ast.ColumnDefinition{
		Name: self.parseIdentifier(),
	}
	fieldType := token.GetFieldType(self.token)
	if fieldType == 0 {
		self.errorUnexpectedMsg(fmt.Sprintf("Unexpected column type: %v", self.token))
		return nil
	}
	self.expectToken(self.token)
	columnDefinition.Type = fieldType
	if self.expectEqualsToken(token.LEFT_PARENTHESIS) {
		length, _ := utils.ConvertInt(self.parseNumberLiteral().Value)
		columnDefinition.Length = length
		if self.expectEqualsToken(token.COMMA) {
			decimal, _ := utils.ConvertInt(self.parseNumberLiteral().Value)
			columnDefinition.Decimal = decimal
		}
		self.expectToken(token.RIGHT_PARENTHESIS)
	} else {
		lengthAndDecimal := common.GetFieldDefaultLengthAndDecimal(fieldType)
		columnDefinition.Length = lengthAndDecimal.Length
		columnDefinition.Decimal = lengthAndDecimal.Decimal
	}

	if self.expectEqualsToken(token.PRIMARY) && self.expectEqualsToken(token.KEY) {
		columnDefinition.Flag |= common.PRIMARY_KEY_FLAG
	}
	if self.expectEqualsToken(token.UNIQUE) {
		columnDefinition.Flag |= common.UNIQUE_KEY_FLAG
	}
	if self.expectEqualsToken(token.UNSIGNED) {
		columnDefinition.Flag |= common.UNSIGNED_FLAG
	}
	if self.expectEqualsToken(token.ZEROFILL) {
		columnDefinition.Flag |= common.ZEROFILL_FLAG
	}
	if self.expectEqualsToken(token.AUTO_INCREMENT) {
		columnDefinition.Flag |= common.AUTO_INCREMENT_FLAG
	}
	if self.expectEqualsToken(token.NOT) && self.expectEqualsToken(token.NULL) {
		columnDefinition.Flag |= common.NOT_NULL_FLAG
	} else if self.expectEqualsToken(token.NULL) {

	}

	if self.expectEqualsToken(token.DEFAULT) {
		columnDefinition.DefaultValue = self.parseExpression()
	}
	if self.expectEqualsToken(token.COLUMN_COMMENT) {
		columnDefinition.Comment = self.parseExpression()
	}
	return columnDefinition
}

func (self *Parser) parseCreateIndexStatement(createIndex uint64) ast.Statement {
	indexType := ast.IndexTypeNone
	if self.token != token.INDEX {
		switch self.token {
		case token.UNIQUE:
			indexType = ast.IndexTypeUnique
		case token.SPATIAL:
			indexType = ast.IndexTypeSpatial
		case token.FULLTEXT:
			indexType = ast.IndexTypeFullText
		default:
			self.errorUnexpectedMsg(fmt.Sprintf("Unexpected index type: %v", self.token))
		}
		self.expectToken(self.token)
	}
	self.expectToken(token.INDEX)
	createIndexStatement := &ast.CreateIndexStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.NOT) && self.expectEqualsToken(token.EXISTS),
		Name:        self.parseIdentifier(),
		Type:        indexType,
	}
	self.expectToken(token.ON)
	createIndexStatement.TableName = self.parseTableName()
	self.expectToken(token.LEFT_PARENTHESIS)
	createIndexStatement.ColumnNames = self.parseColumnNames()
	self.expectToken(token.RIGHT_PARENTHESIS)
	return createIndexStatement
}

func (self *Parser) parseDropStatement() ast.Statement {
	dropIndex := self.expect(token.DROP)
	switch self.token {
	case token.DATABASE:
		return self.parseDropDatabaseStatement(dropIndex)
	case token.TABLE:
		return self.parseDropTableStatement(dropIndex)
	case token.INDEX:
		return self.parseDropIndexStatement(dropIndex)
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseDropDatabaseStatement(dropIndex uint64) ast.Statement {
	self.expectToken(token.DATABASE)
	return &ast.DropDatabaseStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.EXISTS),
		Name:      self.parseIdentifier(),
	}
}

func (self *Parser) parseDropTableStatement(dropIndex uint64) ast.Statement {
	self.expectToken(token.TABLE)
	return &ast.DropTableStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.EXISTS),
		Names:     self.parseTableNames(),
	}
}

func (self *Parser) parseDropIndexStatement(dropIndex uint64) ast.Statement {
	self.expectToken(token.INDEX)
	dropIndexStatement := &ast.DropIndexStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(token.IF) && self.expectEqualsToken(token.EXISTS),
		Name:      self.parseIdentifier(),
	}
	self.expectToken(token.ON)
	dropIndexStatement.TableName = self.parseTableName()
	return dropIndexStatement
}

func (self *Parser) parseInsertStatement() ast.Statement {
	insertStatement := &ast.InsertStatement{
		InsertIndex: self.expect(token.INSERT),
	}
	self.expectToken(token.INTO)
	insertStatement.TableName = self.parseTableName()
	self.expectToken(token.LEFT_PARENTHESIS)
	insertStatement.ColumnNames = self.parseColumnNames()
	self.expectToken(token.RIGHT_PARENTHESIS)
	self.expectToken(token.VALUES)
	for {
		self.expectToken(token.LEFT_PARENTHESIS)
		var values []ast.Expression
		for self.token != token.RIGHT_PARENTHESIS && self.token != token.EOF {
			values = append(values, self.parseExpression())
			self.expectEqualsToken(token.COMMA)
		}
		insertStatement.Values = append(insertStatement.Values, values)
		self.expectToken(token.RIGHT_PARENTHESIS)
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	return insertStatement
}

func (self *Parser) parseDeleteStatement() ast.Statement {
	deleteStatement := &ast.DeleteStatement{
		DeleteIndex: self.expect(token.DELETE),
	}
	self.expectToken(token.FROM)
	deleteStatement.TableName = self.parseTableName()
	if self.expectEqualsToken(token.WHERE) {
		deleteStatement.Where = self.parseWhereExpression()
	}
	if self.token == token.ORDER {
		deleteStatement.Order = self.parseOrderByClause()
	}
	if self.token == token.LIMIT {
		deleteStatement.Limit = self.parseLimit()
	}
	return deleteStatement
}

func (self *Parser) parseUpdateStatement() ast.Statement {
	updateStatement := &ast.UpdateStatement{
		UpdateIndex: self.expect(token.UPDATE),
		TableName:   self.parseTableName(),
	}
	self.expectToken(token.SET)
	for {
		updateStatement.AssignExpressions = append(updateStatement.AssignExpressions, self.parseAssignExpression())
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	if self.expectEqualsToken(token.WHERE) {
		updateStatement.Where = self.parseWhereExpression()
	}
	if self.token == token.ORDER {
		updateStatement.Order = self.parseOrderByClause()
	}
	if self.token == token.LIMIT {
		updateStatement.Limit = self.parseLimit()
	}
	return updateStatement
}

func (self *Parser) parseSelectStatement() *ast.SelectStatement {
	defer func() { self.scope.inSelect = false }()
	self.scope.inSelect = true
	selectStatement := &ast.SelectStatement{
		SelectIndex: self.expect(token.SELECT),
		Fields:      self.parseSelectFields(),
	}
	if self.expectEqualsToken(token.FROM) {
		selectStatement.From = self.parseResultSet()
	}
	if self.expectEqualsToken(token.WHERE) {
		selectStatement.Where = self.parseWhereExpression()
	}
	if self.token == token.ORDER {
		selectStatement.Order = self.parseOrderByClause()
	}
	if self.token == token.LIMIT {
		selectStatement.Limit = self.parseLimit()
	}
	return selectStatement
}

func (self *Parser) parseSelectField() *ast.SelectField {
	defer func() { self.scope.inSelectField = false }()
	self.scope.inSelectField = true
	selectField := &ast.SelectField{}
	switch self.token {
	case token.MULTIPLY:
		selectField.Expr = self.parseKeyWordIdentifier(token.MULTIPLY)
	default:
		selectField.Expr = self.parseExpression()
	}
	if self.expectEqualsToken(token.AS) || self.token == token.IDENTIFIER {
		selectField.AsName = self.parseIdentifier()
	}
	return selectField
}

func (self *Parser) parseSelectFields() (selectFields []*ast.SelectField) {
	for {
		selectFields = append(selectFields, self.parseSelectField())
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	return
}

func (self *Parser) parseOrderByClause() *ast.OrderByClause {
	orderByClause := &ast.OrderByClause{
		OrderByIndex: self.expect(token.ORDER),
	}
	self.expectToken(token.BY)
	orderByClause.Items = self.parseOrderItems()
	return orderByClause
}

func (self *Parser) parseOrderItem() *ast.OrderItem {
	orderItem := &ast.OrderItem{
		ColumnName: self.parseColumnName(),
		Desc:       false,
	}
	if self.token == token.AES || self.token == token.DESC {
		orderItem.Desc = self.token == token.DESC
		orderItem.Order = self.parseKeyWordIdentifier(self.token)
	}
	return orderItem
}

func (self *Parser) parseOrderItems() (orderItems []*ast.OrderItem) {
	for {
		orderItems = append(orderItems, self.parseOrderItem())
		if self.token != token.COMMA {
			break
		}
		self.expectToken(token.COMMA)
	}
	return
}

func (self *Parser) parseLimit() *ast.Limit {
	limit := &ast.Limit{
		LimitIndex: self.expect(token.LIMIT),
		Count:      self.parseExpression(),
	}
	if self.expectEqualsToken(token.COMMA) {
		limit.Offset = limit.Count
		limit.Count = self.parseExpression()
	}
	return limit
}

func (self *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Expr: self.parseExpression(),
	}
}
