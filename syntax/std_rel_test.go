package syntax

import "testing"

func TestRelUnion(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, 2, 3, 4}`, `//rel.union({{1, 2}, {2, 3}, {3, 4}})`)
	AssertCodesEvalToSameValue(t, `{1, 2, 3}   `, `//rel.union({{1}, {2}, {3}})         `)
	AssertCodesEvalToSameValue(t, `{1}         `, `//rel.union({{1}, {1}, {1}})         `)
	AssertCodesEvalToSameValue(t, `{}          `, `//rel.union({})                      `)
}

func TestRelUnionError(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, `arg to //rel.union must be set, not tuple`, `//rel.union(())`)
	AssertCodeErrors(t, `elems of set arg to //rel.union must be sets, not tuple`, `//rel.union({()})`)
	// FIXME: This should error with "arg to //rel.union must be set, not closure".
	AssertCodePanics(t, `//rel.union(\x x)`)
}
