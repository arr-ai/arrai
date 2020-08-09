package syntax

import "testing"

func TestSafeTail(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `1 `, `(a: 1).a?:42                                                     `)
	AssertCodesEvalToSameValue(t, `42`, `(a: 1).b?:42                                                     `)
	AssertCodesEvalToSameValue(t, `1 `, `{"a": 1}("a")?:42                                                `)
	AssertCodesEvalToSameValue(t, `42`, `{"a": 1}("b")?:42                                                `)
	AssertCodesEvalToSameValue(t, `1 `, `(a: (b: 1)).a?.b:42                                              `)
	AssertCodesEvalToSameValue(t, `1 `, `{"a": {"b": 1}}("a")?("b"):42                                    `)
	AssertCodesEvalToSameValue(t, `1 `, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d.e?:42                    `)
	AssertCodesEvalToSameValue(t, `1 `, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d?.e:42                    `)
	AssertCodesEvalToSameValue(t, `42`, `let a = (b: (c: (d: (e: 1)))); a.b.c?.d.f?:42                    `)
	AssertCodesEvalToSameValue(t, `42`, `let a = (b: (c: (d: (e: 1)))); a.b.c?.f?.e:42                    `)
	AssertCodesEvalToSameValue(t, `1 `, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b", "c", "d", "e")?:42 `)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b", "c", "x", "e")?:42 `)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b", "c", "x")?("e"):42 `)
	AssertCodesEvalToSameValue(t, `1 `, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")("e")?:42`)
	AssertCodesEvalToSameValue(t, `1 `, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")?("e"):42`)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("d")("f")?:42`)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": {"c": {"d": {"e": 1}}}}; a("b")("c")?("f")?("e"):42`)
	AssertCodesEvalToSameValue(t, `1 `, `let a = {"b": (c: (d: {"e": 1}))}; a("b").c?.d("e")?:42          `)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": (c: (d: {"e": 1}))}; a("b").c?.d("f")?:42          `)
	AssertCodesEvalToSameValue(t, `42`, `let a = {"b": (c: (d: {"e": 1}))}; a("b").c?.e?("f")?:42         `)

	AssertCodeErrors(t, `(1).c: lhs must be a Tuple, not rel.Number`, `(a: 1).a?.c:42               `)
	AssertCodeErrors(t, `Missing attr "c" (available: |b|)`, `(a: (b: 1)).a?.c:42          `)
	AssertCodeErrors(t, `Call: no return values for input c from set {b: 1}`, `{"a": {"b": 1}}("a")?("c"):42`)
}
