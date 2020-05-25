package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// ArrowExpr returns the tuple applied to a function.
type ArrowExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewArrowExpr returns a new ArrowExpr.
func NewArrowExpr(scanner parser.Scanner, lhs, fn Expr) Expr {
	return &ArrowExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
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
func (e *ArrowExpr) Eval(local Scope) (_ Value, err error) {
	defer wrapPanic(e, &err)
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	return e.fn.body.Eval(local.Update(e.fn.arg.Bind(local, value)))
}
