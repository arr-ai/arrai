package syntax

import (
	"testing"
)

func TestWhereExpr(t *testing.T) {
	t.Parallel()
	s := `{|a,b| (3,41), (2,42), (1,43)}`
	// defer trace().revert()
	AssertCodesEvalToSameValue(t, `{(a:3, b:41)}`, s+` where .a=3`)
}

func TestRelationCall(t *testing.T) {
	t.Parallel()
	s := `{"key": "val"}("key")`
	AssertCodesEvalToSameValue(t, `"val"`, s)
}

func TestOpsAddArrow(t *testing.T) {
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2) +> (c: 3, d: 4)`, `(a: 1, b: 2, c: 3, d: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2) +> (b: 3, c: 4)`, `(a: 1, b: 3, c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 2)) +> (b: 3, c: 4)`, `(a: 1, b: 3, c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 2)) +> (b: (c: 4), c: 4)`, `(a: 1, b: (c: 4), c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 2)) +> (a: (b: 1), b: (c: 4))`, `(a: (b: 1), b: (c: 4))`)
}
