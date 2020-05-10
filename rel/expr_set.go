package rel

import (
	"bytes"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SetExpr returns the tuple or set with a single field replaced by an
// expression.
type SetExpr struct {
	ExprScanner
	elements []Expr
}

// NewSetExpr returns a new TupleExpr.
func NewSetExpr(scanner parser.Scanner, elements ...Expr) Expr {
	values := make([]Value, len(elements))
	for i, expr := range elements {
		value, ok := expr.(Value)
		if !ok {
			return &SetExpr{ExprScanner{scanner}, elements}
		}
		values[i] = value
	}
	return NewSet(values...)
}

// Elements returns a Set's elements.
func (e *SetExpr) Elements() []Expr {
	elements := make([]Expr, len(e.elements))
	copy(elements, e.elements)
	return elements
}

// String returns a string representation of the expression.
func (e *SetExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, expr := range e.elements {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(expr.String())
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the subject
func (e *SetExpr) Eval(local Scope) (Value, error) {
	values := make([]Value, 0, len(e.elements))
	for _, expr := range e.elements {
		value, err := EvalExpr(expr, local)
		if err != nil {
			return nil, wrapContext(err, e)
		}
		values = append(values, value)
	}
	return NewSet(values...), nil
}

// NewIntersectExpr evaluates a <&> b.
func NewIntersectExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "<&>", "(%s <&> %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Intersect(x, y), nil
				}
				return nil, errors.Errorf("<&> rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("<&> lhs must be a Set, not %T", a)
		})
}
