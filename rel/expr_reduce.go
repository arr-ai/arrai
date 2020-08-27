package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/pkg/errors"
)

// ReduceExpr represents a range of operators.
type ReduceExpr struct {
	ExprScanner
	a      Expr
	f      *Function
	format string
	init   func(s Set) (interface{}, error)
	reduce func(acc interface{}, v Value) (interface{}, error)
	output func(acc interface{}) (Value, error)
}

// NewReduceExpr evaluates a reduce f, given a set lhs.
func NewReduceExpr(scanner parser.Scanner,
	a Expr,
	f *Function,
	format string,
	init func(s Set) (interface{}, error),
	reduce func(acc interface{}, v Value) (interface{}, error),
	output func(acc interface{}) (Value, error),
) Expr {
	return &ReduceExpr{ExprScanner{scanner}, a, f, format, init, reduce, output}
}

// NewSumExpr evaluates to the sum of expr over all elements in a.
func NewSumExpr(scanner parser.Scanner, a, b Expr) Expr {
	return NewReduceExpr(
		scanner, a, ExprAsFunction(b), "%s sum ???",
		func(s Set) (interface{}, error) {
			return 0.0, nil
		},
		func(acc interface{}, v Value) (interface{}, error) {
			switch v := v.(type) {
			case Number:
				return acc.(float64) + v.Float64(), nil
			}
			return nil, errors.Errorf("Non-numeric value used in sum")
		},
		func(acc interface{}) (Value, error) {
			return NewNumber(acc.(float64)), nil
		},
	)
}

// NewMaxExpr evaluates to the max of expr over all elements in a.
func NewMaxExpr(scanner parser.Scanner, a, b Expr) Expr {
	return NewReduceExpr(
		scanner, a, ExprAsFunction(b), "%s max ???",
		func(s Set) (interface{}, error) {
			if s.IsTrue() {
				return nil, nil
			}
			return nil, errors.Errorf("Empty set has no max")
		},
		func(acc interface{}, v Value) (interface{}, error) {
			if acc == nil || acc.(Value).Less(v) {
				return v, nil
			}
			return acc, nil
		},
		func(acc interface{}) (Value, error) {
			if acc != nil {
				return acc.(Value), nil
			}
			return nil, errors.Errorf("Empty input to max")
		},
	)
}

// NewMeanExpr evaluates to the mean of expr over all elements in a.
func NewMeanExpr(scanner parser.Scanner, a, b Expr) Expr {
	type Agg struct {
		sum float64
		n   int
	}
	return NewReduceExpr(
		scanner, a, ExprAsFunction(b), "%s mean ???",
		func(s Set) (interface{}, error) {
			if n := s.Count(); n > 0 {
				return Agg{n: n}, nil
			}
			return nil, errors.Errorf("Empty set has no mean")
		},
		func(acc interface{}, v Value) (interface{}, error) {
			agg := acc.(Agg)
			if v, ok := v.(Number); ok {
				return Agg{sum: agg.sum + v.Float64(), n: agg.n}, nil
			}
			return nil, errors.Errorf("Non-numeric value used in mean")
		},
		func(acc interface{}) (Value, error) {
			agg := acc.(Agg)
			if agg.n != 0 {
				return NewNumber(agg.sum / float64(agg.n)), nil
			}
			return nil, errors.Errorf("Non-numeric value used in mean")
		},
	)
}

// NewMinExpr evaluates to the min of expr over all elements in a.
func NewMinExpr(scanner parser.Scanner, a, b Expr) Expr {
	return NewReduceExpr(
		scanner, a, ExprAsFunction(b), "%s min ???",
		func(s Set) (interface{}, error) {
			if s.IsTrue() {
				return nil, nil
			}
			return nil, errors.Errorf("Empty set has no min")
		},
		func(acc interface{}, v Value) (interface{}, error) {
			if acc == nil || v.Less(acc.(Value)) {
				return v, nil
			}
			return acc, nil
		},
		func(acc interface{}) (Value, error) {
			if acc != nil {
				return acc.(Value), nil
			}
			return nil, errors.Errorf("Empty input to min")
		},
	)
}

// String returns a string representation of the expression.
func (e *ReduceExpr) String() string {
	return fmt.Sprintf(e.format, e.a)
}

// Eval returns the subject
func (e *ReduceExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	a, err := e.a.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if s, ok := a.(Set); ok {
		acc, err := e.init(s)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		for i := s.Enumerator(); i.MoveNext(); {
			f, err := e.f.Eval(ctx, local)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			v, err := SetCall(ctx, f.(Closure), i.Current())
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			acc, err = e.reduce(acc, v)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
		}
		return e.output(acc)
	}
	return nil, errors.Errorf("Not a set")
}
