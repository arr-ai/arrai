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

// Pattern returns the PatternExprPair's Pattern.
func (p PatternExprPair) Pattern() Pattern {
	return p.pattern
}

// Expr returns the PatternExprPair's Expr.
func (p PatternExprPair) Expr() Expr {
	return p.expr
}

// String returns a string representation of a PatternPair.
func (p PatternExprPair) String() string {
	return fmt.Sprintf("%s: %s", p.pattern, p.expr)
}

// Bind implements Pattern.Bind.
func (p PatternExprPair) Bind(
	ctx context.Context,
	local Scope,
	value Value,
) (context.Context, Scope, error) {
	return p.pattern.Bind(ctx, local, value)
}

func (p PatternExprPair) eval(ctx context.Context, local Scope) (Value, error) {
	return p.expr.Eval(ctx, local)
}
