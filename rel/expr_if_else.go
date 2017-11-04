package rel

import (
	"fmt"
)

// IfElseExpr returns the tuple applied to a function.
type IfElseExpr struct {
	ifTrue, cond, ifFalse Expr
}

// NewIfElseExpr returns a new IfElseExpr.
func NewIfElseExpr(ifTrue, cond, ifFalse Expr) Expr {
	return &IfElseExpr{ifTrue, cond, ifFalse}
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
func (e *IfElseExpr) Eval(local, global *Scope) (Value, error) {
	cond, err := e.cond.Eval(local, global)
	if err != nil {
		return nil, err
	}
	if cond.Bool() {
		return e.ifTrue.Eval(local, global)
	}
	return e.ifFalse.Eval(local, global)
}
