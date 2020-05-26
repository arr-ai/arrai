package syntax

import "testing"

func TestJSONDecode(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `//encoding.json.decode(123)`, "")

	expected := `{
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

	encoding := `'{
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

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(`+encoding+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.json.decode(<<`+encoding+`>>)`)
}
