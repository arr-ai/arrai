package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
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

	AssertCodesEvalToSameValue(t, `"{
  42: 1,
  '42': 2
}"`, `//fmt.pretty({'42':2,42:1})`)
}

func TestFmtPrettyTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"()"`, `//fmt.pretty(( ))`)
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
	AssertCodesEvalToSameValue(t, `"{}"`, `//fmt.pretty([ ])`)
	AssertCodesEvalToSameValue(t, `"[1, 2, 3]"`, `//fmt.pretty([1,2,3])`)
	AssertCodesEvalToSameValue(t, `"[1, , 3]"`, `//fmt.pretty(9\[1,,3])`)
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

	AssertCodesEvalToSameValue(t, `"[
  (
    a: 1
  )
]"`, `//fmt.pretty([(a:1)])`)
}

func TestFmtPrettySet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"{}"`, `//fmt.pretty({ })`)
	AssertCodesEvalToSameValue(t, `"{1, 3, '2'}"`, `//fmt.pretty({3,1,'2'})`)
	AssertCodesEvalToSameValue(t, `"{1, 2, {3, 4, 5}}"`, `//fmt.pretty({1,2,{3,4,5}})`)
}

func TestFmtPrettyString(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"{}"`, `//fmt.pretty('')`)
	AssertCodesEvalToSameValue(t, `"'abc'"`, `//fmt.pretty('abc')`)
	AssertCodesEvalToSameValue(t, `"42\\'abc'"`, `//fmt.pretty(42\'abc')`)
}

func TestIsSimple(t *testing.T) {
	assert.True(t, isSimple(rel.NewString([]rune("a"))))
	assert.True(t, isSimple(rel.NewNumber(12345)))
	assert.True(t, isSimple(rel.NewArray(rel.NewNumber(1))))
	assert.True(t, isSimple(rel.MustNewSet()))
	assert.True(t, isSimple(rel.MustNewSet(rel.NewString([]rune("a")), rel.NewNumber(12345))))
	assert.True(t, isSimple(rel.MustNewDict(false)))

	d := rel.MustNewDict(false, rel.NewDictEntryTuple(rel.NewString([]rune("a")), rel.NewNumber(1)))
	assert.False(t, isSimple(d))
	assert.False(t, isSimple(rel.NewTuple(rel.NewAttr("a", rel.NewNumber(1)))))
}
