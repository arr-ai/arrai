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
func NewMapExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
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
func (e *DArrowExpr) Eval(local Scope) (Value, error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			v, err := e.fn.body.Eval(local.With(e.fn.arg, i.Current()))
			if err != nil {
				return nil, wrapContext(err, e)
			}
			values = append(values, v)
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf("=> not applicable to %T: %[1]v", value), e)
}
