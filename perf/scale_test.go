package perf

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arrai"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

// Scale points for S0. 10k is the floor of the design-note 10k–50k reconstruct
// band; 1e6 is a million-row join. Both skip default CI (see skipScale).
const (
	scaleReconstructApps = 10000
	scaleJoinRows        = 1000000
)

func skipScale(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}
	if os.Getenv("ARRAI_SCALE") == "" {
		t.Skip("scale suite: set ARRAI_SCALE=1 to run 10k-app reconstruct and million-row join")
	}
}

func TestGenerateSyslModelPBShape(t *testing.T) {
	pb := generateSyslModelPB(2)
	require.Greater(t, len(pb), 100)
	s := string(pb)
	require.Contains(t, s, "App0")
	require.Contains(t, s, "App1")
	require.Contains(t, s, "Get")
	require.Contains(t, s, "Data")
}

func TestScaleJoinPipelineSmall(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}
	const n = 7000
	got := runScaleScript(t, filepath.Join("scale", "join.arrai"), strconv.Itoa(n))
	require.Equal(t, strconv.Itoa((n+6)/7), strings.TrimSpace(got.out),
		"join/pipeline count disagrees with the closed-form oracle (n+6)/7")
}

func TestGenerateSyslModelPBMatchesGolden(t *testing.T) {
	skipScale(t)
	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(dir, "expected.arrai"))
	require.NoError(t, err)

	model := filepath.Join(t.TempDir(), "model.sysl.pb")
	require.NoError(t, os.WriteFile(model, generateSyslModelPB(1000), 0o600))
	got := runReconstruct(t, model)
	require.Equal(t, string(expected), got.out,
		"generated 1000-app protobuf does not reconstruct to the v0.321.0 golden")
}

func TestScaleReconstruct(t *testing.T) {
	skipScale(t)
	model := filepath.Join(t.TempDir(), "model.sysl.pb")
	require.NoError(t, os.WriteFile(model, generateSyslModelPB(scaleReconstructApps), 0o600))
	got := runReconstruct(t, model)
	t.Logf("scale reconstruct %d apps: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated",
		scaleReconstructApps,
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB)

	sum := sha256.Sum256([]byte(got.out))
	digest := hex.EncodeToString(sum[:])
	t.Logf("scale reconstruct sha256: %s (%d bytes)", digest, len(got.out))
	require.Contains(t, got.out, "App0", "reconstruct output missing App0")
	require.Contains(t, got.out, "App"+strconv.Itoa(scaleReconstructApps-1),
		"reconstruct output missing last app")
	require.Equal(t, scaleReconstructSHA256, digest,
		"reconstruct output hash moved; if this is a deliberate change, update scaleReconstructSHA256")
	if !fastPathsEnabled || raceEnabled || runtime.GOOS == "windows" {
		t.Log("fast paths disabled, -race, or windows: output was checked, allocation budget skipped")
		return
	}
	if got.allocsM > scaleReconstructAllocBudgetM {
		t.Errorf("allocation budget exceeded: %.2fM > %.2fM", got.allocsM, scaleReconstructAllocBudgetM)
	}
	if got.allocsM < scaleReconstructAllocBudgetM*0.9 {
		t.Errorf("allocations are %.2fM, well under the %.2fM budget: lower the budget in the commit that earned it",
			got.allocsM, scaleReconstructAllocBudgetM)
	}
}

func TestScaleSeqPipelineSmall(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}
	const n = 20000
	got := runScaleScript(t, filepath.Join("scale", "pipeline.arrai"), strconv.Itoa(n))
	gotN, err := strconv.ParseFloat(strings.TrimSpace(got.out), 64)
	require.NoError(t, err)
	require.Equal(t, float64(n), gotN, "n-stage >> count must equal the base length")
}

func TestScaleSeqPipeline(t *testing.T) {
	skipScale(t)
	got := runScaleScript(t, filepath.Join("scale", "pipeline.arrai"), strconv.Itoa(scaleJoinRows))
	t.Logf("scale pipeline %d elems: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated, out=%s",
		scaleJoinRows,
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB,
		got.out)
	gotN, err := strconv.ParseFloat(strings.TrimSpace(got.out), 64)
	require.NoError(t, err)
	require.Equal(t, float64(scaleJoinRows), gotN, "n-stage >> count must equal the base length")
}

func TestScalePushdownSmall(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}
	const n = 7000
	got := runScaleScript(t, filepath.Join("scale", "pushdown.arrai"), strconv.Itoa(n))
	want := n / 1000
	if n%1000 > 42 {
		want++
	}
	gotN, err := strconv.ParseFloat(strings.TrimSpace(got.out), 64)
	require.NoError(t, err)
	require.Equal(t, float64(want), gotN, "eq-key where after project disagrees with k=42 count")
}

func TestScalePushdown(t *testing.T) {
	skipScale(t)
	got := runScaleScript(t, filepath.Join("scale", "pushdown.arrai"), strconv.Itoa(scaleJoinRows))
	want := scaleJoinRows / 1000
	if scaleJoinRows%1000 > 42 {
		want++
	}
	t.Logf("scale pushdown %d rows: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated, out=%s",
		scaleJoinRows,
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB,
		got.out)
	gotN, err := strconv.ParseFloat(strings.TrimSpace(got.out), 64)
	require.NoError(t, err)
	require.Equal(t, float64(want), gotN, "eq-key where after project disagrees with k=42 count")
}

func TestScaleJoinPipeline(t *testing.T) {
	skipScale(t)
	got := runScaleScript(t, filepath.Join("scale", "join.arrai"), strconv.Itoa(scaleJoinRows))
	t.Logf("scale join %d rows: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated, out=%s",
		scaleJoinRows,
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB,
		got.out)
	want := strconv.Itoa((scaleJoinRows + 6) / 7)
	require.Equal(t, want, strings.TrimSpace(got.out),
		"join/pipeline count disagrees with the closed-form oracle (n+6)/7")
}

func runScaleScript(t *testing.T, script string, args ...string) reconstructResult {
	t.Helper()
	source, err := os.ReadFile(script)
	require.NoError(t, err)
	abs, err := filepath.Abs(script)
	require.NoError(t, err)

	ctxArgs := append([]string{abs}, args...)
	ctx := importcache.WithNewImportCache(
		arraictx.WithArgs(arraictx.InitRunCtx(context.Background()), ctxArgs...))

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	start := time.Now()
	expr, err := syntax.Compile(ctx, abs, string(source))
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

// scaleReconstructSHA256 is the pinned digest of reconstruct output at
// scaleReconstructApps, produced by generateSyslModelPB + vendor/run.arrai.
const scaleReconstructSHA256 = "3120818a5ba3943d94fb1f0b9bff20983238a89aa7158a2bcce829e57c41c5cc"

// Allocation budget for the 10k-app reconstruct, in millions. Same ratchet
// as reconstructAllocBudgetM: lower it when an optimisation earns it.
const scaleReconstructAllocBudgetM = 200.0
