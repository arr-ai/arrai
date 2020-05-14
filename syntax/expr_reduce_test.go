package syntax

import (
	"testing"
)

var testRel1 = `{|a,b| (3,41), (2,42), (1,43)}`
var testRel2 = `{|a,b| (3,41), (2,42), (1,43), (0, 46)}`

func TestSumExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `6`, testRel1+` sum .a`)
	AssertCodesEvalToSameValue(t, `126`, testRel1+` sum .b`)
}

func TestMaxExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, testRel1+` max .a`)
	AssertCodesEvalToSameValue(t, `-1`, testRel1+` max -.a`)
	AssertCodesEvalToSameValue(t, `43`, testRel1+` max .b`)
}

func TestMeanExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`, testRel1+` mean .a`)
	AssertCodesEvalToSameValue(t, `-2`, testRel1+` mean -.a`)
	AssertCodesEvalToSameValue(t, `42`, testRel1+` mean .b`)
	AssertCodesEvalToSameValue(t, `43`, testRel2+` mean .b`)
}

func TestMinExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, testRel1+` min .a`)
	AssertCodesEvalToSameValue(t, `-3`, testRel1+` min -.a`)
	AssertCodesEvalToSameValue(t, `41`, testRel1+` min .b`)
}

func TestMedianExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`, testRel1+` median .a`)
	AssertCodesEvalToSameValue(t, `-2`, testRel1+` median -.a`)
	AssertCodesEvalToSameValue(t, `42`, testRel1+` median .b`)
	AssertCodesEvalToSameValue(t, `42.5`, testRel2+` median .b`)
	AssertCodesEvalToSameValue(t, `43`, testRel1+` where .a=1 median .b`)
	AssertCodesEvalToSameValue(t, `42.5`, testRel1+` where .a<3 median .b`)
}
