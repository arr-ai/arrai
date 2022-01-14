package syntax

import "testing"

func TestExprSetWith(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `{} without {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{{}} without {}`)
	AssertCodesEvalToSameValue(t, `{1}`, `{1} without {}`)
	AssertCodesEvalToSameValue(t, `{1}`, `{1, {}} without {}`)
	AssertCodesEvalToSameValue(t, `{}`, `{} without 1`)
	AssertCodesEvalToSameValue(t, `{{}}`, `{{}} without 1`)
	AssertCodesEvalToSameValue(t, `{}`, `{1} without 1`)
	AssertCodesEvalToSameValue(t, `{{}}`, `{1, {}} without 1`)
}

func TestExprStrWithout(t *testing.T) {
	t.Parallel()

	// working cases
	AssertCodesEvalToSameValue(t, `1\'bc'                       `, `'abc'   without ('a' single)`)
	AssertCodesEvalToSameValue(t, `2\'bc'                       `, `1\'abc' without ((1\'a') single)`)
	AssertCodesEvalToSameValue(t, `'ab'                         `, `'abc'   without ((2\'c') single)`)
	AssertCodesEvalToSameValue(t, `1\'ab'                       `, `1\'abc' without ((3\'c') single)`)
	AssertCodesEvalToSameValue(t, `{|@, @char| (0, 97), (2, 99)}`, `'abc'   without ((1\'b') single)`)
	AssertCodesEvalToSameValue(t, `{|@, @char| (1, 97), (3, 99)}`, `1\'abc' without ((2\'b') single)`)
	AssertCodesEvalToSameValue(t, `{}                           `, `'a'     without ('a' single)`)
	AssertCodesEvalToSameValue(t, `{}                           `, `1\'a'   without ((1\'a') single)`)

	// test without missing character
	AssertCodesEvalToSameValue(t, `'abc'                        `, `'abc'   without ('d' single)`)
	AssertCodesEvalToSameValue(t, `'abc'                        `, `'abc'   without ((2\'d') single)`)
	AssertCodesEvalToSameValue(t, `'abc'                        `, `'abc'   without ((1\'d') single)`)

	// test without missing index
	AssertCodesEvalToSameValue(t, `'abc'                        `, `'abc'   without ((5\'a') single)`)
}
