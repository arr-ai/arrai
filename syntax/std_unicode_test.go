package syntax

import "testing"

func TestStdUnicodeUTF8Encode(t *testing.T) {
	AssertCodesEvalToSameValue(t, `<<>>          `, `//unicode.utf8.encode("")   `)
	AssertCodesEvalToSameValue(t, `<<97, 98, 99>>`, `//unicode.utf8.encode("abc")`)
}
