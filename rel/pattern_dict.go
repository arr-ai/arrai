package rel

import (
	"bytes"
	"context"
	"fmt"
)

type DictPatternEntry struct {
	at      Expr
	pattern FallbackPattern
}

func NewDictPatternEntry(at Expr, pattern FallbackPattern) DictPatternEntry {
	return DictPatternEntry{
		at:      at,
		pattern: pattern,
	}
}

func (p DictPatternEntry) String() string {
	if p.pattern.fallback == nil {
		return fmt.Sprintf("%s: %s", p.at, p.pattern)
	}
	return fmt.Sprintf("%s?: %s", p.at, p.pattern)
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

func (p DictPattern) Bind(ctx context.Context, local Scope, value Value, b *scopeBuilder) (context.Context, error) {
	dict, is := value.(Dict)
	if !is {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("%s is not a dict", value)
	}

	extraElements := make(map[int]int)
	for i, entry := range p.entries {
		if _, is := entry.pattern.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = dict.Count() - len(p.entries)
		}
		if entry.pattern.fallback != nil {
			if len(extraElements) == 1 {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = dict.Count() - len(p.entries)
		}
	}

	if len(p.entries) > dict.Count()+len(extraElements) {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("length of dict %s shorter than dict pattern %s", dict, p)
	}

	if len(extraElements) == 0 && len(p.entries) < dict.Count() {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("length of dict %s longer than dict pattern %s", dict, p)
	}

	m := dict.m
	for _, entry := range p.entries {
		var dictValue Value
		if _, is := entry.pattern.pattern.(ExtraElementPattern); is {
			if m.IsEmpty() {
				dictValue = None
			} else {
				dictValue = newDict(m)
			}
		} else {
			key := entry.at
			if lit, is := key.(LiteralExpr); is {
				key = lit.Literal()
			}

			dictExpr, found := m.Get(key.(Value))
			if !found {
				if entry.pattern.fallback == nil {
					if !b.explain {
						return ctx, errPatternMismatch
					}
					return ctx, lazyErrorf("couldn't find %s in dict %s", key, m)
				}
				var err error
				dictValue, err = entry.pattern.fallback.Eval(ctx, local)
				if err != nil {
					return ctx, err
				}
			} else {
				dictValue = dictExpr.(Value)
				m = m.Without(key.(Value))
			}
		}

		var err error
		ctx, err = entry.pattern.pattern.Bind(ctx, local, dictValue, b)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
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
		bindings[i] = v.pattern.pattern.String()
	}
	return bindings
}
