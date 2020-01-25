package syntax

import "testing"

func TestPackagePi(t *testing.T) {
	AssertCodesEvalToSameValue(t, `0`, `//.math.pi*0`)
}
