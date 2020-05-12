package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrLower(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.lower("")              `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("this is a test")`)
	assertExprPanics(t, `//str.lower(123)`)
}

func TestStrUpper(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.upper("")              `)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("this is a test")`)
	assertExprPanics(t, `//str.upper(123)`)
}

func TestStrTitle(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.title("")              `)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.title("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"ThIs Is A TeST"`, `//str.title("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"This Is A Test"`, `//str.title("this is a test")`)
	assertExprPanics(t, `//str.title(123)`)
}

func assertExprPanics(t *testing.T, code string) {
	assert.Panics(t, func() { AssertCodesEvalToSameValue(t, `"doesn't matter"`, code) })
}
