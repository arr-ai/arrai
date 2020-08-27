package rel

import (
	"context"
	"fmt"
)

// PatternExprPair a Pattern/Expr pair
type PatternExprPair struct {
	pattern Pattern
	expr    Expr
}

// NewPatternExprPair returns a new PatternExprPair.
func NewPatternExprPair(pattern Pattern, expr Expr) PatternExprPair {
	return PatternExprPair{pattern, expr}
}

// String returns a string representation of a PatternPair.
func (pt PatternExprPair) String() string {
	return fmt.Sprintf("%s:%s", pt.pattern, pt.expr)
}

func (pt PatternExprPair) Bind(ctx context.Context, local Scope, value Value) (Scope, error) {
	return pt.pattern.Bind(ctx, local, value)
}

func (pt PatternExprPair) Eval(ctx context.Context, local Scope) (Value, error) {
	return pt.expr.Eval(ctx, local)
}
