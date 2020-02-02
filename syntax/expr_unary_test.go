package syntax

import (
	"testing"
)

func TestCountExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `{41, 42, 43} count`)
}

func TestPowerSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{{}}`, `^{}`)
	AssertCodesEvalToSameValue(t, `{{}, {1}}`, `^{1}`)
	AssertCodesEvalToSameValue(t, `{{}, {1}, {2}, {1, 2}}`, `^{1, 2}`)
	AssertCodesEvalToSameValue(t,
		`{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}}`,
		`^{1, 2, 3}`)
}
