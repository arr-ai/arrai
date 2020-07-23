package syntax

import (
	"testing"
)

func TestFmtPretty(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"(\n  a: 1,\n  b: 2,\n  c: 3\n)"`,
		`//fmt.pretty((a:1,b:2,c:3))`)
}
