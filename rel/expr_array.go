package rel

import (
	"bytes"

	"github.com/arr-ai/wbnf/parser"
)

// ArrayExpr represents an expr that evaluates to an Array.
type ArrayExpr struct {
	ExprScanner
	elements []Expr
}

// NewArrayExpr returns a new Expr that constructs an Array.
func NewArrayExpr(scanner parser.Scanner, elements ...Expr) Expr {
	values := make([]Value, 0, len(elements))
	for _, expr := range elements {
		if expr != nil {
			if value, ok := expr.(Value); ok {
				values = append(values, value)
				continue
			}
		}
		return ArrayExpr{ExprScanner{scanner}, elements}
	}
	return NewArray(values...)
}

// Elements returns a Set's elements.
func (e ArrayExpr) Elements() []Expr {
	elements := make([]Expr, len(e.elements))
	copy(elements, e.elements)
	return elements
}

// String returns a string representation of the expression.
func (e ArrayExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('[')
	for i, expr := range e.elements {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(expr.String())
	}
	b.WriteByte(']')
	return b.String()
}

// Eval returns the subject.
func (e ArrayExpr) Eval(local Scope) (Value, error) {
	values := make([]Value, 0, len(e.elements))
	for _, expr := range e.elements {
		var value Value
		if expr != nil {
			var err error
			value, err = expr.Eval(local)
			if err != nil {
				return nil, wrapContext(err, e, local)
			}
		}
		values = append(values, value)
	}
	return NewArray(values...), nil
}
