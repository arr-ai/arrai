package rel

import (
	"bytes"
	"fmt"
)

// Pattern can be inside an Expr, Expr can be a Pattern.
type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	Bind(scope Scope, value Value) Scope
}

func ExprAsPattern(expr Expr) Pattern {
	switch t := expr.(type) {
	case IdentExpr:
		return t
	case Number:
		return t
	default:
		panic(fmt.Sprintf("%s is not a Pattern", t))
	}
}

type IdentPattern struct {
	ident string
}

func NewIdentPattern(ident string) IdentPattern {
	return IdentPattern{ident}
}

func (p IdentPattern) Bind(scope Scope, value Value) Scope {
	scope.MustGet(p.ident)
	scope.MatchedWith(p.ident, value)
	return EmptyScope.With(p.ident, value)
}

func (p IdentPattern) String() string {
	return p.ident
}

type ArrayPattern struct {
	items []Pattern
}

func NewArrayPattern(elements ...Pattern) ArrayPattern {
	return ArrayPattern{elements}
}

func (p ArrayPattern) Bind(scope Scope, value Value) Scope {
	if s, is := value.(GenericSet); is {
		if s.set.IsEmpty() {
			return EmptyScope
		}
		panic(fmt.Sprintf("value %s is not an array", value))
	}
	array, is := value.(Array)
	if !is {
		panic(fmt.Sprintf("value %s is not an array", value))
	}

	result := EmptyScope
	for i, item := range p.items {
		if len(array.Values()) < i+1 {
			panic(fmt.Sprintf("length of value %s shorter than array pattern %s", array.Values(), p.items))
		}
		result = result.MatchedUpdate(item.Bind(scope, array.Values()[i]))
	}

	return result
}

func (p ArrayPattern) String() string {
	var b bytes.Buffer
	b.WriteByte('[')
	for i, item := range p.items {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(item.String())
	}
	b.WriteByte(']')
	return b.String()
}

type AttrPattern struct {
	name    string
	pattern Pattern
}

func NewAttrPattern(name string, pattern Pattern) AttrPattern {
	return AttrPattern{
		name:    name,
		pattern: pattern,
	}
}

func (p AttrPattern) Bind(scope Scope, value Value) Scope {
	fmt.Println(value)
	return scope
}

func (p AttrPattern) String() string {
	return fmt.Sprintf("%s:%s", p.name, p.pattern)
}

func (p *AttrPattern) IsWildcard() bool {
	return p.name == "*"
}

type TuplePattern struct {
	attrs []AttrPattern
}

func NewTuplePattern(attrs ...AttrPattern) TuplePattern {
	return TuplePattern{attrs}
}

func (p TuplePattern) Bind(scope Scope, value Value) Scope {
	tuple, is := value.(Tuple)
	if !is {
		panic(fmt.Sprintf("%s is not a tuple", value))
	}

	result := EmptyScope
	for _, attr := range p.attrs {
		tupleExpr := tuple.MustGet(attr.name)
		result = result.MatchedUpdate(attr.pattern.Bind(scope, tupleExpr))
	}

	return result
}

func (p TuplePattern) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	for i, attr := range p.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		if attr.IsWildcard() {
			if attr.pattern != DotIdent {
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
