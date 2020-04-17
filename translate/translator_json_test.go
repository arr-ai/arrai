package translate_test

import (
	"encoding/json"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/translate"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertExpectedTranslation asserts that the translated value is the same as the expected string
func AssertExpectedTranslation(t *testing.T, expected string, value rel.Value) bool {
	var pc syntax.ParseContext
	ast, err := pc.Parse(parser.NewScanner(expected))
	if !assert.NoError(t, err, "parsing expected: %s", expected) {
		return false
	}
	expectedExpr := pc.CompileExpr(ast)
	if !rel.AssertExprsEvalToSameValue(t, expectedExpr, value) {
		return assert.Fail(
			t, "Input should translate to same value", "%s == %s", expected, value)
	}
	return true
}

func AssertExpectedJSONTranslation(t *testing.T, expected, rawJSON string) {
	var data interface{}
	require.NoError(t, json.Unmarshal([]byte(rawJSON), &data))
	trans, err := translate.ToArrai(data)
	require.NoError(t, err)
	AssertExpectedTranslation(t, expected, trans)
}

func TestJSONObjectToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedJSONTranslation(t, `{}`, `{}`)

	// different value types
	AssertExpectedJSONTranslation(t, `{"key": 123}           `, `{"key":123}          `)
	AssertExpectedJSONTranslation(t, `{"key": (null: {})}    `, `{"key":null}         `)
	AssertExpectedJSONTranslation(t, `{"key": (s: "val")}    `, `{"key":"val"}        `)
	AssertExpectedJSONTranslation(t, `{"key": (a: [1, 2, 3])}`, `{"key":[1, 2, 3]}    `)
	AssertExpectedJSONTranslation(t, `{"key": {"foo": (s: "bar")}}`, `{"key":{"foo":"bar"}}`)

	// Multiple key-val pairs
	AssertExpectedJSONTranslation(t, `{"key": (s: "val"), "foo": 123}`, `{"key":"val", "foo":123}`)
}

func TestJSONArrayToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedJSONTranslation(t, `(a: [])`, `[]`)

	// Different value types
	AssertExpectedJSONTranslation(t, `(a: [1])                  `, `[1]            `)
	AssertExpectedJSONTranslation(t, `(a: [(null: {})])         `, `[null]         `)
	AssertExpectedJSONTranslation(t, `(a: [(s: "hello")])       `, `["hello"]      `)
	AssertExpectedJSONTranslation(t, `(a: [(a: [1, 2, 3])])     `, `[[1, 2, 3]]    `)
	AssertExpectedJSONTranslation(t, `(a: [{"foo": (s: "bar")}])`, `[{"foo":"bar"}]`)

	// Multiple values with different types
	AssertExpectedJSONTranslation(t, `(a: [1, (s: "Hello"), (null: {})])`, `[1, "Hello", null]`)
}

func TestJSONNullToNone(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `(null: {})`, `null`)
}

func TestJSONStringToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `(s: {})           `, `""           `)
	AssertExpectedJSONTranslation(t, `(s: "Hello World")`, `"Hello World"`)
}

func TestJSONNumericToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `123 `, `123 `)
	AssertExpectedJSONTranslation(t, `1.23`, `1.23`)
}
