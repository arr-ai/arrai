package syntax

import "testing"

func TestSeqConcat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""`, `//.seq.concat([])`)
	AssertCodesEvalToSameValue(t, `"abc"`, `//.seq.concat(["abc"])`)
	AssertCodesEvalToSameValue(t, `"abc"`, `//.seq.concat(["abc", ""])`)
	AssertCodesEvalToSameValue(t, `"def"`, `//.seq.concat(["", "def"])`)
	AssertCodesEvalToSameValue(t, `"abcdef"`, `//.seq.concat(["abc", "def"])`)
	AssertCodesEvalToSameValue(t, `"abcdefghi"`, `//.seq.concat(["abc", "def", "ghi"])`)
}

func TestSeqRepeat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""`, `//.seq.repeat(0, "ab")`)
	AssertCodesEvalToSameValue(t, `""`, `//.seq.repeat(3, "")`)
	AssertCodesEvalToSameValue(t, `"ab"`, `//.seq.repeat(1, "ab")`)
	AssertCodesEvalToSameValue(t, `"ababab"`, `//.seq.repeat(3, "ab")`)

	AssertCodesEvalToSameValue(t, `[]`, `//.seq.repeat(0, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[]`, `//.seq.repeat(2, [])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `//.seq.repeat(1, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 1, 2, 3]`, `//.seq.repeat(2, [1, 2, 3])`)

	AssertCodePanics(t, `//.seq.repeat(2, 3.4)`)
}
