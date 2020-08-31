package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

type ImportExpr struct {
	rel.ExprScanner
	packageExpr rel.Expr
	path        string
}

func NewImportExpr(scanner parser.Scanner, imported rel.Expr, path string) ImportExpr {
	return ImportExpr{
		ExprScanner: rel.ExprScanner{Src: scanner},
		packageExpr: NewPackageExpr(scanner, imported),
		path:        path,
	}
}

func (i ImportExpr) Eval(ctx context.Context, _ rel.Scope) (rel.Value, error) {
	//TODO: evaluate accessed imports to avoid re-evaluation
	return i.packageExpr.Eval(ctx, rel.EmptyScope)
}

func (i ImportExpr) String() string {
	return fmt.Sprintf("//{%s}", i.path)
}
