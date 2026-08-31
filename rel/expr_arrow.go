package rel

import (
	"context"
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

// String returns a string representation of the expression.
func (e *ArrowExpr) String() string {
	if e.fn.Arg() == "." {
		return fmt.Sprintf("(%s -> %s)", e.lhs, e.fn.Body())
	}
	return fmt.Sprintf("(%s -> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *ArrowExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	// Fast path for `let x = ...; body` and `lhs -> \x body`: see Closure.call.
	if ident, is := e.fn.arg.(IdentPattern); is {
		return e.fn.body.Eval(ctx, local.With(string(ident), value))
	}
	var b scopeBuilder
	ctx, err = e.fn.arg.Bind(ctx, local, value, &b)
	if err != nil {
		if err == errPatternMismatch {
			err = explainBind(ctx, e.fn.arg, local, value)
		}
		return nil, WrapContextErr(err, e, local)
	}
	return e.fn.body.Eval(ctx, local.updateWith(&b))
}
