package syntax

import (
	"strings"
	"testing"
)

func TestXStringSimple(t *testing.T) {
	AssertCodesEvalToSameValue(t, `""               `, `$""`)
	AssertCodesEvalToSameValue(t, `"42"             `, `$":{6*7}:"`)
	AssertCodesEvalToSameValue(t, `"a42z"           `, `$"a:{6*7}:z"`)
	AssertCodesEvalToSameValue(t, `"a00042z"        `, `$"a:{05d:6*7}:z"`)
	AssertCodesEvalToSameValue(t, `"a001, 002, 003z"`, `$"a:{03d*:[1, 2, 3]:, }:z"`)
	AssertCodesEvalToSameValue(t, `"a42k3.142z"     `, `$"a:{6*7}:k:{.3f://.math.pi}:z"`)
}

func TestXStringStrings(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"hello"`, `$":{'hello'}:"`)
}

func TestXStringMap(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`"Getcustid() int"`,
		`(name: "custid", type: "int") -> $"Get:{.name}:() :{.type}:"`,
	)
}

func TestXStringMap2(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`"GetCustid()"`,
		`[(name: "custid", type: "int")] -> $":{*:. >> $"Get:{//.str.title(.name)}:()":\n}:"`,
	)
}

func TestXStringNested(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		strings.ReplaceAll(
			`"type Customer interface {
				GetCustid() int
				GetDob() date
				GetAlias() string
			}"`, "\n\t\t\t", "\n",
		),
		`(name: "Customer", fields: [
			(name: "custid", type: "int"   ),
			(name: "dob",    type: "date"  ),
			(name: "alias",  type: "string"),
		]) -> $"
			type :{.name}: interface {
				:{*:.fields >> $"Get:{//.str.title(.name)}:() :{.type}:":\n\i}:
			}"`,
	)
}
