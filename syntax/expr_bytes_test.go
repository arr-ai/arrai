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
		`<<256>>`,
		"BytesExpr.Eval: Number does not represent a byte: 256")
	AssertCodeErrors(t,
		`<<(-2)>>`,
		"BytesExpr.Eval: Number does not represent a byte: -2")
	AssertCodeErrors(t,
		`<<(2\"offset")>>`,
		"BytesExpr.Eval: offsetted String is not supported: offset")
	AssertCodeErrors(t,
		`<<({1, 2, 3})>>`,
		"BytesExpr.Eval: Set {1, 2, 3} is not supported")
	AssertCodeErrors(t,
		`<<((a: 1))>>`,
		"BytesExpr.Eval: *rel.GenericTuple is not supported")
	AssertCodeErrors(t,
		`<<([1, 2, 3])>>`,
		"BytesExpr.Eval: rel.Array is not supported")
}

func toBytes(s string) string {
	return fmt.Sprintf("%s => (:.@, @byte: .@char)", s)
}
