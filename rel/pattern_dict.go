package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/frozen"
)

type DictPatternEntry struct {
	at    Expr
	value Pattern
}

func NewDictPatternEntry(at Expr, value Pattern) DictPatternEntry {
	return DictPatternEntry{
		at:    at,
		value: value,
	}
}

func (p DictPatternEntry) String() string {
	return fmt.Sprintf("%s:%s", p.at, p.value)
}

type DictPattern struct {
	entries []DictPatternEntry
}

func NewDictPattern(entries ...DictPatternEntry) DictPattern {
	names := make(map[string]bool)
	for _, entry := range entries {
		if names[entry.at.String()] {
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
		if _, is := entry.value.(ExtraElementPattern); is {
			if m.IsEmpty() {
				scope, err := entry.value.Bind(local, None)
				if err != nil {
					return EmptyScope, err
				}
				result, err = result.MatchedUpdate(scope)
				if err != nil {
					return EmptyScope, err
				}
			} else {
				scope, err := entry.value.Bind(local, Dict{m: m})
				if err != nil {
					return EmptyScope, err
				}
				result, err = result.MatchedUpdate(scope)
				if err != nil {
					return EmptyScope, err
				}
			}

			continue
		}
		key := entry.at
		if lit, is := key.(LiteralExpr); is {
			key = lit.Literal()
		}
		dictValue, found := m.Get(key)
		if !found {
			return EmptyScope, fmt.Errorf("couldn't find %s in dict %s", key, m)
		}
		scope, err := entry.value.Bind(local, dictValue.(Value))
		if err != nil {
			return EmptyScope, err
		}
		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return EmptyScope, err
		}
		m = m.Without(frozen.NewSet(key))
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
		fmt.Fprintf(&b, "%v: %v", expr.at.String(), expr.value.String())
	}
	b.WriteByte('}')
	return b.String()
}
