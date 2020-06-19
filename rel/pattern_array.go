package rel

import (
	"bytes"
	"fmt"
)

type PatternFallback struct {
	pattern  Pattern
	fallback Expr
}

func NewPatternFallback(pattern Pattern, fallback Expr) PatternFallback {
	return PatternFallback{
		pattern:  pattern,
		fallback: fallback,
	}
}

func (f PatternFallback) String() string {
	if f.fallback == nil {
		return f.pattern.String()
	}
	return fmt.Sprintf("%s?:%s", f.pattern, f.fallback)
}

type ArrayPattern struct {
	items []PatternFallback
}

func NewArrayPattern(elements ...PatternFallback) ArrayPattern {
	return ArrayPattern{elements}
}

func (p ArrayPattern) Bind(local Scope, value Value) (Scope, error) {
	if s, is := value.(GenericSet); is {
		if s.set.IsEmpty() {
			if len(p.items) == 0 {
				return EmptyScope, nil
			}
			return EmptyScope, fmt.Errorf("value [] is empty but pattern %s is not", p)
		}
		return EmptyScope, fmt.Errorf("value %s is not an array", value)
	}

	array, is := value.(Array)
	if !is {
		return EmptyScope, fmt.Errorf("value %s is not an array", value)
	}

	extraElements := make(map[int]int)
	for i, item := range p.items {
		if _, is := item.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = array.Count() - len(p.items)
		}
		if item.fallback != nil {
			if len(extraElements) == 1 {
				return EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = array.Count() - len(p.items)
		}
	}

	if len(p.items) > array.Count()+len(extraElements) {
		return EmptyScope, fmt.Errorf("length of array %s shorter than array pattern %s", array, p)
	}

	if len(extraElements) == 0 && len(p.items) < array.Count() {
		return EmptyScope, fmt.Errorf("length of array %s longer than array pattern %s", array, p)
	}

	result := EmptyScope
	offset := 0
	for i, item := range p.items {
		var value Value
		if _, is := item.pattern.(ExtraElementPattern); is {
			offset = extraElements[i]
			arr := NewArray()
			if offset >= 0 {
				arr = NewArray(array.Values()[i : i+offset+1]...)
			}
			value = arr
		} else if array.Count() <= i+offset {
			if item.fallback == nil {
				return EmptyScope, fmt.Errorf("length of array %s shorter than array pattern %s", array, p)
			}
			var err error
			value, err = item.fallback.Eval(local)
			if err != nil {
				return EmptyScope, err
			}
		} else {
			value = array.Values()[i+offset]
		}

		scope, err := item.pattern.Bind(local, value)
		if err != nil {
			return EmptyScope, err
		}
		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return Scope{}, err
		}
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

func (p ArrayPattern) Bindings() []string {
	bindings := make([]string, len(p.items))
	for i, v := range p.items {
		bindings[i] = v.String()
	}
	return bindings
}
