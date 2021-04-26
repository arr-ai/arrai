package syntax

import (
	"fmt"
	"testing"
)

func TestBytesExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, toBytes(`"hello\n"`), `<<"hello", 10>>                             `)
	AssertCodesEvalToSameValue(t, toBytes(`"abc"    `), `<<97, 98, 99>>                              `)
	AssertCodesEvalToSameValue(t, toBytes(`"ABC"    `), `<<%A, %B, %C>>                              `)
	AssertCodesEvalToSameValue(t, toBytes(`"hello"  `), `<<"hello">>                                 `)
	AssertCodesEvalToSameValue(t, toBytes(`""       `), `<<>>                                        `)
	AssertCodesEvalToSameValue(t, toBytes(`""       `), `<<''>>                                      `)
	AssertCodesEvalToSameValue(t, toBytes(`"a"      `), `<<'', 97>>                                  `)
	AssertCodesEvalToSameValue(t, toBytes(`"aa"     `), `let x = 97; <<x, x>>                        `)
	AssertCodesEvalToSameValue(t, toBytes(`"bcd"    `), `<<("abc" >> . + 1)>>                        `)
	AssertCodesEvalToSameValue(t, toBytes(`"abc"    `), `<<({|@, @char| (0, 97), (1, 98), (2, 99)})>>`)

	AssertCodeErrors(t,
		"BytesExpr.Eval: Number does not represent a byte: 256",
		`<<256>>`)
	AssertCodeErrors(t,
		"BytesExpr.Eval: Number does not represent a byte: -2",
		`<<(-2)>>`)
	AssertCodeErrors(t,
		"BytesExpr.Eval: offset string is not supported: offset",
		`<<(2\"offset")>>`)
	AssertCodeErrors(t,
		"BytesExpr.Eval: set is not supported",
		`<<({1, 2, 3})>>`)
	AssertCodeErrors(t,
		"BytesExpr.Eval: tuple is not supported",
		`<<((a: 1))>>`)
	AssertCodeErrors(t,
		"BytesExpr.Eval: array is not supported",
		`<<([1, 2, 3])>>`)
}

func toBytes(s string) string {
	return fmt.Sprintf("%s => (:.@, @byte: .@char)", s)
}
