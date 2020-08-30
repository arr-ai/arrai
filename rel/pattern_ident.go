package rel

import (
	"context"
	"strings"
)

type DynIdent string

type IdentPattern string

func NewIdentPattern(name string) Pattern {
	if strings.HasPrefix(name, "${") {
		return DynIdentPattern(name)
	}
	return IdentPattern(name)
}

func (p IdentPattern) Ident() string {
	return string(p)
}

func (p IdentPattern) Bind(ctx context.Context, scope Scope, value Value) (context.Context, Scope, error) {
	return ctx, Scope{}.With(string(p), value), nil
}

func (p IdentPattern) String() string {
	return string(p)
}

func (p IdentPattern) Bindings() []string {
	return []string{string(p)}
}

type DynIdentPattern string

func (p DynIdentPattern) Bind(ctx context.Context, scope Scope, value Value) (context.Context, Scope, error) {
	return context.WithValue(ctx, DynIdent(p), value), Scope{}, nil
}

func (p DynIdentPattern) String() string {
	return string(p)
}

func (p DynIdentPattern) Bindings() []string {
	return []string{string(p)}
}
