package rel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayCall(t *testing.T) {
	f := NewArray(
		NewNumber(0),
		NewNumber(1),
		NewNumber(4),
		NewNumber(9),
		NewNumber(16),
		NewNumber(25),
	)
	for i := 0; i < f.Count(); i++ {
		assert.Equal(t, i*i, int(f.Call(NewNumber(float64(i))).(Number).Float64()))
	}
}

func TestAsString(t *testing.T) {
	t.Parallel()

	generic := NewString([]rune("this is a test")).Map(func(v Value) Value { return v })
	stringified, isString := generic.AsString()
	require.True(t, isString)
	assert.True(t, stringified.Equal(NewString([]rune("this is a test"))), fmt.Sprintf("%s", stringified))

	// generic = NewOffsetString([]rune("this is a test"), 100).Map(func(v Value) Value { return v })
	// stringified, isString = generic.AsString()
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewOffsetString([]rune("this is a test"), 100)))

	// generic = NewString([]rune("")).Map(func(v Value) Value { return v })
	// stringified, isString = generic.AsString()
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewString([]rune(""))))
}
