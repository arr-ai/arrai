package tests

import (
	"testing"
)

var s = `{||a,b| {3,41}, {2,42}, {1,43}|}`
var u = `{||a,b| {3,41}, {2,42}, {1,43}, {0, 46}|}`

func TestSumExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `6`, s+` sum .a`)
	assertCodesEvalToSameValue(t, `126`, s+` sum .b`)
}

func TestMaxExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `3`, s+` max .a`)
	assertCodesEvalToSameValue(t, `-1`, s+` max -.a`)
	assertCodesEvalToSameValue(t, `43`, s+` max .b`)
}

func TestMeanExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `2`, s+` mean .a`)
	assertCodesEvalToSameValue(t, `-2`, s+` mean -.a`)
	assertCodesEvalToSameValue(t, `42`, s+` mean .b`)
	assertCodesEvalToSameValue(t, `43`, u+` mean .b`)
}

func TestMinExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `1`, s+` min .a`)
	assertCodesEvalToSameValue(t, `-3`, s+` min -.a`)
	assertCodesEvalToSameValue(t, `41`, s+` min .b`)
}

func TestMedianExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `2`, s+` median .a`)
	assertCodesEvalToSameValue(t, `-2`, s+` median -.a`)
	assertCodesEvalToSameValue(t, `42`, s+` median .b`)
	assertCodesEvalToSameValue(t, `42.5`, u+` median .b`)
	assertCodesEvalToSameValue(t, `43`, s+` where .a=1 median .b`)
	assertCodesEvalToSameValue(t, `42.5`, s+` where .a<3 median .b`)
}
