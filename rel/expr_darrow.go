package rel

import (
	"context"
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

// String returns a string representation of the expression.
func (e *DArrowExpr) String() string {
	return fmt.Sprintf("(%s => %s)", e.lhs, e.fn)
}

// Eval returns the lhs transformed elementwise by fn.
func (e *DArrowExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContext(err, e, local)
	}
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			scope, err := e.fn.arg.Bind(ctx, local, i.Current())
			if err != nil {
				return nil, WrapContext(err, e, local)
			}
			v, err := e.fn.body.Eval(ctx, local.Update(scope))
			if err != nil {
				return nil, WrapContext(err, e, local)
			}
			values = append(values, v)
		}
		return NewSet(values...), nil
	}
	return nil, WrapContext(errors.Errorf("=> not applicable to %T: %[1]v", value), e, local)
}
