package syntax

import "testing"

func TestStrJoin(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"AB"`, `//seq.join("",['A','B'])                         `)
	AssertCodesEvalToSameValue(t, `"Youme"`, `//seq.join("",["You", "me"])`)
	AssertCodesEvalToSameValue(t, `"AB"`, `//seq.join([],["A","B"])`)
	AssertCodesEvalToSameValue(t, `"AB"`, `//seq.join([],['A','B'])`)
	AssertCodesEvalToSameValue(t, `"this is a test"  `, `//seq.join(" ",["this", "is", "a", "test"])`)
	AssertCodesEvalToSameValue(t, `"this"            `, `//seq.join(",",["this"])                   `)
	AssertCodesEvalToSameValue(t, `"You and me"`, `//seq.join(" and ",["You", "me"])`)
	assertExprPanics(t, `//seq.join("this", 2)`)

	// Following cases are not supported to make sure code is clear and simple.
	// Or it has to check array element is rel.String or rel.Number.
	// And they are imapcted by https://github.com/arr-ai/arrai/issues/268 too.
	// AssertCodesEvalToSameValue(t, `""                `, `//seq.join(",",[])                         `)
	// AssertCodesEvalToSameValue(t, `",,"              `, `//seq.join(",",["", "", ""])               `)
}

func TestArrayJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	AssertCodesEvalToSameValue(t, `[1,0,2,0,3,0,4,0,5]`, `//seq.join([0], [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 0, 3, 4, 0, 5, 6]`, `//seq.join([0], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[2, [3, 4], 0, 5, 6]`, `//seq.join([0], [[2, [3, 4]], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 10, 11, 3, 4, 10, 11, 5, 6]`,
		`//seq.join([10,11], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], [3, 4], 0, [5, 6], [7, 8]]`,
		`//seq.join([0], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])`)

	AssertCodesEvalToSameValue(t, `[1, 2, [10], [11], 3, 4, [10], [11], 5, 6]`,
		`//seq.join([[10],[11]], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], [3, 4], [0], [1], [5, 6], [7, 8]]`,
		`//seq.join([[0],[1]], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])`)

	// Following cases are not supported to make sure code is clear and simple.
	// Or it has to check array element is rel.String or rel.Number.
	// And they are imapcted by https://github.com/arr-ai/arrai/issues/268 too.
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
}

func TestBytesJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('hhehlhlho')`,
		`//seq.join({ |@, @byte| (0, 104)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 104), (1, 108), (2, 111), (3, 101), (4, 108), (5, 111),`+
			`(6, 108), (7, 108), (8, 111), (9, 108), (10, 108), (11, 111), (12, 111) }`,
		`//seq.join({ |@, @byte| (0, 108), (1, 111)},{ |@, @byte| (0, 104), (1, 101),`+
			` (2, 108), (3, 108), (4, 111) })`)
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('hateatlatlato')`,
		`//seq.join(//unicode.utf8.encode('at'),//unicode.utf8.encode('hello'))`)
}