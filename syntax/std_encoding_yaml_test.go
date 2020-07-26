package syntax

import "testing"

func TestYAMLDecode(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, "", `//encoding.yaml.decode(':')`)

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
		"i": (a: [(b: true), (b: false)]),
		"j": (s: {}),
		"k": (),
		"l": ()
	}`

	encoding := `
'a: string
b: 123
c: 123.321
d: [1, string again, [], {}]
e:
  f:
    g: "321"
  h: []
i: [true, false]
j: ""
k: null
l: '`

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decode(`+encoding+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decode(<<`+encoding+`>>)`)
}
