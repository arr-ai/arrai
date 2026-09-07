package rel

import "context"

type ExtraElementPattern struct {
	ident string
}

func NewExtraElementPattern(ident string) ExtraElementPattern {
	return ExtraElementPattern{ident}
}

func (p ExtraElementPattern) Bind(
	ctx context.Context, scope Scope, value Value, b *scopeBuilder,
) (context.Context, error) {
	if p.ident == "" {
		return ctx, nil
	}
	return ctx, b.add(p.ident, value)
}

func (p ExtraElementPattern) String() string {
	return "..." + p.ident
}

func (p ExtraElementPattern) Bindings() []string {
	return []string{p.ident}
}
