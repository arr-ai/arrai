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
		assert.Equal(t, expected, expr.String())
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

func TestEvalCond(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, `1`, `cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `1`, `cond (1 < 2 : 1)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2)`)
	assertEvalOutputs(t, `3`, `cond (* : 1 + 2)`)
	assertEvalOutputs(t, `1`, `cond (1 < 2: 1, * : 1 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2: 1, * : 1 + 2)`)
	assertEvalOutputs(t, `3`, `let a = cond (1 > 2: 1, * : 1 + 2);a`)
	assertEvalOutputs(t, `1`, `let a = cond (1 < 2: 1, * : 1 + 2);a * 1`)

	var sb strings.Builder
	assert.Error(t, evalImpl(`cond (1 > 0 : 1, 2 < 3: 2, *:1 + 2)`, &sb))
	assert.Error(t, evalImpl(`cond ()`, &sb))
}

// EvalCondExpr is used in local env only.
func EvalCondExpr(t *testing.T) {
	t.Parallel()
	assertEvalExprString(t, "{(1 > 0): 1, (2 > 3): 2, *: (1 + 2)}", "cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2)")
	assertEvalExprString(t, "{(1 < 2): 1}", "cond (1 < 2 : 1)")
	assertEvalExprString(t, "{(1 > 2): 1, (2 < 3): 2}", "cond (1 > 2 : 1, 2 < 3: 2)")
	assertEvalExprString(t, "{*: (1 + 2)}", "cond (*: 1 + 2)")
	assertEvalExprString(t, "{(1 < 2): 1, *: (1 + 2)}", "cond (1 < 2: 1, * : 1 + 2)")
}
