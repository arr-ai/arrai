package syntax

import "testing"

// `rel where .attr = key` is answered from the relation's attribute index
// when the predicate has that shape. These cases pin the semantics of the
// indexed path against the scan it replaces.
func TestWhereEqAttrIndex(t *testing.T) {
	t.Parallel()

	const rel = `let r = {(a: 1, b: 'x'), (a: 2, b: 'y'), (a: 2, b: 'z'), (a: 3, b: 'x')};`

	// Both operand orders, identifier and literal keys.
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y'), (a: 2, b: 'z')}`, rel+`let k = 2; r where .a = k`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y'), (a: 2, b: 'z')}`, rel+`let k = 2; r where k = .a`)
	AssertCodesEvalToSameValue(t, `{(a: 1, b: 'x'), (a: 3, b: 'x')}`, rel+`r where .b = 'x'`)

	// No matching key, and a key of the wrong type, both yield the empty set.
	AssertCodesEvalToSameValue(t, `{}`, rel+`let k = 9; r where .a = k`)
	AssertCodesEvalToSameValue(t, `{}`, rel+`let k = 'x'; r where .a = k`)

	// An attribute the relation doesn't have falls back to the scan, which
	// fails per row exactly as before.
	AssertCodeErrors(t, ``, rel+`let k = 2; r where .c = k`)

	// Other comparison operators are never indexed.
	AssertCodesEvalToSameValue(t, `{(a: 1, b: 'x'), (a: 3, b: 'x')}`, rel+`let k = 2; r where .a != k`)
	AssertCodesEvalToSameValue(t, `{(a: 3, b: 'x')}`, rel+`let k = 2; r where .a > k`)

	// The key is resolved in the predicate's own scope, including inside a
	// function and after shadowing.
	AssertCodesEvalToSameValue(t, `[1, 2]`,
		rel+`let f = \k (r where .a = k) count; let k = 1; [f(k), f(2)]`)
	AssertCodesEvalToSameValue(t, `{(a: 3, b: 'x')}`,
		rel+`let k = 1; let k = 3; r where .a = k`)

	// The same relation queried repeatedly (the index is memoised) and the
	// index of a derived relation.
	AssertCodesEvalToSameValue(t, `[1, 2, 1]`,
		rel+`[1, 2, 3] >> \i (r where .a = i) count`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y')}`,
		rel+`let k = 2; (r where .b != 'z') where .a = k`)

	// Non-relation sets still work through the generic path.
	AssertCodesEvalToSameValue(t, `{(a: 2)}`, `let k = 2; {(a: 1), (a: 2), (b: 2)} where (.).a?:0 = k`)
	AssertCodesEvalToSameValue(t, `{2}`, `let k = 2; {1, 2, 3} where . = k`)
}

// Conjunctive `.attr = k && .other = v` is answered by intersecting per-column
// indexes (🎯T29.1). Mixed predicates still fall back to the scan.
func TestWhereConjEqAttrIndex(t *testing.T) {
	t.Parallel()

	const rel = `let r = {(a: 1, b: 'x', c: 1), (a: 2, b: 'y', c: 1), (a: 2, b: 'z', c: 2), (a: 3, b: 'x', c: 1)};`

	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y', c: 1)}`, rel+`r where .a = 2 && .b = 'y'`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y', c: 1)}`, rel+`r where .b = 'y' && .a = 2`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'z', c: 2)}`,
		rel+`let k = 2; let v = 'z'; r where .a = k && .b = v`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y', c: 1)}`,
		rel+`r where .a = 2 && .b = 'y' && .c = 1`)
	AssertCodesEvalToSameValue(t, `{}`, rel+`r where .a = 2 && .b = 'x'`)
	AssertCodesEvalToSameValue(t, `{}`, rel+`let k = 9; r where .a = k && .b = 'y'`)

	// Missing attr still errors on the scan fallback.
	AssertCodeErrors(t, ``, rel+`r where .a = 2 && .d = 1`)

	// A non-eq conjunct is not indexed as a conjunction; the scan is correct.
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y', c: 1)}`, rel+`r where .a = 2 && .b != 'z'`)

	// Pushdown through an unchanged project of both attrs.
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 'y')}`,
		rel+`(r => (a: .a, b: .b)) where .a = 2 && .b = 'y'`)
}

// Column predicates scan the arena (🎯T29.5).
func TestWhereColumnPred(t *testing.T) {
	t.Parallel()

	const rel = `let r = {(a: 1, rest: 1, t: 'x'), (a: 2, rest: 0, t: 'y'), (a: 3, rest: 1, t: 'z')};`
	AssertCodesEvalToSameValue(t, `{(a: 1, rest: 1, t: 'x'), (a: 3, rest: 1, t: 'z')}`, rel+`r where .rest`)
	AssertCodesEvalToSameValue(t, `{(a: 2, rest: 0, t: 'y')}`, rel+`r where !.rest`)
	AssertCodesEvalToSameValue(t, `{(a: 1, rest: 1, t: 'x'), (a: 3, rest: 1, t: 'z')}`, rel+`r where .t != 'y'`)
	AssertCodesEvalToSameValue(t, `{(a: 1, rest: 1, t: 'x'), (a: 3, rest: 1, t: 'z')}`,
		rel+`let k = 'y'; r where k != .t`)
	AssertCodesEvalToSameValue(t, `{(a: 1, rest: 1, t: 'x')}`, rel+`r where .t <: {'x', 'w'}`)
	AssertCodesEvalToSameValue(t, `{(a: 2, rest: 0, t: 'y'), (a: 3, rest: 1, t: 'z')}`,
		rel+`r where .t !<: {'x'}`)
	AssertCodesEvalToSameValue(t, `{(a: 1, rest: 1, t: 'x'), (a: 3, rest: 1, t: 'z')}`,
		rel+`let enum = {(typeName: 'y')}; r where .t !<: (enum => .typeName)`)
}

// Function application returns its single result directly. These cases pin
// the result kinds that used to round-trip through a one-element set.
func TestCallReturnsResultDirectly(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(a: 1, b: 2)`, `(\x (a: x, b: 2))(1)`)
	AssertCodesEvalToSameValue(t, `(@: 1, @item: 'x')`, `(\x (@: x, @item: 'x'))(1)`)
	AssertCodesEvalToSameValue(t, `{}`, `(\x {})(1)`)
	AssertCodesEvalToSameValue(t, `3`, `(\f f(1))(\x x + 2)`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `(\x [x, x + 1])(1)`)
	AssertCodesEvalToSameValue(t, `{'k': 1}`, `(\x {'k': x})(1)`)
	AssertCodesEvalToSameValue(t, `'ab'`, `(\x x ++ 'b')('a')`)
	AssertCodesEvalToSameValue(t, `0`, `//math.sin(0)`)
	AssertCodeErrors(t, ``, `({1: 2} | {1: 1})(1)`)
	AssertCodeErrors(t, ``, `{1: 2}(3)`)
}
