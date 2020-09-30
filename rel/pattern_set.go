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

func (p SetPattern) Bind(ctx context.Context, local Scope, value Value) (context.Context, Scope, error) {
	set, is := value.(Set)
	if !is {
		return ctx, EmptyScope, fmt.Errorf("value %s is not a set", value)
	}
	extraElements := make(map[int]int)
	for i, ptn := range p.patterns {
		switch ptn.(type) {
		case ExtraElementPattern, IdentPattern, DynIdentPattern:
			if len(extraElements) == 1 {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = set.Count() - len(p.patterns)
		}
	}

	if len(p.patterns) > set.Count()+len(extraElements) {
		return ctx, EmptyScope, fmt.Errorf("length of set %s shorter than set pattern %s", set, p)
	}

	if len(extraElements) == 0 && len(p.patterns) < set.Count() {
		return ctx, EmptyScope, fmt.Errorf("length of set %s longer than set pattern %s", set, p)
	}

	result := EmptyScope
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
					return ctx, EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v).(GenericSet)
				continue
			}

			if _, is := t.Expr.(IdentExpr); !is {
				return ctx, EmptyScope, fmt.Errorf("item type %s is not supported yet", t)
			}
		case ExprsPattern:
			// Support cases:
			// AssertCodesEvalToSameValue(t, `{5, 6}`, `let x = 1; let y = 42; let {(x), (y), ...t} = {1, 42, 5, 6}; t`)
			// AssertCodeErrors(t, "", `let x = 1; let y = 42; let {(x), (y)} = {1, 4}; 2`)
			if identExpr, is := t.exprs[0].(IdentExpr); is {
				v, has := local.Get(identExpr.ident)
				if !has {
					return ctx, EmptyScope, fmt.Errorf("%q not in scope", identExpr.ident)
				}
				if !set.Has(v.(Value)) {
					return ctx, EmptyScope, fmt.Errorf("item %s is not included in set %s", v, value)
				}
				set = set.Without(v.(Value)).(GenericSet)
			}
		default:
			panic(fmt.Errorf("pattern type %T not supported yet", t))
		}
	}
	for i := range extraElements {
		var scope Scope
		var err error
		if _, is := p.patterns[i].(ExtraElementPattern); is {
			ctx, scope, err = p.patterns[i].Bind(ctx, local, set)
		} else {
			if set.Count() != 1 {
				return ctx, EmptyScope, fmt.Errorf("the length of set %s is wrong", set)
			}

			e := set.Enumerator()
			if !e.MoveNext() {
				panic("set with count 1 failed to enumerate")
			}
			ctx, scope, err = p.patterns[i].Bind(ctx, local, e.Current())
			if err != nil {
				return ctx, EmptyScope, err
			}
		}
		if err != nil {
			return ctx, EmptyScope, err
		}

		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return ctx, EmptyScope, err
		}
	}

	return ctx, result, nil
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
