package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type unaryEval func(a Value, local Scope) (Value, error)

// UnaryExpr represents a range of operators.
type UnaryExpr struct {
	ExprScanner
	a      Expr
	op     string
	format string
	eval   unaryEval
}

func newUnaryExpr(scanner parser.Scanner, a Expr, op, format string, eval unaryEval) Expr {
	return &UnaryExpr{ExprScanner{scanner}, a, op, format, eval}
}

// NewPosExpr evaluates to a.
func NewPosExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "+", "(+%s)",
		func(a Value, _ Scope) (Value, error) { return a, nil },
	)
}

// NewNegExpr evaluates to -a.
func NewNegExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "-", "(-%s)",
		func(a Value, _ Scope) (Value, error) { return a.Negate(), nil },
	)
}

// NewPowerSetExpr evaluates to ^a.
func NewPowerSetExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "**", "(**%s)",
		func(a Value, _ Scope) (Value, error) {
			if s, ok := a.(Set); ok {
				return PowerSet(s), nil
			}
			return nil, errors.Errorf("eval arg must be a Set, not %T", a)
		},
	)
}

// NewNotExpr evaluates to !a.
func NewNotExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "!", "(!%s)",
		func(a Value, _ Scope) (Value, error) { return NewBool(!a.IsTrue()), nil })
}

// NewEvalExpr evaluates to *a, given a set lhs.
func NewEvalExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "*", "(*%s)",
		func(a Value, local Scope) (Value, error) {
			if x, ok := a.(Closure); ok {
				return SetCall(x, None), nil
			}
			return nil, errors.Errorf("eval arg must be a Function, not %T", a)
		})
}

// NewCountExpr evaluates to the number of elements in a.
func NewCountExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "count", "(%s count)",
		func(a Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return NewNumber(float64(x.Count())), nil
			}
			return nil, errors.Errorf("eval arg must be a Set, not %T", a)
		})
}

// NewSingleExpr evaluates to the single element in a or fails if a count != 1.
func NewSingleExpr(scanner parser.Scanner, a Expr) Expr {
	return newUnaryExpr(scanner, a, "single", "(%s single)",
		func(a Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				var result Value
				for e := x.Enumerator(); e.MoveNext(); {
					if result != nil {
						return nil, fmt.Errorf("single: too many elements")
					}
					result = e.Current()
				}
				if result == nil {
					return nil, fmt.Errorf("single: empty set")
				}
				return result, nil
			}
			return nil, errors.Errorf("eval arg must be a Set, not %T", a)
		})
}

// String returns a string representation of the expression.
func (e *UnaryExpr) String() string {
	return fmt.Sprintf(e.format, e.a)
}

// Eval returns the subject
func (e *UnaryExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	val, err := e.eval(a, local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	return val, nil
}
