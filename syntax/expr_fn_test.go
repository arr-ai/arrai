package syntax

import "testing"

func TestExprFn(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `(\[x, y] x + y)([1, 2])`)
	AssertCodesEvalToSameValue(t, `3`, `(\z \[x, y] z/(x + y))(9, [1, 2])`)
	AssertCodeErrors(t, "", `(\[x, y] 42)([1, 2, 3])`)
	AssertCodeErrors(t, "", `(\[x, y] x + y)([1, 2, 3])`)
}

// TestClosureEqualityCapturedBindings pins the language-level rule for `=`
// on closures (#777): equal iff the same function and equal captured
// bindings, so `mk(1) = mk(1)` and `mk(1) = mk(2)` differ, sets keep
// closures that differ only in a capture, and a repeated pattern variable
// refuses to bind two of them.
func TestClosureEqualityCapturedBindings(t *testing.T) {
	t.Parallel()

	// The three reproductions from #777.
	AssertCodesEvalToSameValue(t, `[false, 11, 12]`,
		`let mk = \k \a a + k; let f = mk(1); let g = mk(2); [f = g, f(10), g(10)]`)
	AssertCodesEvalToSameValue(t, `2`, `let mk = \k \a a + k; {mk(1), mk(2)} count`)
	AssertCodeErrors(t, "the value of x is different in both scopes",
		`let mk = \k \a a + k; let [x, x] = [mk(1), mk(2)]; x(10)`)
	AssertCodesEvalToSameValue(t, `11`, `let mk = \k \a a + k; let [x, x] = [mk(1), mk(1)]; x(10)`)

	// Equal captures in different frames are equal; f = f stays true.
	AssertCodesEvalToSameValue(t, `[true, false]`, `let mk = \k \a a + k; [mk(1) = mk(1), mk(1) = mk(2)]`)
	AssertCodesEvalToSameValue(t, `true`, `let f = \a a + 1; f = f`)

	// Only bindings the body uses count.
	AssertCodesEvalToSameValue(t, `true`, `let mk = \k \a a; mk(1) = mk(2)`)

	// A set (and a dict keyed by closures) keeps closures that differ only
	// in a captured value.
	AssertCodesEvalToSameValue(t, `3`, `let mk = \k \a a + k; {mk(1), mk(2), mk(1), mk(3)} count`)
	AssertCodesEvalToSameValue(t, `[2, 'two']`,
		`let mk = \k \a a + k; let d = {mk(1): "one", mk(2): "two"}; [d count, d(mk(2))]`)

	// Stdlib references are free identifiers too, but bound in neither scope.
	AssertCodesEvalToSameValue(t, `[true, false]`,
		`let s = \k \x //str.upper(x) ++ k; [s("a") = s("a"), s("a") = s("b")]`)
}

// TestClosureEqualityRecursive covers closures that capture a `let rec`
// binding: the self-reference must not loop, and equality still follows the
// other captured values.
func TestClosureEqualityRecursive(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true`, `let rec f = \n cond {(n = 0): 0, _: f(n - 1)}; f = f`)
	AssertCodesEvalToSameValue(t, `[true, false, 3]`,
		`let g = \k (let rec f = \n cond {(n = 0): k, _: f(n - 1)}; f); [g(1) = g(1), g(1) = g(2), g(3)(5)]`)
	AssertCodesEvalToSameValue(t, `2`,
		`let g = \k (let rec f = \n cond {(n = 0): k, _: f(n - 1)}; f); {g(1), g(1), g(2)} count`)
	AssertCodesEvalToSameValue(t, `[true, false, true]`,
		`let g = \k (let rec eo = (
			ev: \n cond {(n = 0): k, _: eo.od(n - 1)},
			od: \n cond {(n = 0): !k, _: eo.ev(n - 1)},
		); eo.ev); [g(true) = g(true), g(true) = g(false), g(true)(4)]`)
}
