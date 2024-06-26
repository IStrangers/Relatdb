package parser

import (
	"Relatdb/parser/ast"
	token "Relatdb/parser/token"
	"fmt"
	"strings"
)

type Scope struct {
	outer *Scope

	inSelect      bool
	inSelectField bool
	inWhere       bool
}

type Parser struct {
	baseOffset     uint64
	skipComment    bool
	skipWhiteSpace bool

	content   string
	length    uint64
	chr       rune
	chrOffset uint64
	offset    uint64
	token     token.Token
	literal   string
	value     string
	index     uint64

	scope *Scope
}

func CreateParser(baseOffset uint64, content string, skipComment bool, skipWhiteSpace bool) *Parser {
	return &Parser{
		baseOffset:     baseOffset,
		skipComment:    skipComment,
		skipWhiteSpace: skipWhiteSpace,
		content:        strings.ToLower(content),
		length:         uint64(len(content)),
		chr:            ' ',
	}
}

func (self *Parser) openScope() {
	self.scope = &Scope{
		outer: self.scope,
	}
	scope := self.scope
	outer := self.scope.outer
	if outer != nil {
		scope.inSelect, scope.inWhere = outer.inSelect, outer.inWhere
	}
}

func (self *Parser) closeScope() {
	self.scope = self.scope.outer
}

func (self *Parser) Parse() []ast.Statement {
	self.next()
	return self.parseStatements()
}

func (self *Parser) ScanNextToken() (token.Token, string, string, uint64) {
	return self.scan()
}

func (self *Parser) next() {
	self.token, self.literal, self.value, self.index = self.scan()
}

func (self *Parser) expect(tkn token.Token) uint64 {
	index := self.index
	if self.token != tkn {
		self.errorUnexpectedToken(tkn)
	}
	self.next()
	return index
}

func (self *Parser) expectToken(tkn token.Token) token.Token {
	if self.token != tkn {
		self.errorUnexpectedToken(tkn)
	}
	self.next()
	return tkn
}

func (self *Parser) expectEqualsToken(tkn token.Token) bool {
	if self.token != tkn {
		return false
	}
	self.next()
	return true
}

func (self *Parser) slice(start, end uint64) string {
	from := start - self.baseOffset
	to := end - self.baseOffset
	if from >= 0 && to <= uint64(len(self.content)) {
		return self.content[from:to]
	}
	return ""
}

type ParseState struct {
	chr        rune
	chrOffset  uint64
	offset     uint64
	token      token.Token
	literal    string
	value      string
	index      uint64
	errorIndex uint64
}

func (self *Parser) markParseState() *ParseState {
	return &ParseState{
		chr:       self.chr,
		chrOffset: self.chrOffset,
		offset:    self.offset,
		token:     self.token,
		literal:   self.literal,
		value:     self.value,
		index:     self.index,
	}
}

func (self *Parser) restoreParseState(parseState *ParseState) {
	self.chr = parseState.chr
	self.chrOffset = parseState.chrOffset
	self.offset = parseState.offset
	self.token = parseState.token
	self.literal = parseState.literal
	self.value = parseState.value
	self.index = parseState.index
}

func (self *Parser) errorUnexpectedToken(tkn token.Token) {
	panic(fmt.Sprintf("Unexpected token %v", tkn))
}

func (self *Parser) errorUnexpectedMsg(msg string) {
	panic(msg)
}
