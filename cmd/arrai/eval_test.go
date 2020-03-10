package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEvalOutputs(t *testing.T, expected, source string) bool { //nolint:unparam
	var sb strings.Builder
	return assert.NoError(t, evalImpl(source, &sb)) &&
		assert.Equal(t, expected, strings.TrimRight(sb.String(), "\n"))
}

func TestEvalNumberULP(t *testing.T) {
	assertEvalOutputs(t, `0.3`, `0.1 + 0.1 + 0.1`)
}

func TestEvalDoubleQuoteTupleAccess(t *testing.T) {
	assertEvalOutputs(t, `42`, `("😅":42)."😅"`)
	assertEvalOutputs(t, `42`, `("😅":42).'😅'`)
	assertEvalOutputs(t, `42`, "('😅':42).`😅`")
	assertEvalOutputs(t, `42`, `(x:42)."x"`)
	assertEvalOutputs(t, `42`, `("":42).""`)
}
