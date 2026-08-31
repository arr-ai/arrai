package rel

import (
	"context"
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

func (p FallbackPattern) Bind(ctx context.Context, local Scope, value Value, b *scopeBuilder) (context.Context, error) {
	if value != nil {
		return p.pattern.Bind(ctx, local, value, b)
	}

	if p.fallback == nil {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("no value and no fallback")
	}

	var err error
	value, err = p.fallback.Eval(ctx, local)
	if err != nil {
		return ctx, err
	}
	return p.pattern.Bind(ctx, EmptyScope, value, b)
}

func (p FallbackPattern) String() string {
	if p.fallback == nil {
		return p.pattern.String()
	}
	return fmt.Sprintf("%s:%s", p.pattern, p.fallback)
}
