// Package ast provides the abstract syntax tree for the language.
package ast

import (
	"fmt"

	"github.com/cmgn/compiler/token"
)

// BinaryOperatorType is used in the BinaryOperator node to represent
// the operator type.
type BinaryOperatorType int

// Binary operator type definitions
const (
	BinaryAdd         BinaryOperatorType = iota // '+'
	BinarySub                                   // '-'
	BinaryMul                                   // '*'
	BinaryDiv                                   // '/'
	BinaryLessThan                              // '<'
	BinaryGreaterThan                           // '>'
	BinaryEquals                                // '='
)

// Node is the interface implemented by all syntax tree nodes.
type Node interface {
	SourceInfo() *token.SourceInformation
	String() string
}

// Statement is the interface implemented by all statement node types.
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface implemented by all expression node types.
type Expression interface {
	Node
	expressionNode()
}

// Empty represents an empty statement. The empty statement is used in2
// cases such as "while (something);".
type Empty struct {
	Source token.SourceInformation
}

// SourceInfo gets the source information for the empty statement. This is
// the location of its semicolon.
func (e *Empty) SourceInfo() *token.SourceInformation {
	return &e.Source
}

func (e *Empty) String() string {
	return "Empty[]"
}

func (e *Empty) statementNode() {}

// ExpressionStatement represents an expression followed by a semicolon.
type ExpressionStatement struct {
	Expression Expression
}

// SourceInfo gets the source information for the expression.
func (e *ExpressionStatement) SourceInfo() *token.SourceInformation {
	return e.Expression.SourceInfo()
}

func (e *ExpressionStatement) String() string {
	return "ExpressionStatement[" + e.Expression.String() + "]"
}

func (e *ExpressionStatement) statementNode() {}

// Assignment is an assignment statement.
type Assignment struct {
	Source token.SourceInformation
	Left   Expression
	Right  Expression
}

// SourceInfo gets the source information for the assignment.
func (a *Assignment) SourceInfo() *token.SourceInformation {
	return &a.Source
}

func (a *Assignment) String() string {
	return fmt.Sprintf("Assignment[%s, %s]", a.Left.String(), a.Right.String())
}

func (a *Assignment) statementNode() {}

// Integer is an integer expression.
type Integer struct {
	Source token.SourceInformation
	Value  string
}

// SourceInfo gets the source information for the integer.
func (i *Integer) SourceInfo() *token.SourceInformation {
	return &i.Source
}

func (i *Integer) String() string {
	return i.Value
}

func (i *Integer) expressionNode() {}

// Variable is a variable expression.
type Variable struct {
	Source token.SourceInformation
	Value  string
}

// SourceInfo gets the source information for the variable.
func (v *Variable) SourceInfo() *token.SourceInformation {
	return &v.Source
}

func (v *Variable) String() string {
	return v.Value
}

func (v *Variable) expressionNode() {}

// BinaryOperator represents an occurrence of a binary operator
// expression.
type BinaryOperator struct {
	Type  BinaryOperatorType
	Left  Expression
	Right Expression
}

// SourceInfo gets the source information for the left operand of the
// operator expression.
func (b *BinaryOperator) SourceInfo() *token.SourceInformation {
	return b.Left.SourceInfo()
}

func (b *BinaryOperator) String() string {
	return fmt.Sprintf(
		"BinaryOperator[%s, %s, %s]",
		b.Type.String(),
		b.Left.String(),
		b.Right.String(),
	)
}

func (b *BinaryOperator) expressionNode() {}
