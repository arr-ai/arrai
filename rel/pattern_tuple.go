package rel

import (
	"bytes"
	"context"
	"fmt"
)

type TuplePatternAttr struct {
	name    string
	pattern FallbackPattern
}

func NewTuplePatternAttr(name string, pattern FallbackPattern) TuplePatternAttr {
	return TuplePatternAttr{
		name:    name,
		pattern: pattern,
	}
}

func (a TuplePatternAttr) String() string {
	if a.name == "" {
		return a.pattern.String()
	}
	if a.pattern.fallback == nil {
		return fmt.Sprintf("%s: %s", a.name, a.pattern)
	}
	return fmt.Sprintf("%s?: %s", a.name, a.pattern)
}

func (a *TuplePatternAttr) IsWildcard() bool {
	return a.name == "*"
}

type TuplePattern struct {
	attrs []TuplePatternAttr
}

func NewTuplePattern(attrs ...TuplePatternAttr) (TuplePattern, error) {
	p := TuplePattern{attrs}
	if err := validTuplePattern(p); err != nil {
		return TuplePattern{}, err
	}
	return p, nil
}

func validTuplePattern(p TuplePattern) error {
	names := make(map[string]struct{})
	for _, attr := range p.attrs {
		if _, has := names[attr.name]; has {
			if _, is := attr.pattern.pattern.(ExtraElementPattern); is && attr.name == "" {
				// Attr name is '' when its pattern is '...' which will crash with other attr whose name is '',
				// e.g. `let ('':value, ...) = ('':2); value`. So skip it in this validation.
				continue
			}
			return fmt.Errorf("duplicate fields found in pattern %s ", p)
		}
		if _, is := attr.pattern.pattern.(ExtraElementPattern); !is {
			names[attr.name] = struct{}{}
		}
	}
	return nil
}

func (p TuplePattern) Bind(ctx context.Context, local Scope, value Value, b *scopeBuilder) (context.Context, error) {
	tuple, is := value.(Tuple)
	if !is {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("%s is not a tuple", value)
	}

	// matched counts the tuple's attributes consumed by the pattern; the
	// leftover (for the `...` pattern, or the too-long check) is everything
	// else. This replaces maintaining a shrinking frozen name set per step.
	matched := 0
	var extraPattern *TuplePatternAttr

	for i, attr := range p.attrs {
		var tupleValue Value
		if _, is := attr.pattern.pattern.(ExtraElementPattern); is {
			// detects a second `...`
			if extraPattern != nil {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("non-deterministic pattern is not supported yet")
			}
			extraPattern = &p.attrs[i]
			continue
		}
		if attr.pattern.fallback == nil && matched == tuple.Count() {
			if !b.explain {
				return ctx, errPatternMismatch
			}
			return ctx, lazyErrorf("length of tuple %s shorter than tuple pattern %s", tuple, p)
		}
		var found bool
		tupleValue, found = tuple.Get(attr.name)
		if !found {
			if attr.pattern.fallback == nil {
				if !b.explain {
					return ctx, errPatternMismatch
				}
				return ctx, lazyErrorf("couldn't find %s in tuple %s", attr.name, tuple)
			}
			var err error
			tupleValue, err = attr.pattern.fallback.Eval(ctx, local)
			if err != nil {
				return ctx, err
			}
		} else {
			matched++
		}

		var err error
		ctx, err = attr.pattern.Bind(ctx, local, tupleValue, b)
		if err != nil {
			return ctx, err
		}
	}

	if extraPattern != nil {
		// Leftover names are computed once, only when a `...` needs them.
		names := tuple.Names()
		for _, attr := range p.attrs {
			if _, is := attr.pattern.pattern.(ExtraElementPattern); !is {
				names = names.Without(attr.name)
			}
		}
		tupleValue := tuple.Project(names)
		if tupleValue == nil {
			if !b.explain {
				return ctx, errPatternMismatch
			}
			return ctx, lazyErrorf("tuple %s cannot match tuple pattern %s", tuple, p)
		}
		var err error
		ctx, err = extraPattern.pattern.Bind(ctx, local, tupleValue, b)
		if err != nil {
			return ctx, err
		}
	} else if matched < tuple.Count() {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("length of tuple %s longer than tuple pattern %s", tuple, p)
	}

	return ctx, nil
}

func (p TuplePattern) String() string { //nolint:dupl
	var b bytes.Buffer
	b.WriteByte('(')
	for i, attr := range p.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		if attr.IsWildcard() {
			isDot := false
			if exprpat, is := attr.pattern.pattern.(ExprPattern); is {
				if ident, is := exprpat.Expr.(IdentExpr); is {
					isDot = ident.Ident() == "."
				}
			}
			if !isDot {
				b.WriteString(attr.pattern.String())
			}
			b.WriteString(".*")
		} else {
			b.WriteString(attr.String())
		}
	}
	b.WriteByte(')')
	return b.String()
}

func (p TuplePattern) Bindings() []string {
	bindings := make([]string, len(p.attrs))
	for i, v := range p.attrs {
		bindings[i] = v.pattern.String()
	}
	return bindings
}
