package tests

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertCodesEvalToSameValue asserts that code evaluate to the same value as
// expected.
func AssertCodesEvalToSameValue(
	t *testing.T, expected string, code string,
) bool {
	expectedExpr, err := syntax.Parse(syntax.NewStringLexer(expected))
	if !assert.NoError(t, err, "parsing expected: %s", expected) {
		return false
	}
	codeExpr, err := syntax.Parse(syntax.NewStringLexer(code))
	if !assert.NoError(t, err, "parsing code: %s", code) {
		return false
	}
	if !rel.AssertExprsEvalToSameValue(t, expectedExpr, codeExpr) {
		return assert.Fail(
			t, "Codes should eval to same value", "%s == %s", expected, code)
	}
	return true
}

// RequireCodesEvalToSameValue requires that code evaluate to the same value as
// expected.
func RequireCodesEvalToSameValue(t *testing.T, expected string, code string) {
	expectedExpr, err := syntax.Parse(syntax.NewStringLexer(expected))
	require.NoError(t, err)
	codeExpr, err := syntax.Parse(syntax.NewStringLexer(code))
	require.NoError(t, err)
	rel.AssertExprsEvalToSameValue(t, expectedExpr, codeExpr)
}

// AssertScan asserts that a lexer's next produced token is as expected.
func AssertScan(
	t *testing.T, l *syntax.Lexer, tok syntax.Token, intf interface{}, lexeme string,
) bool {
	if !assert.True(t, l.Scan()) {
		return false
	}

	if !assert.Equal(
		t, syntax.TokenRepr(tok), syntax.TokenRepr(l.Token()), "%s", l,
	) {
		return false
	}

	if intf == nil {
		if !assert.Nil(t, l.Value()) {
			return false
		}
	} else {
		value, err := rel.NewValue(intf)
		require.NoError(t, err)
		if !assert.True(
			t, value.Equal(l.Value()), "%s == %s", value, l.Value(),
		) {
			return false
		}
	}

	return assert.Equal(t, lexeme, lexeme, l)
}
