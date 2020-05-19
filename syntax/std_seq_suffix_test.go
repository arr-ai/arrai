package syntax

import "testing"

func TestStrSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("ABCDE","ABCDE")`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("E","ABCDE")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("DE","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix("CD", "ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix("D","ABCDE")`)

	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix("D","")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("","")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("","ABCDE")`)

	assertExprPanics(t, `//seq.has_suffix(1,"ABC")`)
}

func TestArraySuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['A','B'],['A','B'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix( ['D','E'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['C','D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['D','E'],[])`)

	assertExprPanics(t, `//seq.has_suffix(1,[1,2])`)
	assertExprPanics(t, `//seq.has_suffix('A',['A','B'])`)
}

func TestBytesSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode('o'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode('lo'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('e'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('ell'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode(''),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('o'),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode(''),//unicode.utf8.encode('hello'))`)
}
