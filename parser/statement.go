package parser

import (
	"Relatdb/parser/ast"
	"fmt"
)

func (self *Parser) parseStatement() ast.Statement {
	switch self.token {
	case CREATE:
		return self.parseCreateStatement()
	case DROP:
		return self.parseDropStatement()
	case INSERT:
		return self.parseInsertStatement()
	case DELETE:
		return self.parseDeleteStatement()
	case UPDATE:
		return self.parseUpdateStatement()
	case SELECT:
		return self.parseSelectStatement()
	default:
		return self.parseExpressionStatement()
	}
}

func (self *Parser) parseStatements() (statements []ast.Statement) {
	for self.token != EOF {
		statements = append(statements, self.parseStatement())
		self.expectEqualsToken(SEMICOLON)
	}
	return
}

func (self *Parser) parseCreateStatement() ast.Statement {
	createIndex := self.expect(CREATE)
	switch self.token {
	case DATABASE:
		return self.parseCreateDatabaseStatement(createIndex)
	case TABLE:
		return self.parseCreateTableStatement(createIndex)
	case INDEX, UNIQUE, SPATIAL, FULLTEXT:
		return self.parseCreateIndexStatement(createIndex)
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseCreateDatabaseStatement(createIndex uint64) ast.Statement {
	self.expectToken(DATABASE)
	return &ast.CreateDatabaseStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:        self.parseStringLiteralOrIdentifier(),
	}
}

func (self *Parser) parseCreateTableStatement(createIndex uint64) ast.Statement {
	self.expectToken(TABLE)
	return &ast.CreateTableStatement{
		CreateIndex:       createIndex,
		IfNotExists:       self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:              self.parseTableName(),
		ColumnDefinitions: self.parseColumnDefinitions(),
		RightParenthesis:  self.expect(RIGHT_PARENTHESIS),
	}
}

func (self *Parser) parseColumnDefinitions() (columnDefinitions []*ast.ColumnDefinition) {
	self.expectToken(LEFT_PARENTHESIS)
	for self.token != RIGHT_PARENTHESIS && self.token != EOF {
		columnDefinitions = append(columnDefinitions, self.parseColumnDefinition())
		self.expectEqualsToken(COMMA)
	}
	return
}

func (self *Parser) parseColumnDefinition() *ast.ColumnDefinition {
	return &ast.ColumnDefinition{}
}

func (self *Parser) parseCreateIndexStatement(createIndex uint64) ast.Statement {
	indexType := ast.IndexTypeNone
	if self.token != INDEX {
		switch self.token {
		case UNIQUE:
			indexType = ast.IndexTypeUnique
		case SPATIAL:
			indexType = ast.IndexTypeSpatial
		case FULLTEXT:
			indexType = ast.IndexTypeFullText
		default:
			self.errorUnexpectedMsg(fmt.Sprintf("Unexpected index type: %v", self.token))
		}
		self.expectToken(self.token)
	}
	self.expectToken(INDEX)
	createIndexStatement := &ast.CreateIndexStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:        self.parseStringLiteralOrIdentifier(),
		Type:        indexType,
	}
	self.expectToken(ON)
	createIndexStatement.TableName = self.parseTableName()
	self.expectToken(LEFT_PARENTHESIS)
	createIndexStatement.ColumnNames = self.parseColumnNames()
	self.expectToken(RIGHT_PARENTHESIS)
	return createIndexStatement
}

func (self *Parser) parseDropStatement() ast.Statement {
	dropIndex := self.expect(DROP)
	switch self.token {
	case DATABASE:
		return self.parseDropDatabaseStatement(dropIndex)
	case TABLE:
		return self.parseDropTableStatement(dropIndex)
	case INDEX:
		return self.parseDropIndexStatement(dropIndex)
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseDropDatabaseStatement(dropIndex uint64) ast.Statement {
	self.expectToken(DATABASE)
	return &ast.DropDatabaseStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Name:      self.parseStringLiteralOrIdentifier(),
	}
}

func (self *Parser) parseDropTableStatement(dropIndex uint64) ast.Statement {
	self.expectToken(TABLE)
	return &ast.DropTableStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Names:     self.parseTableNames(),
	}
}

func (self *Parser) parseDropIndexStatement(dropIndex uint64) ast.Statement {
	self.expectToken(INDEX)
	dropIndexStatement := &ast.DropIndexStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Name:      self.parseStringLiteralOrIdentifier(),
	}
	self.expectToken(ON)
	dropIndexStatement.TableName = self.parseTableName()
	return dropIndexStatement
}

func (self *Parser) parseInsertStatement() ast.Statement {
	insertStatement := &ast.InsertStatement{
		InsertIndex: self.expect(INSERT),
	}
	self.expectToken(INTO)
	insertStatement.TableName = self.parseTableName()
	self.expectToken(LEFT_PARENTHESIS)
	insertStatement.ColumnNames = self.parseColumnNames()
	self.expectToken(RIGHT_PARENTHESIS)
	self.expectToken(VALUES)
	for {
		self.expectToken(LEFT_PARENTHESIS)
		var values []ast.Expression
		for self.token != RIGHT_PARENTHESIS && self.token != EOF {
			values = append(values, self.parseExpression())
			self.expectEqualsToken(COMMA)
		}
		insertStatement.Values = append(insertStatement.Values, values)
		self.expectToken(RIGHT_PARENTHESIS)
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return insertStatement
}

func (self *Parser) parseDeleteStatement() ast.Statement {
	deleteStatement := &ast.DeleteStatement{
		DeleteIndex: self.expect(DELETE),
	}
	self.expectToken(FROM)
	deleteStatement.TableName = self.parseTableName()
	if self.expectEqualsToken(WHERE) {
		deleteStatement.Where = self.parseExpression()
	}
	if self.token == ORDER {
		deleteStatement.Order = self.parseOrderByClause()
	}
	if self.token == LIMIT {
		deleteStatement.Limit = self.parseLimit()
	}
	return deleteStatement
}

func (self *Parser) parseUpdateStatement() ast.Statement {
	return nil
}

func (self *Parser) parseSelectStatement() ast.Statement {
	return nil
}

func (self *Parser) parseOrderByClause() *ast.OrderByClause {
	orderByClause := &ast.OrderByClause{
		OrderByIndex: self.expect(ORDER),
	}
	self.expectToken(BY)
	orderByClause.Items = self.parseOrderItems()
	return orderByClause
}

func (self *Parser) parseOrderItem() *ast.OrderItem {
	orderItem := &ast.OrderItem{
		ColumnName: self.parseColumnName(),
		Desc:       false,
	}
	if self.token == AES || self.token == DESC {
		orderItem.Desc = self.token == DESC
		orderItem.Order = self.parseKeyWordIdentifier(self.token)
	}
	return orderItem
}

func (self *Parser) parseOrderItems() (orderItems []*ast.OrderItem) {
	for {
		orderItems = append(orderItems, self.parseOrderItem())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseLimit() *ast.Limit {
	limit := &ast.Limit{
		LimitIndex: self.expect(LIMIT),
		Count:      self.parseExpression(),
	}
	if self.expectEqualsToken(COMMA) {
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
