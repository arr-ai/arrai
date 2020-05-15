package syntax

import "testing"

// TestStrJoin, joiner is string.
func TestStrJoin(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"AB"`, `//seq.join("",['A','B'])                         `)
	AssertCodesEvalToSameValue(t, `""                `, `//seq.join(",",[])                         `)
	AssertCodesEvalToSameValue(t, `",,"              `, `//seq.join(",",["", "", ""])               `)
	AssertCodesEvalToSameValue(t, `"this is a test"  `, `//seq.join(" ",["this", "is", "a", "test"])`)
	AssertCodesEvalToSameValue(t, `"this"            `, `//seq.join(",",["this"])                   `)
	AssertCodesEvalToSameValue(t, `"You and me"`, `//seq.join(" and ",["You", "me"])`)
	assertExprPanics(t, `//seq.join("this", 2)`)
}

func TestArrayJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	// AssertCodesEvalToSameValue(t, `["You", "me"]`, `//seq.join("",["You", "me"])`)
	// AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join("",[1,2])`)

	// AssertCodesEvalToSameValue(t, `["A","B"]`, `//seq.join([],["A","B"])`)
	// AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// // if joinee is empty, the final value will be empty
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	// AssertCodesEvalToSameValue(t, `["A",",","B"]`, `//seq.join([","],["A","B"])`)
	// AssertCodesEvalToSameValue(t, `[1,0,2,0,3,0,4,0,5]`, `//seq.join([0], [1,2,3,4,5])`)
	// // TODO
	// //AssertCodesEvalToSameValue(t, `[1, 2, 0, 3, 4, 0, 5, 6]`, `//seq.join([0], [[1, 2], [3, 4], [5, 6]])`)
	// AssertCodesEvalToSameValue(t, `['A','A','B','A','C','A','D']`, `//seq.join(['A'], ['A','B','C','D'])`)
}

func TestBytesJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	// AssertCodesEvalToSameValue(t, `{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) }`,
	// 	`//seq.join("",//unicode.utf8.encode('hello'))`)
	// AssertCodesEvalToSameValue(t, `{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) }`,
	// 	`//seq.join([],{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) })`)
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	// AssertCodesEvalToSameValue(t, `["A","B"]`, `//seq.join([],["A","B"])`)
	// AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// // if joinee is empty, the final value will be empty
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	// AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('h'))`)
}
