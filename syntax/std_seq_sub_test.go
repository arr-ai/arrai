package syntax

import "testing"

func TestStrSub(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`" BC"`,
		`//seq.sub( "A", " ","ABC")`)
	AssertCodesEvalToSameValue(t,
		`"this is not a test"`,
		`//seq.sub("aaa", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("is not", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("not ", "","this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"t1his is not1 a t1est1"`,
		`//seq.sub("t", "t1","this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is still a test"`,
		`//seq.sub( "doesn't matter", "hello there","this is still a test")`)
	assertExprPanics(t, `//seq.sub("hello there", "test", 1)`)
	/////////////////
	AssertCodesEvalToSameValue(t,
		`""`,
		`//seq.sub( "","", "")`)
	AssertCodesEvalToSameValue(t,
		`"A"`,
		`//seq.sub( "","A", "")`)
	AssertCodesEvalToSameValue(t,
		`""`,
		`//seq.sub( "A","", "")`)

	AssertCodesEvalToSameValue(t,
		`"ABC"`,
		`//seq.sub( "","", "ABC")`)
	AssertCodesEvalToSameValue(t,
		`"EAEBECE"`,
		`//seq.sub( "", "E","ABC")`)
	AssertCodesEvalToSameValue(t,
		`"BC"`,
		`//seq.sub( "A", "","ABC")`)

	assertExprPanics(t, `//seq.sub(1,'B','BCD')`)
}

func TestArraySub(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`['T', 'B', 'T', 'C', 'D', 'E']`,
		`//seq.sub(['A'], ['T'], ['A', 'B', 'A', 'C', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t,
		`[['A', 'B'], ['T','C'],['A','D']]`,
		`//seq.sub([['A','C']], [['T','C']], [['A', 'B'], ['A','C'],['A','D']])`)
	AssertCodesEvalToSameValue(t,
		`[2, 2, 3]`,
		`//seq.sub([1], [2], [1, 2, 3])`)
	AssertCodesEvalToSameValue(t,
		`[[1,1], [4,4], [3,3]]`,
		`//seq.sub([[2,2]], [[4,4]], [[1,1], [2,2], [3,3]])`)

	assertExprPanics(t, `//seq.sub(1,'B',[1,2,3])`)
	assertExprPanics(t, `//seq.sub(1,'B',['A','B','C'])`)
}

func TestArraySubEdgeCases(t *testing.T) {
	/////////////////
	AssertCodesEvalToSameValue(t,
		`[]`,
		`//seq.sub( [],[], [])`)
	AssertCodesEvalToSameValue(t,
		`[1]`,
		`//seq.sub( [],[1], [])`)
	AssertCodesEvalToSameValue(t,
		`[1,2]`,
		`//seq.sub( [],[1,2], [])`)
	AssertCodesEvalToSameValue(t,
		`[[1,2]]`,
		`//seq.sub( [],[[1,2]], [])`)
	AssertCodesEvalToSameValue(t,
		`[]`,
		`//seq.sub( [1],[], [])`)

	AssertCodesEvalToSameValue(t,
		`[1,2,3]`,
		`//seq.sub( [],[], [1,2,3])`)
	AssertCodesEvalToSameValue(t,
		`[4,1,4,2,4,3,4]`,
		`//seq.sub( [], [4],[1,2,3])`)
	AssertCodesEvalToSameValue(t,
		`[4,[1,2],4,[3,4],4]`,
		`//seq.sub( [], [4],[[1,2],[3,4]])`)
	AssertCodesEvalToSameValue(t,
		`[[4],[1,2],[4],[3,4],[4]]`,
		`//seq.sub( [], [[4]],[[1,2],[3,4]])`)
	AssertCodesEvalToSameValue(t,
		`[1,3]`,
		`//seq.sub([2], [],[1,2,3])`)
}

func TestBytesSub(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('oello')`,
		`//seq.sub({ |@, @byte| (0, 104)},{ |@, @byte| (0, 111)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('hehho')`,
		`//seq.sub({ |@, @byte| (0, 108)},{ |@, @byte| (0, 104)},//unicode.utf8.encode('hello'))`)
	///////////////////
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('')`,
		`//seq.sub(//unicode.utf8.encode(''),//unicode.utf8.encode(''),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('')`,
		`//seq.sub(//unicode.utf8.encode('a'),//unicode.utf8.encode(''),//unicode.utf8.encode(''))`)
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('a')`,
		`//seq.sub(//unicode.utf8.encode(''),//unicode.utf8.encode('a'),//unicode.utf8.encode(''))`)

	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('hello')`,
		`//seq.sub(//unicode.utf8.encode(''),//unicode.utf8.encode(''),//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('ello')`,
		`//seq.sub(//unicode.utf8.encode('h'),//unicode.utf8.encode(''),//unicode.utf8.encode('hello'))`)

	AssertCodesEvalToSameValue(t, `//unicode.utf8.encode('thtetltltot')`,
		`//seq.sub(//unicode.utf8.encode(''),//unicode.utf8.encode('t'),//unicode.utf8.encode('hello'))`)
}
