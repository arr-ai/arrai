package syntax

import "testing"

func TestExprSetWith(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{{}}`, `{} with {}`)
	AssertCodesEvalToSameValue(t, `{{}}`, `{{}} with {}`)
	AssertCodesEvalToSameValue(t, `{{}, 1}`, `{1} with {}`)
	AssertCodesEvalToSameValue(t, `{{}, 1}`, `{1, {}} with {}`)
	AssertCodesEvalToSameValue(t, `{1}`, `{} with 1`)
	AssertCodesEvalToSameValue(t, `{1}`, `{1} with 1`)
	AssertCodesEvalToSameValue(t, `{1, {}}`, `{{}} with 1`)
	AssertCodesEvalToSameValue(t, `{1, {}}`, `{{}, 1} with 1`)
}

func TestExprSetWithout(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `{} without {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{{}} without {}`)
	AssertCodesEvalToSameValue(t, `{1}`, `{1} without {}`)
	AssertCodesEvalToSameValue(t, `{1}`, `{1, {}} without {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} without 1`)
	AssertCodesEvalToSameValue(t, `{{}}`, `{{}} without 1`)
	AssertCodesEvalToSameValue(t, `{}`, `{1} without 1`)
	AssertCodesEvalToSameValue(t, `{{}}`, `{1, {}} without 1`)
}
