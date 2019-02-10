// Package parser provides a parser capable of parsing a slice of tokens into
// the AST provided by package ast.
package parser

import (
	"fmt"

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
	}
	if curr.Type != typ {
		p.err = fmt.Errorf("[%s] expected %s, got %s", curr.Source.String(), typ.String(), curr.Type.String())
		return false
	}
	p.pos++
	return true
}

func (p *parser) unexpectedEnd() bool {
	if p.empty() {
		p.err = fmt.Errorf("[%s] unexpected end of input", p.toks[p.pos-1].Source.String())
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
	}

	expr := p.expression()
	if expr == nil {
		return nil
	} else if p.unexpectedEnd() {
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
	} else if p.expect(token.TokSemiColon) {
		return &ast.ExpressionStatement{
			Expression: expr,
		}
	}
	return nil
}

// expression
// | equality
func (p *parser) expression() ast.Expression {
	return p.equality()
}

// equality
// | comparison '=' comparison
// | comparison
func (p *parser) equality() ast.Expression {
	left := p.comparison()
	if left == nil {
		return nil
	}
	for !p.empty() {
		curr := p.curr()
		if curr.Type != token.TokEquals {
			break
		}
		p.expect(token.TokEquals)
		right := p.comparison()
		if right == nil {
			return nil
		}
		left = &ast.BinaryOperator{
			Type:  ast.BinaryEquals,
			Left:  left,
			Right: right,
		}
	}
	return left
}

// comparison
// | summation ">" summation
// | summation "<" summation
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
		case token.TokMinus:
			p.expect(token.TokMinus)
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
// | product '*' terminal
// | product '/' terminal
// | terminal
func (p *parser) product() ast.Expression {
	term := p.terminal()
	if term == nil {
		return nil
	}
loop:
	for !p.empty() {
		curr := p.curr()
		switch curr.Type {
		case token.TokTimes:
			p.expect(token.TokTimes)
			right := p.terminal()
			if right == nil {
				return nil
			}
			term = &ast.BinaryOperator{
				Type:  ast.BinaryMul,
				Left:  term,
				Right: right,
			}
		case token.TokDivide:
			p.expect(token.TokDivide)
			right := p.terminal()
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

// terminal
// | integer
// | variable
// | "(" expression ")"
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
	}
	p.err = fmt.Errorf("[%s] unexpected %s", curr.Source.String(), curr.String())
	return nil
}
