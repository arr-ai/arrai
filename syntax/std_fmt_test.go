package syntax

import (
	"testing"
)

func TestFmtPrettyDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty({'a':1,'b':2,'c':3})`)
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty({'a':1,'b':2,'c':{'d':11,'e':22}})`)
}

func TestFmtPrettyTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty((a:1,b:2,c:3))`)
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty((a:1,b:2,c:(d:11,e:22)))`)
}

func TestFmtPrettySet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty({1, 2, 3})`)
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty({1, 2, {3, 4, 5}})`)
}
