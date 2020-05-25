package syntax

import "testing"

func TestSafeTail(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `1 `, `(a: 1).a?:42                                 `)
	AssertCodesEvalToSameValue(t, `42`, `(a: 1).b?:42                                 `)
	AssertCodesEvalToSameValue(t, `1 `, `{"a": 1}("a")?:42                            `)
	AssertCodesEvalToSameValue(t, `42`, `{"a": 1}("b")?:42                            `)
	AssertCodesEvalToSameValue(t, `1 `, `(a: (b: 1)).a?.b:42                          `)
	AssertCodesEvalToSameValue(t, `1 `, `{"a": {"b": 1}}("a")?("b"):42                `)
	AssertCodesEvalToSameValue(t, `1 `, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d.e?:42`)
	AssertCodesEvalToSameValue(t, `1 `, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d?.e:42`)
	AssertCodesEvalToSameValue(t, `42`, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d.f?:42`)
	AssertCodesEvalToSameValue(t, `42`, `let a = (b: (c: (d: (e: 1)))); a.b.c?.f?.e:42`)
	AssertCodesEvalToSameValue(t,
		`1`,
		`let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")("e")?:42`)
	AssertCodesEvalToSameValue(t,
		`1`,
		`let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")?("e"):42`)
	AssertCodesEvalToSameValue(t,
		`42`,
		`let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")("f")?:42`)
	AssertCodesEvalToSameValue(t,
		`42`,
		`let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("f")?("e"):42`)
	AssertCodesEvalToSameValue(t,
		`42`,
		`let a = {"b": (c: (d: {"e": 1}))}; a("b").c?.d.e:42`)

	AssertCodeErrors(t, `(a: 1).a?.c:42`, `(1).c: lhs must be a Tuple, not rel.Number`)
}
