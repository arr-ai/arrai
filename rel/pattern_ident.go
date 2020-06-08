package rel

import "fmt"

type IdentPattern struct {
	ident string
}

func NewIdentPattern(ident string) IdentPattern {
	return IdentPattern{ident}
}

func (p IdentPattern) Bind(scope Scope, value Value) (Scope, error) {
	if _, has := scope.Get(p.ident); !has {
		return EmptyScope, fmt.Errorf("%q not in scope", p.ident)
	}
	if _, err := scope.MatchedWith(p.ident, value); err != nil {
		return Scope{}, err
	}
	return EmptyScope.With(p.ident, value), nil
}

func (p IdentPattern) String() string {
	return p.ident
}
