package syntax

import (
	"strings"
	"testing"
)

func TestXStringSimple(t *testing.T) {
	AssertCodesEvalToSameValue(t, `""               `, `$""`)
	AssertCodesEvalToSameValue(t, `"42"             `, `$":{6*7}:"`)
	AssertCodesEvalToSameValue(t, `"a42z"           `, `$"a:{6*7}:z"`)
	AssertCodesEvalToSameValue(t, `"a00042z"        `, `$"a:{6*7:05d}:z"`)
	AssertCodesEvalToSameValue(t, `"a001, 002, 003z"`, `$"a:{[1, 2, 3]:03d:, }:z"`)
	AssertCodesEvalToSameValue(t, `"a42k3.142z"     `, `$"a:{6*7}:k:{//.math.pi:.3f}:z"`)
}

func TestXStringStrings(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"hello"`, `$":{'hello'}:"`)
}

func TestXStringIndent(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"a\nb"`, "$`\n  a\n  b`")
	AssertCodesEvalToSameValue(t, `"a\nb\n  c\nd"`, "$'\n  a\n  b\n    c\n  d'")
}

func TestXStringWS(t *testing.T) {
	AssertCodesEvalToSameValue(t, `"1 2"`, `$":{1}: :{2}:"`)
	AssertCodesEvalToSameValue(t, `"1\n2"`, "$'\n  :{1}:\n  :{2}:'")
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
		`[(name: "custid", type: "int")] -> $":{. >> $"Get:{//.str.title(.name)}:()"::}:"`,
	)
}

func TestXStringNested(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		strings.ReplaceAll(
			`"type Customer interface {
				IsCustomer()
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
				Is:{.name}:()
				:{.fields >> $"Get:{//.str.title(.name)}:() :{.type}:"::\n\i}:
			}"`,
	)
}
