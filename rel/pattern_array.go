package rel

import (
	"bytes"
	"context"
	"fmt"
)

type ArrayPattern struct {
	items []FallbackPattern
}

func NewArrayPattern(elements ...FallbackPattern) ArrayPattern {
	return ArrayPattern{elements}
}

func (p ArrayPattern) Bind(ctx context.Context, local Scope, value Value) (context.Context, Scope, error) {
	if s, is := value.(GenericSet); is {
		if s.set.IsEmpty() {
			if len(p.items) == 0 {
				return ctx, EmptyScope, nil
			}
			return ctx, EmptyScope, fmt.Errorf("value [] is empty but pattern %s is not", p)
		}
		return ctx, EmptyScope, fmt.Errorf("value %s is not an array", value)
	}

	array, is := value.(Array)
	if !is {
		return ctx, EmptyScope, fmt.Errorf("value %s is not an array", value)
	}

	extraElements := make(map[int]int)
	for i, item := range p.items {
		if _, is := item.pattern.(ExtraElementPattern); is {
			if len(extraElements) == 1 {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = array.Count() - len(p.items)
		}
		if item.fallback != nil {
			if len(extraElements) == 1 {
				return ctx, EmptyScope, fmt.Errorf("non-deterministic pattern is not supported yet")
			}
			extraElements[i] = array.Count() - len(p.items)
		}
	}

	if len(p.items) > array.Count()+len(extraElements) {
		return ctx, EmptyScope, fmt.Errorf("length of array %s shorter than array pattern %s", array, p)
	}

	if len(extraElements) == 0 && len(p.items) < array.Count() {
		return ctx, EmptyScope, fmt.Errorf("length of array %s longer than array pattern %s", array, p)
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
				return ctx, EmptyScope, fmt.Errorf("length of array %s shorter than array pattern %s", array, p)
			}
			var err error
			value, err = item.fallback.Eval(ctx, local)
			if err != nil {
				return ctx, EmptyScope, err
			}
		} else {
			value = array.Values()[i+offset]
		}

		var scope Scope
		var err error
		ctx, scope, err = item.pattern.Bind(ctx, local, value)
		if err != nil {
			return ctx, EmptyScope, err
		}
		result, err = result.MatchedUpdate(scope)
		if err != nil {
			return ctx, Scope{}, err
		}
	}

	return ctx, result, nil
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
