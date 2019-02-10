package parser

import (
	"fmt"
	"testing"

	"github.com/cmgn/compiler/ast"
	"github.com/cmgn/compiler/token"
)

func TestTerminalInteger(t *testing.T) {
	in := toks(tok(token.TokInteger, "123"))
	parser := makeParser(in)
	term := parser.terminal()
	if _, ok := term.(*ast.Integer); !ok {
		t.Error(
			"For", "123",
			"expected", "integer",
			"got", term,
		)
	}
}

func TestTerminalVariable(t *testing.T) {
	in := toks(tok(token.TokIdentifier, "abc"))
	parser := makeParser(in)
	term := parser.terminal()
	if _, ok := term.(*ast.Variable); !ok {
		t.Error(
			"For", "123",
			"expected", "variable",
			"got", term,
		)
	}
}

func TestTerminalBrackets(t *testing.T) {
	in := toks(
		tok(token.TokLeftBracket, "("),
		tok(token.TokInteger, "123"),
		tok(token.TokRightBracket, ")"),
	)
	parser := makeParser(in)
	term := parser.terminal()
	if _, ok := term.(*ast.Integer); !ok {
		t.Error(
			"For", "123",
			"expected", "integer",
			"got", term,
		)
	}
}

func TestProductTimes(t *testing.T) {
	in := toks(
		tok(token.TokInteger, "123"),
		tok(token.TokTimes, "*"),
		tok(token.TokInteger, "456"),
	)

	parser := makeParser(in)
	prod := parser.product()
	bin, ok := prod.(*ast.BinaryOperator)
	if !ok {
		t.Error(
			"For", "123 * 456",
			"expected", "binary operator",
			"got", prod,
		)
	} else if bin.Type != ast.BinaryMul {
		t.Error(
			"For", "123 * 456",
			"expected", "BinaryMul",
			"got", prod,
		)
	}
}

func TestProductDivide(t *testing.T) {
	in := toks(
		tok(token.TokInteger, "123"),
		tok(token.TokDivide, "/"),
		tok(token.TokInteger, "456"),
	)

	parser := makeParser(in)
	prod := parser.product()
	bin, ok := prod.(*ast.BinaryOperator)
	if !ok {
		t.Error(
			"For", "123 / 456",
			"expected", "binary operator",
			"got", prod,
		)
	} else if bin.Type != ast.BinaryDiv {
		t.Error(
			"For", "123 / 456",
			"expected", "BinaryDiv",
			"got", bin.Type.String(),
		)
	}
}

func TestAssignmentStatement(t *testing.T) {
	in := toks(
		tok(token.TokIdentifier, "abc"),
		tok(token.TokAssign, "="),
		tok(token.TokInteger, "123"),
		tok(token.TokSemiColon, ";"),
	)
	parser := makeParser(in)
	stmt := parser.statement()
	if _, ok := stmt.(*ast.Assignment); !ok {
		fmt.Println(parser.err)
		t.Error(
			"For", "abc = 123;",
			"expected", "assign",
			"got", stmt,
		)
		return
	}
}

func tok(typ token.Type, val string) *token.Token {
	return &token.Token{Type: typ, Value: val}
}

func toks(tokens ...*token.Token) []*token.Token {
	return tokens
}

func makeParser(input []*token.Token) *parser {
	return &parser{
		toks: input,
	}
}
