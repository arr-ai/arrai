package rel

import (
	"context"

	"github.com/arr-ai/wbnf/parser"
)

// LiteralExpr represents an expression that yields a Literal.
type LiteralExpr struct {
	ExprScanner
	literal Value
}

// NewLiteralExpr returns a new LiteralExpr from pairs.
func NewLiteralExpr(scanner parser.Scanner, literal Value) LiteralExpr {
	return LiteralExpr{ExprScanner: ExprScanner{Src: scanner}, literal: literal}
}

func (e LiteralExpr) Literal() Value {
	return e.literal
}

// String returns a string representation of the expression.
func (e LiteralExpr) String() string {
	return e.literal.String()
}

// Eval returns the subject
func (e LiteralExpr) Eval(ctx context.Context, _ Scope) (Value, error) {
	return e.literal, nil
}
