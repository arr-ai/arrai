package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
)

func TestStdTuple(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `()`, `//tuple({})`)
	AssertCodesEvalToSameValue(t, `(a:1)`, `//tuple({"a":1})`)
	AssertCodesEvalToSameValue(t, `(a:1, b:2)`, `//tuple({"a":1, "b":2})`)

	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({"":1})`)
	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({{}:1})`)
	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({[]:1})`)

	AssertCodeErrors(t, "", `//tuple({0: 0})`)
	AssertCodeErrors(t, "", `//tuple({(): ()})`)
	AssertCodeErrors(t, "", `//tuple({1:2})`)
	AssertCodeErrors(t, "", `//tuple({[1]:2})`)
	AssertCodeErrors(t, "", `//tuple((a:1))`)
	AssertCodeErrors(t, "", `//tuple(42)`)
}

func TestStdDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `//dict(())`)
	AssertCodesEvalToSameValue(t, `{"a":1}`, `//dict((a:1))`)
	AssertCodesEvalToSameValue(t, `{"a":1, "b":2}`, `//dict((a:1, b:2))`)

	AssertCodeErrors(t, "", `//dict({42:43})`)
	AssertCodeErrors(t, "", `//dict(42)`)
}

func TestStdScope(t *testing.T) {
	scope := rel.EmptyScope.With("//",
		rel.NewTuple(rel.NewTupleAttr("a",
			rel.NewTupleAttr("b",
				rel.NewFloatAttr("c", 1),
			),
		)),
	)
	expected := rel.EmptyScope.With("//",
		rel.NewTuple(rel.NewTupleAttr("a",
			rel.NewTupleAttr("b",
				rel.NewFloatAttr("c", 2),
			),
		)),
	)

	out := addToScope(scope, []string{"a", "b", "c"}, func() rel.Value { return rel.NewNumber(2) })

	assert.True(t, expected.MustGet("//").(rel.Tuple).Equal(out.MustGet("//")))
}
