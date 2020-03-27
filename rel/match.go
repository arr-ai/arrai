package rel

import "fmt"

type Matcher interface {
	Match(Value) bool
}

type Let func(Value)

func (m Let) Match(v Value) bool {
	m(v)
	return true
}

type BindMatcher struct {
	target interface{}
}

func Bind(target interface{}) BindMatcher {
	return BindMatcher{target: target}
}

func (m BindMatcher) Match(v Value) bool {
	ok := false
	switch target := m.target.(type) {
	case *float64:
		if n, ok := v.(Number); ok {
			*target = n.Float64()
			return true
		}
		return false
	case *int:
		if n, ok := v.(Number); ok {
			f := n.Float64()
			*target = int(f)
			return float64(*target) == f
		}
		return false
	case *string:
		*target = v.String()
		return true
	case *Number:
		*target, ok = v.(Number)
	case *Tuple:
		*target, ok = v.(Tuple)
	case *Set:
		*target, ok = v.(Set)
	case *Value:
		*target = v
		return true
	default:
		panic(fmt.Errorf("%v %[1]T unsupported", target))
	}
	return ok
}

type MatchFirst []Matcher

func (m MatchFirst) Match(v Value) bool {
	for _, o := range m {
		if o.Match(v) {
			return true
		}
	}
	return false
}

type MatchAll []Matcher

func (m MatchAll) Match(v Value) bool {
	for _, o := range m {
		if !o.Match(v) {
			return false
		}
	}
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

type MatchNum func(n Number)

func (m MatchNum) Match(v Value) bool {
	if n, ok := v.(Number); ok {
		m(n)
		return true
	}
	return false
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

type SetMatcher struct {
	elem Matcher
}

func NewSetMatcher(elem Matcher) SetMatcher {
	return SetMatcher{elem: elem}
}

func (m SetMatcher) Match(v Value) bool {
	if s, ok := v.(Set); ok {
		for e := s.Enumerator(); e.MoveNext(); {
			if !m.elem.Match(e.Current()) {
				return false
			}
		}
		return true
	}
	return false
}

type FuncMatcher struct {
	tuples [][2]Matcher
	rest   Matcher
}

func MatchFunc(tuples [][2]Matcher, rest Matcher) FuncMatcher {
	return FuncMatcher{tuples: tuples, rest: rest}
}

func (m FuncMatcher) Match(v Value) bool {
	if s, ok := v.(Set); ok {
		e := s.Enumerator()
		if !e.MoveNext() {
			return false
		}
		el0 := e.Current()
		t, ok := el0.(Tuple)
		if !ok {
			return false
		}
		if t.Count() != 2 {
			return false
		}
		var at, other Value
		tms := map[string]Matcher{}
		for te := t.Enumerator(); te.MoveNext(); {
			if name, _ := te.Current(); name == "@" {
				tms[name] = Let(func(v Value) { at = v })
			} else {
				tms[name] = Let(func(v Value) { other = v })
			}
		}
		tm := NewTupleMatcher(tms, Lit(EmptyTuple))
		tuples := append([][2]Matcher{}, m.tuples...)
		rest := None
	tuple:
		for e = s.Enumerator(); e.MoveNext(); {
			residue := tuples[:0]
			if !tm.Match(e.Current()) {
				return false
			}
			for _, m := range tuples {
				if m[0].Match(at) && m[1].Match(other) {
					continue tuple
				}
				residue = append(residue, m)
			}
			rest = rest.With(e.Current())
			tuples = residue
		}
		return m.rest.Match(rest)
	}
	return false
}
