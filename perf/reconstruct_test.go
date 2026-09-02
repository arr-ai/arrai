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
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Allocation budget for the reconstruct scenario, in millions. Lower it when
// an optimisation earns it; raising it needs a reason in the commit message.
const reconstructAllocBudgetM = 37.5

func TestReconstruct(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}

	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(dir, "expected.arrai"))
	require.NoError(t, err)

	got := runReconstruct(t, filepath.Join(dir, "model.sysl.pb"))
	t.Logf("reconstruct: %s (compile %s + eval %s), %.2fM allocations, %.0fMB allocated",
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB)

	require.Equal(t, string(expected), got.out,
		"output differs from the v0.321.0 reference; this is a correctness regression")

	if !fastPathsEnabled || raceEnabled || runtime.GOOS == "windows" {
		// The budget is calibrated on the development platform; Windows'
		// runtime allocates measurably more for identical work. Output
		// identity is still enforced above on every platform.
		t.Log("fast paths disabled, -race, or windows: output was checked, allocation budget skipped")
		return
	}
	if got.allocsM > reconstructAllocBudgetM {
		t.Errorf("allocation budget exceeded: %.2fM > %.2fM. If this is a deliberate "+
			"trade, raise reconstructAllocBudgetM and say why in the commit message.",
			got.allocsM, reconstructAllocBudgetM)
	}
	if got.allocsM < reconstructAllocBudgetM*0.9 {
		t.Errorf("allocations are %.2fM, well under the %.2fM budget: lower the budget "+
			"in the commit that earned it, so the ratchet keeps holding.",
			got.allocsM, reconstructAllocBudgetM)
	}
}

func TestReconstructPlan(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}

	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(dir, "expected.arrai"))
	require.NoError(t, err)

	got := runReconstructFromPlan(t, filepath.Join(dir, "model.sysl.pb"))
	t.Logf("reconstruct plan: %s (decode %s + eval %s), %.2fM allocations, %.0fMB allocated",
		got.elapsed.Round(time.Millisecond),
		got.compile.Round(time.Millisecond),
		got.eval.Round(time.Millisecond),
		got.allocsM,
		got.totalMiB)

	require.Equal(t, string(expected), got.out,
		"plan-executed reconstruct differs from the v0.321.0 reference")
}
