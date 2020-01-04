package bootstrap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertUnparse(t *testing.T, expected string, parsers Parsers, v interface{}) bool { //nolint:unparam
	var sb strings.Builder
	_, err := parsers.Unparse(v, &sb)
	return assert.NoError(t, err) && assert.Equal(t, expected, sb.String())
}

var expr = Rule("expr")

var exprGrammarSrc = `
// Simple expression grammar
expr -> expr:/([-+])/
      ^ expr:/([*\/])/
      ^ "-"? expr
	  ^ /(\d+)/ | expr
	  ^ expr<:"**"
      ^ "(" expr ")";
`

var exprGrammar = Grammar{
	expr: Tower{
		Delim{Term: expr, Sep: RE(`([-+])`)},
		Delim{Term: expr, Sep: RE(`([*/])`)},
		Seq{Opt(S("-")), expr},
		Oneof{RE(`(\d+)`), expr},
		R2L(expr, S("**")),
		Seq{S("("), expr, S(")")},
	},
}

func TestGrammarParser(t *testing.T) {
	t.Parallel()

	parsers := exprGrammar.Compile()

	r := parse.NewScanner("1+2*3")
	v, err := parsers.Parse(expr, r)
	require.NoError(t, err)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "1+2*3", parsers, v)
	assert.Equal(t,
		`expr\:(expr#1\:(expr#2\_(?(), expr#3\|("1"))), `+
			`"+", `+
			`expr#1\:(expr#2\_(?(), expr#3\|("2")), "*", expr#2\_(?(), expr#3\|("3"))))`,
		fmt.Sprintf("%q", v),
	)

	r = parse.NewScanner("1+(2-3/4)")
	v, err = parsers.Parse(expr, r)
	assert.NoError(t, err)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "1+(2-3/4)", parsers, v)
	assert.Equal(t,
		`expr\:(`+
			`expr#1\:(expr#2\_(?(), expr#3\|("1"))), `+
			`"+", `+
			`expr#1\:(expr#2\_(?(), expr#3\|(expr#4\:(expr#5\_("(", `+
			`expr\:(expr#1\:(expr#2\_(?(), expr#3\|("2"))), `+
			`"-", `+
			`expr#1\:(expr#2\_(?(), expr#3\|("3")), `+
			`"/", `+
			`expr#2\_(?(), expr#3\|("4")))), `+
			`")"))))))`,
		fmt.Sprintf("%q", v),
	)
}

func TestExprGrammarGrammar(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(exprGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err, "r=%v\nv=%v", r.Context(), v)
	require.Equal(t, len(exprGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t,
		`// Simple expression grammar`+
			`expr->expr:([-+])`+
			`^expr:([*\/])`+
			`^-?expr`+
			`^(\d+)|expr`+
			`^expr<:**`+
			`^(expr);`,
		parsers,
		v,
	)
}

func TestGrammarSnippet(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(`prod+`)
	v, err := parsers.Parse(term, r)
	require.NoError(t, err)
	assert.Equal(t,
		`term\:(term#1\:(term#2\?(term#3\_(?(), term#4\_(atom\|("prod"), ?(quant\|("+")))))))`,
		fmt.Sprintf("%q", v),
	)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "prod+", parsers, v)
}

func assertGrammarsMatch(t *testing.T, expected, actual Grammar) {
	if !assert.True(t, reflect.DeepEqual(expected, actual)) {
		t.Logf("raw expected: %#v", expected)
		t.Logf("raw actual: %#v", actual)

		expectedJSON, err := json.Marshal(expected)
		require.NoError(t, err)
		actualJSON, err := json.Marshal(actual)
		require.NoError(t, err)
		t.Log("JSON(expected): ", string(expectedJSON))
		t.Log("JSON(actual): ", string(actualJSON))
		assert.JSONEq(t, string(expectedJSON), string(actualJSON))
	}
}

func TestTinyGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	tiny := Rule("tiny")
	tinyGrammar := Grammar{tiny: S("x")}
	tinyGrammarSrc := `tiny -> "x";`

	parsers := Core()
	r := parse.NewScanner(tinyGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err)
	e := v.(parse.Node)
	assert.NoError(t, parsers.ValidateParse(v))

	grammar2 := NewFromNode(e)
	assertGrammarsMatch(t, tinyGrammar, grammar2)
}

func TestExprGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(exprGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err)
	e := v.(parse.Node)
	assert.NoError(t, parsers.ValidateParse(v))

	grammar2 := NewFromNode(e)
	assertGrammarsMatch(t, exprGrammar, grammar2)
}
