package rel

import (
	"github.com/arr-ai/wbnf/parser"
)

// ExprExpr represents an expression that yields a Literal.
type ExprExpr struct {
	ExprScanner
	Expr
}

// NewExprExpr returns a new ExprExpr from pairs.
func NewExprExpr(scanner parser.Scanner, expr Expr) ExprExpr {
	return ExprExpr{ExprScanner: ExprScanner{Src: scanner}, Expr: expr}
}

func (e ExprExpr) Source() parser.Scanner {
	return e.ExprScanner.Source()
}
