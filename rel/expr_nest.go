package rel

import (
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
func (e *NestExpr) Eval(local Scope) (Value, error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, WrapContext(err, e, local)
	}
	if set, ok := value.(Set); ok {
		if e.inverse {
			setRelAttrs := mustGetRelationAttrs(set)
			if err := validNestOp(setRelAttrs, e.attrs); err != nil {
				return nil, WrapContext(err, e, local)
			}
			e.attrs = setRelAttrs.Minus(e.attrs)
			if !e.attrs.IsTrue() {
				return nil, WrapContext(
					fmt.Errorf("nest attrs cannot be on all of relation attrs (%v)", setRelAttrs),
					e, local,
				)
			}
		}
		return Nest(set, e.attrs, e.attr), nil
	}
	return nil, WrapContext(errors.Errorf("nest not applicable to %T", value), e, local)
}
