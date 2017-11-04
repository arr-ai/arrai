package rel

import (
	"fmt"
	"math"

	"github.com/go-errors/errors"
)

type binEval func(a, b Value, local, global *Scope) (Value, error)

// BinExpr represents a range of operators.
type BinExpr struct {
	a, b   Expr
	op     string
	format string
	eval   binEval
}

func newBinExpr(a, b Expr, op, format string, eval binEval) Expr {
	return &BinExpr{a, b, op, format, eval}
}

type valueEval func(a, b Value) Value

// MakeBinValExpr returns a function that creates a binExpr for the given
// logical operator.
func MakeBinValExpr(op string, eval valueEval) func(a, b Expr) Expr {
	return func(a, b Expr) Expr {
		return newBinExpr(a, b, op, "(%s "+op+" %s)",
			func(a, b Value, _, _ *Scope) (Value, error) {
				return eval(a, b), nil
			})
	}
}

type eqEval func(a, b Value) bool

// MakeEqExpr returns a function that creates a binExpr for the given operator.
func MakeEqExpr(op string, eval eqEval) func(a, b Expr) Expr {
	return func(a, b Expr) Expr {
		return newBinExpr(a, b, op, "(%s "+op+" %s)",
			func(a, b Value, _, _ *Scope) (Value, error) {
				if eval(a, b) {
					return True, nil
				}
				return False, nil
			})
	}
}

type arithEval func(a, b float64) float64

func newArithExpr(a, b Expr, op string, eval arithEval) Expr {
	return newBinExpr(a, b, op, "(%s "+op+" %s)",
		func(a, b Value, _, _ *Scope) (Value, error) {
			if a, ok := a.(*Number); ok {
				if b, ok := b.(*Number); ok {
					return NewNumber(eval(a.number, b.number)), nil
				}
			}
			return nil, errors.Errorf(
				"Both args to %q must be numbers, not %T and %T", op, a, b)
		})
}

func addValues(a, b Value) (Value, error) {
	if a, ok := a.(*Number); ok {
		if b, ok := b.(*Number); ok {
			return NewNumber(a.number + b.number), nil
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
		"Both args to %q must be numbers or tuples, not %T and %T",
		"+", a, b)
}

// NewAddExpr evaluates a + b, given two Numbers.
func NewAddExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "+", "(%s + %s)",
		func(a, b Value, _, _ *Scope) (Value, error) {
			return addValues(a, b)
		})
}

// NewSubExpr evaluates a - b, given two Numbers.
func NewSubExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "-", func(a, b float64) float64 { return a - b })
}

// NewMulExpr evaluates a * b, given two Numbers.
func NewMulExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "*", func(a, b float64) float64 { return a * b })
}

// NewDivExpr evaluates a / b, given two Numbers.
func NewDivExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "/", func(a, b float64) float64 { return a / b })
}

// NewIdivExpr evaluates ⎣a / b⎦, given two Numbers.
func NewIdivExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "/", func(a, b float64) float64 {
		return math.Floor(a / b)
	})
}

// NewModExpr evaluates a % b, given two Numbers.
func NewModExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "%", func(a, b float64) float64 {
		return math.Mod(a, b)
	})
}

// NewSubModExpr evaluates a % b, given two Numbers.
func NewSubModExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "-%", func(a, b float64) float64 {
		return a - math.Mod(a, b)
	})
}

// NewPowExpr evaluates a to the power of b, given two Numbers.
func NewPowExpr(a, b Expr) Expr {
	return newArithExpr(a, b, "**", func(a, b float64) float64 {
		return math.Pow(a, b)
	})
}

// NewWithExpr evaluates a with b, given a set lhs.
func NewWithExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "with", "(%s with %s)",
		func(a, b Value, _, _ *Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.With(b), nil
			}
			return nil, errors.Errorf("'with' lhs must be a Set, not %T", a)
		})
}

// NewWithoutExpr evaluates a without b, given a set lhs.
func NewWithoutExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "without", "(%s without %s)",
		func(a, b Value, _, _ *Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.Without(b), nil
			}
			return nil, errors.Errorf("'without' lhs must be a Set, not %T", a)
		})
}

// NewWhereExpr evaluates a where pred, given a set lhs.
func NewWhereExpr(a, pred Expr) Expr {
	pred = ExprAsFunction(pred)
	return newBinExpr(a, pred, "where", "(%s where %s)",
		func(a, pred Value, local, global *Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if p, ok := pred.(*Function); ok {
					return x.Where(func(v Value) bool {
						match, err := p.Call(v, local, global)
						if err != nil {
							panic(err)
						}
						return match.Bool()
					}), nil
				}
				return nil, errors.Errorf("'where' rhs must be a Fn, not %T", a)
			}
			return nil, errors.Errorf("'where' lhs must be a Set, not %T", a)
		})
}

// NewOrderExpr evaluates a order key, given a set lhs, returning an array.
func NewOrderExpr(a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(a, key, "order", "(%s order %s)",
		func(a, key Value, local, global *Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if k, ok := key.(*Function); ok {
					values, err := Order(x, func(value Value) (Value, error) {
						return k.Call(value, local, global)
					})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'where' rhs must be a Fn, not %T", a)
			}
			return nil, errors.Errorf("'where' lhs must be a Set, not %T", a)
		})
}

// NewCallExpr evaluates a without b, given a set lhs.
func NewCallExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "call", "(%s %s)",
		func(a, b Value, local, global *Scope) (Value, error) {
			switch x := a.(type) {
			case *Function:
				return x.Call(b, local, global)
			case *NativeFunction:
				return x.Call(b, local, global)
			case Set:
				var match func(at Value) bool
				multi := false
				if y, ok := b.(Set); ok {
					match = func(at Value) bool { return y.Has(at) }
					multi = true
				} else {
					match = func(at Value) bool { return b.Equal(at) }
				}

				outs := None
				for e := x.Enumerator(); e.MoveNext(); {
					if t, ok := e.Current().(Tuple); ok {
						if v, found := t.Get("@"); found && match(v) {
							if t.Count() != 2 {
								return nil, errors.Errorf("Too many outputs")
							}
							rest, _ := t.Without("@")
							for e := rest.Enumerator(); e.MoveNext(); {
								_, value := e.Current()
								outs = outs.With(value)
							}
						}
					}
				}
				if multi {
					return outs, nil
				}
				n := outs.Count()
				if n != 1 {
					if n == 0 {
						return nil, errors.Errorf("No items founds")
					}
					return nil, errors.Errorf("Too many items found")
				}
				if e := outs.Enumerator(); e.MoveNext() {
					return e.Current(), nil
				}
			}
			return nil, errors.Errorf(
				"call lhs must be a Function, not %T", a)
		})
}

// LHS returns the left hand side of the BinExpr.
func (e *BinExpr) LHS() Expr {
	return e.a
}

// RHS returns the right hand side of the BinExpr.
func (e *BinExpr) RHS() Expr {
	return e.b
}

// String returns a string representation of the expression.
func (e *BinExpr) String() string {
	return fmt.Sprintf(e.format, e.a, e.b)
}

// Eval returns the subject
func (e *BinExpr) Eval(local, global *Scope) (Value, error) {
	a, err := e.a.Eval(local, global)
	if err != nil {
		return nil, err
	}

	b, err := e.b.Eval(local, global)
	if err != nil {
		return nil, err
	}

	return e.eval(a, b, local, global)
}
