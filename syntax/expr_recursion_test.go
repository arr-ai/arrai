package syntax

import (
	"testing"
)

func TestRecursionExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`55`,
		`let rec fib = \n cond {(n = 0): 0, (n = 1): 1, (n > 0): fib(n-1) + fib(n-2), _: 0};
		fib(10)`,
	)
	AssertCodesEvalToSameValue(t,
		`true`,
		`let rec eo = (
     		even: \n n = 0 || eo.odd(n - 1),
      		odd:  \n n != 0 && eo.even(n - 1),
   		);
   		eo.even(6)`,
	)
	AssertCodeErrors(t,
		`Recursion requires a tuple of functions: `+
			`(even: (\n ((n = 0)) || («(eo.odd)»((n - 1)))), num: 6, odd: (\n ((n != 0)) && («(eo.even)»((n - 1)))))`,
		`let rec eo = (
     		even: \n n = 0 || eo.odd(n - 1),
			odd:  \n n != 0 && eo.even(n - 1),
			num: 6
   		);
		eo.even(eo.num)`,
	)
	AssertCodeErrors(t,
		`Recursion does not support number`,
		`let rec random = 1; random`,
	)
	AssertCodeErrors(t,
		`let rec parameter must be IDENT, not 1`,
		`let rec 1 = 1; 2`,
	)
	// to test compile variables with the prefix rec
	AssertCodesEvalToSameValue(t, `1`, `let recTest = 1; recTest`)
	// FIXME: requires more complex grammar, or maybe this should be a keyword shouldn't be used as a variable
	// AssertCodesEvalToSameValue(t, `1`, `let rec = 1; rec`)
}

// The direct (non-combinator) implementation ties the knot through a scope
// cell. These cases cover the ways self-reference can bite: formatting,
// hashing/equality, reading the binding before it exists, shadowing, and
// recursive closures escaping into data structures.
func TestRecursionExprDirect(t *testing.T) {
	t.Parallel()

	// Deep recursion is plain calls, not combinator plumbing.
	AssertCodesEvalToSameValue(t, `500500`,
		`let rec sum = \n cond {(n = 0): 0, _: n + sum(n - 1)}; sum(1000)`)

	// A recursive closure can be formatted without looping forever.
	AssertCodesEvalToSameValue(t, `true`,
		`let rec f = \n cond {(n = 0): 0, _: f(n - 1)}; $`+"`${f}`"+` count > 0`)

	// ... and hashed: sets, arrays and dict values of recursive closures.
	AssertCodesEvalToSameValue(t, `3`,
		`let rec f = \n cond {(n = 0): 0, _: f(n - 1)}; ({f} | {f}) count + [f, f] count`)
	AssertCodesEvalToSameValue(t, `0`,
		`let rec f = \n cond {(n = 0): 0, _: f(n - 1)}; {'f': f}('f')(5)`)

	// Mutual recursion through a tuple, with the tuple captured in data.
	AssertCodesEvalToSameValue(t, `[true, false]`,
		`let rec eo = (even: \n n = 0 || eo.odd(n - 1), odd: \n n != 0 && eo.even(n - 1));
		let fns = [eo]; [fns(0).even(4), fns(0).even(3)]`)

	// Using the name before the binding exists is an error, not a crash.
	AssertCodeErrors(t, `recursive binding "x" used before it is defined`,
		`let rec x = (f: \n n) +> (g: x.f); x`)

	// The recursive name shadows an outer binding only inside its own body.
	AssertCodesEvalToSameValue(t, `[10, 3]`,
		`let f = \n 10;
		let g = let rec f = \n cond {(n = 0): 0, _: 1 + f(n - 1)}; f;
		[f(3), g(3)]`)

	// A recursive function returned from a function keeps working after the
	// defining scope is gone.
	AssertCodesEvalToSameValue(t, `120`,
		`let mk = \_ let rec fact = \n cond {(n <= 1): 1, _: n * fact(n - 1)}; fact;
		mk(())(5)`)

	// Recursion inside a where predicate and a >> map.
	AssertCodesEvalToSameValue(t, `[1, 2, 4, 8]`,
		`let rec pow2 = \n cond {(n = 0): 1, _: 2 * pow2(n - 1)}; [0, 1, 2, 3] >> pow2(.)`)
}
