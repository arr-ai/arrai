package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	t.Parallel()

	out, err := eval("", "hello", "world")

	require.NoError(t, err)
	assert.Equal(t, "hello world", out.String())
}
