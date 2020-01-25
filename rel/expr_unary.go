package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

type unaryEval func(a Value, local, global *Scope) (Value, error)

// UnaryExpr represents a range of operators.
type UnaryExpr struct {
	a      Expr
	op     string
	format string
	eval   unaryEval
}

func newUnaryExpr(a Expr, op, format string, eval unaryEval) Expr {
	return &UnaryExpr{a, op, format, eval}
}

// NewPosExpr evaluates to a.
func NewPosExpr(a Expr) Expr {
	return newUnaryExpr(a, "+", "(+%s)",
		func(a Value, _, _ *Scope) (Value, error) {
			return a, nil
		},
	)
}

// NewNegExpr evaluates to -a.
func NewNegExpr(a Expr) Expr {
	return newUnaryExpr(a, "-", "(-%s)",
		func(a Value, _, _ *Scope) (Value, error) {
			return a.Negate(), nil
		},
	)
}

// NewPowerSetExpr evaluates to ^a.
func NewPowerSetExpr(a Expr) Expr {
	return newUnaryExpr(a, "**", "(**%s)",
		func(a Value, _, _ *Scope) (Value, error) {
			if s, ok := a.(Set); ok {
				return PowerSet(s), nil
			}
			return nil, errors.Errorf("eval arg must be a Set, not %T", a)
		},
	)
}

// NewNotExpr evaluates to !a.
func NewNotExpr(a Expr) Expr {
	return newUnaryExpr(a, "!", "(!%s)",
		func(a Value, _, _ *Scope) (Value, error) {
			return NewBool(!a.Bool()), nil
		})
}

// NewEvalExpr evaluates to *a, given a set lhs.
func NewEvalExpr(a Expr) Expr {
	return newUnaryExpr(a, "*", "(*%s)",
		func(a Value, local, global *Scope) (Value, error) {
			if x, ok := a.(*Function); ok {
				return x.Call(None, local, global)
			}
			return nil, errors.Errorf("eval arg must be a Function, not %T", a)
		})
}

// NewCountExpr evaluates to the number of elements in a.
func NewCountExpr(a Expr) Expr {
	return newUnaryExpr(a, "count", "(%s count)",
		func(a Value, local, global *Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return NewNumber(float64(x.Count())), nil
			}
			return nil, errors.Errorf("eval arg must be a Function, not %T", a)
		})
}

// Arg returns the UnaryExpr's arg.
func (e *UnaryExpr) Arg() Expr {
	return e.a
}

// String returns a string representation of the expression.
func (e *UnaryExpr) String() string {
	return fmt.Sprintf(e.format, e.a)
}

// Eval returns the subject
func (e *UnaryExpr) Eval(local, global *Scope) (Value, error) {
	a, err := e.a.Eval(local, global)
	if err != nil {
		return nil, err
	}
	return e.eval(a, local, global)
}
