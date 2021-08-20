package syntax

import "testing"

func TestJSONDecode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "()", `//encoding.json.decode('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.json.decode('{}')`)
	AssertCodesEvalToSameValue(t, "(a: [])", `//encoding.json.decode('[]')`)
	AssertCodeErrors(t, "", `//encoding.json.decode(123)`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.json.decode('123')`)

	json := testJSONString()
	expected := testArraiString()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(`+json+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(<<`+json+`>>)`)
}

func TestJSONDecode_NonStrict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "()", `//encoding.json.decoder((strict: false))('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.json.decoder((strict: false))('{}')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.json.decoder((strict: false))('[]')`)
	AssertCodeErrors(t, "", `//encoding.json.decoder((strict: false))(123)`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.json.decoder((strict: false))('123')`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.json.decoder((strict: false))('123')`)

	json := testJSONString()
	expected := testArraiStringLoose()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decoder((strict: false))(`+json+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decoder((strict: false))(<<`+json+`>>)`)
}

func TestJSONEncode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `<<'null\n'>>`, `//encoding.json.encode(())`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.json.encode({})`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.json.encode({123})`)
	AssertCodesEvalToSameValue(t, `<<'[]\n'>>`, `//encoding.json.encode((a: []))`)
	AssertCodesEvalToSameValue(t, `<<'123\n'>>`, `//encoding.json.encode(123)`)
	AssertCodesEvalToSameValue(t, `<<'[123]\n'>>`, `//encoding.json.encode([123])`)
	AssertCodesEvalToSameValue(t, `<<'{"key":123}\n'>>`, `//encoding.json.encode({"key": 123})`)
	AssertCodesEvalToSameValue(t, `<<'"abcde"\n'>>`, `//encoding.json.encode("abcde")`)
	AssertCodesEvalToSameValue(t, `<<'">"\n'>>`, `//encoding.json.encode(">")`)

	AssertCodeErrors(t, "", `//encoding.json.encode((a: [], s: ""))`)

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

func TestJSONEncodeSetOrder(t *testing.T) {
	t.Parallel()

	// Ensures that json encoder has deterministic output.

	AssertCodesEvalToSameValue(t,
		`<<'[{"a":1,"b":2},{"a":2,"b":1},{"a":2,"b":2}]\n'>>`,
		`//encoding.json.encoder((strict: false))({|a, b| (1, 2), (2, 1), (2, 2)})`,
	)

	AssertCodesEvalToSameValue(t,
		`<<'[1,2,3]\n'>>`,
		`//encoding.json.encoder((strict: false))({1, 2, 3})`,
	)
}

func TestJSONEncode_Config(t *testing.T) {
	t.Parallel()

	configTuple := `(
		prefix: '↘️',
		indent: '➡️',
		escapeHTML: true,
		strict: false,
	)`
	data := `(a: {"b": "c", "d": true}, bool: false, number: 0, set: {}, array: [], html: "<script/>")`
	expected := `<<'{
↘️➡️"a": {
↘️➡️➡️"b": "c",
↘️➡️➡️"d": true
↘️➡️},
↘️➡️"array": null,
↘️➡️"bool": null,
↘️➡️"html": "\\u003cscript/\\u003e",
↘️➡️"number": 0,
↘️➡️"set": null
↘️}
'>>`
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.encoder(`+configTuple+`)(`+data+`)`)
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

func testArraiStringLoose() string {
	return `{
    "a": "string",
    "b": 123,
    "c": 123.321,
    "d": [1, "string again", {}, {}],
    "e": {
      "f": {
        "g": "321"
      },
      "h": []
    },
    "i": (),
    "j": [true, {}],
    "k": {}
  }`
}
