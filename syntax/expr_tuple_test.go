package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestTupleType(t *testing.T) {
	t.Parallel()

	AssertCodeEvalsToType(t, rel.StringCharTuple{}, `(@: 1, @char: 65)`)
	AssertCodeEvalsToType(t, rel.ArrayItemTuple{}, `(@: 1, @item: 2)`)
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
