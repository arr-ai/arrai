package syntax

import "testing"

func TestDict(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{|@,@value| ("x", "y")}`, `{"x": "y"}`)
}
