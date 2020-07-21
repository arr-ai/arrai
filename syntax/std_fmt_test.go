package syntax

import (
	"testing"
)

func TestFmtPrettyDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 2, 'c': 3}`, `//fmt.pretty({'b':2, 'a':1,'c':3})`)
	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 2, 'c': {'d': 11, 'e': 22}}`, `//fmt.pretty({'b':2,'a':1,'c':{'d':11,'e':22}})`)
}

func TestFmtPrettyTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:1,b:2,c:3)`, `//fmt.pretty((a:1,b:2,c:3))`)
	AssertCodesEvalToSameValue(t, `(a:1,b:2,c:(d:11,e:22))`, `//fmt.pretty((a:1,b:2,c:(d:11,e:22)))`)
}

func TestFmtPrettySet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{1, 2, 3}`, `//fmt.pretty({1, 2, 3})`)
	AssertCodesEvalToSameValue(t, `{1, 2, {3, 4, 5}}`, `//fmt.pretty({1, 2, {3, 4, 5}})`)
}
