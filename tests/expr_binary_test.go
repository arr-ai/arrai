package tests

import (
	"testing"
)

func TestWhereExpr(t *testing.T) {
	s := `{||a,b| {3,41}, {2,42}, {1,43}|}`
	assertCodesEvalToSameValue(t, `{|{a:3,b:41}|}`, s+` where .a=3`)
}
