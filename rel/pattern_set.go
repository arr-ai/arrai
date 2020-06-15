package rel

import (
	"bytes"
	"fmt"
)

type SetPattern struct {
	patterns []Pattern
}

func NewSetPattern(patterns ...Pattern) SetPattern {
	m := make(map[string]struct{})
	for _, v := range patterns {
		if _, exists := m[v.String()]; exists {
			panic(fmt.Sprintf("item %s is duplicated", v))
		}

		m[v.String()] = struct{}{}
	}
	return SetPattern{patterns}
}

func (p SetPattern) Bind(local Scope, value Value) (Scope, error) {
	set, is := value.(GenericSet)
	if !is {
		return EmptyScope, fmt.Errorf("value %s is not a set", value)
	}

	extraElements := make(map[int]int)
	for i, ptn := range p.patterns {
		if _, is := ptn.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = set.Count() - len(p.patterns)
			continue
		}
		if t, is := ptn.(ExprPattern); is {
			if _, is = t.Expr.(IdentExpr); is {
				if len(extraElements) == 1 {
					return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
				}
				extraElements[i] = set.Count() - len(p.patterns)
			}
		}
	}

	if len(p.patterns) > set.Count()+len(extraElements) {
		return EmptyScope, fmt.Errorf("length of set %s shorter than set pattern %s", set, p)
	}

	if len(extraElements) == 0 && len(p.patterns) < set.Count() {
		return EmptyScope, fmt.Errorf("length of set %s longer than set pattern %s", set, p)
	}

	result := EmptyScope
	for _, ptn := range p.patterns {
		if _, is := ptn.(ExtraElementPattern); is {
			continue
		}
		switch t := ptn.(type) {
		case ExprPattern:
			v, is := t.Expr.(Value)
			if is {
				if !set.Has(v) {
					return EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v).(GenericSet)
				continue
			}

			if _, is := t.Expr.(IdentExpr); !is {
				return EmptyScope, fmt.Errorf("item type %s is not supported yet", t)
			}
		case IdentPattern:
			v, has := local.Get(t.ident)
			if !has {
				return EmptyScope, fmt.Errorf("%q not in scope", t.ident)
			}
			if !set.Has(v.(Value)) {
				return EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
			}
			set = set.Without(v.(Value)).(GenericSet)
		default:
			return EmptyScope, fmt.Errorf("%s not supported yet", t)
		}
	}
	for i := range extraElements {
		var scope Scope
		var err error
		if _, is := p.patterns[i].(ExtraElementPattern); is {
			scope, err = p.patterns[i].Bind(local, set)
		} else {
			if set.Count() != 1 {
				return EmptyScope, fmt.Errorf("the length of set %s is wrong", set)
			}

			scope, err = p.patterns[i].Bind(local, set.set.Any().(Value))
			if err != nil {
				return EmptyScope, err
			}
		}
		if err != nil {
			return EmptyScope, err
		}

		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return EmptyScope, err
		}
	}

	return result, nil
}

func (p SetPattern) String() string {
	elts := p.patterns
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, value := range elts {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(value.String())
	}
	buf.WriteString("}")
	return buf.String()
}

func (p SetPattern) Bindings() []string {
	bindings := make([]string, len(p.patterns))
	for i, v := range p.patterns {
		bindings[i] = v.String()
	}
	return bindings
}
