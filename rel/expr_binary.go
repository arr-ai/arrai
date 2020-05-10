package rel

import (
	"fmt"
	"math"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type binEval func(a, b Value, local Scope) (Value, error)

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
			func(a, b Value, _ Scope) (Value, error) {
				return eval(a, b), nil
			})
	}
}

type arithEval func(a, b float64) float64

func newArithExpr(scanner parser.Scanner, a, b Expr, op string, eval arithEval) Expr {
	return newBinExpr(scanner, a, b, op, "(%s "+op+" %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if a, ok := a.(Number); ok {
				if b, ok := b.(Number); ok {
					return NewNumber(eval(a.Float64(), b.Float64())), nil
				}
			}
			return nil, errors.Errorf(
				"Both args to %q must be numbers, not %T and %T", op, a, b)
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
		"Both args to %q must be numbers or tuples, not %T and %T",
		"+", a, b)
}

// NewAddExpr evaluates a + b, given two Numbers.
func NewAddExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "+", "(%s + %s)",
		func(a, b Value, _ Scope) (Value, error) {
			return addValues(a, b)
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
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.With(b), nil
			}
			return nil, errors.Errorf("'with' lhs must be a Set, not %T", a)
		})
}

// NewWithoutExpr evaluates a without b, given a set lhs.
func NewWithoutExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "without", "(%s without %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.Without(b), nil
			}
			return nil, errors.Errorf("'without' lhs must be a Set, not %T", a)
		})
}

// NewWhereExpr evaluates a where pred, given a set lhs.
func NewWhereExpr(scanner parser.Scanner, a, pred Expr) Expr {
	pred = ExprAsFunction(pred)
	return newBinExpr(scanner, a, pred, "where", "(%s where %s)",
		func(a, pred Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if p, ok := pred.(Closure); ok {
					return x.Where(func(v Value) bool {
						match, err := p.Call(v, local)
						if err != nil {
							panic(err)
						}
						return match.IsTrue()
					}), nil
				}
				return nil, errors.Errorf("'where' rhs must be a Fn, not %T", a)
			}
			return nil, errors.Errorf("'where' lhs must be a Set, not %T", a)
		})
}

// NewOrderByExpr evaluates a orderby key, given a set lhs, returning an array.
func NewOrderByExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "order", "(%s order %s)",
		func(a, key Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if k, ok := key.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return k.Call(value, local)
						},
						func(a, b Value) bool {
							return a.Less(b)
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'order' rhs must be a Fn, not %T", a)
			}
			return nil, errors.Errorf("'order' lhs must be a Set, not %T", a)
		})
}

// NewOrderExpr evaluates a order less, given a set lhs, returning an array.
func NewOrderExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "order", "(%s order %s)",
		func(a, less Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if l, ok := less.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return value, nil
						},
						func(a, b Value) bool {
							f, err := l.Call(a, local)
							if err != nil {
								panic(err)
							}
							result, err := f.(Closure).Call(b, local)
							if err != nil {
								panic(err)
							}
							return result.IsTrue()
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'order' rhs must be a Fn, not %T", a)
			}
			return nil, errors.Errorf("'order' lhs must be a Set, not %T", a)
		})
}

type Callable interface {
	Call(Expr, Scope) (Value, error)
}

func Call(a, b Value, local Scope) (Value, error) {
	switch x := a.(type) {
	case Callable:
		return x.Call(b, local)
	case Set:
		var out Value
		for e := x.Enumerator(); e.MoveNext(); {
			if t, ok := e.Current().(Tuple); ok {
				// log.Printf("%v %v %[2]T %v %[3]T", t, t.MustGet("@"), b)
				if v, found := t.Get("@"); found && b.Equal(v) {
					if out != nil {
						return nil, errors.Errorf("Too many items found")
					}
					if t.Count() != 2 {
						return nil, errors.Errorf("Too many outputs")
					}
					rest := t.Without("@")
					for e := rest.Enumerator(); e.MoveNext(); {
						_, value := e.Current()
						out = value
					}
				}
			}
		}
		if out == nil {
			return nil, errors.Errorf("No items found: %v", b)
		}
		return out, nil
	}
	return nil, errors.Errorf(
		"call lhs must be a function, not %T", a)
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
func (e *BinExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}

	b, err := e.b.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}

	return e.eval(a, b, local)
}
