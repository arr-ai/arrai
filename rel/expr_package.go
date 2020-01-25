package rel

import (
	"fmt"
)

// PackageExpr represents a range of operators.
type PackageExpr struct {
	a Expr
}

// NewPackageExpr evaluates to !a.
func NewPackageExpr(a Expr) Expr {
	return PackageExpr{a: a}
}

// Arg returns the PackageExpr's arg.
func (e PackageExpr) Arg() Expr {
	return e.a
}

// String returns a string representation of the expression.
func (e PackageExpr) String() string {
	return fmt.Sprintf("(//%s)", e.a)
}

// Eval returns the subject
func (e PackageExpr) Eval(_, _ *Scope) (Value, error) {
	return e.a.Eval(stdScope, EmptyScope)
}
