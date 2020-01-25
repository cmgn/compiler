// Package parser provides a parser capable of parsing a slice of tokens into
// the AST provided by package ast.
package parser

import (
	"fmt"
	"strconv"

	"github.com/cmgn/compiler/ast"
	"github.com/cmgn/compiler/token"
)

// Parse parses a slice of tokens into a syntax tree. If the input is invalid
// then nil, error is returned.
func Parse(tokens []*token.Token) ([]ast.Statement, error) {
	parser := &parser{toks: tokens}
	statements := make([]ast.Statement, 0)
	for !parser.empty() {
		stmt := parser.statement()
		if stmt == nil {
			break
		}
		statements = append(statements, stmt)
	}
	if parser.err != nil {
		return nil, parser.err
	}
	return statements, nil
}

type parser struct {
	toks []*token.Token
	pos  int
	err  error
}

func (p *parser) empty() bool {
	return p.pos >= len(p.toks)
}

func (p *parser) curr() *token.Token {
	if p.empty() {
		return nil
	}
	return p.toks[p.pos]
}

func (p *parser) expect(typ token.Type) bool {
	curr := p.curr()
	if curr == nil {
		curr = p.toks[p.pos-1]
		p.err = fmt.Errorf("[%s] unexpected end of input after %s, expected %s",
			curr.Source.String(), curr.String(), typ.String())
		return false
	}
	if curr.Type != typ {
		p.err = fmt.Errorf("[%s] expected %s, got %s",
			curr.Source.String(), typ.String(), curr.String())
		return false
	}
	p.pos++
	return true
}

func (p *parser) unexpected(curr *token.Token) {
	p.err = fmt.Errorf("[%s] unexpected %s", curr.Source.String(), curr.String())
}

func (p *parser) unexpectedEnd() bool {
	if p.empty() {
		prev := p.toks[p.pos-1]
		p.err = fmt.Errorf("[%s] unexpected end of input after %s", prev.Source.String(), prev.String())
		return true
	}
	return false
}

func (p *parser) next() *token.Token {
	p.pos++
	if p.empty() {
		return nil
	}
	return p.curr()
}

// statement
// | expression '=' expression ';'
// | expression ';'
// | 'var' identifier typedecl ';'
// | 'if' expression statement ['else' statement]
// | block
// | ';'
func (p *parser) statement() ast.Statement {
	if p.unexpectedEnd() {
		return nil
	}

	curr := p.curr()
	switch curr.Type {
	case token.TokSemiColon:
		p.pos++
		return &ast.Empty{Source: curr.Source}
	case token.TokVar:
		p.pos++
		name := p.curr()
		if !p.expect(token.TokIdentifier) {
			return nil
		}
		typ := p.typedecl()
		if typ == nil {
			return nil
		}
		if !p.expect(token.TokSemiColon) {
			return nil
		}
		return &ast.Declaration{
			Source: curr.Source,
			Name:   name.Value,
			Type:   typ,
		}
	case token.TokIf:
		p.expect(token.TokIf)
		cond := p.expression()
		if cond == nil {
			return nil
		}
		stmt1 := p.statement()
		if stmt1 == nil {
			return nil
		}
		if p.empty() || p.curr().Type != token.TokElse {
			return &ast.IfStatement{
				Source:     curr.Source,
				Condition:  cond,
				Statement1: stmt1,
				Statement2: &ast.Empty{},
			}
		}
		p.expect(token.TokElse)
		stmt2 := p.statement()
		if stmt2 == nil {
			return nil
		}
		return &ast.IfStatement{
			Source:     curr.Source,
			Condition:  cond,
			Statement1: stmt1,
			Statement2: stmt2,
		}
	case token.TokWhile:
		p.expect(token.TokWhile)
		cond := p.expression()
		if cond == nil {
			return nil
		}
		stmt := p.statement()
		if stmt == nil {
			return nil
		}
		return &ast.WhileStatement{
			Source:    curr.Source,
			Condition: cond,
			Statement: stmt,
		}
	case token.TokLeftCurly:
		return p.block()
	}

	expr := p.expression()
	if expr == nil || p.unexpectedEnd() {
		return nil
	}

	middle := p.curr()
	if middle.Type == token.TokAssign {
		p.expect(token.TokAssign)
		right := p.expression()
		if right == nil {
			return nil
		}
		if !p.expect(token.TokSemiColon) {
			return nil
		}
		return &ast.Assignment{
			Left:   expr,
			Right:  right,
			Source: middle.Source,
		}
	}
	if p.expect(token.TokSemiColon) {
		return &ast.ExpressionStatement{
			Expression: expr,
		}
	}
	return nil
}

// block
// | '{' {statement} '}'
func (p *parser) block() ast.Statement {
	curr := p.curr()
	if !p.expect(token.TokLeftCurly) {
		return nil
	}
	statements := make([]ast.Statement, 0)
	for !p.empty() && p.curr().Type != token.TokRightCurly {
		stmt := p.statement()
		if stmt == nil {
			return nil
		}
		statements = append(statements, stmt)
	}
	if !p.expect(token.TokRightCurly) {
		return nil
	}
	return &ast.BlockStatement{
		Source:     curr.Source,
		Statements: statements,
	}
}

// typedecl
// | 'int'
// | 'char'
// | 'array' '(' integer ')' 'of' typedecl
// | '(' typedecl ')'
func (p *parser) typedecl() ast.Type {
	if p.unexpectedEnd() {
		return nil
	}
	curr := p.curr()
	switch curr.Type {
	case token.TokLeftBracket:
		p.expect(token.TokLeftBracket)
		typ := p.typedecl()
		if typ == nil {
			return nil
		}
		if !p.expect(token.TokRightBracket) {
			return nil
		}
		return typ
	case token.TokInt:
		p.expect(token.TokInt)
		return &ast.Primitive{
			Type:   ast.IntType,
			Source: curr.Source,
		}
	case token.TokChar:
		p.expect(token.TokChar)
		return &ast.Primitive{
			Type:   ast.CharType,
			Source: curr.Source,
		}
	case token.TokArray:
		p.expect(token.TokArray)
		if !p.expect(token.TokLeftBracket) {
			return nil
		}
		size := p.curr()
		if !p.expect(token.TokInteger) {
			return nil
		}
		if !p.expect(token.TokRightBracket) {
			return nil
		}
		if !p.expect(token.TokOf) {
			return nil
		}
		typ := p.typedecl()
		if typ == nil {
			return nil
		}
		sizeInt, err := strconv.Atoi(size.Value)
		if err != nil {
			p.err = fmt.Errorf("[%s] invalid static array size '%s'",
				size.Source.String(), size.Value)
		}
		return &ast.ArrayType{
			Type:   typ,
			Length: sizeInt,
			Source: curr.Source,
		}
	case token.TokPtr:
		p.expect(token.TokPtr)
		if !p.expect(token.TokTo) {
			return nil
		}
		typ := p.typedecl()
		if typ == nil {
			return nil
		}
		return &ast.PointerType{
			Source: curr.Source,
			Type:   typ,
		}
	}
	p.unexpected(curr)
	return nil
}

// expression
// | equality
func (p *parser) expression() ast.Expression {
	return p.equality()
}

// equality
// | comparison '==' comparison
// | comparison '!=' comparison
// | comparison
func (p *parser) equality() ast.Expression {
	left := p.comparison()
	if left == nil {
		return nil
	}
loop:
	for !p.empty() {
		curr := p.curr()
		switch curr.Type {
		case token.TokEquals:
			p.expect(token.TokEquals)
			right := p.comparison()
			if right == nil {
				return nil
			}
			left = &ast.BinaryOperator{
				Type:  ast.BinaryEqual,
				Left:  left,
				Right: right,
			}
		case token.TokNotEqual:
			p.expect(token.TokNotEqual)
			right := p.comparison()
			if right == nil {
				return nil
			}
			left = &ast.BinaryOperator{
				Type:  ast.BinaryNotEqual,
				Left:  left,
				Right: right,
			}
		default:
			break loop
		}
	}
	return left
}

// comparison
// | summation '>' summation
// | summation '<' summation
// | summation
func (p *parser) comparison() ast.Expression {
	left := p.summation()
	if left == nil {
		return nil
	} else if p.empty() {
		return left
	}
	curr := p.curr()
	switch curr.Type {
	case token.TokLessThan:
		p.expect(token.TokLessThan)
		right := p.summation()
		if right == nil {
			return nil
		}
		return &ast.BinaryOperator{
			Type:  ast.BinaryLessThan,
			Left:  left,
			Right: right,
		}
	case token.TokGreaterThan:
		p.expect(token.TokGreaterThan)
		right := p.summation()
		if right == nil {
			return nil
		}
		return &ast.BinaryOperator{
			Type:  ast.BinaryGreaterThan,
			Left:  left,
			Right: right,
		}
	}
	return left
}

// summation
// | summation '+' product
// | summation '-' product
// | product
func (p *parser) summation() ast.Expression {
	prod := p.product()
	if prod == nil {
		return nil
	}
loop:
	for !p.empty() {
		curr := p.curr()
		switch curr.Type {
		case token.TokPlus:
			p.expect(token.TokPlus)
			right := p.product()
			if right == nil {
				return nil
			}
			prod = &ast.BinaryOperator{
				Type:  ast.BinaryAdd,
				Left:  prod,
				Right: right,
			}
		case token.TokDash:
			p.expect(token.TokDash)
			right := p.product()
			if right == nil {
				return nil
			}
			prod = &ast.BinaryOperator{
				Type:  ast.BinarySub,
				Left:  prod,
				Right: right,
			}
		default:
			break loop
		}
	}
	return prod
}

// product
// | product '*' subscript
// | product '/' subscript
// | subscript
func (p *parser) product() ast.Expression {
	term := p.subscript()
	if term == nil {
		return nil
	}
loop:
	for !p.empty() {
		curr := p.curr()
		switch curr.Type {
		case token.TokStar:
			p.expect(token.TokStar)
			right := p.subscript()
			if right == nil {
				return nil
			}
			term = &ast.BinaryOperator{
				Type:  ast.BinaryMul,
				Left:  term,
				Right: right,
			}
		case token.TokFwdSlash:
			p.expect(token.TokFwdSlash)
			right := p.subscript()
			if right == nil {
				return nil
			}
			term = &ast.BinaryOperator{
				Type:  ast.BinaryDiv,
				Left:  term,
				Right: right,
			}
		default:
			break loop
		}
	}
	return term
}

// subscript
// | subscript '[' expression ']'
// | terminal
func (p *parser) subscript() ast.Expression {
	term := p.terminal()
	for !p.empty() && p.curr().Type == token.TokLeftSquare {
		p.expect(token.TokLeftSquare)
		index := p.expression()
		if !p.expect(token.TokRightSquare) {
			return nil
		}
		term = &ast.Subscript{Value: term, Index: index}
	}
	return term
}

// terminal
// | integer
// | variable
// | '(' expression ')'
// | '-' terminal
// | '*' terminal
// | '&' terminal
func (p *parser) terminal() ast.Expression {
	if p.unexpectedEnd() {
		return nil
	}
	curr := p.curr()
	switch curr.Type {
	case token.TokInteger:
		p.pos++
		return &ast.Integer{
			Source: curr.Source,
			Value:  curr.Value,
		}
	case token.TokIdentifier:
		p.pos++
		return &ast.Variable{
			Source: curr.Source,
			Value:  curr.Value,
		}
	case token.TokLeftBracket:
		if !p.expect(token.TokLeftBracket) {
			return nil
		}
		expr := p.expression()
		if expr == nil {
			return nil
		}
		if !p.expect(token.TokRightBracket) {
			return nil
		}
		return expr
	case token.TokStar:
		p.expect(token.TokStar)
		term := p.terminal()
		if term == nil {
			return nil
		}
		return &ast.UnaryOperator{
			Type:  ast.UnaryDereference,
			Value: term,
		}
	case token.TokDash:
		p.expect(token.TokDash)
		term := p.terminal()
		if term == nil {
			return nil
		}
		return &ast.UnaryOperator{
			Type:  ast.UnaryMinus,
			Value: term,
		}
	case token.TokAmpersand:
		p.expect(token.TokAmpersand)
		term := p.terminal()
		if term == nil {
			return nil
		}
		return &ast.UnaryOperator{
			Type:  ast.UnaryAddress,
			Value: term,
		}
	}
	p.unexpected(curr)
	return nil
}
