package parser

import (
	"Relatdb/parser/token"
	"fmt"
	"strings"
)

func (self *Parser) scan() (tkn token.Token, literal string, value string, index uint64) {
	for {
		skipWhiteSpace := self.skipWhiteSpace
		if skipWhiteSpace {
			self.skipWhiteSpaceChr()
		}
		index = self.chrOffset
		switch chr := self.chr; {
		case !skipWhiteSpace && isWhiteSpaceChr(chr):
			tkn, literal, value = token.WHITE_SPACE, string(chr), string(chr)
			self.readChr()
			break
		case isIdentifierStart(chr):
			literal = self.scanIdentifier()
			value = literal
			keywordToken, exists := token.IsKeyword(literal)
			if exists {
				tkn = keywordToken
			} else {
				tkn = token.IDENTIFIER
			}
			break
		case isStringSymbol(chr):
			self.readChr()
			value = self.scanString()
			literal = string(chr) + value + string(self.chr)
			tkn = token.STRING
			self.readChr()
			break
		case isNumeric(chr):
			literal = self.scanNumericLiteral()
			value = literal
			tkn = token.NUMBER
			break
		default:
			self.readChr()
			switch chr {
			case -1:
				tkn = token.EOF
				break
			case '+':
				tkn = self.switchToken("+,=", token.INCREMENT, token.ADDITION_ASSIGN, token.ADDITION)
				literal = tkn.String()
				value = tkn.String()
				break
			case '-':
				tkn = self.switchToken("-,=", token.DECREMENT, token.SUBTRACT_ASSIGN, token.SUBTRACT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '*':
				tkn = self.switchToken("=", token.MULTIPLY_ASSIGN, token.MULTIPLY)
				literal = tkn.String()
				value = tkn.String()
				break
			case '/':
				tkn = self.switchToken("/,*,=", token.COMMENT, token.MULTI_COMMENT, token.DIVIDE_ASSIGN, token.DIVIDE)
				literal = tkn.String()
				value = tkn.String()
				if tkn == token.COMMENT || tkn == token.MULTI_COMMENT {
					comment := self.scanComment(tkn)
					value = comment
					if tkn == token.COMMENT {
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
				tkn = self.switchToken("=", token.REMAINDER_ASSIGN, token.REMAINDER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '@':
				tkn, literal, value = token.AT_IDENTIFIER, string(chr), string(chr)
				break
			case '(':
				tkn, literal, value = token.LEFT_PARENTHESIS, string(chr), string(chr)
				break
			case ')':
				tkn, literal, value = token.RIGHT_PARENTHESIS, string(chr), string(chr)
				break
			case '{':
				tkn, literal, value = token.LEFT_BRACE, string(chr), string(chr)
				break
			case '}':
				tkn, literal, value = token.RIGHT_BRACE, string(chr), string(chr)
				break
			case '[':
				tkn, literal, value = token.LEFT_BRACKET, string(chr), string(chr)
				break
			case ']':
				tkn, literal, value = token.RIGHT_BRACKET, string(chr), string(chr)
				break
			case '.':
				tkn, literal, value = token.DOT, string(chr), string(chr)
				break
			case ',':
				tkn, literal, value = token.COMMA, string(chr), string(chr)
				break
			case ':':
				tkn, literal, value = token.COLON, string(chr), string(chr)
				break
			case ';':
				tkn, literal, value = token.SEMICOLON, string(chr), string(chr)
				break
			case '!':
				tkn = self.switchToken("=", token.NOT_EQUAL, token.NOT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '=':
				tkn = self.switchToken("=", token.EQUAL, token.ASSIGN)
				literal = tkn.String()
				value = tkn.String()
				break
			case '<':
				tkn = self.switchToken("=", token.LESS_OR_EQUAL, token.LESS)
				literal = tkn.String()
				value = tkn.String()
				break
			case '>':
				tkn = self.switchToken("=", token.GREATER_OR_EQUAL, token.GREATER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '&':
				tkn = self.switchToken("&,=", token.LOGICAL_AND, token.AND_ARITHMETIC_ASSIGN, token.AND_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
				break
			case '|':
				tkn = self.switchToken("|,=", token.LOGICAL_OR, token.OR_ARITHMETIC_ASSIGN, token.OR_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
				break
			default:
				tkn = token.ILLEGAL
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

func (self *Parser) scanComment(tkn token.Token) string {
	if tkn == token.MULTI_COMMENT {
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

func (self *Parser) switchToken(keysStr string, tkns ...token.Token) token.Token {
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
	return chr == '"' || chr == '\'' || chr == '`'
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
