package syntax

import (
	"testing"
)

func TestSetPatternWithTupleSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, `cond {(y:1)} { {(:y)}: y*y }`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {(y:1)}{ {(:a)}: y*y }`)

	AssertCodesEvalToSameValue(t,
		`4`, `let m=(a: {(x: 3)}, b: {(x: 3)}); cond m { (a:{(:x)}, b: {(:x)}): x+1}`)
	AssertCodesEvalToSameValue(t,
		`{}`, `let m=(a: {(x: 3)}, b: {(x: 3)}); cond m { (a:{(:a)}, b: {(:x)}): x+1}`)

	AssertCodesEvalToSameValue(t, `{(2)}`, `cond {(y:1, z:2)} { {(:y, :z)}: {(y*z)} }`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {(y:1, z:2)} { {(:a, :z)}: {(y*z)} }`)

	AssertCodesEvalToSameValue(t,
		`2`, `let w = {(y:0, z:2), (y:1, z:3)} where .y=1; cond w { {(:y, ...)}: 2*y }`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {(y:0, z:2), (y:0, z:3)} { {(:y, :z)}: 2*y }`)

	// TODO: This should be an error
	AssertCodePanics(t, `let x = {(y:0, z:2), (y:0, z:3)}; cond x { {(:y, :z), ...}: 2*y }`)
}
