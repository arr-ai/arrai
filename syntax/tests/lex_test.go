package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/rel/syntax"
)

func lxr(input string) *syntax.Lexer {
	return syntax.NewLexer([]byte(input))
}

func assertScan(
	t *testing.T, l *syntax.Lexer, tok syntax.Token, intf interface{},
	lexeme string,
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

	return assert.Equal(t, lexeme, string(lexeme), l)
}

// TestLexSymbols tests Lexer recognising symbols.
func TestLexSymbols(t *testing.T) {
	assertScan(t, lxr(" \t{aa"), syntax.Token('{'), nil, "{")
	assertScan(t, lxr("\t }|}\n "), syntax.Token('}'), nil, "}")
	assertScan(t, lxr("\t +12\n "), syntax.Token('+'), nil, "+")
	assertScan(t, lxr("\t -1\n "), syntax.Token('-'), nil, "-")
	assertScan(t, lxr("\t ,2\n "), syntax.Token(','), nil, ",")
	assertScan(t, lxr("\t :{|0\n "), syntax.Token(':'), nil, ":")
	assertScan(t, lxr("   {||}"), syntax.OSET, nil, "{|")
	assertScan(t, lxr("   |}}"), syntax.CSET, nil, "|}")
	assertScan(t, lxr("   "), syntax.EOF, nil, "")
}

// TestLexIdent tests Lexer recognising identifiers.
func TestLexIdent(t *testing.T) {
	assertScan(t, lxr(" @,"), syntax.IDENT, nil, "@")
	assertScan(t, lxr(" a,{12"), syntax.IDENT, nil, "a")
	assertScan(t, lxr(" Ab|}"), syntax.IDENT, nil, "Ab")
	assertScan(t, lxr(" \na@b 1"), syntax.IDENT, nil, "a@b")
	assertScan(t, lxr(" \n\t a@b_123__"), syntax.IDENT, nil, "a@b_123__")
}

// TestLexNumber tests Lexer recognising numbers.
func TestLexNumber(t *testing.T) {
	assertScan(t, lxr(" 0}"), syntax.NUMBER, 0, "0")
	assertScan(t, lxr(" 123,"), syntax.NUMBER, 123, "123")
	assertScan(t, lxr(" 0.32 |}"), syntax.NUMBER, 0.32, "0.32")
	assertScan(t, lxr(" 4.5e+123}"), syntax.NUMBER, 4.5e+123, "4.5e+123")
}

// TestLexString tests Lexer recognising strings.
func TestLexString(t *testing.T) {
	assertScan(t, lxr(" \"\"}"), syntax.STRING, nil, "\"\"")
	assertScan(t, lxr(" \"abc\","), syntax.STRING, nil, "\"abc\"")
	assertScan(t, lxr(" \"\\t\\n\"|}"), syntax.STRING, nil, "\"\\t\\n\"")
	assertScan(t, lxr(" \"\\\"\":"), syntax.STRING, nil, "\"\\\"\"")
}

// TestLexSequence tests Lexer recognising a sequence of tokens.
func TestLexSequence(t *testing.T) {
	l := lxr("1 2 3")
	if assertScan(t, l, syntax.NUMBER, 1, "1") {
		if assertScan(t, l, syntax.NUMBER, 2, "2") {
			if assertScan(t, l, syntax.NUMBER, 3, "3") {
				assertScan(t, l, syntax.EOF, nil, "")
			}
		}
	}
}

// // TestLexBadInput tests Lexer detecting bad input.
// func TestLexBadInput(t *testing.T) {
// 	l := lxr("{|<")
// 	assert.Equal(t, syntax.OSET, l.Scan())
// 	assert.Equal(t, syntax.ERROR, l.Scan())
// 	assert.Equal(t, syntax.ERROR, l.Scan())
// }
