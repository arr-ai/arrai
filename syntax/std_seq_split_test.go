package syntax

import "testing"

func TestStrSplit(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//seq.split(" ","this is a test") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//seq.split(",","this is a test") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//seq.split("is","this is a test")`)
	assertExprPanics(t, `//seq.split(1, "this is a test")`)

	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//seq.split("","this is a test")`)

	// As https://github.com/arr-ai/arrai/issues/268, `{}`, `[]` and `""` means empty set in arr.ai
	// And the intent for //seq.split is to return an array, so it should be expressed as such.
	// `""` -> empty string, `[]` -> empty array and `{}` -> empty set
	AssertCodesEvalToSameValue(t, `[]`, `//seq.split("","") `)

	AssertCodesEvalToSameValue(t, `[""]`, `//seq.split(",","") `)

	assertExprPanics(t, `//seq.split(1,"ABC")`)
}

func TestArraySplit(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[['B'],['C', 'D', 'E']]`,
		`//seq.split(['A'],['A', 'B', 'A', 'C', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t, `[['B'],['C'], ['D', 'E']]`,
		`//seq.split(['A'],['B', 'A', 'C', 'A', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t, `[['A', 'B', 'C']]`,
		`//seq.split(['F'],['A', 'B', 'C'])`)
	AssertCodesEvalToSameValue(t, `[[['A','B'], ['C','D'], ['E','F']]]`,
		`//seq.split([['F','F']],[['A','B'], ['C','D'], ['E','F']])`)

	AssertCodesEvalToSameValue(t, `[[1],[3]]`, `//seq.split([2],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[[1,2]],[[5,6]]]`, `//seq.split([[3,4]],[[1,2],[3,4],[5,6]])`)
	AssertCodesEvalToSameValue(t, `[[[1,2]], [[3,4]]]`, `//seq.split([],[[1,2], [3,4]])`)
	AssertCodesEvalToSameValue(t, `[['A'],['B'],['A']]`, `//seq.split([],['A', 'B', 'A'])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.split([],[])`)
	AssertCodesEvalToSameValue(t, `[[]]`, `//seq.split(['A'],[])`)

	assertExprPanics(t, `//seq.split(1,[1,2,3])`)
	assertExprPanics(t, `//seq.split('A',['A','B'])`)
}

func TestBytesSplit(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t,
		`[//unicode.utf8.encode('y'),//unicode.utf8.encode('e'),//unicode.utf8.encode('s')]`,
		`//seq.split(//unicode.utf8.encode(""),//unicode.utf8.encode("yes"))`)
	AssertCodesEvalToSameValue(t,
		`[//unicode.utf8.encode("this"), //unicode.utf8.encode("is"),`+
			` //unicode.utf8.encode("a"), //unicode.utf8.encode("test")]`,
		`//seq.split(//unicode.utf8.encode(" "),//unicode.utf8.encode("this is a test"))`)
	AssertCodesEvalToSameValue(t,
		`[//unicode.utf8.encode("this is a test")]`,
		`//seq.split(//unicode.utf8.encode("get"),//unicode.utf8.encode("this is a test"))`)

	AssertCodesEvalToSameValue(t,
		`//unicode.utf8.encode("")`,
		`//seq.split(//unicode.utf8.encode(""),//unicode.utf8.encode(""))`)
	AssertCodesEvalToSameValue(t,
		`[//unicode.utf8.encode("A"),//unicode.utf8.encode("B"),//unicode.utf8.encode("C")]`,
		`//seq.split(//unicode.utf8.encode(""),//unicode.utf8.encode("ABC"))`)
	AssertCodesEvalToSameValue(t,
		`[//unicode.utf8.encode("")]`,
		`//seq.split(//unicode.utf8.encode(","),//unicode.utf8.encode(""))`)
}
