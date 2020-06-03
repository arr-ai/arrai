package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalFile(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	require.NoError(t, evalFile("../../examples/jsfuncs/jsfuncs.arrai", &buf))
	require.NoError(t, evalFile("../../examples/grpc/app.arrai", &buf))
}

func TestEvalNotExistingFile(t *testing.T) {
	require.Equal(t, `"version": not a command and not found as a file in the current directory`,
		evalFile("version", nil).Error())

	require.Equal(t, `"`+string([]rune{'.', os.PathSeparator})+`version": file not found`, evalFile("./version", nil).Error())
}
