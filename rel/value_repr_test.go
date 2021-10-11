package rel

import (
	"testing"

	"github.com/arr-ai/arrai/pkg/fu"

	"github.com/stretchr/testify/assert"
)

func TestReprBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `true`, fu.Repr(NewBool(true)))
	assert.Equal(t, `{}`, fu.Repr(NewBool(false)))
}

func TestReprString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, fu.Repr(NewString([]rune(""))))
	assert.Equal(t, `'abc'`, fu.Repr(NewString([]rune("abc"))))
	assert.Equal(t, `1\'abc'`, fu.Repr(NewOffsetString([]rune("abc"), 1)))
}

func TestReprBytes(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `<<0, 1, 42, 2, 255>>`, fu.Repr(NewBytes([]byte{0, 1, 42, 2, 255})))
}

func TestReprStringCharTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @char: 97)`, fu.Repr(NewStringCharTuple(1, 'a')))
}

func TestReprArrayItemTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @item: 'a')`, fu.Repr(NewArrayItemTuple(1, NewString([]rune("a")))))
}

func TestReprByteItemTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 0, @byte: 97)`, fu.Repr(NewBytesByteTuple(0, 97)))
}

func TestReprDictEntryTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @value: 'a')`, fu.Repr(NewDictEntryTuple(NewNumber(1), NewString([]rune("a")))))
}

func TestReprSet(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, fu.Repr(None))
	assert.Equal(t, `{1}`, fu.Repr(MustNewSet(NewNumber(1))))
	assert.Equal(t, `{1, 'abc'}`, fu.Repr(MustNewSet(NewNumber(1), NewString([]rune("abc")))))
}

func TestReprArray(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, fu.Repr(NewArray()))
	assert.Equal(t, `[1, 'abc']`, fu.Repr(NewArray(NewNumber(1), NewString([]rune("abc")))))
	assert.Equal(t, `[1, , 'abc']`, fu.Repr(NewArray(NewNumber(1), nil, NewString([]rune("abc")))))
}

func TestReprDict(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{1: 2, {42, 'a'}: 43}`, fu.Repr(MustNewDict(false,
		NewDictEntryTuple(NewNumber(1), NewNumber(2)),
		NewDictEntryTuple(MustNewSet(NewString([]rune("a")), NewNumber(42)), NewNumber(43)),
	)))
}

func TestReprTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		`(a: (b: 'foo'), c: ['bar'], "d'": {})`,
		fu.Repr(NewTuple(
			NewAttr("a", NewTuple(NewAttr("b", NewString([]rune("foo"))))),
			NewAttr("c", NewArray(NewString([]rune("bar")))),
			NewAttr("d'", None),
		)))
}
