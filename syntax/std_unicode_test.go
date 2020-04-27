package syntax

import "testing"

func TestStdUnicodeUTF8Encode(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `//unicode.utf8.encode("")`)
	AssertCodesEvalToSameValue(t,
		`{|@,@byte| (0,97), (1,98), (2,99)}`,
		`//unicode.utf8.encode("abc")`,
	)
}
