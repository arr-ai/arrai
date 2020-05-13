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

type ArrayPattern struct {
	items []Pattern
}

func NewArrayPattern(elements ...Pattern) ArrayPattern {
	return ArrayPattern{elements}
}

func (p ArrayPattern) Bind(scope Scope, value Value) Scope {
	array, is := value.(Array)
	if !is {
		panic(fmt.Sprintf("value %s is not an array", value))
	}

	values := make(map[string]Value)
	patterns := make(map[string]Pattern)
	for i, item := range p.items {
		if len(array.Values()) < i+1 {
			panic(fmt.Sprintf("length of value %s shorter than array pattern %s", array.Values(), p.items))
		}
		// `_` should never appear in scope
		if item.String() == "_" {
			continue
		}

		if expr, exists := scope.Get(item.String()); exists {
			if expr.String() != array.Values()[i].String() {
				panic(fmt.Sprintf("%s is redefined differently", item))
			}
		}

		if v, ok := values[item.String()]; ok {
			if v.Kind() == array.Values()[i].Kind() && v.String() == array.Values()[i].String() {
				continue
			}
			panic(fmt.Sprintf("value %s does not equal to value %s", v, array.Values()[i]))
		}
		values[item.String()] = array.Values()[i]
		patterns[item.String()] = item
	}

	result := EmptyScope
	for s, ptn := range patterns {
		result = result.Update(ptn.Bind(scope, values[s]))
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

	values := make(map[string]Value)
	patterns := make(map[string]Pattern)
	for _, attr := range p.attrs {
		tupleExpr := tuple.MustGet(attr.name)
		if expr, exists := scope.Get(attr.pattern.String()); exists {
			if expr.String() != tupleExpr.String() {
				panic(fmt.Sprintf("%s is redefined differently", attr.pattern))
			}
		}

		if v, ok := values[attr.pattern.String()]; ok {
			if v.Kind() == tupleExpr.Kind() && v.String() == tupleExpr.String() {
				continue
			}
			panic(fmt.Sprintf("value %s does not equal to value %s", v, tupleExpr))
		}
		values[attr.pattern.String()] = tupleExpr
		patterns[attr.pattern.String()] = attr.pattern
	}

	result := EmptyScope
	for s, ptn := range patterns {
		result = result.Update(ptn.Bind(scope, values[s]))
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
