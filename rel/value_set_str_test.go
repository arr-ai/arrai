package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringTuple(t *testing.T) {
	s := "hello"
	for e := NewString([]rune(s)).Enumerator(); e.MoveNext(); {
		tuple, is := e.Current().(StringCharTuple)
		if assert.True(t, is) {
			assert.Equal(t, rune(s[tuple.at]), tuple.char)
		}
	}
}

func TestStringCall(t *testing.T) {
	t.Parallel()
	s := "hello"
	f := NewString([]rune(s))
	for i, c := range s {
		assert.Equal(t, c, rune(f.Call(NewNumber(float64(i))).(Number).Float64()))
	}

	assert.Panics(t, func() { f.Call(NewNumber(6)) })
	assert.Panics(t, func() { f.Call(NewNumber(-1)) })
}

func TestStringCallAll(t *testing.T) {
	t.Parallel()

	abc := NewString([]rune("abc"))

	AssertEqualValues(t, NewSet(NewNumber(float64('a'))), abc.CallAll(NewNumber(0)))
	AssertEqualValues(t, NewSet(NewNumber(float64('b'))), abc.CallAll(NewNumber(1)))
	AssertEqualValues(t, NewSet(NewNumber(float64('c'))), abc.CallAll(NewNumber(2)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(5)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(-1)))

	abc = NewOffsetString([]rune("abc"), -2)
	AssertEqualValues(t, NewSet(NewNumber(float64('a'))), abc.CallAll(NewNumber(-2)))
	AssertEqualValues(t, NewSet(NewNumber(float64('b'))), abc.CallAll(NewNumber(-1)))
	AssertEqualValues(t, NewSet(NewNumber(float64('c'))), abc.CallAll(NewNumber(0)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(1)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(-3)))

	abc = NewOffsetString([]rune("abc"), 2)
	AssertEqualValues(t, NewSet(NewNumber(float64('a'))), abc.CallAll(NewNumber(2)))
	AssertEqualValues(t, NewSet(NewNumber(float64('b'))), abc.CallAll(NewNumber(3)))
	AssertEqualValues(t, NewSet(NewNumber(float64('c'))), abc.CallAll(NewNumber(4)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(1)))
	AssertEqualValues(t, None, abc.CallAll(NewNumber(5)))
}
