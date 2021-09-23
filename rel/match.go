package rel

type Matcher interface {
	Match(Value) bool
}

type Let func(Value)

func (m Let) Match(v Value) bool {
	m(v)
	return true
}

type LiteralMatcher struct {
	value Value
}

func Lit(value Value) LiteralMatcher {
	return LiteralMatcher{value: value}
}

func (m LiteralMatcher) Match(v Value) bool {
	return v.Equal(m.value)
}

type MatchInt func(i int)

func (m MatchInt) Match(v Value) bool {
	if n, ok := v.(Number); ok {
		if i, is := n.Int(); is {
			m(i)
			return true
		}
	}
	return false
}

type TupleMatcher struct {
	attrs map[string]Matcher
	rest  Matcher
}

func NewTupleMatcher(attrs map[string]Matcher, rest Matcher) TupleMatcher {
	return TupleMatcher{attrs: attrs, rest: rest}
}

func (m TupleMatcher) Match(v Value) bool {
	if t, ok := v.(Tuple); ok {
		for name, matcher := range m.attrs {
			if a, has := t.Get(name); has {
				if matcher.Match(a) {
					continue
				}
			}
			return false
		}
		rest := EmptyTuple
		for e := t.Enumerator(); e.MoveNext(); {
			name, attr := e.Current()
			if _, has := m.attrs[name]; !has {
				rest = rest.With(name, attr)
			}
		}
		return m.rest.Match(rest)
	}
	return false
}
