package rel

import (
	"fmt"
)

// TupleMapExpr returns the tuple applied to a function.
type TupleMapExpr struct {
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewTupleMapExpr(lhs Expr, fn Expr) Expr {
	return &TupleMapExpr{lhs, ExprAsFunction(fn)}
}

// LHS returns the LHS of the AtArrowExpr.
func (e *TupleMapExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *TupleMapExpr) Fn() *Function {
	return e.fn
}

// String returns a string representation of the expression.
func (e *TupleMapExpr) String() string {
	return fmt.Sprintf("(%s :> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *TupleMapExpr) Eval(local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				panic(r)
			}
		}
	}()
	return value.(Tuple).Map(func(v Value) Value {
		v, err = e.fn.body.Eval(local.With(e.fn.arg, v))
		if err != nil {
			panic(err)
		}
		return v
	}), nil
}
