package rel

type IdentPattern struct {
	ident string
}

func NewIdentPattern(ident string) IdentPattern {
	return IdentPattern{ident}
}

func (p IdentPattern) Bind(scope Scope, value Value) (Scope, error) {
	scope.MustGet(p.ident)
	scope.MatchedWith(p.ident, value)
	return EmptyScope.With(p.ident, value), nil
}

func (p IdentPattern) String() string {
	return p.ident
}
