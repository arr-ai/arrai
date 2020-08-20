package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// NestExpr returns the relation with names grouped into a nested relation.
type SingleNestExpr struct {
	ExprScanner
	lhs  Expr
	attr string
}

// NewNestExpr returns a new NestExpr.
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
		return nil, WrapContext(err, e, local)
	}
	if set, ok := value.(Set); ok {
		return SingleAttrNest(set, e.attr), nil
	}
	return nil, WrapContext(errors.Errorf("nest not applicable to %T", value), e, local)
}
