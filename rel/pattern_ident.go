package rel

type IdentPattern struct {
	ident string
}

func NewIdentPattern(ident string) IdentPattern {
	return IdentPattern{ident}
}

func (p IdentPattern) Bind(scope Scope, value Value) (Scope, error) {
	return EmptyScope.With(p.ident, value), nil
}

func (p IdentPattern) String() string {
	return p.ident
}

func (p IdentPattern) Bindings() []string {
	return []string{p.ident}
}
