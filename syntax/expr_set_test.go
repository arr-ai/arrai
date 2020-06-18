package syntax

import (
	"testing"
)

func TestSetCount(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `0`, `{} count`)
	AssertCodesEvalToSameValue(t, `1`, `{1} count`)
	AssertCodesEvalToSameValue(t, `2`, `{1, 2} count`)
	AssertCodesEvalToSameValue(t, `3`, `{1, 2, {3, 4}} count`)
}

func TestSetSingle(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `{42} single`)
	AssertCodeErrors(t, "", `{} single`)
	AssertCodeErrors(t, "", `{1, 2} single`)
	AssertCodeErrors(t, "", `{1, 2, {3, 4}} single`)
}
