package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// NestExpr returns the relation with names grouped into a nested relation.
type NestExpr struct {
	ExprScanner
	inverse bool
	lhs     Expr
	attrs   Names
	attr    string
}

// NewNestExpr returns a new NestExpr.
func NewNestExpr(scanner parser.Scanner, inverse bool, lhs Expr, attrs Names, attr string) Expr {
	return &NestExpr{ExprScanner{scanner}, inverse, lhs, attrs, attr}
}

// String returns a string representation of the expression.
func (e *NestExpr) String() string {
	return fmt.Sprintf("(%s nest %s %s)", e.lhs, e.attrs, e.attr)
}

// Eval returns e.lhs with e.attrs grouped under e.attr.
func (e *NestExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if set, ok := value.(Set); ok {
		relAttrs, err := RelationAttrs(set)
		if err != nil {
			return nil, err
		}
		attrs := e.attrs
		if e.inverse {
			if err := validNestOp(relAttrs, attrs); err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			attrs = relAttrs.Minus(attrs)
			if !attrs.IsTrue() {
				return nil, WrapContextErr(
					fmt.Errorf("nest attrs cannot be on all of relation attrs (%v)", relAttrs),
					e, local,
				)
			}
		}
		return Nest(set, relAttrs, attrs, e.attr), nil
	}
	return nil, WrapContextErr(errors.Errorf("nest lhs must be relation, not %s", ValueTypeAsString(value)), e, local)
}
