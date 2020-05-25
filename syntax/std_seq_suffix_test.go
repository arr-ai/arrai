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

	AssertCodeErrors(t, `//seq.has_suffix(1,"ABC")`, "")
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

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([3, 4],[1 ,2, 3, 4])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([3, 4],[[1 ,2], 3, 4])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([3, 4],[3, 4])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([[3, 4]],[[1 ,2], [3, 4]])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['D','E'],[])`)

	AssertCodeErrors(t, `//seq.has_suffix(1,[1,2])`, "")
	AssertCodeErrors(t, `//seq.has_suffix('A',['A','B'])`, "")
}

func TestBytesSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.has_suffix(<<'hello'>>,<<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<'o'>>,<<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<'lo'>>,<<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'e'>>,<<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'ell'>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'h'>>,<<'hello'>>)  `)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<>>,<<>>)       `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'o'>>,<<>>)    `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<>>,<<'hello'>>)`)
}
