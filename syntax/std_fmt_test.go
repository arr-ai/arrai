package syntax

import (
	"github.com/alecthomas/assert"
	"github.com/arr-ai/arrai/rel"
	"testing"
)

func TestFmtPrettyDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `"{
  'a': 1,
  'b': 2,
  'c': 3
}"`, `//fmt.pretty({'b':2, 'a':1,'c':3})`)

	AssertCodesEvalToSameValue(t, `"{
  'a': 1,
  'b': 2,
  'c': {
    'd': 11,
    'e': 22
  }
}"`, `//fmt.pretty({'b':2,'a':1,'c':{'d':11,'e':22}})`)
}

func TestFmtPrettyTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"(
  a: 1,
  b: 2,
  c: 3
)"`, `//fmt.pretty((a:1,c:3,b:2))`)
	AssertCodesEvalToSameValue(t, `"(
  a: 1,
  b: 2,
  c: (
    d: 11,
    e: 22
  )
)"`, `//fmt.pretty((a:1,b:2,c:(d:11,e:22)))`)
}

func TestFmtPrettyArraySimple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"[1, 2, 3]"`, `//fmt.pretty([1,2,3])`)
	AssertCodesEvalToSameValue(t, `"[1, 2, ['a', 'b', 'c']]"`, `//fmt.pretty([1,2,["a",'b',  "c"]])`)
}

func TestFmtPrettyArrayComplex(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"[
  (
    a: 1,
    b: 2
  ),
  3,
  {
    'c': 4,
    [5, 6]: [7, 8]
  }
]"`, `//fmt.pretty([(a:1, b:2),3,{"c":4,[5,6]:[7,8]}])`)
}

func TestFmtPrettySet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"{
  1,
  2,
  3
}"`, `//fmt.pretty({1, 2, 3})`)
	AssertCodesEvalToSameValue(t, `"{
  1,
  2,
  {
    3,
    4,
    5
  }
}"`, `//fmt.pretty({1, 2, {3, 4, 5}})`)
}

func TestIsSimple(t *testing.T) {
	assert.True(t, isSimple(rel.NewString([]rune("a"))))
	assert.True(t, isSimple(rel.NewNumber(12345)))
	assert.True(t, isSimple(rel.NewArray(rel.NewNumber(1))))
	assert.True(t, isSimple(rel.NewSet()))
	assert.True(t, isSimple(rel.NewSet(rel.NewString([]rune("a")), rel.NewNumber(12345))))
	assert.True(t, isSimple(rel.NewDict(false)))

	d := rel.NewDict(false, rel.NewDictEntryTuple(rel.NewString([]rune("a")), rel.NewNumber(1)))
	assert.False(t, isSimple(d))
	assert.False(t, isSimple(rel.NewTuple(rel.NewAttr("a", rel.NewNumber(1)))))
}
