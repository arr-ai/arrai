package syntax

import "testing"

func TestJSONDecode(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, "()", `//encoding.json.decode('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.json.decode('{}')`)
	AssertCodesEvalToSameValue(t, "(a: [])", `//encoding.json.decode('[]')`)
	AssertCodeErrors(t, "", `//encoding.json.decode(123)`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.json.decode('123')`)

	expected := testArraiString()
	encoding := testJSONString()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(`+encoding+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(<<`+encoding+`>>)`)
}

func TestJSONEncode(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `<<'null\n'>>`, `//encoding.json.encode(())`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.json.encode({})`)
	AssertCodesEvalToSameValue(t, `<<'[]\n'>>`, `//encoding.json.encode((a: []))`)
	AssertCodesEvalToSameValue(t, `<<'123\n'>>`, `//encoding.json.encode(123)`)
	AssertCodesEvalToSameValue(t, `<<'[123]\n'>>`, `//encoding.json.encode([123])`)
	AssertCodesEvalToSameValue(t, `<<'"abcde"\n'>>`, `//encoding.json.encode("abcde")`)
	AssertCodesEvalToSameValue(t, `<<'">"\n'>>`, `//encoding.json.encode(">")`)

	encoding := testArraiString()
	expected := `<<'{"a":"string","b":123,"c":123.321,"d":[1,"string again",[],{}],"e":{"f":{"g":"321"},"h":[]},"i":null,"j":[true,false],"k":""}\n'>>` //nolint:lll

	AssertCodesEvalToSameValue(t, expected, `//encoding.json.encode(`+encoding+`)`)
}

func TestJSONEncodeIndent(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `<<'null\n'>>`, `//encoding.json.encode_indent(())`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.json.encode_indent({})`)
	AssertCodesEvalToSameValue(t, `<<'[]\n'>>`, `//encoding.json.encode_indent((a: []))`)
	AssertCodesEvalToSameValue(t, `<<'123\n'>>`, `//encoding.json.encode_indent(123)`)
	AssertCodesEvalToSameValue(t, `<<'">"\n'>>`, `//encoding.json.encode_indent(">")`)

	encoding := testArraiString()
	expected := `<<` + testJSONString() + `>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.json.encode_indent(`+encoding+`)`)
}

func testJSONString() string {
	return `'{
  "a": "string",
  "b": 123,
  "c": 123.321,
  "d": [
    1,
    "string again",
    [],
    {}
  ],
  "e": {
    "f": {
      "g": "321"
    },
    "h": []
  },
  "i": null,
  "j": [
    true,
    false
  ],
  "k": ""
}\n'`
}

func testArraiString() string {
	return `{
    "a": (s: "string"),
    "b": 123,
    "c": 123.321,
    "d": (a: [1, (s: "string again"), (a: []), {}]),
    "e": {
      "f": {
        "g": (s: "321")
      },
      "h": (a: [])
    },
    "i": (),
    "j": (a: [(b: true), (b: false)]),
    "k": (s: {})
  }`
}
