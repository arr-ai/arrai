package syntax

import "testing"

func TestXString(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"a42z"`, `$"a:{6*7}:z"`)
	AssertCodesEvalToSameValue(t, `"a00042z"`, `$"a:{05d:6*7}:z"`)
	AssertCodesEvalToSameValue(t, `"a001, 002, 003z"`, `$"a:{03d*:[1, 2, 3]:, }:z"`)
	AssertCodesEvalToSameValue(t, `"a42k3.142z"`, `$"a:{6*7}:k:{.3f://.math.pi}:z"`)
}
