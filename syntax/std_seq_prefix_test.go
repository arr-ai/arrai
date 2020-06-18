package syntax

import "testing"

func TestStrHasPrefix(t *testing.T) {
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

	AssertCodeErrors(t, "", `//seq.has_prefix(1,"ABC")`)
}

func TestArrayHasPrefix(t *testing.T) { //nolint:dupl
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

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([1, 2],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([[1, 2]],[[1, 2], [3]])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E','F'],[])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix([],['A','B','C','D','E'])`)

	AssertCodeErrors(t, "", `//seq.has_prefix(1,[1,2,3])`)
	AssertCodeErrors(t, "", `//seq.has_prefix('A',['A','B','C'])`)
}

func TestBytesHasPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'hello'>>, <<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<'h'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<'he'>>, <<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'e'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'l'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'o'>>, <<'hello'>>) `)

	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'h'>>, <<'hello'>>)    `)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'he'>>, <<'hello'>>)   `)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(<<'hello'>>, <<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<>>, <<>>)       `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(<<'o'>>, <<>>)    `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.has_prefix(<<>>, <<'hello'>>)`)
}

func TestStrTrimPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_prefix("ABC","ABC")`)

	AssertCodesEvalToSameValue(t, `"BCDE"`, `//seq.trim_prefix("A","ABCDE")`)
	AssertCodesEvalToSameValue(t, `"CDE"`, `//seq.trim_prefix("AB","ABCDE")`)
	AssertCodesEvalToSameValue(t, `"ABCDE"`, `//seq.trim_prefix("BCD","ABCDE")`)
	AssertCodesEvalToSameValue(t, `"ABCDE"`, `//seq.trim_prefix("CD","ABCDE")`)
	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_prefix("CD","")`)
	AssertCodesEvalToSameValue(t, `"ABCD"`, `//seq.trim_prefix("","ABCD")`)

	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_prefix("","")`)
	AssertCodesEvalToSameValue(t, `""`, `//seq.trim_prefix("A","")`)
	AssertCodesEvalToSameValue(t, `"A"`, `//seq.trim_prefix("","A")`)

	AssertCodeErrors(t, "", `//seq.trim_prefix(1,"ABC")`)
}

func TestArrayTrimPrefix(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_prefix(['A'],['A'])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_prefix(['A','B'],['A','B'])`)

	AssertCodesEvalToSameValue(t, `['B','C','D','E']`, `//seq.trim_prefix(['A']            ,['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `    ['C','D','E']`, `//seq.trim_prefix(['A','B']        ,['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `        ['D','E']`, `//seq.trim_prefix(['A','B','C']    ,['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `            ['E']`, `//seq.trim_prefix(['A','B','C','D'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_prefix(['B'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_prefix(['B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t,
		`['A','B','C','D','E']`,
		`//seq.trim_prefix(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `[3]`, `//seq.trim_prefix([1, 2],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[3]]`, `//seq.trim_prefix([[1, 2]],[[1, 2], [3]])`)

	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.trim_prefix(['A','B','C','D','E','F'],[])`)
	AssertCodesEvalToSameValue(t, `['A','B','C','D','E']`, `//seq.trim_prefix([],['A','B','C','D','E'])`)

	AssertCodeErrors(t, "", `//seq.trim_prefix(1,[1,2,3])`)
	AssertCodeErrors(t, "", `//seq.trim_prefix('A',['A','B','C'])`)
}

func TestBytesTrimPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `<<>>`, `//seq.trim_prefix(<<'hello'>>, <<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `<<'ello'>> `, `//seq.trim_prefix(<<'h'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `<<'llo'>> `, `//seq.trim_prefix(<<'he'>>, <<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_prefix(<<'e'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_prefix(<<'l'>>, <<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.trim_prefix(<<'o'>>, <<'hello'>>) `)

	AssertCodesEvalToSameValue(t, `<<>> `, `//seq.trim_prefix(<<>>, <<>>)       `)
	AssertCodesEvalToSameValue(t, `<<>>`, `//seq.trim_prefix(<<'o'>>, <<>>)    `)
	AssertCodesEvalToSameValue(t, `<<'hello'>> `, `//seq.trim_prefix(<<>>, <<'hello'>>)`)

	AssertCodeErrors(t, "", `//seq.trim_prefix(1, <<'hello'>>)`)
	AssertCodeErrors(t, "", `//seq.trim_prefix("he", <<'hello'>>)`)
}
