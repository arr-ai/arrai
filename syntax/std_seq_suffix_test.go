package syntax

import "testing"

func TestStrHasSuffix(t *testing.T) {
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

func TestArrayHasSuffix(t *testing.T) {
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

func TestBytesHasSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(<<'hello'>>, <<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<'o'>>, <<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<'lo'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'e'>>, <<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'ell'>>, <<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'h'>>, <<'hello'>>)  `)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<>>, <<>>)       `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(<<'o'>>, <<>>)    `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_suffix(<<>>, <<'hello'>>)`)
}

func TestStrTrimSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_suffix("ABCDE","ABCDE")`)

	AssertCodesEvalToSameValue(t, `"ABCD"`, `//seq.trim_suffix("E","ABCDE")`)
	AssertCodesEvalToSameValue(t, `"ABC"`, `//seq.trim_suffix("DE","ABCDE")`)
	AssertCodesEvalToSameValue(t, `"ABCDE"`, `//seq.trim_suffix("CD", "ABCDE")`)
	AssertCodesEvalToSameValue(t, `"ABCDE"`, `//seq.trim_suffix("D","ABCDE")`)

	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_suffix("D","")`)
	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_suffix("","")`)
	AssertCodesEvalToSameValue(t, `"ABCDE"`, `//seq.trim_suffix("","ABCDE")`)

	AssertCodeErrors(t, `//seq.trim_suffix(1,"ABC")`, "")
}

func TestArrayTrimSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_suffix(['A','B'],['A','B'])`)

	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_suffix([],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `['A','B','C','D']`, `//seq.trim_suffix(['E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `['A','B','C']`, `//seq.trim_suffix( ['D','E'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_suffix(['D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_suffix(['C','D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t,
		`['A','B','C','D','E']`,
		`//seq.trim_suffix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `[1 ,2]  `, `//seq.trim_suffix([3, 4],[1 ,2, 3, 4])`)
	AssertCodesEvalToSameValue(t, `[[1 ,2]]`, `//seq.trim_suffix([3, 4],[[1 ,2], 3, 4])`)
	AssertCodesEvalToSameValue(t, `[]      `, `//seq.trim_suffix([3, 4],[3, 4])`)
	AssertCodesEvalToSameValue(t, `[[1, 2]]`, `//seq.trim_suffix([[3, 4]],[[1 ,2], [3, 4]])`)

	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_suffix([],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_suffix([], [])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_suffix(['D','E'],[])`)

	AssertCodeErrors(t, `//seq.trim_suffix(1,[1,2])`, "")
	AssertCodeErrors(t, `//seq.trim_suffix('A',['A','B'])`, "")
}

func TestBytesTrimSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `<<>>`, `//seq.trim_suffix(<<'hello'>>, <<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `<<'hell'>>`, `//seq.trim_suffix(<<'o'>>, <<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `<<'hel'>>`, `//seq.trim_suffix(<<'lo'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_suffix(<<'e'>>, <<'hello'>>)  `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_suffix(<<'ell'>>, <<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_suffix(<<'h'>>, <<'hello'>>)  `)

	AssertCodesEvalToSameValue(t, `<<>>`, `//seq.trim_suffix(<<>>, <<>>)       `)
	AssertCodesEvalToSameValue(t, `<<>>`, `//seq.trim_suffix(<<'o'>>, <<>>)    `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_suffix(<<>>, <<'hello'>>)`)
}
