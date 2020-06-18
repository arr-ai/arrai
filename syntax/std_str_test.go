package syntax

import (
	"testing"
)

func TestStrLower(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.lower("")              `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//str.lower("this is a test")`)
	AssertCodeErrors(t, "", `//str.lower(123)`)
}

func TestStrUpper(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.upper("")              `)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.upper("this is a test")`)
	AssertCodeErrors(t, "", `//str.upper(123)`)
}

func TestStrTitle(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""              `, `//str.title("")              `)
	AssertCodesEvalToSameValue(t, `"THIS IS A TEST"`, `//str.title("THIS IS A TEST")`)
	AssertCodesEvalToSameValue(t, `"ThIs Is A TeST"`, `//str.title("ThIs is A TeST")`)
	AssertCodesEvalToSameValue(t, `"This Is A Test"`, `//str.title("this is a test")`)
	AssertCodeErrors(t, "", `//str.title(123)`)
}

func TestStrExprStr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"/a+b+c/"`, `let arr = ['a', 'b', 'c']; $'/${arr::+}/'`)
	AssertCodesEvalToSameValue(t, `"/a++b+c/"`, `let arr = ['a', , 'b', 'c']; $'/${arr::+}/'`)
}
