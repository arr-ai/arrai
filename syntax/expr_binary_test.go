//nolint:dupl
package syntax

import (
	"testing"
)

func TestWhereExpr(t *testing.T) {
	t.Parallel()
	s := `{|a,b| (3,41), (2,42), (1,43)}`
	// defer trace().revert()
	AssertCodesEvalToSameValue(t, `{(a:3, b:41)}`, s+` where .a=3`)
}

func TestRelationCall(t *testing.T) {
	t.Parallel()
	s := `{"key": "val"}("key")`
	AssertCodesEvalToSameValue(t, `"val"`, s)
}

func TestOpsAddArrowForTuples(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2, c: 3, d: 4)`, `(a: 1, b: 2) +> (c: 3, d: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 2), c: (b: 1), d: (c: 4))`, `(a: 1, b: (c: 2)) +> (c: (b: 1), d: (c: 4))`)

	AssertCodesEvalToSameValue(t, `(a: 1, b: 3, c: 4)`, `(a: 1, b: 2) +> (b: 3, c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 3, c: 4)`, `(a: 1, b: (c: 2)) +> (b: 3, c: 4)`)

	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 4), c: 4)`, `(a: 1, b: (c: 2)) +> (b: (c: 4), c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: (b: 1), b: (c: 4))`, `(a: 1, b: (c: 2)) +> (a: (b: 1), b: (c: 4))`)
}

func TestOpsAddArrowForDicts(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `{} +> {}`)
	AssertCodesEvalToSameValue(t, `{1: 2}`, `{1: 2} +> {}`)
	AssertCodesEvalToSameValue(t, `{3: 4}`, `{} +> {3: 4}`)

	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 3, 'd': 4}`, `{'a': 1, 'b': 2} +> {'b': 3, 'd': 4}`)
	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 3, 'd': 4}`, `{'a': 1, 'b': {'a': 2}} +> {'b': 3, 'd': 4}`)
	AssertCodesEvalToSameValue(t, `{'a': {'c': 2}}`, `{'a': {'b': 1}} +> {'a': {'c': 2}}`)

	AssertCodesEvalToSameValue(t, `{'a': 'A1', 'b': 'A3', 'd': 'A4'}`, `{'a': 'A1', 'b': 'A2'} +> {'b': 'A3', 'd': 'A4'}`)
	AssertCodesEvalToSameValue(t, `{'a': {'c': 'ABC2'}}`, `{'a': {'b': 'ABC1'}} +> {'a': {'c': 'ABC2'}}`)

	AssertCodesEvalToSameValue(t, `{'a': 'A1', 'b': 'A2', 'c': 'A3', 'd': 'A4'}`,
		`{'a': 'A1', 'b': 'A2'} +> {'c': 'A3', 'd': 'A4'}`)
	AssertCodesEvalToSameValue(t, `{'a': {'b': 'ABC1'}, 'b': {'c': 'ABC2'}}`,
		`{'a': {'b': 'ABC1'}} +> {'b': {'c': 'ABC2'}}`)
}

func TestOpsNestedAddArrowMixTupleDict(t *testing.T) {
	t.Parallel()
	// |
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2})))`,
		`(a: (b: (c: {1}))) +> (a+>: (b+>: (c|: {2})))`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {1, 2}}}}`,
		`{'a': {'b': {'c': {1}}}} +> {'a'+>: {'b'+>: {'c'|: {2}}}}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: {'b': (c: {1, 2})})`,
		`(a: {'b': (c: {1})}) +> (a+>: {'b'+>: (c|: {2})})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': (b: {'c': {1, 2}})}`,
		`{'a': (b: {'c': {1}})} +> {'a'+>: ('b'+>: {'c'|: {2}})}`,
	)

	// ++
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: [1, 2])))`,
		`(a: (b: (c: [1]))) +> (a+>: (b+>: (c++: [2])))`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': [1, 2]}}}`,
		`{'a': {'b': {'c': [1]}}} +> {'a'+>: {'b'+>: {'c'++: [2]}}}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: {'b': (c: [1, 2])})`,
		`(a: {'b': (c: [1])}) +> (a+>: {'b'+>: (c++: [2])})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': (b: {'c': [1, 2]})}`,
		`{'a': (b: {'c': [1]})} +> {'a'+>: ('b'+>: {'c'++: [2]})}`,
	)

	// very deep
	AssertCodesEvalToSameValue(t,
		`(a: (b: {'c': {'d': (e: {'f': {'g': (h: (i: (j: 2)))}})}}))`,
		`
		let x = (a: (b: {'c': {'d': (e: {'f': {'g': (h: (i: (j: 1)))}})}}));
		x +> (a+>: (b+>: {'c'+>: {'d'+>: (e+>: {'f'+>: {'g'+>: (h+>: (i+>: (j: 2)))}})}}))
		`,
	)

	// multiple nested ops at the same level
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2}, d: {1}), e: (f: 1)), g: (h: 1))`,
		`(a: (b: (c: {1}))) +> (a+>: (b+>: (c|: {2}, d|: {1}), e+>: (f: 1)), g+>: (h: 1))`,
	)

	// chain of +>
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2})), e: (f: 1))`,
		`(a: (b: (c: {1}))) +> (a+>: (b+>: (c|: {2}))) +> (e+>: (f: (g: 1))) +> (e+>: (f: 1))`,
	)
}

func TestOpsNestedAddArrowMixNestedOp(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2}, d: {2})), e: 2)`,
		`(a: (b: (c: {1}, d: {1})), e: 1) +> (a+>: (b+>: (c|: {2}, d: {2})), e: 2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {1, 2}, 'd': {2}}}, 'e': 2}`,
		`{'a': {'b': {'c': {1}, 'd': {1}}}, 'e': 1} +> {'a'+>: {'b'+>: {'c'|: {2}, 'd': {2}}}, 'e': 2}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {2}, d: {2})), e: 2)`,
		`(a: (b: (c: {1}, d: {1})), e: 1) +> (a+>: (b: (c: {2}, d: {2})), e: 2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {2}, 'd': {2}}}, 'e': 2}`,
		`{'a': {'b': {'c': {1}, 'd': {1}}}, 'e': 1} +> {'a'+>: {'b': {'c': {2}, 'd': {2}}}, 'e': 2}`,
	)
}

func TestOpsNestedAddArrowNestedOpsWithGap(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2})))`,
		`(a: (b: (c: {1}, d: {1}))) +> (a+>: (b: (c|: {2})))`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {1, 2}}}}`,
		`{'a': {'b': {'c': {1}, 'd': {1}}}} +> {'a'+>: {'b': {'c'|: {2}}}}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: [1, 2])))`,
		`(a: (b: (c: [1], d: [1]))) +> (a+>: (b: (c++: [2])))`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': [1, 2]}}}`,
		`{'a': {'b': {'c': [1], 'd': [1]}}} +> {'a'+>: {'b': {'c'++: [2]}}}`,
	)
}

func TestOpsNestedAddArrowWithMissingKeys(t *testing.T) {
	t.Parallel()
	// missing at the start
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: 1)))`,
		`() +> (a+>: (b+>: (c+>: 1)))`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: 1, d: 1), e: 1), f: 1)`,
		`() +> (a+>: (b+>: (c+>: 1, d: 1), e: 1), f: 1)`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': 1}}}`,
		`{} +> {'a'+>: {'b'+>: {'c'+>: 1}}}`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': 1, 'd': 1}, 'e': 1}, 'f': 1}`,
		`{} +> {'a'+>: {'b'+>: {'c'+>: 1, 'd': 1}, 'e': 1}, 'f': 1}`,
	)
	// missing at the end
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: 1, d: {1})))`,
		`(a: (b: (c: 1))) +> (a+>: (b+>: (d|: {1})))`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': 1, 'd': {1}}}}`,
		`{'a': {'b': {'c': 1}}} +> {'a'+>: {'b'+>: {'d'|: {1}}}}`,
	)
	// missing everywhere
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2}, d: {1}), e: {1}), f: {1})`,
		`(a: (b: (c: {1}))) +> (a+>: (b+>: (c|: {2}, d+>: {1}), e|: {1}), f|: {1})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {1, 2}, 'd': 1}, 'e': {1}}, 'f': {1}}`,
		`{'a': {'b': {'c': {1}}}} +> {'a'+>: {'b'+>: {'c'|: {2}, 'd'+>: 1}, 'e'|: {1}}, 'f'|: {1}}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: {'b': (c: {1, 2}, d: {1}), 'e': {1}}, f: {1})`,
		`(a: {'b': (c: {1})}) +> (a+>: {'b'+>: (c|: {2}, d|: {1}), 'e'|: {1}}, f|: {1})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': (b: {'c': {1, 2}, 'd': {1}}, e: {1}), 'f': {1}}`,
		`{'a': (b: {'c': {1}})} +> {'a'+>: (b+>: {'c'|: {2}, 'd'|: {1}}, e|: {1}), 'f'|: {1}}`,
	)
	// missing everywhere with nested ops gaps
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: {1, 2}, d: {1}), e: {1}), f: {1})`,
		`(a: (b: (c: {1}, d: 2))) +> (a+>: (b: (c|: {2}, d: {1}), e|: {1}), f|: {1})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {'b': {'c': {1, 2}, 'd': {1}}, 'e': {1}}, 'f': {1}}`,
		`{'a': {'b': {'c': {1}, 'd': 2}}} +> {'a'+>: {'b': {'c'|: {2}, 'd': {1}}, 'e'|: {1}}, 'f'|: {1}}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: {'b': (c: {1, 2}, d: {1}), 'e': {1}}, f: {1})`,
		`(a: {'b': (c: {1}, d: 2)}) +> (a+>: {'b': (c|: {2}, d: {1}), 'e'|: {1}}, f|: {1})`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': (b: {'c': {1, 2}, 'd': {1}}, e: {1}), 'f': {1}}`,
		`{'a': (b: {'c': {1}, 'd': 2})} +> {'a'+>: (b: {'c'|: {2}, 'd': {1}}, e|: {1}), 'f'|: {1}}`,
	)
}

func TestOpsNestedAddArrowMixedWithOtherExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`(a: (b: (c: 2)), x: 2)`,
		`
			(a: (b: (c: 1))) +>
				let x = 2;
				(
					a+>: (b+>: (c: x)),
					:x,
				)
		`,
	)
}

func TestOpsAddArrowForError(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, "", `{1, 2, 3} +> {4, 5, 6}`)
	AssertCodeErrors(t, "", `1 +> 4`)
	AssertCodeErrors(t,
		"attr/key operation only allowed in rhs of a merge operation: a+>: (b: 1)",
		`(a+>: (b: 1)) +> (a+>: (b: 2))`,
	)
	AssertCodeErrors(t,
		"attr/key operation only allowed in rhs of a merge operation: 'a'+>: {'b': 1}",
		`{'a'+>: {'b': 1}} +> {'a'+>: {'b': 2}}`,
	)
	AssertCodeErrors(t,
		"| lhs must be a set, not number",
		`(a: (b: 1)) +> (a+>: (b|: 2))`,
	)
	AssertCodeErrors(t,
		"| lhs must be a set, not number",
		`{'a': {'b': 1}} +> {'a'+>: {'b'|: 2}}`,
	)
	AssertCodeErrors(t,
		"++ lhs must be a set, not number",
		`(a: (b: 1)) +> (a+>: (b++: 2))`,
	)
	AssertCodeErrors(t,
		"++ lhs must be a set, not number",
		`{'a': {'b': 1}} +> {'a'+>: {'b'++: 2}}`,
	)
	AssertCodeErrors(t,
		"attr name must be explicitly defined for attr operation: "+
			"\n\x1b[1;37m:3:19:\x1b[0m\n\n\t\tlet x = (a: (b: 1));\n\t\t(a: (b: 2)) +> (\x1b[1;31m+>:x.a.b\x1b[0m)\n\t\t",
		`
		let x = (a: (b: 1));
		(a: (b: 2)) +> (+>:x.a.b)
		`,
	)
}
