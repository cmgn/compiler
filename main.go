package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cmgn/compiler/lexer"
	"github.com/cmgn/compiler/parser"
)

func runString(filename, str string) {
	tokens, err := lexer.Lex(filename, str)
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

func mustRead(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(contents)
}

func main() {
	if len(os.Args) == 1 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			runString("<stdin>", scanner.Text())
		}
		return
	}

	for _, filename := range os.Args[1:] {
		runString(filename, mustRead(filename))
	}
}
