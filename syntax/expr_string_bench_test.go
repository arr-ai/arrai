package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
)

// String-workload benchmarks: the operations text-rendering programs (like
// sysl reconstruct) lean on. Each evaluates a whole expression so the
// numbers include real evaluator dispatch.

// Building a string by repeated concatenation — the shape template
// rendering reduces to.
const benchStrConcatCode = `let rec f = \n \acc cond {n = 0: acc, _: f(n - 1, acc ++ $"line ${n}\n")};
f(2000, "")`

// Template rendering per element, then joining.
const benchStrTemplateCode = `//seq.join("\n",
	//seq.repeat(2000, [0]) => (@: .@, @item: $"item ${.@}: value=${.@ % 97}"))`

// Sorting strings: Less-heavy.
const benchStrSortCode = `(//seq.repeat(3000, [0]) => $"key-${(.@ * 7919) % 3000}") orderby .`

// Set-of-strings membership and dedupe: Hash128- and Equal-heavy.
const benchStrSetCode = `(//seq.repeat(5000, [0]) => $"name-${.@ % 1250}") count`

func BenchmarkStrConcat(b *testing.B)   { benchEvalLoop(b, benchStrConcatCode) }
func BenchmarkStrTemplate(b *testing.B) { benchEvalLoop(b, benchStrTemplateCode) }
func BenchmarkStrSort(b *testing.B)     { benchEvalLoop(b, benchStrSortCode) }
func BenchmarkStrSet(b *testing.B)      { benchEvalLoop(b, benchStrSetCode) }

// benchEvalLoop compiles once and measures evaluation only, under b.Loop:
// re-parsing per iteration would drown the evaluator in wbnf noise, and
// b.N-style loops allow inlining artifacts.
func benchEvalLoop(b *testing.B, code string) {
	b.Helper()
	ctx := importcache.WithNewImportCache(arraictx.InitRunCtx(context.Background()))
	expr, err := Compile(ctx, NoPath, code)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := expr.Eval(ctx, rel.Scope{}); err != nil {
			b.Fatal(err)
		}
	}
}
