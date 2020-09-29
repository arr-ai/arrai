package rel

import (
	"bytes"
	"context"

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
func NewSetExpr(scanner parser.Scanner, elements ...Expr) (Expr, error) {
	values := make([]Value, len(elements))
	for i, expr := range elements {
		value, is := exprIsValue(expr)
		if !is {
			return &SetExpr{ExprScanner{scanner}, elements}, nil
		}
		values[i] = value
	}
	s, err := NewSet(values...)
	if err != nil {
		return nil, err
	}
	return NewLiteralExpr(scanner, s), nil
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
func (e *SetExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	values := make([]Value, 0, len(e.elements))
	for _, expr := range e.elements {
		value, err := expr.Eval(ctx, local)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		values = append(values, value)
	}
	s, err := NewSet(values...)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return s, nil
}

// NewIntersectExpr evaluates a <&> b.
func NewIntersectExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "<&>", "(%s <&> %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Intersect(x, y), nil
				}
				return nil, errors.Errorf("<&> rhs must be a set, not %s", ValueTypeAsString(b))
			}
			return nil, errors.Errorf("<&> lhs must be a set, not %s", ValueTypeAsString(a))
		})
}
