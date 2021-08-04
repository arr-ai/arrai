package translate_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/translate"
)

func TestJSONTranslator_NonStrict(t *testing.T) {
	t.Parallel()

	approx := translate.NewTranslator(false)
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.None)
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.MustNewSet())
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.NewArray())
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.MustNewDict(false))
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.NewString([]rune("")))
	AssertArraiValueToJSONTranslationWith(t, approx, `null`, rel.NewBool(false))

	AssertArraiToJSONTranslationWith(t, approx, `42`, `42`)
	AssertArraiToJSONTranslationWith(t, approx, `true`, `true`)
	AssertArraiToJSONTranslationWith(t, approx, `"foo"`, `"foo"`)
	AssertArraiToJSONTranslationWith(t, approx, `["foo","bar"]`, `["foo", "bar"]`)
	AssertArraiToJSONTranslationWith(t, approx, `["bar","foo"]`, `{"foo", "bar"}`)
	AssertArraiToJSONTranslationWith(t, approx, `{"foo":"bar"}`, `{"foo": "bar"}`)
}

func TestJSONObjectToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertJSONToArraiTranslation(t, `{}`, `{}`)

	// different value types
	AssertJSONToArraiTranslation(t, `{"key": 123}                `, `{"key":123}          `)
	AssertJSONToArraiTranslation(t, `{"key": ()}                 `, `{"key":null}         `)
	AssertJSONToArraiTranslation(t, `{"key": (s: "val")}         `, `{"key":"val"}        `)
	AssertJSONToArraiTranslation(t, `{"key": (a: [1, 2, 3])}     `, `{"key":[1, 2, 3]}    `)
	AssertJSONToArraiTranslation(t, `{"key": {"foo": (s: "bar")}}`, `{"key":{"foo":"bar"}}`)

	// Multiple key-val pairs
	AssertJSONToArraiTranslation(t, `{"key": (s: "val"), "foo": 123}`, `{"key":"val", "foo":123}`)
}

func TestJSONArrayToArrai(t *testing.T) {
	t.Parallel()

	// Empty
	AssertJSONToArraiTranslation(t, `(a: [])`, `[]`)

	// Different value types
	AssertJSONToArraiTranslation(t, `(a: [1])                  `, `[1]            `)
	AssertJSONToArraiTranslation(t, `(a: [()])                 `, `[null]         `)
	AssertJSONToArraiTranslation(t, `(a: [(s: "hello")])       `, `["hello"]      `)
	AssertJSONToArraiTranslation(t, `(a: [(a: [1, 2, 3])])     `, `[[1, 2, 3]]    `)
	AssertJSONToArraiTranslation(t, `(a: [{"foo": (s: "bar")}])`, `[{"foo":"bar"}]`)

	// Multiple values with different types
	AssertJSONToArraiTranslation(t, `(a: [1, (s: "Hello"), ()])`, `[1, "Hello", null]`)
}

func TestJSONNullToNone(t *testing.T) {
	t.Parallel()
	AssertJSONToArraiTranslation(t, `()`, `null`)
}

func TestJSONStringToArrai(t *testing.T) {
	t.Parallel()
	AssertJSONToArraiTranslation(t, `(s: {})           `, `""           `)
	AssertJSONToArraiTranslation(t, `(s: "Hello World")`, `"Hello World"`)
}

func TestJSONNumericToArrai(t *testing.T) {
	t.Parallel()
	AssertJSONToArraiTranslation(t, `123 `, `123 `)
	AssertJSONToArraiTranslation(t, `1.23`, `1.23`)
}

// AssertExpectedTranslation asserts that the translated value is the same as the expected string
func AssertExpectedTranslation(t *testing.T, expected string, value rel.Value) bool {
	if !rel.AssertExprsEvalToSameValue(t, compileExpr(t, expected), value) {
		return assert.Fail(
			t, "Input should translate to same value", "%s == %s", expected, value)
	}
	return true
}

func AssertJSONToArraiTranslation(t *testing.T, expected, rawJSON string) {
	AssertJSONToArraiTranslationWith(t, translate.StrictTranslator(), expected, rawJSON)
}

func AssertArraiToJSONTranslationWith(t *testing.T, translator translate.Translator, expectedJSON, arrai string) {
	AssertArraiValueToJSONTranslationWith(t, translator, expectedJSON, eval(t, arrai))
}

func AssertJSONToArraiTranslationWith(t *testing.T, translator translate.Translator, expected, rawJSON string) {
	var data interface{}
	require.NoError(t, json.Unmarshal([]byte(rawJSON), &data))
	trans, err := translator.ToArrai(data)
	require.NoError(t, err)
	AssertExpectedTranslation(t, expected, trans)
}

func AssertArraiValueToJSONTranslationWith(
	t *testing.T,
	translator translate.Translator,
	expectedJSON string,
	value rel.Value,
) {
	v, err := translator.FromArrai(value)
	require.NoError(t, err)
	jb, err := json.Marshal(v)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, string(jb))
}

func compileExpr(t *testing.T, src string) rel.Expr {
	var pc syntax.ParseContext
	ctx := arraictx.InitRunCtx(context.Background())
	ast, err := pc.Parse(ctx, parser.NewScanner(src))
	require.NoError(t, err)
	expr, err := pc.CompileExpr(ctx, ast)
	require.NoError(t, err)
	return expr
}

func eval(t *testing.T, src string) rel.Value {
	ctx := arraictx.InitRunCtx(context.Background())
	expr := compileExpr(t, src)
	value, err := expr.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	return value
}
