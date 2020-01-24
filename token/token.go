// Package token provides the definitions for token types.
package token

// To add a new token you should first make a type for it. The
// comment on the right of the type is what it will be displayed as
// in errors. If the token is a constant token (i.e. it will always have
// the same string value) then add it to the ConstantTokens map. If it is
// a keyword then add it to the keywords map.

import "strconv"

// Type represents the type of a token.
type Type int

// Definitions for token types.
const (
	TokInteger      Type = iota // integer
	TokIdentifier               // identifier
	TokAssign                   // '='
	TokEquals                   // '=='
	TokLessThan                 // '<'
	TokGreaterThan              // '>'
	TokPlus                     // '+'
	TokDash                     // '-'
	TokStar                     // '*'
	TokFwdSlash                 // '/'
	TokAmpersand                // '&'
	TokIf                       // 'if'
	TokElse                     // 'else'
	TokWhile                    // 'while'
	TokLeftBracket              // '('
	TokRightBracket             // ')'
	TokLeftCurly                // '{'
	TokRightCurly               // '}'
	TokLeftSquare               // ']'
	TokRightSquare              // ']'
	TokSemiColon                // ';'
	TokVar                      // 'var'
	TokInt                      // 'int'
	TokArray                    // 'array'
	TokOf                       // 'of'
	TokPtr                      // 'ptr'
	TokTo                       // 'to'
	TokChar                     // 'char'
	TokNotEqual                 // '!='
	TokNot                      // '!'
)

// SourceInformation holds the source information for a token.
type SourceInformation struct {
	FileName string
	Line     int
}

func (si *SourceInformation) String() string {
	return si.FileName + ":" + strconv.Itoa(si.Line)
}

// Token represents a token.
type Token struct {
	// Type holds the type of the token.
	Type Type
	// Value holds the string value of the token.
	Value string
	// Source holds the source information for the token.
	Source SourceInformation
}

func (t *Token) String() string {
	if t.Type == TokInteger || t.Type == TokIdentifier {
		return "'" + t.Value + "'"
	}
	return t.Type.String()
}

// ConstantTokens contains a mapping of constant tokens to their
// string equivalent. A constant token is a token that will always have
// the same value e.g. '+'
var ConstantTokens = map[Type]string{
	TokAssign:       "=",
	TokEquals:       "==",
	TokLessThan:     "<",
	TokGreaterThan:  ">",
	TokPlus:         "+",
	TokDash:         "-",
	TokStar:         "*",
	TokFwdSlash:     "/",
	TokAmpersand:    "&",
	TokIf:           "if",
	TokElse:         "else",
	TokWhile:        "while",
	TokLeftBracket:  "(",
	TokRightBracket: ")",
	TokLeftCurly:    "{",
	TokRightCurly:   "}",
	TokLeftSquare:   "[",
	TokRightSquare:  "]",
	TokSemiColon:    ";",
	TokVar:          "var",
	TokInt:          "int",
	TokArray:        "array",
	TokOf:           "of",
	TokPtr:          "ptr",
	TokTo:           "to",
	TokChar:         "char",
	TokNotEqual:     "!=",
	TokNot:          "!",
}

// Keywords contains identifiers that are language-level keywords.
var Keywords = map[string]Type{
	"if":    TokIf,
	"while": TokWhile,
	"else":  TokElse,
	"var":   TokVar,
	"int":   TokInt,
	"array": TokArray,
	"of":    TokOf,
	"ptr":   TokPtr,
	"to":    TokTo,
	"char":  TokChar,
}
