package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalFile(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	require.NoError(t, evalFile("../../examples/jsfuncs/jsfuncs.arrai", &buf))
	require.NoError(t, evalFile("../../examples/grpc/app.arrai", &buf))
}
