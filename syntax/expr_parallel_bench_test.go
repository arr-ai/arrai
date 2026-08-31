package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/frozen"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

// benchEval measures one arr.ai expression end to end, optionally with
// parallel evaluation disabled, so the pair quantifies what the fan-out in
// `where` and `=>` buys for real closure bodies.
func benchEval(b *testing.B, parallel bool, code string) {
	b.Helper()
	if !parallel {
		defer frozen.SetMinParallelChunk(frozen.DisableParallel())
	}
	ctx := arraictx.InitRunCtx(context.Background())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := EvaluateExpr(ctx, "", code); err != nil {
			b.Fatal(err)
		}
	}
}

const benchWhereCode = `(//seq.repeat(50000, [0]) ` +
	`=> (a: .@, b: .@ % 97) where .a % 7 = 0 && .b < 50) count`

const benchDArrowCode = `(//seq.repeat(50000, [0]) ` +
	`=> (a: .@, b: .@ * .@ % 101, c: $"${.@}")) count`

func BenchmarkEvalWhereParallel(b *testing.B)   { benchEval(b, true, benchWhereCode) }
func BenchmarkEvalWhereSequential(b *testing.B) { benchEval(b, false, benchWhereCode) }
func BenchmarkEvalDArrowParallel(b *testing.B)  { benchEval(b, true, benchDArrowCode) }
func BenchmarkEvalDArrowSequential(b *testing.B) {
	benchEval(b, false, benchDArrowCode)
}
