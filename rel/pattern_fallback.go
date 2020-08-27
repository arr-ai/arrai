package rel

import (
	"context"
	"errors"
	"fmt"
)

type FallbackPattern struct {
	pattern  Pattern
	fallback Expr
}

func NewFallbackPattern(pattern Pattern, fallback Expr) FallbackPattern {
	return FallbackPattern{
		pattern:  pattern,
		fallback: fallback,
	}
}

func (p FallbackPattern) Bind(ctx context.Context, local Scope, value Value) (Scope, error) {
	if value != nil {
		return p.pattern.Bind(ctx, EmptyScope, value)
	}

	if p.fallback == nil {
		return EmptyScope, errors.New("no value and no fallback")
	}

	var err error
	value, err = p.fallback.Eval(ctx, local)
	if err != nil {
		return EmptyScope, err
	}
	return p.pattern.Bind(ctx, EmptyScope, value)
}

func (p FallbackPattern) String() string {
	if p.fallback == nil {
		return p.pattern.String()
	}
	return fmt.Sprintf("%s:%s", p.pattern, p.fallback)
}
