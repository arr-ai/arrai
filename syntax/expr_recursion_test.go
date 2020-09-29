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
		`Recursion does not support rel.Number`,
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
