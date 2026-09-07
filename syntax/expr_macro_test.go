package syntax

import "testing"

// A macro's result is a Value evaluated eagerly at parse time, not a compiled
// Expr, so it carried whatever generic Source() its own type defaults to
// instead of the embed's real location - the same bug class as the bytes-literal
// scanner bug, just at a third call site (compileMacro).
func TestMacroResultHasRealSource(t *testing.T) {
	t.Parallel()
	code := `
		let date = (
			@grammar: {://grammar.lang.wbnf:date -> y=\d{4} "-" m=\d{2} "-" d=\d{2};:},
			@transform: (date: \ast ast -> (year: .y, month: .m, day: .d) :> //eval.value(.''))
		);
		{:date:2020-06-09:} = (day: 9, month: 6, year: 2020)
	`
	AssertCodesEvalToSameValue(t, `true`, code)
}
