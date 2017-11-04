package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

// UnnestExpr returns the relation with names grouped into a nested relation.
type UnnestExpr struct {
	lhs  Expr
	attr string
}

// NewUnnestExpr returns a new UnnestExpr.
func NewUnnestExpr(lhs Expr, attr string) Expr {
	return &UnnestExpr{lhs, attr}
}

// LHS returns the LHS of the UnnestExpr.
func (e *UnnestExpr) LHS() Expr {
	return e.lhs
}

// AttrToUnnest returns the attr name to unnest.
func (e *UnnestExpr) AttrToUnnest() string {
	return e.attr
}

// String returns a string representation of the expression.
func (e *UnnestExpr) String() string {
	return fmt.Sprintf("(%s unnest %s)", e.lhs, e.attr)
}

// Eval returns e.lhs with e.attrs grouped under e.attr.
func (e *UnnestExpr) Eval(local, global *Scope) (Value, error) {
	value, err := e.lhs.Eval(local, global)
	if err != nil {
		return nil, err
	}
	if set, ok := value.(Set); ok {
		return Unnest(set, e.attr), nil
	}
	return nil, errors.Errorf("=> not applicable to %T", value)
}
