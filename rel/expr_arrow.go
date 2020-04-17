package rel

import (
	"fmt"
)

// ArrowExpr returns the tuple applied to a function.
type ArrowExpr struct {
	lhs Expr
	fn  *Function
}

// NewArrowExpr returns a new ArrowExpr.
func NewArrowExpr(lhs, fn Expr) Expr {
	return &ArrowExpr{lhs, ExprAsFunction(fn)}
}

// LHS returns the LHS of the ArrowExpr.
func (e *ArrowExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *ArrowExpr) Fn() *Function {
	return e.fn
}

// String returns a string representation of the expression.
func (e *ArrowExpr) String() string {
	if e.fn.Arg() == "." {
		return fmt.Sprintf("(%s -> %s)", e.lhs, e.fn.Body())
	}
	return fmt.Sprintf("(%s -> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *ArrowExpr) Eval(local Scope) (Value, error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, err
	}
	return e.fn.body.Eval(local.With(e.fn.arg, value))
}
