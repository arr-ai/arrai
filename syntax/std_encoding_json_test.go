package syntax

import "testing"

func TestJSONDecode(t *testing.T) {
	t.Parallel()
	AssertCodePanics(t, `//.encoding.json.decode(123)`)
	AssertCodesEvalToSameValue(t,
		`{
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
			"i": (null: {}),
			"j": (a: [(b: {()}), (b: {})]),
			"k": (s: {})
		}`,
		`//.encoding.json.decode(
			'{
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
			}'
		)`,
	)
}
