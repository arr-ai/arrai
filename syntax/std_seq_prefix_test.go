package syntax

import "testing"

func TestStrPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("A","ABCDE")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("AB","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("BCD","ABCDE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("CD","ABCDE")`)
}

func TestArrayPrefix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix('A',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix([],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix('B',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['B'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],[])`)
}

func TestBytesPrefix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('he'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('e'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('l'),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('o'),//unicode.utf8.encode('hello'))`)
}
