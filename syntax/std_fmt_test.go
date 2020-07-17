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

func TestFmtPrettyMixed(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty({'a':1,'b':2,'c':(d:11, e:22)})`)
	AssertCodesEvalToSameValue(t, `""`, `//fmt.pretty((a:1, b:2, c:{'d':11, 'e':22}))`)
}
