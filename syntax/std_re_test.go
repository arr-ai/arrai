package syntax

import "testing"

func TestStdRe(t *testing.T) {
	AssertCodeErrors(t, "//re.compile", "//re.compile(`x(y))`)")
}

func TestStdReMatch(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a1', 1\'1'], [2\'b2', 3\'2'], [4\'c3', 5\'3']]`,
		"//re.compile(`.(\\d)`).match('a1b2c3')")
	AssertCodesEvalToSameValue(t,
		`{}`,
		"//re.compile(`x(\\d)`).match('a1b2c3')")
	AssertCodesEvalToSameValue(t,
		`{}`,
		"//re.compile(`x(\\d)`).match('a1b2c3')")
}

func TestStdReMismatch(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a1'] | 2\[1\'1'], [2\'b~2', 3\'~', 4\'2'], [5\'c3'] | 2\[6\'3']]`,
		"//re.compile(`[a-z](~)?(\\d)`).match('a1b~2c3')")
}

func TestStdReSub(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`'-1-2-3'`,
		"//re.compile(`.(\\d)`).sub(`-$1`, 'a1b2c3')")
	AssertCodesEvalToSameValue(t,
		`'a1b2c3'`,
		"//re.compile(`x(\\d)`).sub(`-$1`, 'a1b2c3')")
}

func TestStdReSubf(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`'A1B2C3'`,
		"//re.compile(`.(\\d)`).subf(//str.upper, 'a1b2c3')")
	AssertCodesEvalToSameValue(t,
		`'a1b2c3'`,
		"//re.compile(`x(\\d)`).subf(//str.upper, 'a1b2c3')")
}
