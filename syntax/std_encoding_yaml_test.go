package syntax

import "testing"

func TestYAMLDecode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "()", `//encoding.yaml.decode('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.yaml.decode('{}')`)
	AssertCodesEvalToSameValue(t, "(a: [])", `//encoding.yaml.decode('[]')`)
	AssertCodeErrors(t, "", `//encoding.yaml.decode(123)`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.yaml.decode('123')`)

	expected := testArraiString()
	yaml := testYAMLString()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decode(`+yaml+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decode(<<`+yaml+`>>)`)
}

func TestYAMLDecode_NonStrict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "()", `//encoding.yaml.decoder((strict: false))('null')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.yaml.decoder((strict: false))('{}')`)
	AssertCodesEvalToSameValue(t, "{}", `//encoding.yaml.decoder((strict: false))('[]')`)
	AssertCodeErrors(t, "", `//encoding.yaml.decoder((strict: false))(123)`)
	AssertCodesEvalToSameValue(t, "123", `//encoding.yaml.decoder((strict: false))('123')`)

	yaml := testYAMLString()
	expected := testArraiStringLoose()

	// String
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decoder((strict: false))(`+yaml+`)`)
	// Bytes
	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.decoder((strict: false))(<<`+yaml+`>>)`)
}

func TestYAMLEncode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `<<'null\n'>>`, `//encoding.yaml.encode(())`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.yaml.encode({})`)
	AssertCodesEvalToSameValue(t, `<<'{}\n'>>`, `//encoding.yaml.encode({123})`)
	AssertCodesEvalToSameValue(t, `<<'[]\n'>>`, `//encoding.yaml.encode((a: []))`)
	AssertCodesEvalToSameValue(t, `<<'123\n'>>`, `//encoding.yaml.encode(123)`)
	AssertCodesEvalToSameValue(t, `<<'- 123\n'>>`, `//encoding.yaml.encode([123])`)
	AssertCodesEvalToSameValue(t, `<<'key: 123\n'>>`, `//encoding.yaml.encode({"key": 123})`)
	AssertCodesEvalToSameValue(t, `<<'abcde\n'>>`, `//encoding.yaml.encode("abcde")`)
	AssertCodesEvalToSameValue(t, `<<"'>'\n">>`, `//encoding.yaml.encode(">")`)

	AssertCodeErrors(t, "", `//encoding.yaml.encode((a: [], s: ""))`)

	encoding := testArraiString()
	expected := `<<` + testYAMLString() + `>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.yaml.encoder((indent: 2))(`+encoding+`)`)
}

func TestYAMLEncodeIndent(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `<<'a:\n    b: 1\n'>>`, `//encoding.yaml.encode({"a": {"b": 1}})`) // 4 spaces by default
	AssertCodesEvalToSameValue(t, `<<'a:\n  b: 1\n'>>`, `//encoding.yaml.encoder((indent:2))({"a": {"b": 1}})`)
	AssertCodesEvalToSameValue(t, `<<'a:\n    b: 1\n'>>`, `//encoding.yaml.encoder((indent:4))({"a": {"b": 1}})`)
}

func testYAMLString() string {
	return `
'a: string
b: 123
c: 123.321
d:
  - 1
  - string again
  - []
  - {}
e:
  f:
    g: "321"
  h: []
i: null
j:
  - true
  - false
k: ""
'`
}
