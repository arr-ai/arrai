// Package perf holds end-to-end performance regression scenarios: whole
// arr.ai programs, run in process so their cost can be measured and pinned.
//
// The scenario here is the one that prompted the 2026 evaluator work. It is
// the anz-bank/sysl reconstruct pipeline — load a Sysl model from protobuf,
// normalise it into relations, and render it back to Sysl source — over a
// generated model of 1000 synthetic applications. It exercises the paths
// that dominate real arr.ai programs: relation construction, indexed
// lookups, joins, nesting, pattern matching, closure calls, string building
// and grammar parsing.
//
// Two things are asserted, and they are asserted differently on purpose:
//
//   - Output is compared byte for byte against expected.arrai, which was
//     produced by arrai v0.321.0, before the frozen generics migration. Any
//     change to it is a correctness regression, full stop.
//   - Allocation count is compared against a committed budget. Allocation is
//     deterministic where wall-clock is not, and it tracks the work done
//     closely enough to catch a regression. The budget is a ratchet in both
//     directions: an improvement should lower it in the same commit that
//     earns it, so the number only ever moves deliberately.
//
// Wall-clock time is reported but never asserted; it is far too sensitive to
// the machine and to other load to gate on.
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

// Allocation budget for the reconstruct scenario, in millions. Lower it when
// an optimisation earns it; raising it needs a reason in the commit message.
const reconstructAllocBudgetM = 40.0

func TestReconstruct(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}

	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	model := filepath.Join(dir, "model.sysl.pb")
	script := filepath.Join(dir, "vendor", "run.arrai")

	source, err := os.ReadFile(script)
	require.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(dir, "expected.arrai"))
	require.NoError(t, err)

	// run.arrai resolves its imports relative to the working directory.
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(filepath.Join(dir, "vendor")))
	defer func() { require.NoError(t, os.Chdir(wd)) }()

	ctx := importcache.WithNewImportCache(
		arraictx.WithArgs(arraictx.InitRunCtx(context.Background()), script, model))

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	start := time.Now()

	// Compile and evaluate as separate phases: compile is the cost of
	// parsing the program and the sysl library, eval is the algorithm. The
	// split is what makes comparisons with the native implementation fair —
	// its compilation happened at build time.
	expr, err := syntax.Compile(ctx, script, string(source))
	require.NoError(t, err)
	compiled := time.Now()

	value, err := expr.Eval(ctx, rel.Scope{})
	require.NoError(t, err)
	var out bytes.Buffer
	require.NoError(t, arrai.OutputValue(ctx, value, &out, ""))

	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)

	allocsM := float64(after.Mallocs-before.Mallocs) / 1e6
	t.Logf("reconstruct: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated",
		elapsed.Round(time.Millisecond),
		compiled.Sub(start).Round(time.Millisecond),
		time.Since(compiled).Round(time.Millisecond),
		allocsM,
		float64(after.TotalAlloc-before.TotalAlloc)/(1<<20))

	require.Equal(t, string(expected), out.String(),
		"output differs from the v0.321.0 reference; this is a correctness regression")

	if !fastPathsEnabled || raceEnabled {
		t.Log("fast paths disabled or -race: output was checked, allocation budget skipped")
		return
	}
	if allocsM > reconstructAllocBudgetM {
		t.Errorf("allocation budget exceeded: %.2fM > %.2fM. If this is a deliberate "+
			"trade, raise reconstructAllocBudgetM and say why in the commit message.",
			allocsM, reconstructAllocBudgetM)
	}
	if allocsM < reconstructAllocBudgetM*0.9 {
		t.Errorf("allocations are %.2fM, well under the %.2fM budget: lower the budget "+
			"in the commit that earned it, so the ratchet keeps holding.",
			allocsM, reconstructAllocBudgetM)
	}
}
