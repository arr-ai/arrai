package syntax

import "testing"

func TestStrJoin(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""`, `//seq.join(",",[])`)

	AssertCodesEvalToSameValue(t, `"this is a test" `, `//seq.join(" ",["this", "is", "a", "test"])`)
	AssertCodesEvalToSameValue(t, `"this"`, `//seq.join(",",["this"])`)
	AssertCodesEvalToSameValue(t, `"You and me"`, `//seq.join(" and ",["You", "me"])`)
	AssertCodesEvalToSameValue(t, `"AB"`, `//seq.join("",['A','B'])`)

	AssertCodesEvalToSameValue(t, `"Youme"`, `//seq.join("",["You", "me"])`)
	AssertCodesEvalToSameValue(t, `""`, `//seq.join(",",[])`)
	AssertCodesEvalToSameValue(t, `",,"`, `//seq.join(",",["", "", ""])`)

	// It is not supported
	// AssertCodesEvalToSameValue(t, `""`, `//seq.join("",[])`)

	AssertCodeErrors(t, "", `//seq.join("this", 2)`)
}

//nolint:dupl
func TestArrayJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `[1, 2, 0, 3, 4, 0, 5, 6]`, `//seq.join([0], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[2, [3, 4], 0, 5, 6]`, `//seq.join([0], [[2, [3, 4]], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 10, 11, 3, 4, 10, 11, 5, 6]`,
		`//seq.join([10,11], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], [3, 4], 0, [5, 6], [7, 8]]`,
		`//seq.join([0], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])`)

	AssertCodesEvalToSameValue(t, `[1, 2, [10], [11], 3, 4, [10], [11], 5, 6]`,
		`//seq.join([[10],[11]], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], [3, 4], [0], [1], [5, 6], [7, 8]]`,
		`//seq.join([[0],[1]], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])`)

	AssertCodesEvalToSameValue(t, `['AA', 'AB', 'BB', 'AB', 'CC', 'DD']`,
		`//seq.join(['AB'], [['AA'], ['BB'], ['CC' , 'DD']])`)
	AssertCodesEvalToSameValue(t, `['AA', 'AB', 'BB', ['CC', 'DD']]`,
		`//seq.join(['AB'], [['AA'], ['BB' ,['CC' , 'DD']]])`)

	// Test cases the delimiter is []
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[[1],[2]])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[[1],[],[2]])`)
	AssertCodesEvalToSameValue(t, `[1,3,3,2]`, `//seq.join([3],[[1],[],[2]])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join([],[])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)

	AssertCodesEvalToSameValue(t, `[1, 2, 3, 4]`, `//seq.join([], [[1, 2], [3, 4]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], 3, 4]`, `//seq.join([], [[[1, 2]], [3, 4]])`)
	AssertCodesEvalToSameValue(t, `[[1, 2], [3, 4], 5]`, `//seq.join([], [[[1, 2]], [[3,4], 5]])`)

	AssertCodeErrors(t, "", `//seq.join(1, [1,2,3,4,5])`)
	AssertCodeErrors(t, "", `//seq.join('A', [1,2])`)
	AssertCodeErrors(t, "", `//seq.join([],[1,2])`)
	AssertCodeErrors(t, "", `//seq.join([1],[1,2])`)
	AssertCodeErrors(t, "", `//seq.join([0], [1,2,3,4,5])`)
	AssertCodeErrors(t, "", `//seq.join(['A'], ['AT','BB', 'CD'])`)
}

func TestBytesJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `<<'hhehlhlho'>>    `, `//seq.join(<<104>>,<<'hello'>>) `)
	AssertCodesEvalToSameValue(t, `<<'hateatlatlato'>>`, `//seq.join(<<'at'>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t,
		`<<104, 108, 111, 101, 108, 111, 108, 108, 111, 108, 108, 111, 111>>`,
		`//seq.join(<<108, 111>>, <<104, 101, 108, 108, 111>>)`)

	AssertCodesEvalToSameValue(t, `<<>>       `, `//seq.join(<<>>,<<>>)       `)
	AssertCodesEvalToSameValue(t, `<<'hello'>>`, `//seq.join(<<>>,<<'hello'>>)`)
	AssertCodesEvalToSameValue(t, `<<>>       `, `//seq.join(<<'h'>>,<<>>)    `)
}
