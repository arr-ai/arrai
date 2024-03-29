package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

func TestFmtPrettyDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `"{
  'a': 1,
  'b': 2,
  'c': 3,
}"`, `//fmt.pretty({'b':2, 'a':1,'c':3})`)

	AssertCodesEvalToSameValue(t, `"{
  'a': 1,
  'b': 2,
  'c': {
    'd': 11,
    'e': 22,
  },
}"`, `//fmt.pretty({'b':2,'a':1,'c':{'d':11,'e':22}})`)

	AssertCodesEvalToSameValue(t, `"{
  42: 1,
  '42': 2,
}"`, `//fmt.pretty({'42':2,42:1})`)
}

func TestFmtPrettyTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"()"`, `//fmt.pretty(( ))`)
	AssertCodesEvalToSameValue(t, `"(
  a: 1,
  b: 2,
  c: 3,
)"`, `//fmt.pretty((a:1,c:3,b:2))`)
	AssertCodesEvalToSameValue(t, `"(
  a: 1,
  b: 2,
  c: (
    d: 11,
    e: 22,
  ),
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
    b: 2,
  ),
  3,
  {
    'c': 4,
    [5, 6]: [7, 8],
  },
]"`, `//fmt.pretty([(a:1, b:2),3,{"c":4,[5,6]:[7,8]}])`)

	AssertCodesEvalToSameValue(t, `"[
  (
    a: 1,
  ),
]"`, `//fmt.pretty([(a:1)])`)
}

func TestFmtPrettySet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"{}"`, `//fmt.pretty({ })`)
	AssertCodesEvalToSameValue(t, `"{1, 3, '2'}"`, `//fmt.pretty({3,1,'2'})`)
	AssertCodesEvalToSameValue(t, `"{1, 2, {3, 4, 5}}"`, `//fmt.pretty({1,2,{3,4,5}})`)
}

func TestFmtPrettyUnionSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"{
  (
    a: 'abc',
  ),
  (
    b: 2,
  ),
}"`,
		`//fmt.pretty({(a: 'abc'), (b: 2)})`,
	)
	AssertCodesEvalToSameValue(t,
		`"{
  1,
  2,
  3,
  'a',
  [1, 2, 3],
  (
    a: 1,
  ),
}"`,
		`//fmt.pretty({(a: 1), 1, 2, 3, "a", [1, 2, 3]})`,
	)
	AssertCodesEvalToSameValue(t,
		`"{
  (
    a: (
      b: (
        c: [
          (
            d: 1,
          ),
        ],
      ),
    ),
  ),
  (
    e: {
      (
        a: {
          (
            a: 1,
          ),
          (
            b: 1,
          ),
        },
      ),
      (
        b: 2,
      ),
    },
  ),
}"`,
		`//fmt.pretty({(a: (b: (c: [(d: 1)]))), (e: {(a: {(a: 1), (b: 1)}), (b: 2)})})`,
	)
}

func TestFmtPrettyString(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"{}"`, `//fmt.pretty('')`)
	AssertCodesEvalToSameValue(t, `"'abc'"`, `//fmt.pretty('abc')`)
	AssertCodesEvalToSameValue(t, `"42\\'abc'"`, `//fmt.pretty(42\'abc')`)
}

func TestFmtPrettyRelation(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`"{
  |a, b, c|
  (1, 1, 1),
}"`,
		`//fmt.pretty({(a: 1, b: 1, c: 1)})`,
	)

	AssertCodesEvalToSameValue(t,
		`"{
  |a, b, c|
  (1, 1, 1),
}"`,
		`//fmt.pretty({|a, b, c| (1, 1, 1)})`,
	)

	AssertCodesEvalToSameValue(t,
		`"{
  |a, b, c|
  (1, 1, 1),
  (2, 2, 2),
  (3, 3, 3),
}"`,
		`//fmt.pretty({|a, b, c| (1, 1, 1), (2, 2, 2), (3, 3, 3)})`,
	)

	AssertCodesEvalToSameValue(t,
		`"{
  |a, b, c|
  (
    1,
    1,
    {
      |a|
      (1),
    },
  ),
  (
    1,
    1,
    {
      |a|
      (
        (
          b: 1,
        ),
      ),
    },
  ),
  (
    1,
    [
      (
        a: 1,
      ),
    ],
    1,
  ),
  (
    (
      b: 1,
      c: 1,
    ),
    1,
    1,
  ),
}"`,
		`//fmt.pretty({|a, b, c| ((b: 1, c: 1), 1, 1), (1, [(a: 1)], 1), (1, 1, {(a: 1)}), (1, 1, {(a: (b:1))})})`,
	)
}

func TestIsSimple(t *testing.T) {
	assert.True(t, isSimple(rel.NewString([]rune("a"))))
	assert.True(t, isSimple(rel.NewNumber(12345)))
	assert.True(t, isSimple(rel.NewArray(rel.NewNumber(1))))
	assert.True(t, isSimple(rel.None))
	assert.True(t, isSimple(rel.MustNewSet(rel.NewString([]rune("a")), rel.NewNumber(12345))))
	assert.True(t, isSimple(rel.MustNewDict(false)))

	d := rel.MustNewDict(false, rel.NewDictEntryTuple(rel.NewString([]rune("a")), rel.NewNumber(1)))
	assert.False(t, isSimple(d))
	tp := rel.NewTuple(rel.NewAttr("a", rel.NewNumber(1)))
	assert.False(t, isSimple(tp))

	sb := rel.NewSetBuilder()
	sb.Add(tp)
	r, err := sb.Finish()
	require.NoError(t, err)
	require.IsType(t, rel.Relation{}, r)
	assert.False(t, isSimple(r))
}
