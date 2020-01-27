package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

// PackageExpr represents a range of operators.
type PackageExpr struct {
	a rel.Expr
}

// NewPackageExpr evaluates to !a.
func NewPackageExpr(a rel.Expr) rel.Expr {
	return PackageExpr{a: a}
}

// Arg returns the PackageExpr's arg.
func (e PackageExpr) Arg() rel.Expr {
	return e.a
}

// String returns a string representation of the expression.
func (e PackageExpr) String() string {
	return fmt.Sprintf("(//%s)", e.a)
}

// Eval returns the subject
func (e PackageExpr) Eval(_, _ rel.Scope) (rel.Value, error) {
	return e.a.Eval(stdScope, rel.EmptyScope)
}
