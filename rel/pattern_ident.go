package rel

import (
	"context"
)

type IdentPattern string

func (p IdentPattern) Bind(ctx context.Context, scope Scope, value Value) (context.Context, Scope, error) {
	// Bind value for identexpr in Pattern, like `let (a: x, b: y) = (a: 4, b: 7); x`
	return ctx, Scope{}.With(string(p), value), nil
}

func (p IdentPattern) String() string {
	return string(p)
}

func (p IdentPattern) Bindings() []string {
	return []string{string(p)}
}
