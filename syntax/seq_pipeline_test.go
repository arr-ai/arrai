package syntax

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
)

func TestSeqPipelineCountPassthrough(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`,
		`([0, 1, 2] >> . + 1 >> . * 2) count`)
}

func TestSeqPipelineValuesMatchNaive(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`[2, 4, 6]`,
		`[0, 1, 2] >> . + 1 >> . * 2`)
}

func TestProjectDotsCountPassthrough(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`,
		`({|a,b| (1, 10), (2, 20)} => (a: .a)) count`)
}

func TestPruneStackedProject(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{|a| (1)}`,
		`{|a, b, c| (1, 2, 3)} => (a: .a, b: .b, c: .c) => (a: .a)`)
}

func TestWhereIndexCacheKeepsCallerAttrs(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`[{|k| (1)}, {|a, k| (10, 1)}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .k); [p where .k = 1, r where .k = 1]`)
	AssertCodesEvalToSameValue(t,
		`[{|a, k| (10, 1)}, {|k| (1)}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .k); [r where .k = 1, p where .k = 1]`)
}

func TestWhereIndexCacheKeysByStoreColumn(t *testing.T) {
	t.Parallel()
	// Dest name k, store column a: p's .k is r's .a. Cache must not key by dest name.
	AssertCodesEvalToSameValue(t,
		`[{|k| (10)}, {}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .a); [p where .k = 10, r where .k = 10]`)
	AssertCodesEvalToSameValue(t,
		`[{}, {|k| (10)}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .a); [r where .k = 10, p where .k = 10]`)
	AssertCodesEvalToSameValue(t,
		`[{}, {|a, k| (10, 1)}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .a); [p where .k = 1, r where .k = 1]`)
	AssertCodesEvalToSameValue(t,
		`[{|a, k| (10, 1)}, {}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .a); [r where .k = 1, p where .k = 1]`)
	AssertCodesEvalToSameValue(t,
		`[{|x| (1)}, {}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (x: .k); let q = r => (x: .a); [p where .x = 1, q where .x = 1]`)
	AssertCodesEvalToSameValue(t,
		`[{}, {|x| (1)}]`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (x: .k); let q = r => (x: .a); [q where .x = 1, p where .x = 1]`)
}

func TestWherePushdownThroughProject(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{|a| (1)}`,
		`{|a,b| (1, 10), (2, 20)} => (a: .a) where .a = 1`)
}

// A `>>` whose body is itself a `>>` leaves an unforced seqPipeline as an
// element of the outer array; that must still compare equal to a literal.
func TestSeqPipelineNestedEqualsArrayLiteral(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`[['bar']]`,
		`[['bar']] >> (. >> .)`,
	)
}

func TestSeqPipelineIsSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`,
		`(([0, 1, 2] >> . + 1 >> . * 2) >> . + 0) count`)
	_ = rel.EmptyScope
}

func TestDictPipelineCountPassthrough(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`,
		`({'a': 0, 'b': 1, 'c': 2} >> . + 1 >> . * 2) count`)
}

func TestDictPipelineValuesMatchNaive(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{'a': 2, 'b': 4, 'c': 6}`,
		`{'a': 0, 'b': 1, 'c': 2} >> . + 1 >> . * 2`)
}

func TestWherePushdownErrorMatchesNaive(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, "single: too many elements",
		`{|a, b| (1, 10), (2, 20)} => (a: .a) where .a = ({1, 2} single)`)
}

func TestLetFanoutTwoConsumers(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`[2, 4, 6, 3, 6, 9]`,
		`let x = [0, 1, 2] >> . + 1; (x >> . * 2) ++ (x >> . * 3)`)
}

// TestSeqPipelineEffectfulMapRunsWhenDiscarded is the standing oracle for a
// regression where >> suspended every map in a seqPipeline/dictPipeline,
// including one whose body calls something effectful (e.g. //log.print). A
// pipeline that is never forced — because its result is bound to `_` and
// never used again, or is only ever asked for its count — silently dropped
// those calls. A call is assumed effectful because the callee might be, so
// >> must run it eagerly instead of suspending it.
func TestSeqPipelineEffectfulMapRunsWhenDiscarded(t *testing.T) {
	t.Parallel()

	v, calls := evalCounting(t, `let _ = [1, 2, 3] >> tick(.); 0`)
	require.Equal(t, int64(3), calls, "discarded >> result must still run its effectful body")
	want, err := EvaluateExpr(arraictx.InitRunCtx(context.Background()), NoPath, `0`)
	require.NoError(t, err)
	require.True(t, want.Equal(v))
}

func TestSeqPipelineEffectfulMapRunsUnderCount(t *testing.T) {
	t.Parallel()

	v, calls := evalCounting(t, `([1, 2, 3] >> tick(.)) count`)
	require.Equal(t, int64(3), calls, "count must not answer from metadata for an effectful >>")
	want, err := EvaluateExpr(arraictx.InitRunCtx(context.Background()), NoPath, `3`)
	require.NoError(t, err)
	require.True(t, want.Equal(v))
}

func TestDictPipelineEffectfulMapRunsWhenDiscarded(t *testing.T) {
	t.Parallel()

	_, calls := evalCounting(t, `let _ = {'a': 1, 'b': 2} >> tick(.); 0`)
	require.Equal(t, int64(2), calls, "discarded >> result over a dict must still run its effectful body")
}

// TestSeqPipelineEffectfulChainStillMaterialisesOnce checks that gating the
// fast path on purity does not reintroduce double evaluation for a chained
// effectful >>: each stage still runs its body exactly once per element.
func TestSeqPipelineEffectfulChainStillMaterialisesOnce(t *testing.T) {
	t.Parallel()

	// tick ignores its argument and returns the shared counter, so only the
	// call count (not the resulting values, whose assignment across the two
	// stages races) is asserted here.
	v, calls := evalCounting(t, `[1, 2] >> tick(.) >> tick(.)`)
	require.Equal(t, int64(4), calls, "each of 2 elements through 2 effectful stages")
	s, ok := v.(rel.Set)
	require.True(t, ok, "%v", v)
	require.Equal(t, 2, s.Count())
}
