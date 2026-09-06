package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

// TestPlanRoundtripCorpus compiles every *_test.arrai file under testdata/planroundtrip and
// evaluates it twice: once directly, and once after lowering it to a compiled plan and decoding
// it back (the same LowerPlan/EncodePlan/DecodePlan round trip a .arraiz bundle with plan.bin
// goes through). Every test case in the corpus must pass both ways; a mismatch means the plan
// encode/decode path is losing information that direct evaluation preserves.
func TestPlanRoundtripCorpus(t *testing.T) {
	ctx := arraictx.InitRunCtx(context.Background())

	files, err := getTestFiles(ctx, "testdata/planroundtrip")
	require.NoError(t, err)

	for i := range files {
		file := files[i]
		t.Run(relPath(file.Path), func(t *testing.T) {
			t.Parallel()

			fileCtx := importcache.WithNewImportCache(ctx)
			expr, err := syntax.Compile(fileCtx, file.Path, file.Source)
			require.NoError(t, err, "compile")

			direct, err := expr.Eval(fileCtx, rel.Scope{})
			require.NoError(t, err, "direct eval")
			requireAllPassed(t, "direct", classifyResults(direct))

			p, err := rel.LowerPlan(expr)
			require.NoError(t, err, "lower plan")
			b, err := rel.EncodePlan(p)
			require.NoError(t, err, "encode plan")
			p2, err := rel.DecodePlan(b)
			require.NoError(t, err, "decode plan")

			viaPlan, err := p2.Eval(fileCtx, rel.Scope{})
			require.NoError(t, err, "plan eval")
			requireAllPassed(t, "plan", classifyResults(viaPlan))
		})
	}
}

func requireAllPassed(t *testing.T, label string, results []Result) {
	t.Helper()
	for _, r := range results {
		if r.Outcome != Passed {
			t.Errorf("[%s] %s: outcome=%d (%s)", label, r.Name, r.Outcome, r.Message)
		}
	}
}
