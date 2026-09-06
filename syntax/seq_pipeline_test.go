package syntax

import (
	"testing"

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
