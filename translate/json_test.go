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
	trans, err := translate.JSONToArrai(data)
	require.NoError(t, err)
	AssertExpectedTranslation(t, expected, trans)
}

func TestJSONObjectToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedJSONTranslation(t, `{}`, `{}`)

	// different value types
	AssertExpectedJSONTranslation(t, `{"key": "val"}`, `{"key":"val"}`)
	AssertExpectedJSONTranslation(t, `{"key": 123}`, `{"key":123}`)
	AssertExpectedJSONTranslation(t, `{"key": {"foo": "bar"}}`, `{"key":{"foo":"bar"}}`)
	AssertExpectedJSONTranslation(t, `{"key": [1, 2, 3]}`, `{"key":[1, 2, 3]}`)
	AssertExpectedJSONTranslation(t, `{"key": {}}`, `{"key":null}`)

	// Multiple key-val pairs
	AssertExpectedJSONTranslation(t, `{"key": "val", "foo": 123}`, `{"key":"val", "foo":123}`)
}

func TestJSONArrayToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedJSONTranslation(t, `[]`, `[]`)

	// Different value types
	AssertExpectedJSONTranslation(t, `[1]`, `[1]`)
	AssertExpectedJSONTranslation(t, `["hello"]`, `["hello"]`)
	AssertExpectedJSONTranslation(t, `[{"foo": "bar"}]`, `[{"foo":"bar"}]`)
	AssertExpectedJSONTranslation(t, `[[1, 2, 3]]`, `[[1, 2, 3]]`)
	AssertExpectedJSONTranslation(t, `[{}]`, `[null]`)

	// Multiple values with different types
	AssertExpectedJSONTranslation(t, `[1, "Hello", {}]`, `[1, "Hello", null]`)
}

func TestJSONNullToNone(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `{}`, `null`)
}

func TestJSONStringToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `""`, `""`)
	AssertExpectedJSONTranslation(t, `"Hello World"`, `"Hello World"`)
}

func TestJSONNumericToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedJSONTranslation(t, `123`, `123`)
	AssertExpectedJSONTranslation(t, `1.23`, `1.23`)
}
