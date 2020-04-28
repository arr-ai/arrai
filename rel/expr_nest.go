package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// NestExpr returns the relation with names grouped into a nested relation.
type NestExpr struct {
	ExprScanner
	lhs   Expr
	attrs Names
	attr  string
}

// NewNestExpr returns a new NestExpr.
func NewNestExpr(scanner parser.Scanner, lhs Expr, attrs Names, attr string) Expr {
	return &NestExpr{ExprScanner{scanner}, lhs, attrs, attr}
}

// LHS returns the LHS of the NestExpr.
func (e *NestExpr) LHS() Expr {
	return e.lhs
}

// AttrsToNest returns the attrs from the LHS to be nested.
func (e *NestExpr) AttrsToNest() Names {
	return e.attrs
}

// NestedAttr returns the attr name for the nested relations.
func (e *NestExpr) NestedAttr() Names {
	return e.attrs
}

// String returns a string representation of the expression.
func (e *NestExpr) String() string {
	return fmt.Sprintf("(%s nest %s %s)", e.lhs, e.attrs, e.attr)
}

// Eval returns e.lhs with e.attrs grouped under e.attr.
func (e *NestExpr) Eval(local Scope) (Value, error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	if set, ok := value.(Set); ok {
		return Nest(set, e.attrs, e.attr), nil
	}
	return nil, wrapContext(errors.Errorf("nest not applicable to %T", value), e)
}
