package rel

import (
	"context"
	"fmt"
	"math"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type binEval func(ctx context.Context, a, b Value, local Scope) (Value, error)

// BinExpr represents a range of operators.
type BinExpr struct {
	ExprScanner
	a, b   Expr
	op     string
	format string
	eval   binEval
}

func newBinExpr(scanner parser.Scanner, a, b Expr, op, format string, eval binEval) Expr {
	return &BinExpr{ExprScanner{scanner}, a, b, op, format, eval}
}

type valueEval func(a, b Value) Value

// MakeBinValExpr returns a function that creates a binExpr for the given
// logical operator.
func MakeBinValExpr(op string, eval valueEval) func(scanner parser.Scanner, a, b Expr) Expr {
	return func(scanner parser.Scanner, a, b Expr) Expr {
		return newBinExpr(scanner, a, b, op, "(%s "+op+" %s)",
			func(ctx context.Context, a, b Value, _ Scope) (Value, error) {
				return eval(a, b), nil
			})
	}
}

type arithEval func(a, b float64) float64

func newArithExpr(scanner parser.Scanner, a, b Expr, op string, eval arithEval) Expr {
	return newBinExpr(scanner, a, b, op, "(%s "+op+" %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if a, ok := a.(Number); ok {
				if b, ok := b.(Number); ok {
					return NewNumber(eval(a.Float64(), b.Float64())), nil
				}
			}
			return nil, errors.Errorf(
				"Both args to %q must be numbers, not %s and %s",
				op, ValueTypeAsString(a), ValueTypeAsString(b))
		})
}

func addValues(a, b Value) (Value, error) {
	if a, ok := a.(Number); ok {
		if b, ok := b.(Number); ok {
			return NewNumber(a.Float64() + b.Float64()), nil
		}
	}
	if a, ok := a.(Tuple); ok {
		if b, ok := b.(Tuple); ok {
			return MergeLeftToRight(a, b), nil
		}
	}
	if a, ok := a.(Set); ok {
		if b, ok := b.(Set); ok {
			return Concatenate(a, b)
		}
	}
	return nil, errors.Errorf(
		"Both args to + must be numbers or tuples, not %s and %s",
		ValueTypeAsString(a), ValueTypeAsString(b))
}

// NewAddExpr evaluates a + b, given two Numbers.
func NewAddExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "+", "(%s + %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			return addValues(a, b)
		})
}

// NewAddArrowExpr returns a new BinExpr which supports operator `+>`.
func NewAddArrowExpr(scanner parser.Scanner, lhs, rhs Expr) Expr {
	return newBinExpr(scanner, lhs, rhs, "+>", "(%s +> %s)",
		func(_ context.Context, lhs, rhs Value, _ Scope) (Value, error) {
			return evalValForAddArrow(lhs, rhs)
		})
}

// NewSubExpr evaluates a - b, given two Numbers.
func NewSubExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "-", func(a, b float64) float64 { return a - b })
}

// NewMulExpr evaluates a * b, given two Numbers.
func NewMulExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "*", func(a, b float64) float64 { return a * b })
}

// NewDivExpr evaluates a / b, given two Numbers.
func NewDivExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "/", func(a, b float64) float64 { return a / b })
}

// NewIdivExpr evaluates ⎣a / b⎦, given two Numbers.
func NewIdivExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "/", func(a, b float64) float64 {
		return math.Floor(a / b)
	})
}

// NewModExpr evaluates a % b, given two Numbers.
func NewModExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "%%", func(a, b float64) float64 {
		return math.Mod(a, b)
	})
}

// NewSubModExpr evaluates a % b, given two Numbers.
func NewSubModExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "-%", func(a, b float64) float64 {
		return a - math.Mod(a, b)
	})
}

// NewPowExpr evaluates a to the power of b, given two Numbers.
func NewPowExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "^", func(a, b float64) float64 {
		return math.Pow(a, b)
	})
}

// NewWithExpr evaluates a with b, given a set lhs.
func NewWithExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "with", "(%s with %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.With(b), nil
			}
			return nil, errors.Errorf("'with' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewWithoutExpr evaluates a without b, given a set lhs.
func NewWithoutExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "without", "(%s without %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.Without(b), nil
			}
			return nil, errors.Errorf("'without' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewWhereExpr evaluates a where pred, given a set lhs.
func NewWhereExpr(scanner parser.Scanner, a, pred Expr) Expr {
	pred = ExprAsFunction(pred)
	return newBinExpr(scanner, a, pred, "where", "(%s where %s)",
		func(ctx context.Context, a, pred Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if p, ok := pred.(Closure); ok {
					return x.Where(func(v Value) (bool, error) {
						r, err := SetCall(ctx, p, v)
						if err != nil {
							return false, err
						}
						return r.IsTrue(), nil
					})
				}
				return nil, errors.Errorf("'where' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'where' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewOrderByExpr evaluates a orderby key, given a set lhs, returning an array.
func NewOrderByExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "orderby", "(%s orderby %s)",
		func(ctx context.Context, a, key Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if k, ok := key.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return SetCall(ctx, k, value)
						},
						func(a, b Value) bool {
							return a.Less(b)
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'orderby' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'orderby' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewOrderExpr evaluates a order less, given a set lhs, returning an array.
func NewOrderExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "order", "(%s orderby %s)",
		func(ctx context.Context, a, less Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if l, ok := less.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return value, nil
						},
						func(a, b Value) bool {
							c, err := SetCall(ctx, l, a)
							if err != nil {
								panic(err)
							}
							less, err := SetCall(ctx, c.(Closure), b)
							if err != nil {
								panic(err)
							}
							return less.IsTrue()
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'order' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'order' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewRankExpr evaluates a rank tuplef, given a relation lhs, returning a new
// relation with each lhs tuple augmented by the tuplef attrs containing the
// corresponding rank.
func NewRankExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "rank", "(%s rank %s)",
		func(ctx context.Context, a, tuplef Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if l, ok := tuplef.(Closure); ok {
					return Rank(x, func(v Tuple) (Tuple, error) {
						result, err := SetCall(ctx, l, v)
						if err != nil {
							return nil, err
						}
						return result.(Tuple), nil
					})
				}
				return nil, errors.Errorf("'rank' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'rank' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

func Call(ctx context.Context, a, b Value, _ Scope) (Value, error) {
	if x, ok := a.(Set); ok {
		return SetCall(ctx, x, b)
	}
	return nil, errors.Errorf(
		"call lhs must be a function, not %s", ValueTypeAsString(a))
}

// NewCallExpr evaluates a without b, given a set lhs.
func NewCallExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "call", "«%s»(%s)", Call)
}

func NewCallExprCurry(scanner parser.Scanner, f Expr, args ...Expr) Expr {
	for _, arg := range args {
		f = NewCallExpr(scanner, f, arg)
	}
	return f
}

// String returns a string representation of the expression.
func (e *BinExpr) String() string {
	return fmt.Sprintf(e.format, e.a, e.b)
}

// Eval returns the subject
func (e *BinExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	a, err := e.a.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}

	b, err := e.b.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	val, err := e.eval(ctx, a, b, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return val, nil
}

// evalValForAddArrow evaluates operator `+>`.
func evalValForAddArrow(lhs, rhs Value) (Value, error) {
	switch lhs := lhs.(type) {
	case Tuple:
		if rhs, ok := rhs.(Tuple); ok {
			return MergeLeftToRight(lhs, rhs), nil
		}
	case Dict:
		switch rhs := rhs.(type) {
		case Dict:
			return mergeDicts(lhs, rhs), nil
		case Set:
			if !rhs.IsTrue() {
				return lhs, nil
			}
		}
	case Set:
		if !lhs.IsTrue() {
			switch rhs := rhs.(type) {
			case Dict:
				return rhs, nil
			case Set:
				if !rhs.IsTrue() {
					return lhs, nil
				}
			}
		}
	}

	return nil, errors.Errorf(
		"Args to +> must be both tuples or both dicts, not %s and %s",
		ValueTypeAsString(lhs), ValueTypeAsString(rhs))
}

func mergeDicts(lhs Dict, rhs Dict) Dict {
	tempMap := lhs.m

	for e := rhs.DictEnumerator(); e.MoveNext(); {
		key, value := e.Current()
		tempMap = tempMap.With(key, value)
	}

	return Dict{m: tempMap}
}
