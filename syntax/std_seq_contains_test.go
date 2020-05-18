package syntax

import "testing"

func TestStrContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains("", "A")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("", "")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("A", "")`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("", "this is a test")             `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("is a test", "this is a test")    `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("is not a test", "this is a test")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("a is", "this is a test")`)
}

func TestArrayContains(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(1, [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains([], [])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([], [1])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(1, [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(3, [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(5, [1,2,3,4,5])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('A',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('E',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('C',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C','D','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['B','C','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C'], ['A', 'A', 'B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B','C'],['A', 'A', 'B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([['B','C']],[['A', 'B'], ['B','C'],['D','E']])`)
}

func TestBytesContains(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 104)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 111)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 108),(0, 108)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)
}
