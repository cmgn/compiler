package main

import (
	"fmt"

	"github.com/cmgn/compiler/lexer"
	"github.com/cmgn/compiler/parser"
)

func main() {
	input := `y = 3 < 2;`
	tokens, err := lexer.Lex("stdin", input)
	if err != nil {
		fmt.Println(err)
		return
	}
	stmts, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, stmt := range stmts {
		fmt.Println(stmt.String())
	}
}
