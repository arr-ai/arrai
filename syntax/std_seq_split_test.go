package syntax

import "testing"

func TestStrSplit(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[[], "B", "CD"]     `, `//seq.split("A","ABACD")`)
	AssertCodesEvalToSameValue(t, `["ABAC", []]     `, `//seq.split("D","ABACD")`)

	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//seq.split(" ","this is a test") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//seq.split(",","this is a test") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//seq.split("is","this is a test")`)
	AssertCodeErrors(t, `//seq.split(1, "this is a test")`, "")

	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//seq.split("","this is a test")`)

	// As https://github.com/arr-ai/arrai/issues/268, `{}`, `[]` and `""` means empty set in arr.ai
	// And the intent for //seq.split is to return an array, so it should be expressed as such.
	// `""` -> empty string, `[]` -> empty array and `{}` -> empty set
	AssertCodesEvalToSameValue(t, `[]`, `//seq.split("","") `)

	AssertCodesEvalToSameValue(t, `[""]`, `//seq.split(",","") `)

	AssertCodeErrors(t, `//seq.split(1,"ABC")`, "")
}

func TestArraySplit(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[['A'], ['B']]`,
		`//seq.split([],['A', 'B'])`)
	AssertCodesEvalToSameValue(t, `[[], ['B']]`,
		`//seq.split(['A'],['A', 'B'])`)
	AssertCodesEvalToSameValue(t, `[['A'], []]`,
		`//seq.split(['B'],['A', 'B'])`)
	AssertCodesEvalToSameValue(t, `[[],['B'],['C', 'D', 'E']]`,
		`//seq.split(['A'],['A', 'B', 'A', 'C', 'D', 'E'])`)

	AssertCodesEvalToSameValue(t, `[['B'],['C'], ['D', 'E']]`,
		`//seq.split(['A'],['B', 'A', 'C', 'A', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t, `[['A', 'B', 'C']]`,
		`//seq.split(['F'],['A', 'B', 'C'])`)
	AssertCodesEvalToSameValue(t, `[[['A','B'], ['C','D'], ['E','F']]]`,
		`//seq.split([['F','F']],[['A','B'], ['C','D'], ['E','F']])`)
	AssertCodesEvalToSameValue(t, `[[['A','B']], [['E','F']]]`,
		`//seq.split([['C','D']],[['A','B'], ['C','D'], ['E','F']])`)
	AssertCodesEvalToSameValue(t, `[[['A','B']], [['E','F'],['G']]]`,
		`//seq.split([['C','D']],[['A','B'], ['C','D'], ['E','F'], ['G']])`)

	AssertCodesEvalToSameValue(t, `[[['A','B']], [['G']]]`,
		`//seq.split([['C','D'],['E','F']],[['A','B'], ['C','D'], ['E','F'], ['G']])`)
	AssertCodesEvalToSameValue(t, `[[['A','B'], ['C','D'], ['E','F'], ['G']]]`,
		`//seq.split([['C','D'],['E','T']],[['A','B'], ['C','D'], ['E','F'], ['G']])`)

	AssertCodesEvalToSameValue(t, `[[],[2,3]]`, `//seq.split([1],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[1,2],[]]`, `//seq.split([3],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[1],[3]]`, `//seq.split([2],[1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[[[1,2]],[[5,6]]]`, `//seq.split([[3,4]],[[1,2],[3,4],[5,6]])`)
	AssertCodesEvalToSameValue(t, `[[[1,2]], [[3,4]]]`, `//seq.split([],[[1,2], [3,4]])`)
	AssertCodesEvalToSameValue(t, `[['A'],['B'],['A']]`, `//seq.split([],['A', 'B', 'A'])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.split([],[])`)
	AssertCodesEvalToSameValue(t, `[[]]`, `//seq.split(['A'],[])`)

	AssertCodeErrors(t, `//seq.split(1,[1,2,3])`, "")
	AssertCodeErrors(t, `//seq.split('A',['A','B'])`, "")
}

func TestBytesSplit(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `[<<'y'>>,<<'e'>>,<<'s'>>]`, `//seq.split(<<>>,<<"yes">>)`)
	AssertCodesEvalToSameValue(t,
		`[<<"this">>, <<"is">>, <<"a">>, <<"test">>]`,
		`//seq.split(<<" ">>, <<"this is a test">>)`)
	AssertCodesEvalToSameValue(t, `[<<"this is a test">>]`, `//seq.split(<<"get">>, <<"this is a test">>)`)

	AssertCodesEvalToSameValue(t, `[[], <<"B">>, <<"CD">>]`, `//seq.split(<<"A">>,<<"ABACD">>)`)
	AssertCodesEvalToSameValue(t, `[<<"ABAC">>, []]`, `//seq.split(<<"D">>,<<"ABACD">>)       `)

	AssertCodesEvalToSameValue(t, `<<>>                     `, `//seq.split(<<>>,<<>>)     `)
	AssertCodesEvalToSameValue(t, `[<<"A">>,<<"B">>,<<"C">>]`, `//seq.split(<<>>,<<"ABC">>)`)
	AssertCodesEvalToSameValue(t, `[<<>>]                   `, `//seq.split(<<",">>,<<>>)  `)

	AssertCodeErrors(t, `//seq.split(",", <<"hello">>)`,
		"unexpected panic: delimiter and subject have to be of the same type, "+
			"currently: delimiter: rel.String, subject: rel.Bytes")
}
