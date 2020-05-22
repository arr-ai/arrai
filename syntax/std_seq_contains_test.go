package syntax

import "testing"

func TestStrContains(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains("A", "A")`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("", "this is a test")             `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("is a test", "this is a test")    `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("is not a test", "this is a test")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("a is", "this is a test")`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains("", "A")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("A", "")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains("", "")`)

	assertExprPanics(t, `//seq.contains(1, "ABC")`)
}

func TestArrayContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A'],['A', 'D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['E'],['A',C','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C'],[B','C','D'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['L','M','N'],['L','M','N','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B','C'],['T','B','C','X','Y'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C','D','E'],['1','3','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A1','B1','C1','D1','E1'],['A1','B1','C1','D1','E1'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['B1','C2','E3'],['A','B1','C2','D','E3'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A2','B3','C4','E5'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A1','B2','C3'], ['A', 'A1', 'B2','C3','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B4','C5'],['A', 'A', 'B4','C5','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([['B1','C1']],[['A', 'B'], ['B1','C1'],['D','E']])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([1,2], [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([1,2,3,4,5], [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([[1,2],[3,4],[5]], [[1,2],[3,4],[5]])`)

	AssertCodesEvalToSameValue(t, `false`, `//seq.contains([1], [])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([], [1])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([], [])`)

	assertExprPanics(t, `//seq.contains(1, [1,2,3,4,5])`)
	assertExprPanics(t, `//seq.contains('A',['A','B','C','D','E'])`)
}

func TestBytesContains(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode('hello'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 104)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 111)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 108),(0, 108)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `false`,
		`//seq.contains(//unicode.utf8.encode('A'),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode(''),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode(''),//unicode.utf8.encode('hello'))`)
}
