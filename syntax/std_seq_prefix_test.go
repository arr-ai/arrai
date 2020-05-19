package syntax

import "testing"

func TestStrPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("ABC","ABC")`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("A","ABCDE")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("AB","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("BCD","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("CD","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("CD","")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("","ABCD")`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("","")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("A","")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("","A")`)

	assertExprPanics(t, `//seq.has_prefix(1,"ABC")`)
}

func TestArrayPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A'],['A'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B'],['A','B'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['B'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],[])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([],['A','B','C','D','E'])`)

	assertExprPanics(t, `//seq.has_prefix(1,[1,2,3])`)
	assertExprPanics(t, `//seq.has_prefix('A',['A','B','C'])`)
}

func TestBytesPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('he'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('e'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('l'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('o'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(//unicode.utf8.encode('he'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true `,
		`//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode(''),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('o'),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode(''),//unicode.utf8.encode('hello'))`)
}
