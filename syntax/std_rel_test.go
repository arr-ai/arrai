package syntax

import "testing"

func TestRelUnion(t *testing.T) {
	t.Parallel()

	AssertCodeEvalsToType(t, `{1, 2, 3, 4}`, `//rel.union({{1, 2}, {2, 3}, {3, 4}})`)
	AssertCodeEvalsToType(t, `{1, 2, 3}   `, `//rel.union({{1}, {2}, {3}})         `)
	AssertCodeEvalsToType(t, `{1}         `, `//rel.union({{1}, {1}, {1}})         `)
	AssertCodeEvalsToType(t, `{}          `, `//rel.union({})                      `)
	AssertCodeEvalsToType(t, `{1, 2, 3, 4}`, `//rel.union([{1, 2}, {2, 3}, {3, 4}])`)
	AssertCodeEvalsToType(t, `{1, 2, 3}   `, `//rel.union([{1}, {2}, {3}])         `)
	AssertCodeEvalsToType(t, `{1}         `, `//rel.union([{1}, {1}, {1}])         `)
	AssertCodeEvalsToType(t, `{}          `, `//rel.union([])                      `)
}
