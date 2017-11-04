package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/rel/syntax"
)

// intSet returns a new set from the given elements.
func intSet(elts ...interface{}) rel.Set {
	result, err := rel.NewSetFrom(elts...)
	if err != nil {
		panic(err)
	}
	return result
}

// assertEqualValues asserts that the two values are Equal.
func assertEqualValues(t *testing.T, expected, actual rel.Value) bool {
	return assert.True(t, expected.Equal(actual), "%s ==\n%s", expected, actual)
}

// requireEqualValues requires that the two values are Equal.
func requireEqualValues(t *testing.T, expected, actual rel.Value) {
	require.True(t, expected.Equal(actual), "%s ==\n%s", expected, actual)
}

// AssertExprsEvalToSameValue asserts that the exprs evaluate to the same value.
func AssertExprsEvalToSameValue(
	t *testing.T, expected, expr rel.Expr,
) bool {
	expectedValue, err := expected.Eval(rel.EmptyScope, rel.EmptyScope)
	if !assert.NoError(t, err, "evaluating expected: %s", expected) {
		return false
	}
	value, err := expr.Eval(rel.EmptyScope, rel.EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	if !assertEqualValues(t, expectedValue, value) {
		return assert.Fail(t, "exprs !=", "%s ==\n%s", expected, expr)
	}
	return true
}

// RequireExprsEvalToSameValue requires that the exprs evaluate to the same
// value.
func RequireExprsEvalToSameValue(
	t *testing.T, expected, expr rel.Expr,
) {
	expectedValue, err := expected.Eval(rel.EmptyScope, rel.EmptyScope)
	require.NoError(t, err)
	value, err := expr.Eval(rel.EmptyScope, rel.EmptyScope)
	require.NoError(t, err)
	requireEqualValues(t, expectedValue, value)
}

// assertCodesEvalToSameValue asserts that code evaluate to the same value as
// expected.
func assertCodesEvalToSameValue(
	t *testing.T, expected string, code string,
) bool {
	expectedExpr, err := syntax.Parse([]byte(expected))
	if !assert.NoError(t, err, "parsing expected: %s", expected) {
		return false
	}
	codeExpr, err := syntax.Parse([]byte(code))
	if !assert.NoError(t, err, "parsing code: %s", code) {
		return false
	}
	if !AssertExprsEvalToSameValue(t, expectedExpr, codeExpr) {
		return assert.Fail(t,
			"Codes should eval to same value", "%s == %s", expected, code)
	}
	return true
}

// assertCodesEvalToSameValue requires that code evaluate to the same value as
// expected.
func requireCodesEvalToSameValue(t *testing.T, expected string, code string) {
	expectedExpr, err := syntax.Parse([]byte(expected))
	require.NoError(t, err)
	codeExpr, err := syntax.Parse([]byte(code))
	require.NoError(t, err)
	AssertExprsEvalToSameValue(t, expectedExpr, codeExpr)
}
