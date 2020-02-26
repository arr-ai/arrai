package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSub(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.sub("this is not a test", "is not", "is")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is still a test"`,
		`//.str.sub("this is still a test", "doesn't matter", "hello there")`,
	)
	assertExprPanics(t, `//.str.sub("hello there", "test", 1)`)
}

func TestSplit(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`["this", "is", "a", "test"]`,
		`//.str.split("this is a test", " ")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`["this is a test"]`,
		`//.str.split("this is a test", ",")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`["th", " ", " a test"]`,
		`//.str.split("this is a test", "is")`,
	)
	assertExprPanics(t, `//.str.split("this is a test", 1)`)
}

func TestLower(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.lower("THIS IS A TEST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.lower("ThIs is A TeST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.lower("this is a test")`,
	)
	assertExprPanics(t, `//.str.lower(123)`)
}

func TestUpper(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.upper("THIS IS A TEST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.upper("ThIs is A TeST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.upper("this is a test")`,
	)
	assertExprPanics(t, `//.str.upper(123)`)
}

func TestTitle(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.title("THIS IS A TEST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"ThIs Is A TeST"`,
		`//.str.title("ThIs is A TeST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"This Is A Test"`,
		`//.str.title("this is a test")`,
	)
	assertExprPanics(t, `//.str.title(123)`)
}

func TestContains(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`true`,
		`//.str.contains("this is a test", "is a test")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`false`,
		`//.str.contains("this is a test", "is not a test")`,
	)
	assertExprPanics(t, `//.str.contains(123, 124)`)
}

func TestConcat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`""`,
		`//.str.concat([])`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.concat(["this", " is", " a", " test"])`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this"`,
		`//.str.concat(["this"])`,
	)
	assertExprPanics(t, `//.str.concat("this")`)
}

func TestJoin(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`""`,
		`//.str.join([], ",")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.join(["this", "is", "a", "test"], " ")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this"`,
		`//.str.join(["this"], ",")`,
	)
	assertExprPanics(t, `//.str.join("this", 2)`)
}

func assertExprPanics(t *testing.T, code string) {
	assert.Panics(t, func() {
		AssertCodesEvalToSameValue(
			t,
			`"doesn't matter"`,
			code,
		)
	})
}
