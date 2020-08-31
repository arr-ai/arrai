package rel

import "context"

type ExtraElementPattern struct {
	ident string
}

func NewExtraElementPattern(ident string) ExtraElementPattern {
	return ExtraElementPattern{ident}
}

func (p ExtraElementPattern) Bind(ctx context.Context, scope Scope, value Value) (context.Context, Scope, error) {
	if p.ident == "" {
		return ctx, EmptyScope, nil
	}
	return ctx, EmptyScope.With(p.ident, value), nil
}

func (p ExtraElementPattern) String() string {
	return "..." + p.ident
}

func (p ExtraElementPattern) Bindings() []string {
	return []string{p.ident}
}
