package main

import (
	"strings"
	"testing"

	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/assert"
)

func assertEvalOutputs(t *testing.T, expected, source string) bool { //nolint:unparam
	var sb strings.Builder
	return assert.NoError(t, evalImpl(source, &sb)) &&
		assert.Equal(t, expected, strings.TrimRight(sb.String(), "\n"))
}

func assertEvalExprString(t *testing.T, expected, source string) bool { //nolint:unparam
	expr, err := syntax.Compile(".", source)
	return assert.True(t, err == nil) &&
		assert.True(t, expr != nil) &&
		assert.Equal(t, expected, strings.Replace(expr.String(), ` `, ``, -1))
}

func TestEvalNumberULP(t *testing.T) {
	assertEvalOutputs(t, `0.3`, `0.1 + 0.1 + 0.1`)
}

func TestEvalString(t *testing.T) {
	assertEvalOutputs(t, ``, `""`)
	assertEvalOutputs(t, ``, `{}`)
	assertEvalOutputs(t, `abc`, `"abc"`)
}

func TestEvalComplex(t *testing.T) {
	assertEvalOutputs(t, `[42, 'abc']`, `[42, "abc"]`)
	assertEvalOutputs(t, `{42, 'abc'}`, `{"abc", 42}`)
}
