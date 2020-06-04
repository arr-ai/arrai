package rel

import (
	"bytes"
	"fmt"
)

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
		if _, is := item.(ExtraElementPattern); is {
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
