package rel

import (
	"fmt"
)

// MemberExpr returns the set applied elementwise to a function.
type MemberExpr struct {
	lhs Expr
	rhs Expr
}

// NewMemberExpr returns a new MemberExpr.
func NewMemberExpr(lhs, rhs Expr) Expr {
	return &MemberExpr{lhs, rhs}
}

// LHS returns the LHS of the MemberExpr.
func (e *MemberExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *MemberExpr) RHS() Expr {
	return e.rhs
}

// String returns a string representation of the expression.
func (e *MemberExpr) String() string {
	return fmt.Sprintf("(%s <: %s)", e.lhs, e.rhs)
}

// Eval returns the lhs transformed elementwise by fn.
func (e *MemberExpr) Eval(local Scope) (Value, error) {
	lhs, err := e.lhs.Eval(local)
	if err != nil {
		return nil, err
	}
	rhs, err := e.rhs.Eval(local)
	if err != nil {
		return nil, err
	}
	if rhs.(Set).Has(lhs) {
		return True, nil
	}
	return False, nil
}
