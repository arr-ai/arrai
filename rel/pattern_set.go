package rel

import (
	"bytes"
	"context"
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

func (p SetPattern) Bind(ctx context.Context, local Scope, value Value, b *scopeBuilder) (context.Context, error) {
	set, is := value.(Set)
	if !is {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("value %s is not a set", value)
	}
	extraElements := make(map[int]int)
	for i, ptn := range p.patterns {
		switch ptn.(type) {
		case ExtraElementPattern, IdentPattern, DynIdentPattern:
			if len(extraElements) == 1 {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = set.Count() - len(p.patterns)
		}
	}

	if len(p.patterns) > set.Count()+len(extraElements) {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("length of set %v shorter than set pattern %s", set, p)
	}

	if len(extraElements) == 0 && len(p.patterns) < set.Count() {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("length of set %v longer than set pattern %s", set, p)
	}

	for _, ptn := range p.patterns {
		if _, is := ptn.(ExtraElementPattern); is {
			continue
		}
		switch t := ptn.(type) {
		case IdentPattern:
		case DynIdentPattern:
		case ExprPattern:
			if v, is := t.Expr.(Value); is {
				if !set.Has(v) {
					if !b.explain {
						return ctx, errPatternMismatch
					}
					return ctx, lazyErrorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v)
				continue
			}

			if _, is := t.Expr.(IdentExpr); !is {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("item type %s is not supported yet", t)
			}
		case ExprsPattern:
			// Support cases:
			// AssertCodesEvalToSameValue(t, `{5, 6}`, `let x = 1; let y = 42; let {(x), (y), ...t} = {1, 42, 5, 6}; t`)
			// AssertCodeErrors(t, "", `let x = 1; let y = 42; let {(x), (y)} = {1, 4}; 2`)
			if identExpr, is := t.exprs[0].(IdentExpr); is {
				v, has := local.Get(identExpr.ident)
				if !has {
					if !b.explain {
						return ctx, errPatternMismatch
					}
					return ctx, lazyErrorf("%q not in scope", identExpr.ident)
				}
				if !set.Has(v.(Value)) {
					if !b.explain {
						return ctx, errPatternMismatch
					}
					return ctx, lazyErrorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v.(Value))
			}
		default:
			if len(p.patterns) == 1 {
				for e := set.Enumerator(); e.MoveNext(); {
					return t.Bind(ctx, local, e.Current(), b)
				}
			}
			// TODO: This is should return an error
			panic(fmt.Errorf("pattern type %T not supported yet", t))
		}
	}
	return p.bindExtras(ctx, local, set, extraElements, b)
}

// bindExtras binds the `...`/ident leftovers once the explicit elements are
// consumed.
func (p SetPattern) bindExtras(
	ctx context.Context, local Scope, set Set, extraElements map[int]int, b *scopeBuilder,
) (context.Context, error) {
	for i := range extraElements {
		var err error
		if _, is := p.patterns[i].(ExtraElementPattern); is {
			ctx, err = p.patterns[i].Bind(ctx, local, set, b)
		} else {
			if set.Count() != 1 {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("the length of set %v is wrong", set)
			}

			e := set.Enumerator()
			if !e.MoveNext() {
				panic("set with count 1 failed to enumerate")
			}
			ctx, err = p.patterns[i].Bind(ctx, local, e.Current(), b)
		}
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
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
