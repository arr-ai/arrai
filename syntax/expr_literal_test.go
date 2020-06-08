package syntax

import "testing"

func TestExprRelationLiteral(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{(a: "hello", b:1), (a: "world", b:2)}`, `{|a,b| ("hello",1), ("world",2)}`)
}
