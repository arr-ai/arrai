package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// DArrowExpr returns the set applied elementwise to a function.
type DArrowExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewDArrowExpr returns a new DArrowExpr.
func NewDArrowExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &DArrowExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// LHS returns the LHS of the DArrowExpr.
func (e *DArrowExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *DArrowExpr) Fn() *Function {
	return e.fn
}

// String returns a string representation of the expression.
func (e *DArrowExpr) String() string {
	return fmt.Sprintf("(%s => %s)", e.lhs, e.fn)
}

// Eval returns the lhs transformed elementwise by fn.
func (e *DArrowExpr) Eval(local Scope) (_ Value, err error) {
	defer wrapPanic(e, &err, local)
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			scope, err := e.fn.arg.Bind(local, i.Current())
			if err != nil {
				return nil, wrapContext(err, e, local)
			}
			v, err := e.fn.body.Eval(local.Update(scope))
			if err != nil {
				return nil, wrapContext(err, e, local)
			}
			values = append(values, v)
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf("=> not applicable to %T: %[1]v", value), e, local)
}
