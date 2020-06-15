package syntax

import "testing"

func TestStrSub(t *testing.T) { //nolint:dupl
	t.Parallel()

	AssertCodesEvalToSameValue(t, `" BC"`, `//seq.sub( "A", " ","ABC")`)
	AssertCodesEvalToSameValue(t, `"this is not a test"`, `//seq.sub("aaa", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//seq.sub("is not", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//seq.sub("not ", "","this is not a test")`)
	AssertCodesEvalToSameValue(t, `"t1his is not1 a t1est1"`, `//seq.sub("t", "t1","this is not a test")`)
	AssertCodesEvalToSameValue(t, `"this is still a test"`,
		`//seq.sub( "doesn't matter", "hello there","this is still a test")`)
	AssertCodeErrors(t, `//seq.sub("hello there", "test", 1)`, "")
	/////////////////
	AssertCodesEvalToSameValue(t, `""`, `//seq.sub( "","", "")`)
	AssertCodesEvalToSameValue(t, `"A"`, `//seq.sub( "","A", "")`)
	AssertCodesEvalToSameValue(t, `""`, `//seq.sub( "A","", "")`)
	AssertCodesEvalToSameValue(t, `"ABC"`, `//seq.sub( "","", "ABC")`)
	AssertCodesEvalToSameValue(t, `"EAEBECE"`, `//seq.sub( "", "E","ABC")`)
	AssertCodesEvalToSameValue(t, `"BC"`, `//seq.sub( "A", "","ABC")`)

	AssertCodeErrors(t, `//seq.sub(1,'B','BCD')`, "")
}

func TestArraySub(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `['T', 'B', 'T', 'C', 'D', 'E']`,
		`//seq.sub(['A'], ['T'], ['A', 'B', 'A', 'C', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t, `[['A', 'B'], ['T','C'],['A','D']]`,
		`//seq.sub([['A','C']], [['T','C']], [['A', 'B'], ['A','C'],['A','D']])`)
	AssertCodesEvalToSameValue(t, `[2, 2, 3]`, `//seq.sub([1], [2], [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[1,1], [4,4], [3,3]]`, `//seq.sub([[2,2]], [[4,4]], [[1,1], [2,2], [3,3]])`)

	AssertCodeErrors(t, `//seq.sub(1,'B',[1,2,3])`, "")
	AssertCodeErrors(t, `//seq.sub(1,'B',['A','B','C'])`, "")
}

func TestArraySubEdgeCases(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[]`, `//seq.sub( [],[], [])`)
	AssertCodesEvalToSameValue(t, `[1]`, `//seq.sub( [],[1], [])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.sub( [],[1,2], [])`)
	AssertCodesEvalToSameValue(t, `[[1,2]]`, `//seq.sub( [],[[1,2]], [])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.sub( [1],[], [])`)
	AssertCodesEvalToSameValue(t, `[1,2,3]`, `//seq.sub( [],[], [1,2,3])`)
	AssertCodesEvalToSameValue(t, `[[1,2],3]`, `//seq.sub( [],[], [[1,2],3])`)
	AssertCodesEvalToSameValue(t, `[4,1,4,2,4,3,4]`, `//seq.sub( [], [4],[1,2,3])`)
	AssertCodesEvalToSameValue(t, `[4,[1,2],4,[3,4],4]`, `//seq.sub( [], [4],[[1,2],[3,4]])`)
	AssertCodesEvalToSameValue(t, `[[4],[1,2],[4],[3,4],[4]]`, `//seq.sub( [], [[4]],[[1,2],[3,4]])`)
	AssertCodesEvalToSameValue(t, `[1,3]`, `//seq.sub([2], [],[1,2,3])`)
}

func TestBytesSub(t *testing.T) {
	t.Parallel()

	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `<<'oello'>>`, `//seq.sub(<<104>>,<<111>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `<<'hehho'>>`, `//seq.sub(<<108>>,<<104>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `<<>>       `, `//seq.sub(<<>>,<<>>,<<>>)             `)
	AssertCodesEvalToSameValue(t, `<<>>       `, `//seq.sub(<<'a'>>,<<>>,<<>>)          `)
	AssertCodesEvalToSameValue(t, `<<'a'>>    `, `//seq.sub(<<>>,<<'a'>>,<<>>)          `)

	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.sub(<<>>,<<>>,<<'hello'>>)   `)
	AssertCodesEvalToSameValue(t, `<<'ello'>> `, `//seq.sub(<<'h'>>,<<>>,<<'hello'>>)`)

	AssertCodesEvalToSameValue(t, `<<'thtetltltot'>>`, `//seq.sub(<<>>,<<'t'>>,<<'hello'>>)`)
}
