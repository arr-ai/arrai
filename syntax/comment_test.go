package syntax

import "testing"

func TestCommentLocations(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "(a: 1)", `
		(
			# this is a comment
			a: 1
		)
	`)
	AssertCodesEvalToSameValue(t, "{'a': 1}", `
		{
			# this is a comment
			'a': 1
		}
	`)
	AssertCodesEvalToSameValue(t, "{(a: 1, b: 1)}", `
		{
			# this is a comment
			|a, b|
			# this is another comment
			(1, 1)
		}
	`)
}
