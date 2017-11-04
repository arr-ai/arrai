package tests

import (
	"testing"
)

func TestCountExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `3`, `{|41, 42, 43|} count`)
}
