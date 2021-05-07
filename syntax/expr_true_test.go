package syntax

import "testing"

func TestExprSetIntoTrue(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, sTrue, `({(), 1} without 1)                                   `)
	AssertCodesEvalToSameValue(t, sTrue, `({(), (@: 1, @value: 1)} without (@: 1, @value: 1))   `)
	AssertCodesEvalToSameValue(t, sTrue, `({(), 1} where . != 1)                                `)
	AssertCodesEvalToSameValue(t, sTrue, `({(), (@: 1, @value: 1)} where . != (@: 1, @value: 1))`)
	AssertCodesEvalToSameValue(t, sTrue, `{} with ()                                            `)
}

func TestExprTrueSetWith(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, sTrue, `true with ()`)
	AssertCodesEvalToSameValue(t, `{(), 'abc'}`, `true with 'abc'`)
}

func TestExprTrueSetMap(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, sTrue, `    true => ()               `)
	AssertCodesEvalToSameValue(t, `'a'    `, `true => (@: 0, @char: 97)`)
	AssertCodesEvalToSameValue(t, `<<'a'>>`, `true => (@: 0, @byte: 97)`)
	AssertCodesEvalToSameValue(t, `{0: ()}`, `true => (@: 0, @value: .)`)
	AssertCodesEvalToSameValue(t, `[()]   `, `true => (@: 0, @item: .) `)
}
