package rel

import "fmt"

type PatternExpr struct {
	pattern Pattern
	expr    Expr
}

// NewPatternExpr returns a new PatternExpr.
func NewPatternExpr(pattern Pattern, expr Expr) PatternExpr {
	return PatternExpr{pattern, expr}
}

// String returns a string representation of a PatternPair.
func (pt PatternExpr) String() string {
	return fmt.Sprintf("%s:%s", pt.pattern, pt.expr)
}

func (pt PatternExpr) Bind(local Scope, value Value) (Scope, error) {
	return pt.pattern.Bind(local, value)
}

func (pt PatternExpr) Eval(local Scope) (Value, error) {
	return pt.expr.Eval(local)
}
