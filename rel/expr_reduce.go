package rel

import (
	"fmt"

	"github.com/pkg/errors"
)

// ReduceExpr represents a range of operators.
type ReduceExpr struct {
	a      Expr
	f      *Function
	format string
	init   func(s Set) (interface{}, error)
	reduce func(acc interface{}, v Value) (interface{}, error)
	output func(acc interface{}) (Value, error)
}

// NewReduceExpr evaluates a reduce f, given a set lhs.
func NewReduceExpr(
	a Expr,
	f *Function,
	format string,
	init func(s Set) (interface{}, error),
	reduce func(acc interface{}, v Value) (interface{}, error),
	output func(acc interface{}) (Value, error),
) Expr {
	return &ReduceExpr{a, f, format, init, reduce, output}
}

// NewSumExpr evaluates to the sum of expr over all elements in a.
func NewSumExpr(a, b Expr) Expr {
	return NewReduceExpr(
		a, ExprAsFunction(b), "%s sum ???",
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
func NewMaxExpr(a, b Expr) Expr {
	return NewReduceExpr(
		a, ExprAsFunction(b), "%s max ???",
		func(s Set) (interface{}, error) {
			if s.Bool() {
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
func NewMeanExpr(a, b Expr) Expr {
	type Agg struct {
		sum float64
		n   int
	}
	return NewReduceExpr(
		a, ExprAsFunction(b), "%s mean ???",
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
func NewMinExpr(a, b Expr) Expr {
	return NewReduceExpr(
		a, ExprAsFunction(b), "%s min ???",
		func(s Set) (interface{}, error) {
			if s.Bool() {
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
func (e *ReduceExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, err
	}
	if s, ok := a.(Set); ok {
		acc, err := e.init(s)
		if err != nil {
			return nil, err
		}
		for i := s.Enumerator(); i.MoveNext(); {
			b, err := e.f.Call(i.Current(), local)
			if err != nil {
				return nil, err
			}
			acc, err = e.reduce(acc, b)
			if err != nil {
				return nil, err
			}
		}
		return e.output(acc)
	}
	return nil, errors.Errorf("Not a set")
}
