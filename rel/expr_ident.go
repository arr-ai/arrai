package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// IdentExpr returns the variable referenced by ident.
type IdentExpr struct {
	ExprScanner
	ident string
}

// DotIdent represents the special identifier '.'.
var DotIdent = IdentExpr{ExprScanner: ExprScanner{Src: *parser.NewScanner(".")}, ident: "."}

// NewIdentExpr returns a new identifier.
func NewIdentExpr(scanner parser.Scanner, ident string) IdentExpr {
	if ident == "." {
		return DotIdent
	}
	return IdentExpr{ExprScanner{scanner}, ident}
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
func (e IdentExpr) Eval(local Scope) (Value, error) {
	if a, found := local.Get(e.ident); found && a != nil {
		return a.Eval(local)
	}
	return nil, wrapContext(fmt.Errorf("name %q not found in %v", e.ident, local.m.Keys()), e, local)
}

func (e IdentExpr) Bind(scope Scope, value Value) (Scope, error) {
	return EmptyScope.With(e.ident, value), nil
}
