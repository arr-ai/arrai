package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestTupleType(t *testing.T) {
	t.Parallel()

	AssertCodeEvalsToType(t, rel.StringCharTuple{}, `(@: 1, @char: 65)`)
	AssertCodePanics(t, `(@: 1, @char: "x")`)
	AssertCodePanics(t, `(@: {}, @char: 65)`)

	AssertCodeEvalsToType(t, rel.ArrayItemTuple{}, `(@: 1, @item: 2)`)
	AssertCodePanics(t, `(@: {}, @item: 2)`)

	AssertCodeEvalsToType(t, rel.DictEntryTuple{}, `(@: {1, 2}, @value: 2)`)
}

func TestTupleGet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(a: 1, b: 42).b`)
	AssertCodesEvalToSameValue(t, `42`, `(a: 1, 'b': 42)."b"`)
	AssertCodesEvalToSameValue(t, `42`, `(a: 1, "b": 42).'b'`)
	AssertCodesEvalToSameValue(t, `42`, "(a: 1, `b`: 42).`b`")
	AssertCodesEvalToSameValue(t, `42`, `(a: 1, '👋': 42)."👋"`)
	AssertCodesEvalToSameValue(t, `42`, `(a: 1, '': 42).""`)
}

func TestTupleCallGet(t *testing.T) {
	AssertCodesEvalToSameValue(t, `2`, `(a: \x (b: x)).a(2).b`)
	AssertCodesEvalToSameValue(t, `2`, `let t = (a: \x (b: x)); t.a(2).b`)
}

func TestTupleLiteral(t *testing.T) {
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `let x = 1; let y = 2; (x: x, y: y)`)
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `let x = 1; let y = 2; (:x, :y)`)
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `let t = (x: 1, y: 2); (:t.x, :t.y)`)
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `(x: 1, y: 2) -> (:.x, :.y)`)
}
