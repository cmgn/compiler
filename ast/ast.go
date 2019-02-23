// Package ast provides the abstract syntax tree for the language.
package ast

import (
	"fmt"
	"strings"

	"github.com/cmgn/compiler/token"
)

// UnaryOperatorType is used in the UnaryOperator node to represent
// the operator type.
type UnaryOperatorType int

// Unary operator type definitions.
const (
	UnaryDereference UnaryOperatorType = iota // '*'
	UnaryMinus                                // '-'
	UnaryAddress                              // '&'
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
	BinaryEqual                                 // '=='
	BinaryNotEqual                              // '!='
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

// Type is the interface implemented by all type node types.
type Type interface {
	Node
	typeNode()
	// Size gets the number of bytes occupied by a type on the stack.
	Size() int
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

// Declaration represents a variable declaration statement.
type Declaration struct {
	Source token.SourceInformation
	Name   string
	Type   Type
}

func (d *Declaration) String() string {
	return fmt.Sprintf(
		"Declaration[%s, %s]",
		d.Name,
		d.Type.String(),
	)
}

// SourceInfo retrieves the source information for the 'var' keyword
// in the delcaration.
func (d *Declaration) SourceInfo() *token.SourceInformation {
	return &d.Source
}

func (d *Declaration) statementNode() {}

// IfStatement represents an occurrence of an if statement. Both ifs with &
// without an else are represented by this, in the latter case Statement2 will
// be the empty statement.
type IfStatement struct {
	Source     token.SourceInformation
	Condition  Expression
	Statement1 Statement
	Statement2 Statement
}

// SourceInfo gets the source information for the 'if' token part of
// the if statment.
func (i *IfStatement) SourceInfo() *token.SourceInformation {
	return &i.Source
}

func (i *IfStatement) String() string {
	return fmt.Sprintf(
		"If[%s, %s, %s]",
		i.Condition.String(),
		i.Statement1.String(),
		i.Statement2.String(),
	)
}

func (i *IfStatement) statementNode() {}

// WhileStatement is a 'while' statement.
type WhileStatement struct {
	Source    token.SourceInformation
	Condition Expression
	Statement Statement
}

// SourceInfo gets the source information for the 'while' keyword part
// of the while statement.
func (w *WhileStatement) SourceInfo() *token.SourceInformation {
	return &w.Source
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf(
		"While[%s, %s]",
		w.Condition.String(),
		w.Statement.String(),
	)
}

func (w *WhileStatement) statementNode() {}

// BlockStatement is a series of statements surrounded by curly brackets.
type BlockStatement struct {
	Source     token.SourceInformation
	Statements []Statement
}

// SourceInfo gets the source information for the opening bracket
// of the block.
func (b *BlockStatement) SourceInfo() *token.SourceInformation {
	return &b.Source
}

func (b *BlockStatement) String() string {
	strs := make([]string, len(b.Statements))
	for i, statement := range b.Statements {
		strs[i] = statement.String()
	}
	return fmt.Sprintf(
		"Block[%s]",
		strings.Join(strs, ", "),
	)
}

func (b *BlockStatement) statementNode() {}

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

// UnaryOperator represents an occurrence of a unary operator
// expression.
type UnaryOperator struct {
	Type  UnaryOperatorType
	Value Expression
}

// SourceInfo gets the source information for the operator inside the
// unary operator node.
func (u *UnaryOperator) SourceInfo() *token.SourceInformation {
	return u.Value.SourceInfo()
}

func (u *UnaryOperator) String() string {
	return fmt.Sprintf(
		"UnaryOperator[%s, %s]",
		u.Type.String(),
		u.Value.String(),
	)
}

func (u *UnaryOperator) expressionNode() {}

// PrimitiveType is used in the Primitive node to represent which primitive
// type is contained in it.
type PrimitiveType int

// Primitive type definitions.
const (
	IntType  PrimitiveType = iota // 'int'
	CharType                      // 'char'
)

// Primitive is the type for primitive machine types such as integers
// and characters.
type Primitive struct {
	Source token.SourceInformation
	Type   PrimitiveType
}

// SourceInfo gets the source information for where the primitive type occurred.
func (p *Primitive) SourceInfo() *token.SourceInformation {
	return &p.Source
}

func (p *Primitive) String() string {
	return p.Type.String()
}

// Size gets the size of the contained primitive type.
func (p *Primitive) Size() int {
	switch p.Type {
	case IntType:
		return 8
	case CharType:
		return 1
	}
	return 0
}

func (p *Primitive) typeNode() {}

// ArrayType is the type for fixed-length statically allocated arrays.
type ArrayType struct {
	Source token.SourceInformation
	Length int
	Type   Type
}

// SourceInfo gets the source information for where the array type is defined.
func (a *ArrayType) SourceInfo() *token.SourceInformation {
	return &a.Source
}

func (a *ArrayType) String() string {
	return fmt.Sprintf(
		"Array[%d, %s]",
		a.Length,
		a.Type.String(),
	)
}

// Size gets the size of the array in bytes, which is the length times the size of
// the the array's type.
func (a *ArrayType) Size() int {
	return a.Type.Size() * a.Length
}

func (a *ArrayType) typeNode() {}

// PointerType represents an occurrence of a pointer type in the program.
type PointerType struct {
	Source token.SourceInformation
	Type   Type
}

// SourceInfo gets the source information for the 'ptr' keyword part of the
// occurrence.
func (p *PointerType) SourceInfo() *token.SourceInformation {
	return &p.Source
}

func (p *PointerType) String() string {
	return fmt.Sprintf("Pointer[%s]", p.Type.String())
}

// Size gets the size of a pointer in bytes, which is always eight bytes.
func (p *PointerType) Size() int {
	return 8
}

func (p *PointerType) typeNode() {}
