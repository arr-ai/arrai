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

	AssertCodeErrors(t, `//seq.has_prefix(1,"ABC")`, "")
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

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([1, 2],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([[1, 2]],[[1, 2], [3]])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],[])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([],['A','B','C','D','E'])`)

	AssertCodeErrors(t, `//seq.has_prefix(1,[1,2,3])`, "")
	AssertCodeErrors(t, `//seq.has_prefix('A',['A','B','C'])`, "")
}

func TestBytesPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.has_prefix(<<'hello'>>,<<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<'h'>>,<<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<'he'>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'e'>>,<<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'l'>>,<<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'o'>>,<<'hello'>>) `)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'h'>>,<<'hello'>>)    `)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'he'>>,<<'hello'>>)   `)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'hello'>>,<<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<>>,<<>>)       `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'o'>>,<<>>)    `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<>>,<<'hello'>>)`)
}
