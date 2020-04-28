package rel

import (
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

// LHS returns the LHS of the IfElseExpr.
func (e *IfElseExpr) LHS() Expr {
	return e.ifTrue
}

// Cond returns the condition of the IfElseExpr.
func (e *IfElseExpr) Cond() Expr {
	return e.cond
}

// RHS returns the RHS of the IfElseExpr.
func (e *IfElseExpr) RHS() Expr {
	return e.ifFalse
}

// String returns a string representation of the expression.
func (e *IfElseExpr) String() string {
	return fmt.Sprintf("(%s if %s else %s)", e.ifTrue, e.cond, e.ifFalse)
}

// Eval returns the ifTrue
func (e *IfElseExpr) Eval(local Scope) (Value, error) {
	cond, err := e.cond.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	if cond.IsTrue() {
		return e.ifTrue.Eval(local)
	}
	return e.ifFalse.Eval(local)
}
