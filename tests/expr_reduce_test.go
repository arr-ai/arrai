package tests

import (
	"testing"
)

var s = `{ |a,b| (3,41), (2,42), (1,43)}`
var u = `{ |a,b| (3,41), (2,42), (1,43), (0, 46)}`

func TestSumExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `6`, s+` sum .a`)
	AssertCodesEvalToSameValue(t, `126`, s+` sum .b`)
}

func TestMaxExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, s+` max .a`)
	AssertCodesEvalToSameValue(t, `-1`, s+` max -.a`)
	AssertCodesEvalToSameValue(t, `43`, s+` max .b`)
}

func TestMeanExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`, s+` mean .a`)
	AssertCodesEvalToSameValue(t, `-2`, s+` mean -.a`)
	AssertCodesEvalToSameValue(t, `42`, s+` mean .b`)
	AssertCodesEvalToSameValue(t, `43`, u+` mean .b`)
}

func TestMinExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, s+` min .a`)
	AssertCodesEvalToSameValue(t, `-3`, s+` min -.a`)
	AssertCodesEvalToSameValue(t, `41`, s+` min .b`)
}

func TestMedianExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`, s+` median .a`)
	AssertCodesEvalToSameValue(t, `-2`, s+` median -.a`)
	AssertCodesEvalToSameValue(t, `42`, s+` median .b`)
	AssertCodesEvalToSameValue(t, `42.5`, u+` median .b`)
	AssertCodesEvalToSameValue(t, `43`, s+` where .a=1 median .b`)
	AssertCodesEvalToSameValue(t, `42.5`, s+` where .a<3 median .b`)
}
