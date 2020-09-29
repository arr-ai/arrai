package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SingleNestExpr returns the relation with names grouped into a nested relation.
type SingleNestExpr struct {
	ExprScanner
	lhs  Expr
	attr string
}

// NewSingleNestExpr returns a new SingleNestExpr.
func NewSingleNestExpr(scanner parser.Scanner, lhs Expr, attr string) Expr {
	return &SingleNestExpr{ExprScanner{scanner}, lhs, attr}
}

// String returns a string representation of the expression.
func (e *SingleNestExpr) String() string {
	return fmt.Sprintf("(%s nest %s)", e.lhs, e.attr)
}

// Eval returns e.lhs with e.attrs grouped under e.attr.
func (e *SingleNestExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if set, ok := value.(Set); ok {
		relAttrs, err := RelationAttrs(set)
		if err != nil {
			return nil, err
		}
		return SingleAttrNest(set, relAttrs, e.attr), nil
	}
	return nil, WrapContextErr(errors.Errorf("nest lhs must be relation, not %s", ValueTypeAsString(value)), e, local)
}
