package bootstrap

import (
	"fmt"
	"log"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	t.Parallel()

	expr := Rule("expr")
	g := Grammar{
		expr: Tower{
			Delim{Term: expr, Sep: RE(`([-+])`)},
			Delim{Term: expr, Sep: RE(`([*/])`)},
			Oneof{expr, Seq{Opt(S("-")), expr}},
			Oneof{RE(`(\d+)`), expr},
			Seq{S("("), expr, S(")")},
		},
	}.Compile()

	r := parse.NewScanner("1+2*3")
	var v interface{}
	assert.True(t, g[expr].Parse(r, &v))
	log.Printf("%#v", v)
	assert.Equal(t,
		`expr(expr#1(expr#2(expr#3("1"))), "+", expr#1(expr#2(expr#3("2")), "*", expr#2(expr#3("3"))))`,
		fmt.Sprintf("%q", v),
	)

	r = parse.NewScanner("1+(2-3/4)")
	assert.True(t, g[expr].Parse(r, &v))
	assert.Equal(t,
		`expr(`+
			`expr#1(expr#2(expr#3("1"))), `+
			`"+", `+
			`expr#1(expr#2(expr#3(expr#4("(", `+
			(`expr(expr#1(expr#2(expr#3("2"))), `+
				`"-", `+
				(`expr#1(`+
					`expr#2(expr#3("3")), `+
					`"/", `+
					`expr#2(expr#3("4"))`+
					`)), `))+
			`")")))))`,
		fmt.Sprintf("%q", v),
	)
}

func TestGrammarGrammar(t *testing.T) {
	t.Parallel()

	src := `
		// Simple expression grammar
		expr -> expr:/([-+])/
		      ^ expr:/[\/*]/
		      ^ expr | "-" expr
		      ^ /(\d+)/ | expr
		      ^ "(" expr ")";
	`
	g := GrammarGrammar.Compile()
	r := parse.NewScanner(src)
	var v interface{}
	assert.True(t, g["grammar"].Parse(r, &v), "r=%v\nv=%v", r.Context(), v)
	assert.Equal(t, len(src), r.Offset(), "r=%v\nv=%v", r.Context(), v)
}

func TestGrammarExpr(t *testing.T) {
	t.Parallel()

	g := GrammarGrammar.Compile()
	r := parse.NewScanner(`prod+`)
	var v interface{}
	assert.True(t, g["expr"].Parse(r, &v))
	assert.Equal(t,
		`expr(expr#1(expr#2(expr#3(expr#4(_(atom("prod"), ?(quant("+"))))))))`,
		fmt.Sprintf("%q", v),
	)
}

// Non-terminals
var grammarGrammarSrc = `
	grammar -> prod+;
	stmt    -> comment | prod;
	comment -> /(//.*$)/;
	prod    -> ident "->" expr+ ";";
	expr    -> expr:"^"
			^ expr:"|"
			^ expr+
			^ expr | ("<" expr ">")? expr
			^ atom quant | atom;
	atom    -> ident | str | re | "(" expr ")";
	quant   -> /([?*+])/ | "{" int? "," int? "}" | ":" atom;

	// Terminals
	ident   -> /([A-Za-z_\.]\w*)/;
	str     -> /"([^"\\]|\\.)*"/;
	i       -> /(\d+)/;
	re      -> /\/([^\/\\]|\\.)\//;
	.wrapRE -> /\s*()\s*/;
`

func TestGrammarGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	g := GrammarGrammar.Compile()
	r := parse.NewScanner(grammarGrammarSrc)
	var v interface{}
	assert.True(t, g["grammar"].Parse(r, &v))
}
