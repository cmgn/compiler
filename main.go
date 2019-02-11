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
	fmt.Print("> ")
	for scanner.Scan() {
		tokens, err := lexer.Lex("stdin", scanner.Text())
		if err != nil {
			fmt.Println(err)
		} else {
			stmts, err := parser.Parse(tokens)
			if err != nil {
				fmt.Println(err)
			} else {
				for _, stmt := range stmts {
					fmt.Println(stmt.String())
				}
			}
		}
		fmt.Print("> ")
	}
}
