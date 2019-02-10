package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cmgn/compiler/lexer"
	"github.com/cmgn/compiler/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens, err := lexer.Lex("stdin", scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		stmts, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, stmt := range stmts {
			fmt.Println(stmt.String())
		}
	}
}
