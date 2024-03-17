package parser

import (
	"Relatdb/parser/ast"
	"fmt"
)

type Parser struct {
	baseOffset     uint64
	skipComment    bool
	skipWhiteSpace bool

	content   string
	length    uint64
	chr       rune
	chrOffset uint64
	offset    uint64
	token     Token
	literal   string
	value     string
	index     uint64
}

func CreateParser(baseOffset uint64, content string, skipComment bool, skipWhiteSpace bool) *Parser {
	return &Parser{
		baseOffset:     baseOffset,
		skipComment:    skipComment,
		skipWhiteSpace: skipWhiteSpace,
		content:        content,
		length:         uint64(len(content)),
		chr:            ' ',
	}
}

func (self *Parser) Parse() *ast.Node {
	self.next()
	return nil
}

func (self *Parser) ScanNextToken() (Token, string, string, uint64) {
	return self.scan()
}

func (self *Parser) next() {
	self.token, self.literal, self.value, self.index = self.scan()
}

func (self *Parser) expect(tkn Token) uint64 {
	index := self.index
	if self.token != tkn {
		self.errorUnexpectedToken(tkn)
	}
	self.next()
	return index
}

func (self *Parser) expectToken(tkn Token) Token {
	if self.token != tkn {
		self.errorUnexpectedToken(tkn)
	}
	self.next()
	return tkn
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
	token      Token
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

func (self *Parser) errorUnexpectedToken(tkn Token) {
	panic(fmt.Sprintf("Unexpected token %v", tkn))
}
