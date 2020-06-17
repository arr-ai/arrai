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
		`let rec eo = (
     		even: \n n = 0 || eo.odd(n - 1),
			odd:  \n n != 0 && eo.even(n - 1),
			num: 6
   		);
		eo.even(eo.num)`,
		"Recursion requires a tuple of functions: "+
			"(even: ⦇(\\n ((n  = 0)) || («(eo.odd)»((n - 1))))⦈, num: 6, odd: ⦇(\\n ((n  != 0)) && («(eo.even)»((n - 1))))⦈)",
	)
	AssertCodeErrors(t,
		`let rec random = 1; random`,
		`Recursion does not support rel.Number`,
	)
	AssertCodeErrors(t,
		`let rec 1 = 1; 2`,
		`Recursion does not support rel.Number`,
	)
}
