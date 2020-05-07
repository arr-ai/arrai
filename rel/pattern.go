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
	return scope.With(p.ident, value)
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

	return scope.With(p.String(), value)
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
	for i, v := range value.(Array).Values() {
		scope = p.items[i].Bind(scope, v)
	}
	return scope
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
