package translate_test

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/arr-ai/arrai/translate"
	"github.com/stretchr/testify/require"
)

func AssertExpectedYAMLTranslation(t *testing.T, expected, rawYAML string) {
	var data interface{}
	require.NoError(t, yaml.Unmarshal([]byte(rawYAML), &data))
	trans, err := translate.ToArrai(data)
	require.NoError(t, err)
	AssertExpectedTranslation(t, expected, trans)
}

func TestYAMLObjectToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedYAMLTranslation(t, `{}`, `--- {}`)

	// different value types
	AssertExpectedYAMLTranslation(t, `{"key": (123)}              `, `key: 123       `)
	AssertExpectedYAMLTranslation(t, `{"key": (null: {})}         `, `key:           `)
	AssertExpectedYAMLTranslation(t, `{"key": (s: "val")}         `, `key: val       `)
	AssertExpectedYAMLTranslation(t, `{"key": (a: [1, 2, 3])}     `, `key: [1,2,3]   `)
	AssertExpectedYAMLTranslation(t, `{"key": {"foo": (s: "bar")}}`, `key: {foo: bar}`)

	// Multiple key-val pairs
	AssertExpectedYAMLTranslation(t, `{"key": (s: "val"), "foo": 123}  `, `{"key":"val", "foo":123}`)
}

func TestYAMLArrayToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertExpectedYAMLTranslation(t, `(a: [])`, `[]`)

	// Different value types
	AssertExpectedYAMLTranslation(t, `(a: [1])                            `, `[1]            `)
	AssertExpectedYAMLTranslation(t, `(a: [(null: {})])                   `, `[null]         `)
	AssertExpectedYAMLTranslation(t, `(a: [(s: "hello")])                 `, `["hello"]      `)
	AssertExpectedYAMLTranslation(t, `(a: [(a: [1, 2, 3])])               `, `[[1, 2, 3]]    `)
	AssertExpectedYAMLTranslation(t, `(a: [{"foo": (s: "bar")}])          `, `[{"foo":"bar"}]`)

	// Multiple values with different types
	AssertExpectedYAMLTranslation(t, `(a: [1, (s: 'Hello'), (null: {})])`, `[1, "Hello", null]`)
}

func TestYAMLNullToNone(t *testing.T) {
	t.Parallel()
	AssertExpectedYAMLTranslation(t, `(null: {})`, `null`)
}

func TestYAMLStringToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedYAMLTranslation(t, `(s: {})           `, `""           `)
	AssertExpectedYAMLTranslation(t, `(s: "Hello World")`, `"Hello World"`)
}

func TestYAMLNumericToArrai(t *testing.T) {
	t.Parallel()
	AssertExpectedYAMLTranslation(t, `123 `, `123 `)
	AssertExpectedYAMLTranslation(t, `1.23`, `1.23`)
}
