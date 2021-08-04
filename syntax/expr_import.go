package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/rel"
)

type ImportExpr struct {
	rel.ExprScanner
	importedExpr rel.Expr
	path         string
}

func NewImportExpr(scanner parser.Scanner, imported rel.Expr, path string) ImportExpr {
	return ImportExpr{
		ExprScanner:  rel.ExprScanner{Src: scanner},
		importedExpr: imported,
		path:         path,
	}
}

func (i ImportExpr) Eval(ctx context.Context, _ rel.Scope) (rel.Value, error) {
	//TODO: evaluate accessed imports to avoid re-evaluation
	return i.importedExpr.Eval(ctx, rel.EmptyScope)
}

func (i ImportExpr) String() string {
	return fmt.Sprintf("//{%s}", i.path)
}
