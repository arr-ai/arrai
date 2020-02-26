package rel

import "github.com/arr-ai/frozen"

const ArrayItemAttr = "@item"

// NewArray constructs an array as a relation.
func NewArray(values ...Value) Set {
	tuples := make([]Value, len(values))
	for i, value := range values {
		tuples[i] = NewTuple(
			Attr{"@", NewNumber(float64(i))},
			Attr{ArrayItemAttr, value},
		)
	}
	return NewSet(tuples...)
}

func ArrayEnumerator(set Set) ValueEnumerator {
	s := set.(*genericSet)
	return &arrayEnumerator{s.set.OrderedRange(func(a, b interface{}) bool {
		return a.(Tuple).MustGet("@").(Number) < b.(Tuple).MustGet("@").(Number)
	})}
}

type arrayEnumerator struct {
	i frozen.Iterator
}

func (a *arrayEnumerator) MoveNext() bool {
	return a.i.Next()
}

func (a *arrayEnumerator) Current() Value {
	return a.i.Value().(Tuple).MustGet(ArrayItemAttr)
}

// func isArrayTuple(v Value) (index int, item Value, is bool) {
// 	is = NewTupleMatcher(
// 		map[string]Matcher{
// 			"@":           MatchInt(func(i int) { index = i }),
// 			ArrayItemAttr: Bind(&item),
// 		},
// 		Lit(EmptyTuple),
// 	).Match(v)
// 	return
// }
