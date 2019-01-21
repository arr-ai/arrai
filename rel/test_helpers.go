package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// intSet returns a new set from the given elements.
func intSet(elts ...interface{}) Set {
	result, err := NewSetFrom(elts...)
	if err != nil {
		panic(err)
	}
	return result
}

// AssertEqualValues asserts that the two values are Equal.
func AssertEqualValues(t *testing.T, expected, actual Value) bool {
	return assert.True(t, expected.Equal(actual), "%s ==\n%s", expected, actual)
}

// requireEqualValues requires that the two values are Equal.
func requireEqualValues(t *testing.T, expected, actual Value) {
	require.True(t, expected.Equal(actual), "%s ==\n%s", expected, actual)
}

// AssertExprsEvalToSameValue asserts that the exprs evaluate to the same value.
func AssertExprsEvalToSameValue(
	t *testing.T, expected, expr Expr,
) bool {
	expectedValue, err := expected.Eval(EmptyScope, EmptyScope)
	if !assert.NoError(t, err, "evaluating expected: %s", expected) {
		return false
	}
	value, err := expr.Eval(EmptyScope, EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	if !AssertEqualValues(t, expectedValue, value) {
		return assert.Fail(t, "exprs !=", "%s ==\n%s", expected, expr)
	}
	return true
}

// RequireExprsEvalToSameValue requires that the exprs evaluate to the same
// value.
func RequireExprsEvalToSameValue(
	t *testing.T, expected, expr Expr,
) {
	expectedValue, err := expected.Eval(EmptyScope, EmptyScope)
	require.NoError(t, err)
	value, err := expr.Eval(EmptyScope, EmptyScope)
	require.NoError(t, err)
	requireEqualValues(t, expectedValue, value)
}
