package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/frozen"
)

type DictPatternEntry struct {
	at       Expr
	value    Pattern
	fallback Expr
}

func NewDictPatternEntry(at Expr, value Pattern, fallback Expr) DictPatternEntry {
	return DictPatternEntry{
		at:       at,
		value:    value,
		fallback: fallback,
	}
}

func (p DictPatternEntry) String() string {
	if p.fallback == nil {
		return fmt.Sprintf("%s: %s", p.at, p.value)
	}
	return fmt.Sprintf("%s?: %s:%s", p.at, p.value, p.fallback)
}

type DictPattern struct {
	entries []DictPatternEntry
}

func NewDictPattern(entries ...DictPatternEntry) DictPattern {
	names := make(map[string]bool)
	for _, entry := range entries {
		if entry.at != nil && names[entry.at.String()] {
			// TODO: Return a runtime error
			panic(fmt.Sprintf("name %s is duplicated in dict", entry.at))
		}
	}

	return DictPattern{entries}
}

func (p DictPattern) Bind(local Scope, value Value) (Scope, error) {
	dict, is := value.(Dict)
	if !is {
		return EmptyScope, fmt.Errorf("%s is not a dict", value)
	}

	extraElements := make(map[int]int)
	for i, entry := range p.entries {
		if _, is := entry.value.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = dict.Count() - len(p.entries)
		}
	}

	if len(p.entries) > dict.Count()+len(extraElements) {
		return EmptyScope, fmt.Errorf("length of dict %s shorter than dict pattern %s", dict, p)
	}

	if len(extraElements) == 0 && len(p.entries) < dict.Count() {
		return EmptyScope, fmt.Errorf("length of dict %s longer than dict pattern %s", dict, p)
	}

	result := EmptyScope
	m := dict.m
	for _, entry := range p.entries {
		var dictValue Value
		if _, is := entry.value.(ExtraElementPattern); is {
			if m.IsEmpty() {
				dictValue = None
			} else {
				dictValue = Dict{m: m}
			}
		} else {
			key := entry.at
			if lit, is := key.(LiteralExpr); is {
				key = lit.Literal()
			}

			dictExpr, found := m.Get(key)
			if !found {
				if entry.fallback == nil {
					return EmptyScope, fmt.Errorf("couldn't find %s in dict %s", key, m)
				}
				var err error
				dictValue, err = entry.fallback.Eval(local)
				if err != nil {
					return EmptyScope, err
				}
			} else {
				dictValue = dictExpr.(Value)
				m = m.Without(frozen.NewSet(key))
			}
		}

		scope, err := entry.value.Bind(local, dictValue)
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

func (p DictPattern) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, expr := range p.entries {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s", expr)
	}
	b.WriteByte('}')
	return b.String()
}

func (p DictPattern) Bindings() []string {
	bindings := make([]string, len(p.entries))
	for i, v := range p.entries {
		bindings[i] = v.value.String()
	}
	return bindings
}
