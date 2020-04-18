package syntax

import "testing"

func TestReMatch(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`let o = \off \s s => (@:.@+off, :.@char);
		 [['a1', o(1, '1')], [o(2, 'b2'), o(3, '2')], [o(4, 'c3'), o(5, '3')]]`,
		"//.re.compile(`.(\\d)`).match('a1b2c3')")
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`{}`,
		"//.re.compile(`x(\\d)`).match('a1b2c3')")
}

func TestReSub(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`'-1-2-3'`,
		"//.re.compile(`.(\\d)`).sub(`-$1`, 'a1b2c3')")
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`'a1b2c3'`,
		"//.re.compile(`x(\\d)`).sub(`-$1`, 'a1b2c3')")
}

func TestReSubf(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`'A1B2C3'`,
		"//.re.compile(`.(\\d)`).subf(//.str.upper, 'a1b2c3')")
	AssertCodesEvalToSameValue(t,
		// TODO: Use n\'abc' syntax when implemented.
		`'a1b2c3'`,
		"//.re.compile(`x(\\d)`).subf(//.str.upper, 'a1b2c3')")
}
