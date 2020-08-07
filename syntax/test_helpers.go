package syntax

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertCodesEvalToSameValue asserts that code evaluate to the same value as
// expected.
func AssertCodesEvalToSameValue(t *testing.T, expected, code string) bool {
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(expected))
	if !assert.NoError(t, err, "parsing expected: %s", expected) {
		return false
	}
	expectedExpr := pc.CompileExpr(ast)
	ast, err = pc.Parse(parser.NewScanner(code))
	if !assert.NoError(t, err, "parsing code: %s", code) {
		return false
	}
	codeExpr := pc.CompileExpr(ast)
	value, err := codeExpr.Eval(rel.Scope{})
	if !assert.NoError(t, err, "evaluating expected: %s", expected) {
		return false
	}
	// log.Printf("code=%v, codeExpr=%v", code, codeExpr)
	if !rel.AssertExprsEvalToSameValue(t, expectedExpr, value) {
		t.Errorf("\nexpected: %s\ncode:     %s", expected, code)
		return false
	}
	return true
}

// RequireCodesEvalToSameValue requires that code evaluates to the same value as
// expected.
func RequireCodesEvalToSameValue(t *testing.T, expected string, code string) {
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(expected))
	require.NoError(t, err)
	expectedExpr := pc.CompileExpr(ast)
	ast, err = pc.Parse(parser.NewScanner(code))
	require.NoError(t, err)
	codeExpr := pc.CompileExpr(ast)
	rel.AssertExprsEvalToSameValue(t, expectedExpr, codeExpr)
}

// AssertCodeEvalsToType asserts that code evaluates to the same type as expected.
func AssertCodeEvalsToType(t *testing.T, expected interface{}, code string) bool {
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(code))
	if !assert.NoError(t, err, "parsing code: %s", code) {
		return false
	}
	codeExpr := pc.CompileExpr(ast)
	if !rel.AssertExprEvalsToType(t, expected, codeExpr) {
		t.Errorf("\nexpected: %T\ncode:     %s", expected, code)
		return false
	}
	return true
}

// AssertCodeEvalsToGrammar asserts that code evaluates to a grammar equal to expected.
func AssertCodeEvalsToGrammar(t *testing.T, expected parser.Grammar, code string) {
	pc := ParseContext{SourceDir: ".."}
	astElt := pc.MustParseString(code)
	astExpr := pc.CompileExpr(astElt)
	astValue, err := astExpr.Eval(rel.EmptyScope)
	assert.NoError(t, err, "parsing code: %s", code)
	astNode := rel.ASTNodeFromValue(astValue).(ast.Branch)
	astGrammar := wbnf.NewFromAst(astNode)

	assert.EqualValues(t, expected, astGrammar)
}

// AssertCodePanics asserts that code panics when executed.
// TODO: Remove this. Should only intentionally panic for implementation bugs.
func AssertCodePanics(t *testing.T, code string) bool {
	return assert.Panics(t, func() {
		pc := ParseContext{SourceDir: ".."}
		ast, err := pc.Parse(parser.NewScanner(code))
		if assert.NoError(t, err, "parsing code: %s", code) {
			codeExpr := pc.CompileExpr(ast)
			codeExpr.Eval(rel.EmptyScope) //nolint:errcheck
		}
	})
}

// AssertCodeErrors asserts that code fails with a certain
// message when executed.
func AssertCodeErrors(t *testing.T, errString, code string) bool {
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(code))
	if assert.NoError(t, err, "parsing code: %s", code) {
		codeExpr := pc.CompileExpr(ast)
		_, err := codeExpr.Eval(rel.EmptyScope)
		if err == nil {
			panic(fmt.Sprintf("the code `%s` didn't generate any error", code))
		}
		assert.EqualError(t, errors.New(err.Error()[:len(errString)]), errString)
	}
	return false
}

// AssertScan asserts that a lexer's next produced token is as expected.
func AssertScan(t *testing.T, l *Lexer, tok Token, intf interface{}, lexeme string) bool {
	if !assert.True(t, l.Scan()) {
		return false
	}

	if !assert.Equal(
		t, TokenRepr(tok), TokenRepr(l.Token()), "%s", l,
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

// AssertEvalExprString asserts Expr string.
func AssertEvalExprString(t *testing.T, expected, source string) bool {
	expr, err := Compile(".", source)
	return assert.NoError(t, err) &&
		assert.NotNil(t, expr) &&
		assert.Equal(t, expected, strings.Replace(expr.String(), ` `, ``, -1))
}
