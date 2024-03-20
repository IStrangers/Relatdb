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
	return &ast.DropIndexStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		TableName: self.parseTableName(),
		Name:      self.parseStringLiteralOrIdentifier(),
	}
}

func (self *Parser) parseInsertStatement() ast.Statement {
	return nil
}

func (self *Parser) parseDeleteStatement() ast.Statement {
	return nil
}

func (self *Parser) parseUpdateStatement() ast.Statement {
	return nil
}

func (self *Parser) parseSelectStatement() ast.Statement {
	return nil
}

func (self *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Expr: self.parseExpression(),
	}
}
