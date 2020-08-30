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
	names := make(map[string]bool)
	for _, attr := range attrs {
		if names[attr.name] {
			panic(fmt.Sprintf("name %s is duplicated in tuple", attr.name))
		}
	}
	return TuplePattern{attrs}
}

func (p TuplePattern) Bind(ctx context.Context, local Scope, value Value) (context.Context, Scope, error) {
	tuple, is := value.(Tuple)
	if !is {
		return ctx, EmptyScope, fmt.Errorf("%s is not a tuple", value)
	}

	extraElements := make(map[int]int)
	for i, attr := range p.attrs {
		if _, is := attr.pattern.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = tuple.Count() - len(p.attrs)
		}
		if attr.pattern.fallback != nil {
			if len(extraElements) == 1 {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = tuple.Count() - len(p.attrs)
		}
	}

	if len(p.attrs) > tuple.Count()+len(extraElements) {
		return ctx, EmptyScope, fmt.Errorf("length of tuple %s shorter than tuple pattern %s", tuple, p)
	}

	if len(extraElements) == 0 && len(p.attrs) < tuple.Count() {
		return ctx, EmptyScope, fmt.Errorf("length of tuple %s longer than tuple pattern %s", tuple, p)
	}

	result := EmptyScope
	names := tuple.Names()
	for _, attr := range p.attrs {
		var tupleValue Value
		if _, is := attr.pattern.pattern.(ExtraElementPattern); is {
			tupleValue = tuple.Project(names)
			if tupleValue == nil {
				return ctx, EmptyScope, fmt.Errorf("tuple %s cannot match tuple pattern %s", tuple, p)
			}
		} else {
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
		var scope Scope
		var err error
		ctx, scope, err = attr.pattern.Bind(ctx, local, tupleValue)
		if err != nil {
			return ctx, EmptyScope, err
		}
		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return ctx, EmptyScope, err
		}
		names = names.Without(attr.name)
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
