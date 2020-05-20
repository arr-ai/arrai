package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsString(t *testing.T) {
	t.Parallel()

	generic := NewString([]rune("this is a test")).Map(func(v Value) Value { return v })
	stringified, isString := AsString(generic)
	require.True(t, isString)
	assert.True(t, stringified.Equal(NewString([]rune("this is a test"))), stringified.String())

	// generic = NewOffsetString([]rune("this is a test"), 100).Map(func(v Value) Value { return v })
	// stringified, isString = AsString(generic)
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewOffsetString([]rune("this is a test"), 100)))

	// generic = NewString([]rune("")).Map(func(v Value) Value { return v })
	// stringified, isString = AsString(generic)
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewString([]rune(""))))
}
