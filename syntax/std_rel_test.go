package syntax

import "testing"

func TestRelUnion(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, 2, 3, 4}`, `//rel.union({{1, 2}, {2, 3}, {3, 4}})`)
	AssertCodesEvalToSameValue(t, `{1, 2, 3}   `, `//rel.union({{1}, {2}, {3}})         `)
	AssertCodesEvalToSameValue(t, `{1}         `, `//rel.union({{1}, {1}, {1}})         `)
	AssertCodesEvalToSameValue(t, `{}          `, `//rel.union({})                      `)
}
