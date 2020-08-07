package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalFile(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	require.NoError(t, evalFile(context.Background(), "../../examples/jsfuncs/jsfuncs.arrai", &buf, ""))
	require.NoError(t, evalFile(context.Background(), "../../examples/grpc/app.arrai", &buf, ""))
}

func TestEvalNotExistingFile(t *testing.T) {
	t.Parallel()
	require.Equal(t, `"version": not a command and not found as a file in the current directory`,
		evalFile(context.Background(), "version", nil, "").Error())

	require.Equal(t, `"`+string([]rune{'.', os.PathSeparator})+`version": file not found`,
		evalFile(context.Background(), string([]rune{'.', os.PathSeparator})+"version", nil, "").Error())
}
