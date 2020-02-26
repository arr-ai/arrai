package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSub(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"javascript should be illegal to use"`,
		`//.str.sub("javascript is easy to use", "is easy", "should be illegal")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"javascript shouldn't exist"`,
		`//.str.sub("javascript shouldn't exist", "doesn't matter", "hello there")`,
	)
	assert.Panics(t, func() {
		AssertCodesEvalToSameValue(
			t,
			`"hello there"`,
			`//.str.sub("hello there", "test", 1)`,
		)
	})
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
	assert.Panics(t, func() {
		AssertCodesEvalToSameValue(
			t,
			`["this is a test"]`,
			`//.str.split("this is a test", 1)`,
		)
	})
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
	assert.Panics(t, func() {
		AssertCodesEvalToSameValue(
			t,
			`"doesn't matter"`,
			`//.str.lower(123)`,
		)
	})
}

func TestUpper(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.lower("THIS IS A TEST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"THIS IS A TEST"`,
		`//.str.lower("ThIs is A TeST")`,
	)
	AssertCodesEvalToSameValue(
		t,
		`"this is a test"`,
		`//.str.lower("this is a test")`,
	)
	assert.Panics(t, func() {
		AssertCodesEvalToSameValue(
			t,
			`"doesn't matter"`,
			`//.str.lower(123)`,
		)
	})
}

func TestTitle(t *testing.T) {
	t.Parallel()

}

func TestContains(t *testing.T) {
	t.Parallel()

}

func TestConcat(t *testing.T) {
	t.Parallel()

}

func TestJoin(t *testing.T) {
	t.Parallel()

}
