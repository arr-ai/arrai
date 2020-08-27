package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// IfElseExpr returns the tuple applied to a function.
// Deprecated: IfElseExpr will be removed, it should use CondExpr instead.
type IfElseExpr struct {
	ExprScanner
	ifTrue, cond, ifFalse Expr
}

// NewIfElseExpr returns a new IfElseExpr.
// Deprecated: NewIfElseExpr will be removed, it should use NewCondExpr instead.
func NewIfElseExpr(scanner parser.Scanner, ifTrue, cond, ifFalse Expr) Expr {
	return &IfElseExpr{ExprScanner{scanner}, ifTrue, cond, ifFalse}
}

// String returns a string representation of the expression.
func (e *IfElseExpr) String() string {
	return fmt.Sprintf("(%s if %s else %s)", e.ifTrue, e.cond, e.ifFalse)
}

// Eval returns the ifTrue
func (e *IfElseExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	cond, err := e.cond.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if cond.IsTrue() {
		return e.ifTrue.Eval(ctx, local)
	}
	return e.ifFalse.Eval(ctx, local)
}
