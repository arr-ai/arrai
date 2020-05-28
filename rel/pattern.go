package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/frozen"
	"github.com/go-errors/errors"
)

// Pattern can be inside an Expr, Expr can be a Pattern.
type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	Bind(scope Scope, value Value) (Scope, error)
}

func ExprAsPattern(expr Expr) Pattern {
	switch t := expr.(type) {
	case IdentExpr:
		return t
	case Number:
		return t
	case Array:
		return t
	case ArrayExpr:
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

func (p IdentPattern) Bind(scope Scope, value Value) (Scope, error) {
	scope.MustGet(p.ident)
	scope.MatchedWith(p.ident, value)
	return EmptyScope.With(p.ident, value), nil
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

func (p ExtraElementPattern) Bind(scope Scope, value Value) (Scope, error) {
	if p.ident == "" {
		return EmptyScope, nil
	}
	return EmptyScope.With(p.ident, value), nil
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

func (p ArrayPattern) Bind(local Scope, value Value) (Scope, error) {
	if s, is := value.(GenericSet); is {
		if s.set.IsEmpty() {
			if len(p.items) == 0 {
				return EmptyScope, nil
			}
			panic(fmt.Sprintf("value [] is empty but pattern %s is not", p))
		}
		panic(fmt.Sprintf("value %s is not an array", value))
	}

	array, is := value.(Array)
	if !is {
		panic(fmt.Sprintf("value %s is not an array", value))
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

	if len(p.items) > array.Count()+len(extraElements) {
		panic(fmt.Sprintf("length of array %s shorter than array pattern %s", array, p))
	}

	if len(extraElements) == 0 && len(p.items) < array.Count() {
		panic(fmt.Sprintf("length of array %s longer than array pattern %s", array, p))
	}

	result := EmptyScope
	offset := 0
	for i, item := range p.items {
		if _, is := item.(ExtraElementPattern); is {
			offset = extraElements[i]
			arr := NewArray()
			if offset >= 0 {
				arr = NewArray(array.Values()[i : i+offset+1]...)
			}
			scope, err := item.Bind(local, arr)
			if err != nil {
				return EmptyScope, err
			}
			result = result.MatchedUpdate(scope)
			continue
		}
		scope, err := item.Bind(local, array.Values()[i+offset])
		if err != nil {
			return EmptyScope, err
		}
		result = result.MatchedUpdate(scope)
	}

	return result, nil
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

func (p TuplePattern) Bind(local Scope, value Value) (Scope, error) {
	tuple, is := value.(Tuple)
	if !is {
		panic(fmt.Sprintf("%s is not a tuple", value))
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

	if len(p.attrs) > tuple.Count()+len(extraElements) {
		panic(fmt.Sprintf("length of tuple %s shorter than tuple pattern %s", tuple, p))
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
			scope, err := attr.pattern.Bind(local, tupleExpr)
			if err != nil {
				return EmptyScope, err
			}
			result = result.MatchedUpdate(scope)
			continue
		}
		tupleExpr := tuple.MustGet(attr.name)
		scope, err := attr.pattern.Bind(local, tupleExpr)
		if err != nil {
			return EmptyScope, err
		}
		result = result.MatchedUpdate(scope)
		names = names.Without(attr.name)
	}

	return result, nil
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

func (p DictPattern) Bind(local Scope, value Value) (Scope, error) {
	dict, is := value.(Dict)
	if !is {
		panic(fmt.Sprintf("%s is not a dict", value))
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

	if len(p.entries) > dict.Count()+len(extraElements) {
		panic(fmt.Sprintf("length of dict %s shorter than dict pattern %s", dict, p))
	}

	if len(extraElements) == 0 && len(p.entries) < dict.Count() {
		panic(fmt.Sprintf("length of dict %s longer than dict pattern %s", dict, p))
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
				result = result.MatchedUpdate(scope)
			} else {
				scope, err := entry.value.Bind(local, Dict{m: m})
				if err != nil {
					return EmptyScope, err
				}
				result = result.MatchedUpdate(scope)
			}

			continue
		}
		dictValue := m.MustGet(entry.at)
		scope, err := entry.value.Bind(local, dictValue.(Value))
		if err != nil {
			return EmptyScope, err
		}
		result = result.MatchedUpdate(scope)
		m = m.Without(frozen.NewSet(entry.at))
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

type ExprsPattern struct {
	exprs []Expr
}

func NewExprsPattern(exprs []Expr) ExprsPattern {
	return ExprsPattern{exprs: exprs}
}

func (ep ExprsPattern) Bind(scope Scope, value Value) (Scope, error) {
	incomingVal, err := value.Eval(scope)
	if err != nil {
		return EmptyScope, err
	}

	for _, e := range ep.exprs {
		val, err := e.Eval(scope)
		if err != nil {
			return EmptyScope, err
		}
		if incomingVal.Equal(val) {
			return scope, nil
		}
	}

	return EmptyScope, errors.Errorf("didn't find matched value")
}

func (ep ExprsPattern) String() string {
	var b bytes.Buffer
	b.WriteByte('[')

	for i, e := range ep.exprs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v", e.String())
	}

	b.WriteByte(']')
	return b.String()
}
