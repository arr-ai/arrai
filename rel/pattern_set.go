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
		panic(fmt.Sprintf("value %s is not a set", value))
	}

	extraElements := make(map[int]int)
	for i, ptn := range p.patterns {
		if _, is := ptn.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				panic("multiple ... not supported yet")
			}
			extraElements[i] = set.Count() - len(p.patterns)
			continue
		}
		if t, is := ptn.(ExprPattern); is {
			if _, is = t.expr.(IdentExpr); is {
				if len(extraElements) == 1 {
					panic("multiple idents not supported yet")
				}
				extraElements[i] = set.Count() - len(p.patterns)
			}
		}
	}

	if len(p.patterns) > set.Count()+len(extraElements) {
		panic(fmt.Sprintf("length of set %s shorter than set pattern %s", set, p))
	}

	if len(extraElements) == 0 && len(p.patterns) < set.Count() {
		panic(fmt.Sprintf("length of set %s longer than set pattern %s", set, p))
	}

	result := EmptyScope
	for _, ptn := range p.patterns {
		if _, is := ptn.(ExtraElementPattern); is {
			continue
		}
		switch t := ptn.(type) {
		case ExprPattern:
			v, is := t.expr.(Value)
			if is {
				if !set.Has(v) {
					return EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v).(GenericSet)
				continue
			}

			if _, is := t.expr.(IdentExpr); !is {
				return EmptyScope, fmt.Errorf("item type %s is not supported yet", t)
			}
		case IdentPattern:
			v := local.MustGet(t.ident).(Value)
			if !set.Has(v) {
				return EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
			}
			set = set.Without(v).(GenericSet)
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

		result = result.MatchedUpdate(scope)
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
