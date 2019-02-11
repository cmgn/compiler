package lexer

import (
	"strconv"
	"testing"

	"github.com/cmgn/compiler/token"
)

func TestIntegerLex(t *testing.T) {
	in := "123 456 7 9"
	out := []*token.Token{
		tok(token.TokInteger, "123"),
		tok(token.TokInteger, "456"),
		tok(token.TokInteger, "7"),
		tok(token.TokInteger, "9"),
	}
	runTests(in, out, t)
}

func TestIdentifierLex(t *testing.T) {
	in := "abc def g hi if while else var of array ptr int to char"
	out := []*token.Token{
		tok(token.TokIdentifier, "abc"),
		tok(token.TokIdentifier, "def"),
		tok(token.TokIdentifier, "g"),
		tok(token.TokIdentifier, "hi"),
		tok(token.TokIf, "if"),
		tok(token.TokWhile, "while"),
		tok(token.TokElse, "else"),
		tok(token.TokVar, "var"),
		tok(token.TokOf, "of"),
		tok(token.TokArray, "array"),
		tok(token.TokPtr, "ptr"),
		tok(token.TokInt, "int"),
		tok(token.TokTo, "to"),
		tok(token.TokChar, "char"),
	}
	runTests(in, out, t)
}

func TestSymbolLex(t *testing.T) {
	in := "+-{}=*/==><;&"
	out := []*token.Token{
		tok(token.TokPlus, "+"),
		tok(token.TokDash, "-"),
		tok(token.TokLeftCurly, "{"),
		tok(token.TokRightCurly, "}"),
		tok(token.TokAssign, "="),
		tok(token.TokStar, "*"),
		tok(token.TokFwdSlash, "/"),
		tok(token.TokEquals, "=="),
		tok(token.TokGreaterThan, ">"),
		tok(token.TokLessThan, "<"),
		tok(token.TokSemiColon, ";"),
		tok(token.TokAmpersand, "&"),
	}
	runTests(in, out, t)
}

func TestComplexExpression(t *testing.T) {
	in := "1 + ((2 * abc) - (def / 743))"
	out := []*token.Token{
		tok(token.TokInteger, "1"),
		tok(token.TokPlus, "+"),
		tok(token.TokLeftBracket, "("),
		tok(token.TokLeftBracket, "("),
		tok(token.TokInteger, "2"),
		tok(token.TokStar, "*"),
		tok(token.TokIdentifier, "abc"),
		tok(token.TokRightBracket, ")"),
		tok(token.TokDash, "-"),
		tok(token.TokLeftBracket, "("),
		tok(token.TokIdentifier, "def"),
		tok(token.TokFwdSlash, "/"),
		tok(token.TokInteger, "743"),
		tok(token.TokRightBracket, ")"),
		tok(token.TokRightBracket, ")"),
	}
	runTests(in, out, t)
}

func TestSimpleProgram(t *testing.T) {
	in := `a = 0;
	b = 1;
	while (a < b) {
		a = a + b;
		b = a - b;
	}`
	out := []*token.Token{
		tok(token.TokIdentifier, "a"),
		tok(token.TokAssign, "="),
		tok(token.TokInteger, "0"),
		tok(token.TokSemiColon, ";"),
		tok(token.TokIdentifier, "b"),
		tok(token.TokAssign, "="),
		tok(token.TokInteger, "1"),
		tok(token.TokSemiColon, ";"),
		tok(token.TokWhile, "while"),
		tok(token.TokLeftBracket, "("),
		tok(token.TokIdentifier, "a"),
		tok(token.TokLessThan, "<"),
		tok(token.TokIdentifier, "b"),
		tok(token.TokRightBracket, ")"),
		tok(token.TokLeftCurly, "{"),
		tok(token.TokIdentifier, "a"),
		tok(token.TokAssign, "="),
		tok(token.TokIdentifier, "a"),
		tok(token.TokPlus, "+"),
		tok(token.TokIdentifier, "b"),
		tok(token.TokSemiColon, ";"),
		tok(token.TokIdentifier, "b"),
		tok(token.TokAssign, "="),
		tok(token.TokIdentifier, "a"),
		tok(token.TokDash, "-"),
		tok(token.TokIdentifier, "b"),
		tok(token.TokSemiColon, ";"),
	}
	runTests(in, out, t)
}

func TestLex(t *testing.T) {
	source := "x = 100;"
	expectedOut := []*token.Token{
		tok(token.TokIdentifier, "x"),
		tok(token.TokAssign, "="),
		tok(token.TokInteger, "100"),
		tok(token.TokSemiColon, ";"),
	}
	tokens, err := Lex("test", source)
	if err != nil {
		t.Error("error should not have occurred")
	}
	if len(expectedOut) != len(tokens) {
		t.Errorf(
			"%s %s %d %s %d",
			"For token's length",
			"expected",
			len(expectedOut),
			"got",
			len(tokens),
		)
	}
	for i := 0; i < len(expectedOut); i++ {
		if tokens[i].Source.Line != 1 {
			t.Error(
				"For token's line",
				"expected 1",
				"got", strconv.Itoa(tokens[i].Source.Line),
			)
		} else if tokens[i].Source.FileName != "test" {
			t.Error(
				"For token's file name",
				"expected test",
				"got", tokens[i].Source.FileName,
			)
		} else if !tokenMatches(expectedOut[i], tokens[i]) {
			t.Error(
				"For token's contens",
				"expected ", expectedOut[i].String(),
				"got", tokens[i].String(),
			)
		}
	}
}

func TestInvalidLex(t *testing.T) {
	tokens, err := Lex("test", "@")
	if err == nil {
		t.Error(
			"For invalid input",
			"expected error",
			"got nil",
		)
	} else if tokens != nil {
		t.Error(
			"For invalid input",
			"expected nil",
			"got slice of tokens",
		)
	}
}

func TestLineNumbering(t *testing.T) {
	in := "12\n34\n56"
	lexer := makeLexer(in)
	for i := 0; i < 3; i++ {
		lexer.next()
	}
	if lexer.line != 3 {
		t.Error(
			"For", "12\\n45\\n56",
			"expected", "3",
			"got", strconv.Itoa(lexer.line),
		)
	}
}
func TestMakesError(t *testing.T) {
	in := "@"
	lexer := makeLexer(in)
	lexer.next()
	if lexer.err == nil {
		t.Error(
			"For", in,
			"expected", "error",
			"got", "nil",
		)
	}
}

func runTests(in string, out []*token.Token, t *testing.T) {
	lexer := makeLexer(in)
	for _, token := range out {
		next := lexer.next()
		if !tokenMatches(next, token) {
			t.Error(
				"For", in,
				"expected", token,
				"got", next,
			)
			return
		}
	}
}

func makeLexer(source string) *lexerState {
	return &lexerState{
		source: source,
		line:   1,
	}
}

func tokenMatches(a, b *token.Token) bool {
	return a.Type == b.Type && a.Value == b.Value
}

func tok(typ token.Type, val string) *token.Token {
	return &token.Token{
		Type:  typ,
		Value: val,
	}
}
