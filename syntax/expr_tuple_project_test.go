package syntax

import "testing"

func TestTupleProjectExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(a: 1)            `, `(a: 1, b: 2, c: 3).|a|       `)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2)      `, `(a: 1, b: 2, c: 3).|a, b|    `)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2, c: 3)`, `(a: 1, b: 2, c: 3).|a, b, c| `)
	AssertCodesEvalToSameValue(t, `(b: 2, c: 3)      `, `(a: 1, b: 2, c: 3).~|a|      `)
	AssertCodesEvalToSameValue(t, `(c: 3)            `, `(a: 1, b: 2, c: 3).~|a, b|   `)
	AssertCodesEvalToSameValue(t, `()                `, `(a: 1, b: 2, c: 3).~|a, b, c|`)
}

func TestTupleProjectExprError(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, `lhs does not evaluate to tuple: {1: 1}`, `{1: 1}.|a|                      `)
	AssertCodeErrors(t, `lhs does not evaluate to tuple: {1: 1}`, `{1: 1}.~|a|                     `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).|d|          `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).|a, b, d|    `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).|a, b, c, d| `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).~|d|         `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).~|a, b, d|   `)
	AssertCodeErrors(t, `names are not subset of lhs: |a, b, c|`, `(a: 1, b: 2, c: 3).~|a, b, c, d|`)
}
