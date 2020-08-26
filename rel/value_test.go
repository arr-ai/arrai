package rel

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCall(t *testing.T) {
	t.Parallel()

	foo := func(at int, v Value) Tuple {
		return NewTuple(NewAttr("@", NewNumber(float64(at))), NewAttr("@foo", v))
	}

	set := NewSet(
		foo(1, NewNumber(42)),
		foo(1, NewNumber(24)),
	)

	result, err := SetCall(set, NewNumber(1))
	assert.Error(t, err, "%v", result)
	result, err = SetCall(set, NewNumber(0))
	assert.Error(t, err, "%v", result)

	set = NewSet(
		foo(1, NewNumber(42)),
		foo(2, NewNumber(24)),
	)

	result, err = SetCall(set, NewNumber(1))
	require.NoError(t, err)
	assert.True(t, result.Equal(NewNumber(42)))
	result, err = SetCall(set, NewNumber(2))
	require.NoError(t, err)
	assert.True(t, result.Equal(NewNumber(24)))
}

func TestReflect(t *testing.T) {
	fset := token.NewFileSet()
	bs, err := ioutil.ReadFile("value_test.go")
	assert.NoError(t, err)

	f, err := parser.ParseFile(fset, "", string(bs), 0)
	assert.NoError(t, err)

	v, err := NewValue(f)
	assert.NoError(t, err)

	assert.NotNil(t, v)
}
