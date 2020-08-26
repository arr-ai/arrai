package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

// PackageExpr represents a range of operators.
type PackageExpr struct {
	rel.ExprScanner
	a rel.Expr
}

// NewPackageExpr evaluates to !a.
func NewPackageExpr(scanner parser.Scanner, a rel.Expr) rel.Expr {
	return PackageExpr{ExprScanner: rel.ExprScanner{Src: scanner}, a: a}
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
func (e PackageExpr) Eval(ctx context.Context, _ rel.Scope) (rel.Value, error) {
	return e.a.Eval(ctx, StdScope())
}
