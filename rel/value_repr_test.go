package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReprBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `true`, Repr(NewBool(true)))
	assert.Equal(t, `{}`, Repr(NewBool(false)))
}

func TestReprString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, Repr(NewString([]rune(""))))
	assert.Equal(t, `'abc'`, Repr(NewString([]rune("abc"))))
	assert.Equal(t, `1\'abc'`, Repr(NewOffsetString([]rune("abc"), 1)))
}

func TestReprBytes(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "\n", Repr(NewBytes([]byte{'\n'})))
}

func TestReprStringCharTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @char: 97)`, Repr(NewStringCharTuple(1, 'a')))
}

func TestReprArrayItemTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @item: 'a')`, Repr(NewArrayItemTuple(1, NewString([]rune("a")))))
}

func TestReprDictEntryTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(@: 1, @value: 'a')`, Repr(NewDictEntryTuple(NewNumber(1), NewString([]rune("a")))))
}

func TestReprSet(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, Repr(NewSet()))
	assert.Equal(t, `{1}`, Repr(NewSet(NewNumber(1))))
	assert.Equal(t, `{1, 'abc'}`, Repr(NewSet(NewNumber(1), NewString([]rune("abc")))))
}

func TestReprArray(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, Repr(NewArray()))
	assert.Equal(t, `[1, 'abc']`, Repr(NewArray(NewNumber(1), NewString([]rune("abc")))))
	assert.Equal(t, `[1, , 'abc']`, Repr(NewArray(NewNumber(1), nil, NewString([]rune("abc")))))
}

func TestReprDict(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{1: 2, {42, 'a'}: 43}`, Repr(NewDict(false,
		NewDictEntryTuple(NewNumber(1), NewNumber(2)),
		NewDictEntryTuple(NewSet(NewString([]rune("a")), NewNumber(42)), NewNumber(43)),
	)))
}

func TestReprTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		`(a: (b: 'foo'), c: ['bar'], "d'": {})`,
		Repr(NewTuple(
			NewAttr("a", NewTuple(NewAttr("b", NewString([]rune("foo"))))),
			NewAttr("c", NewArray(NewString([]rune("bar")))),
			NewAttr("d'", None),
		)))
}

func TestReprUnknown(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() { Repr(nil) })
}
