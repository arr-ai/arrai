package tests

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tests"
)

// assertCodesEvalToSameValue asserts that code evaluate to the same value as
// expected.
func assertCodesEvalToSameValue(
	t *testing.T, expected string, code string,
) bool {
	expectedExpr, err := syntax.Parse(bytes.NewBufferString(expected))
	if !assert.NoError(t, err, "parsing expected: %s", expected) {
		return false
	}
	codeExpr, err := syntax.Parse(bytes.NewBufferString(code))
	if !assert.NoError(t, err, "parsing code: %s", code) {
		return false
	}
	if !tests.AssertExprsEvalToSameValue(t, expectedExpr, codeExpr) {
		return assert.Fail(
			t, "Codes should eval to same value", "%s == %s", expected, code)
	}
	return true
}

// assertCodesEvalToSameValue requires that code evaluate to the same value as
// expected.
func requireCodesEvalToSameValue(t *testing.T, expected string, code string) {
	expectedExpr, err := syntax.Parse(bytes.NewBufferString(expected))
	require.NoError(t, err)
	codeExpr, err := syntax.Parse(bytes.NewBufferString(code))
	require.NoError(t, err)
	tests.AssertExprsEvalToSameValue(t, expectedExpr, codeExpr)
}
