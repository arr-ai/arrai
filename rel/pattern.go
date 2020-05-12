package rel

import (
	"bytes"
	"fmt"
)

type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	Bind(scope Scope, value Value) Scope
}

type IdentPattern struct {
	ident string
}

func NewIdentPattern(ident string) IdentPattern {
	return IdentPattern{ident}
}

func (p IdentPattern) Bind(scope Scope, value Value) Scope {
	return EmptyScope.With(p.ident, value)
}

func (p IdentPattern) String() string {
	return p.ident
}

type ValuePattern struct {
	value Value
}

func NewValuePattern(val Value) ValuePattern {
	return ValuePattern{val}
}

func (p ValuePattern) Bind(scope Scope, value Value) Scope {
	switch v := p.value.(type) {
	case Number:
		if !v.Equal(value) {
			panic(fmt.Sprintf("%s doesn't equal to %s", v, value))
		}
	}

	return EmptyScope
}

func (p ValuePattern) String() string {
	return p.value.String()
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
