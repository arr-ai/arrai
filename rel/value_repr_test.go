package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepr(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `{}`, Repr(NewString([]rune(""))))
	assert.Equal(t, `'abc'`, Repr(NewString([]rune("abc"))))

	assert.Equal(t, `{}`, Repr(NewSet()))
	assert.Equal(t, `{1}`, Repr(NewSet(NewNumber(1))))
	assert.Equal(t, `{1, 'abc'}`, Repr(NewSet(NewNumber(1), NewString([]rune("abc")))))

	assert.Equal(t, `[1, 'abc']`, Repr(NewArray(NewNumber(1), NewString([]rune("abc")))))
	assert.Equal(t, `{1: 2, {42, 'a'}: 43}`, Repr(NewDict(false,
		NewDictEntryTuple(NewNumber(1), NewNumber(2)),
		NewDictEntryTuple(NewSet(NewString([]rune("a")), NewNumber(42)), NewNumber(43)),
	)))

	assert.Equal(t, `(a: (b: 'foo'), c: ['bar'], d: {})`, Repr(NewTuple(
		NewAttr("a", NewTuple(NewAttr("b", NewString([]rune("foo"))))),
		NewAttr("c", NewArray(NewString([]rune("bar")))),
		NewAttr("d", None),
	)))
}
