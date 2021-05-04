package syntax

import "testing"

func TestExprSetWithoutIntoEmpty(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({1} without 1)`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ([1] without (@: 0, @item: 1))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = (<<'a'>> without (@: 0, @byte: 97))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ('a' without (@: 0, @char: 97))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({0: 1} without (@: 0, @value: 1))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({} without {})`)
}

func TestExprSetWhereIntoEmpty(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({1} where \. . != 1)`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ([1] where \. . != (@: 0, @item: 1))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = (<<'a'>> where \. . != (@: 0, @byte: 97))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ('a' where \. . != (@: 0, @char: 97))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({0: 1} where \. . != (@: 0, @value: 1))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({} where \. . != {})`)
	AssertCodesEvalToSameValue(t, sTrue, `{} = ({(@: 0, @byte: 97), (@: 0, @char: 97)} where false)`)
}
