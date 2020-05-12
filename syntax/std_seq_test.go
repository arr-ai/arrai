package syntax

import "testing"

func TestSeqConcat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""              `, `//seq.concat([])                            `)
	AssertCodesEvalToSameValue(t, `""              `, `//seq.concat(["", "", "", ""])              `)
	AssertCodesEvalToSameValue(t, `"hello"         `, `//seq.concat(["", "", "", "", "hello", ""]) `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//seq.concat(["this", " is", " a", " test"])`)
	AssertCodesEvalToSameValue(t, `"this"          `, `//seq.concat(["this"])                      `)
	assertExprPanics(t, `//seq.concat("this")`)

	AssertCodesEvalToSameValue(t, `[]`, `//seq.concat([[]])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.concat([[], []])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `//seq.concat([[1, 2, 3]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `//seq.concat([[1, 2, 3], []])`)
	AssertCodesEvalToSameValue(t, `[4, 5, 6]`, `//seq.concat([[], [4, 5, 6]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 4, 5, 6]`, `//seq.concat([[1, 2, 3], [4, 5, 6]])`)
	AssertCodePanics(t, `//seq.concat(42)`)
}

func TestSeqRepeat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""      `, `//seq.repeat(0, "ab")`)
	AssertCodesEvalToSameValue(t, `""      `, `//seq.repeat(3, "")`)
	AssertCodesEvalToSameValue(t, `"ab"    `, `//seq.repeat(1, "ab")`)
	AssertCodesEvalToSameValue(t, `"ababab"`, `//seq.repeat(3, "ab")`)

	AssertCodesEvalToSameValue(t, `[]                `, `//seq.repeat(0, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[]                `, `//seq.repeat(2, [])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]         `, `//seq.repeat(1, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 1, 2, 3]`, `//seq.repeat(2, [1, 2, 3])`)

	AssertCodePanics(t, `//seq.repeat(2, 3.4)`)
}

func TestSeqContains(t *testing.T) {
	t.Parallel()
}

///////////////////
func TestStrSub(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//str.sub("this is not a test", "is not", "is")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//str.sub("this is not a test", "not ", "")`)
	AssertCodesEvalToSameValue(t,
		`"this is still a test"`,
		`//str.sub("this is still a test", "doesn't matter", "hello there")`)
	assertExprPanics(t, `//str.sub("hello there", "test", 1)`)
}

func TestStrSplit(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//str.split("this is a test", "")`)
	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//str.split("this is a test", " ") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//str.split("this is a test", ",") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//str.split("this is a test", "is")`)
	assertExprPanics(t, `//str.split("this is a test", 1)`)
}

func TestStrContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true `, `//str.contains("this is a test", "")             `)
	AssertCodesEvalToSameValue(t, `true `, `//str.contains("this is a test", "is a test")    `)
	AssertCodesEvalToSameValue(t, `false`, `//str.contains("this is a test", "is not a test")`)
	assertExprPanics(t, `//str.contains(123, 124)`)
}

func TestStrJoin(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""                `, `//str.join([], ",")                         `)
	AssertCodesEvalToSameValue(t, `",,"              `, `//str.join(["", "", ""], ",")               `)
	AssertCodesEvalToSameValue(t, `"this is a test"  `, `//str.join(["this", "is", "a", "test"], " ")`)
	AssertCodesEvalToSameValue(t, `"this"            `, `//str.join(["this"], ",")                   `)
	assertExprPanics(t, `//str.join("this", 2)`)
}
