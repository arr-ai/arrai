//nolint:dupl
package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const hello = "hello"

func TestIsStringTuple(t *testing.T) {
	for e := NewString([]rune(hello)).Enumerator(); e.MoveNext(); {
		tuple, is := e.Current().(StringCharTuple)
		if assert.True(t, is) {
			assert.Equal(t, rune(hello[tuple.at]), tuple.char)
		}
	}
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
