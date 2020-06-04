package rel

type ExtraElementPattern struct {
	ident string
}

func NewExtraElementPattern(ident string) ExtraElementPattern {
	return ExtraElementPattern{ident}
}

func (p ExtraElementPattern) Bind(scope Scope, value Value) (Scope, error) {
	if p.ident == "" {
		return EmptyScope, nil
	}
	return EmptyScope.With(p.ident, value), nil
}

func (p ExtraElementPattern) String() string {
	return "..." + p.ident
}
