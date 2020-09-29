package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// UnnestExpr returns the relation with names grouped into a nested relation.
type UnnestExpr struct {
	ExprScanner
	lhs  Expr
	attr string
}

// NewUnnestExpr returns a new UnnestExpr.
func NewUnnestExpr(scanner parser.Scanner, lhs Expr, attr string) Expr {
	return &UnnestExpr{ExprScanner{scanner}, lhs, attr}
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
func (e *UnnestExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if set, ok := value.(Set); ok {
		return Unnest(set, e.attr)
	}
	return nil, WrapContextErr(errors.Errorf("unnest lhs must be relation, not %s", ValueTypeAsString(value)), e, local)
}
