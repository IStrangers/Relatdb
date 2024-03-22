package parser

type Node interface {
	StartIndex() uint64
	EndIndex() uint64
}

type Expression interface {
	Node
	expression()
}

type Statement interface {
	Node
	statement()
}

type _Expression_ struct {
	Expression
}

func (self *_Expression_) expression() {
}

type _Statement_ struct {
	Statement
}

func (self *_Statement_) statement() {
}

type ResultSet interface {
	Expression

	resultSet()
}

type _ResultSet_ struct {
	_Expression_
}

func (self *_ResultSet_) resultSet() {
}

type ExpressionStatement struct {
	_Statement_

	Expr Expression
}

func (self *ExpressionStatement) StartIndex() uint64 {
	return self.Expr.StartIndex()
}

func (self *ExpressionStatement) EndIndex() uint64 {
	return self.Expr.EndIndex()
}

type NumberLiteral struct {
	_Expression_
	Index   uint64
	Literal string
	Value   any
}

func (self *NumberLiteral) StartIndex() uint64 {
	return self.Index
}

func (self *NumberLiteral) EndIndex() uint64 {
	return self.Index + uint64(len(self.Literal))
}

type StringLiteral struct {
	_Expression_
	Index   uint64
	Literal string
	Value   string
}

func (self *StringLiteral) StartIndex() uint64 {
	return self.Index
}

func (self *StringLiteral) EndIndex() uint64 {
	return self.Index + uint64(len(self.Literal))
}

type BooleanLiteral struct {
	_Expression_
	Index uint64
	Value bool
}

func (self *BooleanLiteral) StartIndex() uint64 {
	return self.Index
}

func (self *BooleanLiteral) EndIndex() uint64 {
	if self.Value {
		return self.Index + 4
	}
	return self.Index + 5
}

type NullLiteral struct {
	_Expression_
	Index uint64
}

func (self *NullLiteral) StartIndex() uint64 {
	return self.Index
}

func (self *NullLiteral) EndIndex() uint64 {
	return self.Index + 4
}

type Identifier struct {
	_Expression_

	Index uint64
	Name  string
}

func (self *Identifier) StartIndex() uint64 {
	return self.Index
}

func (self *Identifier) EndIndex() uint64 {
	return self.Index + uint64(len(self.Name))
}

type AssignExpression struct {
	_Expression_

	Left     Expression
	Operator Token
	Right    Expression
}

func (self *AssignExpression) StartIndex() uint64 {
	return self.Left.StartIndex()
}

func (self *AssignExpression) EndIndex() uint64 {
	return self.Right.EndIndex()
}

type BinaryExpression struct {
	_Expression_
	Left       Expression
	Operator   Token
	Right      Expression
	Comparison bool
}

func (self *BinaryExpression) StartIndex() uint64 {
	return self.Left.StartIndex()
}

func (self *BinaryExpression) EndIndex() uint64 {
	return self.Right.EndIndex()
}

type UnaryExpression struct {
	_Expression_
	Index    uint64
	Operator Token
	Operand  Expression
	Postfix  bool
}

func (self *UnaryExpression) StartIndex() uint64 {
	return self.Index
}

func (self *UnaryExpression) EndIndex() uint64 {
	if self.Postfix {
		return self.Operand.EndIndex() + 2
	}
	return self.Operand.EndIndex()
}

type CallExpression struct {
	_Expression_
	Callee           Expression
	LeftParenthesis  uint64
	Arguments        []Expression
	RightParenthesis uint64
}

func (self *CallExpression) StartIndex() uint64 {
	return self.Callee.StartIndex()
}

func (self *CallExpression) EndIndex() uint64 {
	return self.RightParenthesis + 1
}
