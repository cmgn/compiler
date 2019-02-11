// Package lexer implements a lexer capable of transforming a string
// into the token types contained in package token.
package lexer

import (
	"errors"

	"github.com/cmgn/compiler/token"
)

// Lex lexes a string and returns the tokens encountered, or nil and an error
// if it is an invalid string. The filename parameter is used in creating the
// source information for the tokens.
func Lex(filename string, contents string) ([]*token.Token, error) {
	tokens := make([]*token.Token, 0)
	lexer := &lexerState{
		fname:  filename,
		source: contents,
		line:   1,
	}
	for !lexer.empty() {
		tok := lexer.next()
		if tok == nil {
			break
		}
		tokens = append(tokens, tok)
	}
	if lexer.err != nil {
		return nil, lexer.err
	}
	return tokens, nil
}

// lexerState represents the state of a lexer.
type lexerState struct {
	// fname is the name of the source file.
	fname string
	// source is the source string.
	source string
	// line is the current line number.
	line int
	// pos is the current position in the string.
	pos int
	// err is the error if one has been countered, nil otherwise.
	err error
}

// curr returns the current byte.
func (l *lexerState) curr() byte {
	return l.source[l.pos]
}

// empty checks if there's more bytes.
func (l *lexerState) empty() bool {
	return l.pos >= len(l.source)
}

// sourceInfo creates the source information for the current position.
func (l *lexerState) sourceInfo() token.SourceInformation {
	return token.SourceInformation{
		FileName: l.fname,
		Line:     l.line,
	}
}

// buildToken builds a token with a given value and type, using the current
// position's source info.
func (l *lexerState) buildToken(typ token.Type, val string) *token.Token {
	return &token.Token{
		Type:   typ,
		Value:  val,
		Source: l.sourceInfo(),
	}
}

// buildConstantToken builds a constant token using the buildToken method.
func (l *lexerState) buildConstantToken(typ token.Type) *token.Token {
	val, ok := token.ConstantTokens[typ]
	// This isn't an error we should handle gracefully, it's a logic error.
	if !ok {
		panic("called with non-constant token")
	}
	return &token.Token{
		Type:   typ,
		Value:  val,
		Source: l.sourceInfo(),
	}
}

// error sets the error field.
func (l *lexerState) error(msg string) {
	l.err = errors.New(msg)
}

func (l *lexerState) readIdentifier() *token.Token {
	start := l.pos
	for !l.empty() && (isAlpha(l.curr()) || isDigit(l.curr())) {
		l.pos++
	}
	ident := l.source[start:l.pos]
	if typ, ok := token.Keywords[ident]; ok {
		return l.buildConstantToken(typ)
	}
	return l.buildToken(token.TokIdentifier, ident)
}

func (l *lexerState) readInteger() *token.Token {
	start := l.pos
	for !l.empty() && isDigit(l.curr()) {
		l.pos++
	}
	return l.buildToken(token.TokInteger, l.source[start:l.pos])
}

// next gets the next token, it returns nil and sets the err field to an error
// if it encounters an invalid character.
func (l *lexerState) next() *token.Token {
loop:
	for l.pos < len(l.source) {
		curr := l.curr()
		if isSpace(curr) {
			if curr == '\n' {
				l.line++
			}
			l.pos++
			continue
		} else if isAlpha(curr) {
			return l.readIdentifier()
		} else if isDigit(curr) {
			return l.readInteger()
		} else if typ, ok := byteTokens[curr]; ok {
			l.pos++
			return l.buildConstantToken(typ)
		}
		switch curr {
		case '=':
			l.pos++
			if l.curr() == '=' {
				l.pos++
				return l.buildConstantToken(token.TokEquals)
			}
			return l.buildConstantToken(token.TokAssign)
		default:
			l.error("invalid character: " + string(curr))
			break loop
		}
	}
	return nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func isAlpha(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// NB: tokens such as '=' are not in here as they could potentially
// be a multibyte token.
var byteTokens = map[byte]token.Type{
	'+': token.TokPlus,
	'-': token.TokDash,
	'*': token.TokStar,
	';': token.TokSemiColon,
	'/': token.TokFwdSlash,
	'(': token.TokLeftBracket,
	')': token.TokRightBracket,
	'{': token.TokLeftCurly,
	'}': token.TokRightCurly,
	'<': token.TokLessThan,
	'>': token.TokGreaterThan,
}
