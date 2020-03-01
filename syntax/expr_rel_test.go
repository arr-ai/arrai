package syntax

import "testing"

func TestConcat(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"abcdef"`, `"abc" ++ "def"`)
	AssertCodesEvalToSameValue(t, `   "def"`, `""    ++ "def"`)
	AssertCodesEvalToSameValue(t, `"abc"   `, `"abc" ++ ""   `)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 4, 5, 6]`, `[1, 2, 3] ++ [4, 5, 6]`)
	AssertCodesEvalToSameValue(t, `[         4, 5, 6]`, `[       ] ++ [4, 5, 6]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3         ]`, `[1, 2, 3] ++ [       ]`)
}
