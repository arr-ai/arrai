package rel

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
