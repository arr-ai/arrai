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

func NewTuplePattern(attrs ...TuplePatternAttr) TuplePattern {
	p := TuplePattern{attrs}
	//if err := validTuplePattern(p); err != nil {
	//	panic(err)
	//}
	return p
}

func validTuplePattern(p TuplePattern) error {
	names := make(map[string]struct{})
	for _, attr := range p.attrs {
		if _, has := names[attr.name]; has {
			return fmt.Errorf("duplicate fields found in pattern %s ", p)
		} else {
			if _, is := attr.pattern.pattern.(ExtraElementPattern); !is {
				names[attr.name] = struct{}{}
			}
		}
	}
	return nil
}

func (p TuplePattern) Bind(ctx context.Context, local Scope, value Value) (context.Context, Scope, error) {
	tuple, is := value.(Tuple)
	if !is {
		return ctx, EmptyScope, fmt.Errorf("%s is not a tuple", value)
	}

	if err := validTuplePattern(p); err != nil {
		return ctx, EmptyScope, err
	}
	bind := func(
		ctx context.Context,
		attr TuplePatternAttr,
		base Scope,
		tupleValue Value,
	) (context.Context, Scope, error) {
		var scope Scope
		var err error
		ctx, scope, err = attr.pattern.Bind(ctx, local, tupleValue)
		if err != nil {
			return ctx, EmptyScope, err
		}
		result, err := base.MatchedUpdate(scope)
		if err != nil {
			return ctx, EmptyScope, err
		}
		return ctx, result, nil
	}

	result := EmptyScope
	names := tuple.Names()
	var extraPattern *TuplePatternAttr

	for i, attr := range p.attrs {
		var tupleValue Value
		if _, is := attr.pattern.pattern.(ExtraElementPattern); is {
			// detects a second `...`
			if extraPattern != nil {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraPattern = &p.attrs[i]
			continue
		} else {
			if attr.pattern.fallback == nil && !names.IsTrue() {
				return ctx, EmptyScope, fmt.Errorf("length of tuple %s shorter than tuple pattern %s", tuple, p)
			}
			var found bool
			tupleValue, found = tuple.Get(attr.name)
			if !found {
				if attr.pattern.fallback == nil {
					return ctx, EmptyScope, fmt.Errorf("couldn't find %s in tuple %s", attr.name, tuple)
				}
				var err error
				tupleValue, err = attr.pattern.fallback.Eval(ctx, local)
				if err != nil {
					return ctx, EmptyScope, err
				}
			}
		}
		var err error
		ctx, result, err = bind(ctx, attr, result, tupleValue)
		if err != nil {
			return ctx, EmptyScope, err
		}
		names = names.Without(attr.name)
	}

	if extraPattern != nil {
		tupleValue := tuple.Project(names)
		if tupleValue == nil {
			return ctx, EmptyScope, fmt.Errorf("tuple %s cannot match tuple pattern %s", tuple, p)
		}
		var err error
		ctx, result, err = bind(ctx, *extraPattern, result, tupleValue)
		if err != nil {
			return ctx, EmptyScope, err
		}
		names = EmptyNames
	}

	if names.IsTrue() {
		return ctx, EmptyScope, fmt.Errorf("length of tuple %s longer than tuple pattern %s", tuple, p)
	}

	return ctx, result, nil
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
