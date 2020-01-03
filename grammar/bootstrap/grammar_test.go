package bootstrap

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertUnparse(t *testing.T, expected string, g Grammar, v interface{}) bool { //nolint:unparam
	var sb strings.Builder
	_, err := g.Unparse(v, &sb)
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

	parsers, grammar := exprGrammar.Compile()

	r := parse.NewScanner("1+2*3")
	v, ok := parsers.Parse(expr, r)
	require.True(t, ok)
	assert.NoError(t, grammar.ValidateParse(v))
	assertUnparse(t, "1+2*3", grammar, v)
	assert.Equal(t,
		`expr\:(expr#1\:(expr#2\_(?(), expr#3\|("1"))), `+
			`"+", `+
			`expr#1\:(expr#2\_(?(), expr#3\|("2")), "*", expr#2\_(?(), expr#3\|("3"))))`,
		fmt.Sprintf("%q", v),
	)

	r = parse.NewScanner("1+(2-3/4)")
	v, ok = parsers.Parse(expr, r)
	assert.True(t, ok)
	assert.NoError(t, grammar.ValidateParse(v))
	assertUnparse(t, "1+(2-3/4)", grammar, v)
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

	parsers, grammar := GrammarGrammar.Compile()
	r := parse.NewScanner(exprGrammarSrc)
	v, ok := parsers.Parse(grammarR, r)
	require.True(t, ok, "r=%v\nv=%v", r.Context(), v)
	require.Equal(t, len(exprGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	assert.NoError(t, grammar.ValidateParse(v))
	log.Printf("%#v", v)
	assertUnparse(t,
		`// Simple expression grammar`+
			`expr->expr:([-+])`+
			`^expr:([*\/])`+
			`^-?expr`+
			`^(\d+)|expr`+
			`^expr<:**`+
			`^(expr);`,
		grammar,
		v,
	)
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

func TestGrammarSnippet(t *testing.T) {
	t.Parallel()

	parsers, grammar := GrammarGrammar.Compile()
	r := parse.NewScanner(`prod+`)
	v, ok := parsers.Parse(term, r)
	require.True(t, ok)
	assert.Equal(t,
		`term\:(term#1\:(term#2\?(term#3\_(?(), term#4\_(atom\|("prod"), ?(quant\|("+")))))))`,
		fmt.Sprintf("%q", v),
	)
	assert.NoError(t, grammar.ValidateParse(v))
	assertUnparse(t, "prod+", grammar, v)
}

func TestTinyGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	tiny := Rule("tiny")
	tinyGrammar := Grammar{tiny: S("x")}
	tinyGrammarSrc := `tiny -> "x";`

	parsers, grammar := GrammarGrammar.Compile()
	r := parse.NewScanner(tinyGrammarSrc)
	v, ok := parsers.Parse(grammarR, r)
	require.True(t, ok)
	log.Print(v)
	assert.Equal(t, len(tinyGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	e := v.(parse.Node)
	assert.NoError(t, grammar.ValidateParse(v))

	grammar2 := CompileGrammarNode(e)
	assert.NotNil(t, grammar2)
	assertGrammarsMatch(t, tinyGrammar, grammar2)
}

func TestExprGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	parsers, grammar := GrammarGrammar.Compile()
	r := parse.NewScanner(exprGrammarSrc)
	v, ok := parsers.Parse(grammarR, r)
	require.True(t, ok)
	assert.Equal(t, len(exprGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	e := v.(parse.Node)
	assert.NoError(t, grammar.ValidateParse(v))

	grammar2 := CompileGrammarNode(e)
	assert.NotNil(t, grammar2)
	assertGrammarsMatch(t, exprGrammar, grammar2)
}

const grammarGrammarSrc = `
grammar -> stmt+;
stmt    -> comment | prod;
comment -> /(\/\/.*$|(?s:\/\*(?:[^*]|\*+[^*\/])\*\/))/;
prod    -> ident "->" term+ ";";
term    -> term:"^"
         ^ term:"|"
         ^ term+
         ^ ("<" ident ">")? term
         ^ atom quant?;
atom    -> ident | str | re | "(" term ")";
quant   -> /([?*+])/
         | "{" int? "," int? "}"
         | /(<:|:>?)/ atom;
ident   -> /([A-Za-z_\.]\w*)/;
str     -> /"((?:[^"\\]|\\.)*)"/;
int     -> /(\d+)/;
re      -> /\/((?:[^\/\\]|\\.)*)\//;
.wrapRE -> /\s*()\s*/;
`

func TestGrammarGrammarGrammarGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	// 0. Have a boostrap Grammar schema defined in Go.

	// 1. Build a grammar for itself in that structure.
	parsers, grammar := GrammarGrammar.Compile()

	for i := 0; i < 2; i++ {
		// 2. Use that grammar grammar to parse the syntax for the grammar.
		r := parse.NewScanner(grammarGrammarSrc)
		v, ok := parsers.Parse(grammarR, r)
		require.True(t, ok)
		g := v.(parse.Node)
		require.Equal(t, len(grammarGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
		assert.NoError(t, grammar.ValidateParse(g))

		// 3. Compile the resulting AST back into a Grammar structure.
		grammar2 := CompileGrammarNode(g)
		assert.NotNil(t, grammar2)

		// 4. Check that that structure matches the previous structure.
		assertGrammarsMatch(t, GrammarGrammar, grammar2)

		// 5. Repeat steps 2 to 4 to confirm that we've reached a fixed point.
		parsers, grammar = grammar2.Compile()
	}
}
