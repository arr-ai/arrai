package rel

import (
	"errors"
	"fmt"
)

type FallbackPattern struct {
	pattern  Pattern
	fallback Pattern
}

func NewFallbackPattern(pattern, fallback Pattern) FallbackPattern {
	return FallbackPattern{
		pattern:  pattern,
		fallback: fallback,
	}
}

func (p FallbackPattern) Bind(local Scope, value Value) (Scope, error) {
	if value != nil {
		return p.pattern.Bind(EmptyScope, value)
	}

	if p.fallback == nil {
		return EmptyScope, errors.New("no value and no fallback")
	}
	scope, err := p.fallback.Bind(local, value)
	if err != nil {
		return EmptyScope, err
	}
	return p.pattern.Bind(scope, value)
}

func (f FallbackPattern) String() string {
	if f.fallback == nil {
		return f.pattern.String()
	}
	return fmt.Sprintf("%s?:%s", f.pattern, f.fallback)
}
