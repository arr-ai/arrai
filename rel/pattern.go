package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/frozen"
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

type ExtraElementPattern struct {
	ident string
}

func NewExtraElementPattern(ident string) ExtraElementPattern {
	return ExtraElementPattern{ident}
}

func (p ExtraElementPattern) Bind(scope Scope, value Value) Scope {
	if p.ident == "" {
		return EmptyScope
	}
	return EmptyScope.With(p.ident, value)
}

func (p ExtraElementPattern) String() string {
	return "..." + p.ident
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

	if len(p.items) > array.Count() {
		panic(fmt.Sprintf("length of array %s shorter than array pattern %s", array, p))
	}

	extraElements := make(map[int]int)
	for i, item := range p.items {
		if _, is := item.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				panic("multiple ... not supported yet")
			}
			extraElements[i] = array.Count() - len(p.items)
		}
	}

	if len(extraElements) == 0 && len(p.items) < array.Count() {
		panic(fmt.Sprintf("length of array %s longer than array pattern %s", array, p))
	}

	result := EmptyScope
	offset := 0
	for i, item := range p.items {
		if _, is := item.(ExtraElementPattern); is {
			offset = extraElements[i]
			result = result.MatchedUpdate(item.Bind(scope, NewArray(array.Values()[i:i+offset+1]...)))
			continue
		}
		result = result.MatchedUpdate(item.Bind(scope, array.Values()[i+offset]))
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

func (p TuplePattern) Bind(scope Scope, value Value) Scope {
	tuple, is := value.(Tuple)
	if !is {
		panic(fmt.Sprintf("%s is not a tuple", value))
	}

	if len(p.attrs) > tuple.Count() {
		panic(fmt.Sprintf("length of tuple %s shorter than tuple pattern %s", tuple, p))
	}

	extraElements := make(map[int]int)
	for i, attr := range p.attrs {
		if _, is := attr.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				panic("multiple ... not supported yet")
			}
			extraElements[i] = tuple.Count() - len(p.attrs)
		}
	}

	if len(extraElements) == 0 && len(p.attrs) < tuple.Count() {
		panic(fmt.Sprintf("length of tuple %s longer than tuple pattern %s", tuple, p))
	}

	result := EmptyScope
	names := tuple.Names()
	for _, attr := range p.attrs {
		if _, is := attr.pattern.(ExtraElementPattern); is {
			tupleExpr := tuple.Project(names)
			if tupleExpr == nil {
				panic(fmt.Sprintf("tuple %s cannot match tuple pattern %s", tuple, p))
			}
			result = result.MatchedUpdate(attr.pattern.Bind(scope, tupleExpr))
			continue
		}
		tupleExpr := tuple.MustGet(attr.name)
		result = result.MatchedUpdate(attr.pattern.Bind(scope, tupleExpr))
		names = names.Without(attr.name)
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

func (p DictPattern) Bind(scope Scope, value Value) Scope {
	dict, is := value.(Dict)
	if !is {
		panic(fmt.Sprintf("%s is not a dict", value))
	}

	if len(p.entries) > dict.Count() {
		panic(fmt.Sprintf("length of dict %s shorter than dict pattern %s", dict, p))
	}

	extraElements := make(map[int]int)
	for i, entry := range p.entries {
		if _, is := entry.value.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				panic("multiple ... not supported yet")
			}
			extraElements[i] = dict.Count() - len(p.entries)
		}
	}

	if len(extraElements) == 0 && len(p.entries) < dict.Count() {
		panic(fmt.Sprintf("length of dict %s longer than dict pattern %s", dict, p))
	}

	result := EmptyScope
	m := dict.m
	for _, entry := range p.entries {
		if _, is := entry.value.(ExtraElementPattern); is {
			result = result.MatchedUpdate(entry.value.Bind(scope, Dict{m: m}))
			continue
		}
		dictValue := m.MustGet(entry.at)
		result = result.MatchedUpdate(entry.value.Bind(scope, dictValue.(Value)))
		m = m.Without(frozen.NewSet(entry.at))
	}

	return result
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
