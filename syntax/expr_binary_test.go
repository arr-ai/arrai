package syntax

import (
	"testing"
)

func TestWhereExpr(t *testing.T) {
	t.Parallel()
	s := `{|a,b| (3,41), (2,42), (1,43)}`
	// defer trace().revert()
	AssertCodesEvalToSameValue(t, `{(a:3, b:41)}`, s+` where .a=3`)
}

func TestRelationCall(t *testing.T) {
	t.Parallel()
	s := `{"key": "val"}("key")`
	AssertCodesEvalToSameValue(t, `"val"`, s)
}

func TestOpsAddArrowForTuples(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a: 1, b: 2, c: 3, d: 4)`, `(a: 1, b: 2) +> (c: 3, d: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 2), c: (b: 1), d: (c: 4))`, `(a: 1, b: (c: 2)) +> (c: (b: 1), d: (c: 4))`)

	AssertCodesEvalToSameValue(t, `(a: 1, b: 3, c: 4)`, `(a: 1, b: 2) +> (b: 3, c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 3, c: 4)`, `(a: 1, b: (c: 2)) +> (b: 3, c: 4)`)

	AssertCodesEvalToSameValue(t, `(a: 1, b: (c: 4), c: 4)`, `(a: 1, b: (c: 2)) +> (b: (c: 4), c: 4)`)
	AssertCodesEvalToSameValue(t, `(a: (b: 1), b: (c: 4))`, `(a: 1, b: (c: 2)) +> (a: (b: 1), b: (c: 4))`)
}

func TestOpsAddArrowForDicts(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 3, 'd': 4}`, `{'a': 1, 'b': 2} +> {'b': 3, 'd': 4}`)
	AssertCodesEvalToSameValue(t, `{'a': 1, 'b': 3, 'd': 4}`, `{'a': 1, 'b': {'a': 2}} +> {'b': 3, 'd': 4}`)
	AssertCodesEvalToSameValue(t, `{'a': {'c': 2}}`, `{'a': {'b': 1}} +> {'a': {'c': 2}}`)

	AssertCodesEvalToSameValue(t, `{'a': 'A1', 'b': 'A3', 'd': 'A4'}`, `{'a': 'A1', 'b': 'A2'} +> {'b': 'A3', 'd': 'A4'}`)
	AssertCodesEvalToSameValue(t, `{'a': {'c': 'ABC2'}}`, `{'a': {'b': 'ABC1'}} +> {'a': {'c': 'ABC2'}}`)

	AssertCodesEvalToSameValue(t, `{'a': 'A1', 'b': 'A2', 'c': 'A3', 'd': 'A4'}`, `{'a': 'A1', 'b': 'A2'} +> {'c': 'A3', 'd': 'A4'}`)
	AssertCodesEvalToSameValue(t, `{'a': {'b': 'ABC1'}, 'b': {'c': 'ABC2'}}`, `{'a': {'b': 'ABC1'}} +> {'b': {'c': 'ABC2'}}`)
}

func TestOpsAddArrowForError(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, "", `{1, 2, 3} +> {4, 5, 6}`)
	AssertCodeErrors(t, "", `1 +> 4`)
}
