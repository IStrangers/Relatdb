package parser

import (
	"fmt"
	"strings"
)

func (self *Parser) scan() (tkn Token, literal string, value string, index uint64) {
	for {
		skipWhiteSpace := self.skipWhiteSpace
		if skipWhiteSpace {
			self.skipWhiteSpaceChr()
		}
		index = self.chrOffset
		switch chr := self.chr; {
		case !skipWhiteSpace && isWhiteSpaceChr(chr):
			tkn, literal, value = WHITE_SPACE, string(chr), string(chr)
			self.readChr()
			break
		case isIdentifierStart(chr):
			literal = self.scanIdentifier()
			value = literal
			keywordToken, exists := IsKeyword(literal)
			if exists {
				tkn = keywordToken
			} else {
				tkn = IDENTIFIER
			}
			break
		case isStringSymbol(chr):
			self.readChr()
			value = self.scanString()
			literal = string(chr) + value + string(self.chr)
			tkn = STRING
			self.readChr()
			break
		case isNumeric(chr):
			literal = self.scanNumericLiteral()
			value = literal
			tkn = NUMBER
			break
		default:
			self.readChr()
			switch chr {
			case -1:
				tkn = EOF
				break
			case '+':
				tkn = self.switchToken("+,=", INCREMENT, ADDITION_ASSIGN, ADDITION)
				literal = tkn.String()
				value = tkn.String()
				break
			case '-':
				tkn = self.switchToken("-,=", DECREMENT, SUBTRACT_ASSIGN, SUBTRACT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '*':
				tkn = self.switchToken("=", MULTIPLY_ASSIGN, MULTIPLY)
				literal = tkn.String()
				value = tkn.String()
				break
			case '/':
				tkn = self.switchToken("/,*,=", COMMENT, MULTI_COMMENT, DIVIDE_ASSIGN, DIVIDE)
				literal = tkn.String()
				value = tkn.String()
				if tkn == COMMENT || tkn == MULTI_COMMENT {
					comment := self.scanComment(tkn)
					value = comment
					if tkn == COMMENT {
						comment = "//" + comment
					} else {
						comment = "/*" + comment
						if self.chr == '/' {
							comment += "/"
							self.readChr()
						}
					}
					literal = comment
					if self.skipComment {
						continue
					}
				}
				break
			case '%':
				tkn = self.switchToken("=", REMAINDER_ASSIGN, REMAINDER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '(':
				tkn, literal, value = LEFT_PARENTHESIS, string(chr), string(chr)
				break
			case ')':
				tkn, literal, value = RIGHT_PARENTHESIS, string(chr), string(chr)
				break
			case '{':
				tkn, literal, value = LEFT_BRACE, string(chr), string(chr)
				break
			case '}':
				tkn, literal, value = RIGHT_BRACE, string(chr), string(chr)
				break
			case '[':
				tkn, literal, value = LEFT_BRACKET, string(chr), string(chr)
				break
			case ']':
				tkn, literal, value = RIGHT_BRACKET, string(chr), string(chr)
				break
			case '.':
				tkn, literal, value = DOT, string(chr), string(chr)
				break
			case ',':
				tkn, literal, value = COMMA, string(chr), string(chr)
				break
			case ':':
				tkn, literal, value = COLON, string(chr), string(chr)
				break
			case ';':
				tkn, literal, value = SEMICOLON, string(chr), string(chr)
				break
			case '!':
				tkn = self.switchToken("=", NOT_EQUAL, NOT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '=':
				tkn = self.switchToken("=", EQUAL, ASSIGN)
				literal = tkn.String()
				value = tkn.String()
				break
			case '<':
				tkn = self.switchToken("=", LESS_OR_EQUAL, LESS)
				literal = tkn.String()
				value = tkn.String()
				break
			case '>':
				tkn = self.switchToken("=", GREATER_OR_EQUAL, GREATER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '&':
				tkn = self.switchToken("&,=", LOGICAL_AND, AND_ARITHMETIC_ASSIGN, AND_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
				break
			case '|':
				tkn = self.switchToken("|,=", LOGICAL_OR, OR_ARITHMETIC_ASSIGN, OR_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
				break
			default:
				tkn = ILLEGAL
				panic(fmt.Sprintf("Unexpected end of input: %v", index))
				break
			}
		}
		return
	}
}

func (self *Parser) skipWhiteSpaceChr() {
	for isWhiteSpaceChr(self.chr) {
		self.readChr()
	}
}

func (self *Parser) readChr() rune {
	if self.offset < self.length {
		self.chrOffset = self.offset
		self.chr = rune(self.content[self.offset])
		self.offset += 1
		return self.chr
	}
	self.chrOffset = self.length
	self.chr = -1
	return self.chr
}

func (self *Parser) scanByFilter(filter func(rune) bool) string {
	chrOffset := self.chrOffset
	for filter(self.chr) {
		self.readChr()
	}
	return self.content[chrOffset:self.chrOffset]
}

func (self *Parser) scanIdentifier() string {
	return self.scanByFilter(isIdentifierPart)
}

func (self *Parser) scanNumericLiteral() string {
	return self.scanByFilter(isNumericPart)
}

func (self *Parser) scanString() string {
	return self.scanByFilter(isNotStringSymbol)
}

func (self *Parser) scanComment(tkn Token) string {
	if tkn == MULTI_COMMENT {
		multiCommentCount := 1
		multiComment := self.scanByFilter(func(chr rune) bool {
			if chr == '/' && self.readChr() == '*' {
				multiCommentCount++
			}
			if chr == '*' && self.readChr() == '/' {
				multiCommentCount--
			}
			if chr == -1 {
				self.errorUnexpectedToken(self.token)
			}
			return multiCommentCount > 0 && chr != -1
		})
		return multiComment
	} else {
		return self.scanByFilter(isNotLineTerminator)
	}
}

func (self *Parser) switchToken(keysStr string, tkns ...Token) Token {
	keys := strings.Split(keysStr, ",")
	for i, key := range keys {
		if self.chr == rune(key[0]) {
			self.readChr()
			return tkns[i]
		}
	}
	return tkns[len(tkns)-1]
}

func isWhiteSpaceChr(chr rune) bool {
	return chr == ' ' || chr == '\t' || chr == '\r' || chr == '\n' || chr == '\f'
}

func isIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || (chr >= 'A' && chr <= 'Z') || (chr >= 'a' && chr <= 'z')
}
func isIdentifierPart(chr rune) bool {
	return isIdentifierStart(chr) || isNumeric(chr)
}

func isNumeric(chr rune) bool {
	return chr >= '0' && chr <= '9'
}
func isNumericPart(chr rune) bool {
	return chr == '.' || isNumeric(chr)
}

func isStringSymbol(chr rune) bool {
	return chr == '"' || chr == '\''
}
func isNotStringSymbol(chr rune) bool {
	return !isStringSymbol(chr)
}

func isLineTerminator(chr rune) bool {
	switch chr {
	case '\u000a', '\u000d', '\u2028', '\u2029', -1:
		return true
	}
	return false
}
func isNotLineTerminator(chr rune) bool {
	return !isLineTerminator(chr)
}
