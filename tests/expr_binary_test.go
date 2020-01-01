package tests

import (
	"testing"
)

func TestWhereExpr(t *testing.T) {
	t.Parallel()
	s := `{ |a,b| (3,41), (2,42), (1,43) }`
	AssertCodesEvalToSameValue(t, `{(a:3,b:41)}`, s+` where .a=3`)
}
