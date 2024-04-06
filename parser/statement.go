package parser

import (
	"Relatdb/common"
	"Relatdb/utils"
	"fmt"
)

func (self *Parser) parseStatement() Statement {
	defer self.closeScope()
	self.openScope()

	switch self.token {
	case SHOW:
		return self.parseShowStatement()
	case USE:
		return self.parseUseStatement()
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

func (self *Parser) parseStatements() (statements []Statement) {
	for self.token != EOF {
		statements = append(statements, self.parseStatement())
		self.expectEqualsToken(SEMICOLON)
	}
	return
}

func (self *Parser) parseShowStatement() Statement {
	showIndex := self.expect(SHOW)
	switch self.token {
	case DATABASES, TABLES, VARIABLES:
		showType := map[Token]ShowStatementType{
			DATABASES: ShowDatabases,
			TABLES:    ShowTables,
			VARIABLES: ShowVariables,
		}[self.token]
		return &ShowStatement{
			ShowIndex: showIndex,
			Type:      showType,
			KeyWord:   self.parseKeyWordIdentifier(self.token),
		}
	default:
		self.errorUnexpectedToken(self.token)
		return nil
	}
}

func (self *Parser) parseUseStatement() Statement {
	useIndex := self.expect(USE)
	return &UseStatement{
		UseIndex: useIndex,
		Database: self.parseIdentifier(),
	}
}

func (self *Parser) parseCreateStatement() Statement {
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

func (self *Parser) parseCreateDatabaseStatement(createIndex uint64) Statement {
	self.expectToken(DATABASE)
	return &CreateDatabaseStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:        self.parseIdentifier(),
	}
}

func (self *Parser) parseCreateTableStatement(createIndex uint64) Statement {
	self.expectToken(TABLE)
	return &CreateTableStatement{
		CreateIndex:       createIndex,
		IfNotExists:       self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:              self.parseTableName(),
		ColumnDefinitions: self.parseColumnDefinitions(),
		RightParenthesis:  self.expect(RIGHT_PARENTHESIS),
	}
}

func (self *Parser) parseColumnDefinitions() (columnDefinitions []*ColumnDefinition) {
	self.expectToken(LEFT_PARENTHESIS)
	for self.token != RIGHT_PARENTHESIS && self.token != EOF {
		columnDefinitions = append(columnDefinitions, self.parseColumnDefinition())
		self.expectEqualsToken(COMMA)
	}
	return
}

func (self *Parser) parseColumnDefinition() *ColumnDefinition {
	columnDefinition := &ColumnDefinition{
		Name: self.parseIdentifier(),
	}
	fieldType := GetFieldType(self.token)
	if fieldType == 0 {
		self.errorUnexpectedMsg(fmt.Sprintf("Unexpected column type: %v", self.token))
		return nil
	}
	self.expectToken(self.token)
	columnDefinition.Type = fieldType
	if self.expectEqualsToken(LEFT_PARENTHESIS) {
		length, _ := utils.ConvertInt(self.parseNumberLiteral().Value)
		columnDefinition.Length = length
		if self.expectEqualsToken(COMMA) {
			decimal, _ := utils.ConvertInt(self.parseNumberLiteral().Value)
			columnDefinition.Decimal = decimal
		}
		self.expectToken(RIGHT_PARENTHESIS)
	} else {
		lengthAndDecimal := common.GetFieldDefaultLengthAndDecimal(fieldType)
		columnDefinition.Length = lengthAndDecimal.Length
		columnDefinition.Decimal = lengthAndDecimal.Decimal
	}

	if self.expectEqualsToken(PRIMARY) && self.expectEqualsToken(KEY) {
		columnDefinition.Flag |= common.PRI_KEY_FLAG
	}
	if self.expectEqualsToken(UNIQUE) {
		columnDefinition.Flag |= common.UNIQUE_KEY_FLAG
	}
	if self.expectEqualsToken(UNSIGNED) {
		columnDefinition.Flag |= common.UNSIGNED_FLAG
	}
	if self.expectEqualsToken(ZEROFILL) {
		columnDefinition.Flag |= common.ZEROFILL_FLAG
	}
	if self.expectEqualsToken(AUTO_INCREMENT) {
		columnDefinition.Flag |= common.AUTO_INCREMENT_FLAG
	}
	if self.expectEqualsToken(NOT) && self.expectEqualsToken(NULL) {
		columnDefinition.Flag |= common.NOT_NULL_FLAG
	} else if self.expectEqualsToken(NULL) {

	}

	if self.expectEqualsToken(DEFAULT) {
		columnDefinition.DefaultValue = self.parseExpression()
	}
	if self.expectEqualsToken(COLUMN_COMMENT) {
		columnDefinition.Comment = self.parseStringLiteral()
	}
	return columnDefinition
}

func (self *Parser) parseCreateIndexStatement(createIndex uint64) Statement {
	indexType := IndexTypeNone
	if self.token != INDEX {
		switch self.token {
		case UNIQUE:
			indexType = IndexTypeUnique
		case SPATIAL:
			indexType = IndexTypeSpatial
		case FULLTEXT:
			indexType = IndexTypeFullText
		default:
			self.errorUnexpectedMsg(fmt.Sprintf("Unexpected index type: %v", self.token))
		}
		self.expectToken(self.token)
	}
	self.expectToken(INDEX)
	createIndexStatement := &CreateIndexStatement{
		CreateIndex: createIndex,
		IfNotExists: self.expectEqualsToken(IF) && self.expectEqualsToken(NOT) && self.expectEqualsToken(EXISTS),
		Name:        self.parseIdentifier(),
		Type:        indexType,
	}
	self.expectToken(ON)
	createIndexStatement.TableName = self.parseTableName()
	self.expectToken(LEFT_PARENTHESIS)
	createIndexStatement.ColumnNames = self.parseColumnNames()
	self.expectToken(RIGHT_PARENTHESIS)
	return createIndexStatement
}

func (self *Parser) parseDropStatement() Statement {
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

func (self *Parser) parseDropDatabaseStatement(dropIndex uint64) Statement {
	self.expectToken(DATABASE)
	return &DropDatabaseStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Name:      self.parseIdentifier(),
	}
}

func (self *Parser) parseDropTableStatement(dropIndex uint64) Statement {
	self.expectToken(TABLE)
	return &DropTableStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Names:     self.parseTableNames(),
	}
}

func (self *Parser) parseDropIndexStatement(dropIndex uint64) Statement {
	self.expectToken(INDEX)
	dropIndexStatement := &DropIndexStatement{
		DropIndex: dropIndex,
		IfExists:  self.expectEqualsToken(IF) && self.expectEqualsToken(EXISTS),
		Name:      self.parseIdentifier(),
	}
	self.expectToken(ON)
	dropIndexStatement.TableName = self.parseTableName()
	return dropIndexStatement
}

func (self *Parser) parseInsertStatement() Statement {
	insertStatement := &InsertStatement{
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
		var values []Expression
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

func (self *Parser) parseDeleteStatement() Statement {
	deleteStatement := &DeleteStatement{
		DeleteIndex: self.expect(DELETE),
	}
	self.expectToken(FROM)
	deleteStatement.TableName = self.parseTableName()
	if self.expectEqualsToken(WHERE) {
		deleteStatement.Where = self.parseWhereExpression()
	}
	if self.token == ORDER {
		deleteStatement.Order = self.parseOrderByClause()
	}
	if self.token == LIMIT {
		deleteStatement.Limit = self.parseLimit()
	}
	return deleteStatement
}

func (self *Parser) parseUpdateStatement() Statement {
	updateStatement := &UpdateStatement{
		UpdateIndex: self.expect(UPDATE),
		TableName:   self.parseTableName(),
	}
	self.expectToken(SET)
	for {
		updateStatement.AssignExpressions = append(updateStatement.AssignExpressions, self.parseAssignExpression())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	if self.expectEqualsToken(WHERE) {
		updateStatement.Where = self.parseWhereExpression()
	}
	if self.token == ORDER {
		updateStatement.Order = self.parseOrderByClause()
	}
	if self.token == LIMIT {
		updateStatement.Limit = self.parseLimit()
	}
	return updateStatement
}

func (self *Parser) parseSelectStatement() *SelectStatement {
	defer func() { self.scope.inSelect = false }()
	self.scope.inSelect = true
	selectStatement := &SelectStatement{
		SelectIndex: self.expect(SELECT),
		Fields:      self.parseSelectFields(),
	}
	if self.expectEqualsToken(FROM) {
		selectStatement.From = self.parseResultSet()
	}
	if self.expectEqualsToken(WHERE) {
		selectStatement.Where = self.parseWhereExpression()
	}
	if self.token == ORDER {
		selectStatement.Order = self.parseOrderByClause()
	}
	if self.token == LIMIT {
		selectStatement.Limit = self.parseLimit()
	}
	return selectStatement
}

func (self *Parser) parseSelectField() *SelectField {
	defer func() { self.scope.inSelectField = false }()
	self.scope.inSelectField = true
	selectField := &SelectField{}
	switch self.token {
	case MULTIPLY:
		selectField.Expr = self.parseKeyWordIdentifier(MULTIPLY)
	default:
		selectField.Expr = self.parseExpression()
	}
	if self.expectEqualsToken(AS) || self.token == IDENTIFIER {
		selectField.AsName = self.parseIdentifier()
	}
	return selectField
}

func (self *Parser) parseSelectFields() (selectFields []*SelectField) {
	for {
		selectFields = append(selectFields, self.parseSelectField())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseOrderByClause() *OrderByClause {
	orderByClause := &OrderByClause{
		OrderByIndex: self.expect(ORDER),
	}
	self.expectToken(BY)
	orderByClause.Items = self.parseOrderItems()
	return orderByClause
}

func (self *Parser) parseOrderItem() *OrderItem {
	orderItem := &OrderItem{
		ColumnName: self.parseColumnName(),
		Desc:       false,
	}
	if self.token == AES || self.token == DESC {
		orderItem.Desc = self.token == DESC
		orderItem.Order = self.parseKeyWordIdentifier(self.token)
	}
	return orderItem
}

func (self *Parser) parseOrderItems() (orderItems []*OrderItem) {
	for {
		orderItems = append(orderItems, self.parseOrderItem())
		if self.token != COMMA {
			break
		}
		self.expectToken(COMMA)
	}
	return
}

func (self *Parser) parseLimit() *Limit {
	limit := &Limit{
		LimitIndex: self.expect(LIMIT),
		Count:      self.parseExpression(),
	}
	if self.expectEqualsToken(COMMA) {
		limit.Offset = limit.Count
		limit.Count = self.parseExpression()
	}
	return limit
}

func (self *Parser) parseExpressionStatement() Statement {
	return &ExpressionStatement{
		Expr: self.parseExpression(),
	}
}
