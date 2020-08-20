package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// IdentExpr returns the variable referenced by ident.
type IdentExpr struct {
	ExprScanner
	ident string
}

// NewIdentExpr returns a new identifier.
func NewIdentExpr(scanner parser.Scanner, ident string) IdentExpr {
	return IdentExpr{ExprScanner{scanner}, ident}
}

func NewDotIdent(source parser.Scanner) IdentExpr {
	return NewIdentExpr(source, ".")
}

// Ident returns the ident for the IdentExpr.
func (e IdentExpr) Ident() string {
	return e.ident
}

// String returns a string representation of the expression.
func (e IdentExpr) String() string {
	if e.ident == "." {
		return "(" + e.ident + ")"
	}
	return e.ident
}

// Eval returns the value from scope corresponding to the ident.
func (e IdentExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	if a, found := local.Get(e.ident); found && a != nil {
		return a.Eval(ctx, local)
	}
	return nil, WrapContext(fmt.Errorf("name %q not found in %v", e.ident, local.m.Keys()), e, local)
}
