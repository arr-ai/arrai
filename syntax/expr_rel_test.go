package syntax

import "testing"

func TestConcat(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"abcdef"`, `"abc" ++ "def"`)
	AssertCodesEvalToSameValue(t, `   "def"`, `""    ++ "def"`)
	AssertCodesEvalToSameValue(t, `"abc"   `, `"abc" ++ ""   `)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 4, 5, 6]`, `[1, 2, 3] ++ [4, 5, 6]`)
	AssertCodesEvalToSameValue(t, `[         4, 5, 6]`, `[       ] ++ [4, 5, 6]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3         ]`, `[1, 2, 3] ++ [       ]`)
}

// A sparse array's element count skips its holes, and a non-zero-offset
// array's count doesn't reflect its indices at all, so Concatenate's generic
// fallback (used whenever the fast contiguous-zero-offset path doesn't
// apply) can't use a.Count() as the position to start placing b's elements:
// it silently collided with / reordered around a's own elements instead of
// appending strictly after them.
func TestConcatSparseAndOffsetArrays(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[1, , 3, 4, 5]`, `[1, , 3] ++ [4, 5]`)
	AssertCodesEvalToSameValue(t, `5\[1, 2, 3, 4, 5]`, `5\[1, 2, 3] ++ [4, 5]`)
}
