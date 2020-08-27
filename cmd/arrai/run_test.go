package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/require"
)

func TestEvalFile(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	ctx := arraictx.InitRunCtx(context.Background())
	require.NoError(t, evalFile(ctx, "../../examples/jsfuncs/jsfuncs.arrai", &buf, ""))
	require.NoError(t, evalFile(ctx, "../../examples/grpc/app.arrai", &buf, ""))
}

func TestEvalNotExistingFile(t *testing.T) {
	t.Parallel()
	ctx := arraictx.InitRunCtx(context.Background())
	require.Equal(t, `"version": not a command and not found as a file in the current directory`,
		evalFile(ctx, "version", nil, "").Error())

	require.Equal(t, `"`+string([]rune{'.', os.PathSeparator})+`version": file not found`,
		evalFile(ctx, string([]rune{'.', os.PathSeparator})+"version", nil, "").Error())
}
