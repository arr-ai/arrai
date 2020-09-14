package syntax

import "testing"

func TestJSONDecode(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, "()", `//encoding.json.decode('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.json.decode('{}')`)
	AssertCodesEvalToSameValue(t, "(a: [])", `//encoding.json.decode('[]')`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.json.decode(123)`)

	expected := arraiString()
	encoding := jsonString()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(`+encoding+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(<<`+encoding+`>>)`)
}

func TestJSONEncode(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `<<'null'>>`, `//encoding.json.encode(())`)
	AssertCodesEvalToSameValue(t, `<<'{}'>>`, `//encoding.json.encode({})`)
	AssertCodesEvalToSameValue(t, `<<'[]'>>`, `//encoding.json.encode((a: []))`)
	AssertCodesEvalToSameValue(t, `<<'123'>>`, `//encoding.json.encode(123)`)

	encoding := arraiString()
	expected := `<<'{"a":"string","b":123,"c":123.321,"d":[1,"string again",[],{}],"e":{"f":{"g":"321"},"h":[]},"i":null,"j":[true,false],"k":""}'>>` //nolint:lll

	AssertCodesEvalToSameValue(t, expected, `//encoding.json.encode(`+encoding+`)`)
}

func jsonString() string {
	return `'{
		"a": "string",
		"b": 123,
		"c": 123.321,
		"d": [1, "string again", [], {}],
		"e": {
			"f": {
				"g": "321"
			},
			"h": []
		},
		"i": null,
		"j": [true, false],
		"k": ""
	}'`
}

func arraiString() string {
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
