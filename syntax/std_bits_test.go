package syntax

import "testing"

func TestBitsSetInt(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{42, 56} `, `//bits.set(72061992084439040)`)
	AssertCodesEvalToSameValue(t, `{1, 3, 5}`, `//bits.set(42)               `)
	AssertCodesEvalToSameValue(t, `{7}      `, `//bits.set(128)              `)
	AssertCodesEvalToSameValue(t, `{0}      `, `//bits.set(1)                `)
	AssertCodesEvalToSameValue(t, `{}       `, `//bits.set(0)                `)

	assertExprPanics(t, `//bits.set({})`)
	assertExprPanics(t, `//bits.set(-1)`)
}

func TestBitsMask(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `72061992084439040`, `//bits.mask({42, 56}) `)
	AssertCodesEvalToSameValue(t, `42               `, `//bits.mask({1, 3, 5})`)
	AssertCodesEvalToSameValue(t, `128              `, `//bits.mask({7})      `)
	AssertCodesEvalToSameValue(t, `1                `, `//bits.mask({0})      `)
	AssertCodesEvalToSameValue(t, `0                `, `//bits.mask({})       `)

	assertExprPanics(t, `//bits.mask({"number"})`)
	assertExprPanics(t, `//bits.mask(3)         `)
}
