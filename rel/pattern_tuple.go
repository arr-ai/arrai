package rel

import (
	"bytes"
	"fmt"
)

type TuplePatternAttr struct {
	name    string
	pattern Pattern
}

func NewTuplePatternAttr(name string, pattern Pattern) TuplePatternAttr {
	return TuplePatternAttr{
		name:    name,
		pattern: pattern,
	}
}

func (a TuplePatternAttr) String() string {
	return fmt.Sprintf("%s:%s", a.name, a.pattern)
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

func (p TuplePattern) Bind(local Scope, value Value) (Scope, error) {
	tuple, is := value.(Tuple)
	if !is {
		return EmptyScope, fmt.Errorf("%s is not a tuple", value)
	}

	extraElements := make(map[int]int)
	for i, attr := range p.attrs {
		if _, is := attr.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = tuple.Count() - len(p.attrs)
		}
	}

	if len(p.attrs) > tuple.Count()+len(extraElements) {
		return EmptyScope, fmt.Errorf("length of tuple %s shorter than tuple pattern %s", tuple, p)
	}

	if len(extraElements) == 0 && len(p.attrs) < tuple.Count() {
		return EmptyScope, fmt.Errorf("length of tuple %s longer than tuple pattern %s", tuple, p)
	}

	result := EmptyScope
	names := tuple.Names()
	for _, attr := range p.attrs {
		if _, is := attr.pattern.(ExtraElementPattern); is {
			tupleExpr := tuple.Project(names)
			if tupleExpr == nil {
				return EmptyScope, fmt.Errorf("tuple %s cannot match tuple pattern %s", tuple, p)
			}
			scope, err := attr.pattern.Bind(local, tupleExpr)
			if err != nil {
				return EmptyScope, err
			}
			result = result.MatchedUpdate(scope)
			continue
		}
		tupleExpr, found := tuple.Get(attr.name)
		if !found {
			return EmptyScope, fmt.Errorf("couldn't find %s in tuple %s", attr.name, tuple)
		}
		scope, err := attr.pattern.Bind(local, tupleExpr)
		if err != nil {
			return EmptyScope, err
		}
		result = result.MatchedUpdate(scope)
		names = names.Without(attr.name)
	}

	return result, nil
}

func (p TuplePattern) String() string { //nolint:dupl
	var b bytes.Buffer
	b.WriteByte('(')
	for i, attr := range p.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		if attr.IsWildcard() {
			if ident, is := attr.pattern.(IdentExpr); !is || ident.Ident() != "." {
				b.WriteString(attr.pattern.String())
			}
			b.WriteString(".*")
		} else {
			b.WriteString(attr.name)
			b.WriteString(": ")
			b.WriteString(attr.pattern.String())
		}
	}
	b.WriteByte(')')
	return b.String()
}
