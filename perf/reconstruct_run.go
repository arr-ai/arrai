package perf

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arrai"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

type reconstructResult struct {
	out      string
	elapsed  time.Duration
	compile  time.Duration
	eval     time.Duration
	allocsM  float64
	totalMiB float64
}

func runReconstruct(t *testing.T, modelPath string) reconstructResult {
	t.Helper()
	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	script := filepath.Join(dir, "vendor", "run.arrai")
	source, err := os.ReadFile(script)
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(filepath.Join(dir, "vendor")))
	defer func() { require.NoError(t, os.Chdir(wd)) }()

	ctx := importcache.WithNewImportCache(
		arraictx.WithArgs(arraictx.InitRunCtx(context.Background()), script, modelPath))

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	start := time.Now()

	expr, err := syntax.Compile(ctx, script, string(source))
	require.NoError(t, err)
	compiled := time.Now()

	value, err := expr.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, arrai.OutputValue(ctx, value, &buf, ""))

	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)
	return reconstructResult{
		out:      buf.String(),
		elapsed:  elapsed,
		compile:  compiled.Sub(start),
		eval:     time.Since(compiled),
		allocsM:  float64(after.Mallocs-before.Mallocs) / 1e6,
		totalMiB: float64(after.TotalAlloc-before.TotalAlloc) / (1 << 20),
	}
}

// runReconstructFromPlan compiles once, serializes the operator graph, then
// times decode+eval as the compiled-.arraiz path (🎯T25).
func runReconstructFromPlan(t *testing.T, modelPath string) reconstructResult {
	t.Helper()
	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	script := filepath.Join(dir, "vendor", "run.arrai")
	source, err := os.ReadFile(script)
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(filepath.Join(dir, "vendor")))
	defer func() { require.NoError(t, os.Chdir(wd)) }()

	ctx := importcache.WithNewImportCache(
		arraictx.WithArgs(arraictx.InitRunCtx(context.Background()), script, modelPath))

	expr, err := syntax.Compile(ctx, script, string(source))
	require.NoError(t, err)
	plan, err := rel.LowerPlan(expr)
	require.NoError(t, err)
	encoded, err := rel.EncodePlan(plan)
	require.NoError(t, err)

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	start := time.Now()

	loaded, err := rel.DecodePlan(encoded)
	require.NoError(t, err)
	compiled := time.Now()

	value, err := loaded.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, arrai.OutputValue(ctx, value, &buf, ""))

	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)
	return reconstructResult{
		out:      buf.String(),
		elapsed:  elapsed,
		compile:  compiled.Sub(start),
		eval:     time.Since(compiled),
		allocsM:  float64(after.Mallocs-before.Mallocs) / 1e6,
		totalMiB: float64(after.TotalAlloc-before.TotalAlloc) / (1 << 20),
	}
}
