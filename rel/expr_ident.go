package rel

import (
	"fmt"
)

// IdentLookupFailed represents failure to lookup a scope variable.
type IdentLookupFailed struct {
	scope *Scope
	expr  *IdentExpr
}

func (e *IdentLookupFailed) Error() string {
	return fmt.Sprintf("Name not found: %q", e.expr.ident)
}

// IdentExpr returns the variable referenced by ident.
type IdentExpr struct {
	ident string
}

var DotIdent = &IdentExpr{"."}

// NewIdentExpr returns a new identifier.
func NewIdentExpr(ident string) *IdentExpr {
	if ident == "." {
		return DotIdent
	}
	return &IdentExpr{ident}
}

// Ident returns the ident for the IdentExpr.
func (e *IdentExpr) Ident() string {
	return e.ident
}

// String returns a string representation of the expression.
func (e *IdentExpr) String() string {
	if e.ident == "." {
		return "(" + e.ident + ")"
	}
	return e.ident
}

// Eval returns the value from scope corresponding to the ident.
func (e *IdentExpr) Eval(local, global *Scope) (Value, error) {
	if a, found := local.Get(e.ident); found && a != nil {
		return a.Eval(global, global)
	}
	return nil, &IdentLookupFailed{local, e}
}
