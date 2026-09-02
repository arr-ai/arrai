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

func TestWherePushdownThroughProject(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{|a| (1)}`,
		`{|a,b| (1, 10), (2, 20)} => (a: .a) where .a = 1`)
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
