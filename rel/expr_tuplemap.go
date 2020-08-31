package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// TupleMapExpr returns the tuple applied to a function.
type TupleMapExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewTupleMapExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &TupleMapExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// String returns a string representation of the expression.
func (e *TupleMapExpr) String() string {
	return fmt.Sprintf("(%s :> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *TupleMapExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return value.(Tuple).Map(func(v Value) (Value, error) {
		ctx, scope, err := e.fn.arg.Bind(ctx, local, v)
		if err != nil {
			return nil, err
		}
		ans, err := e.fn.body.Eval(ctx, local.Update(scope))
		if err != nil {
			return nil, err
		}
		return ans, nil
	})
}
