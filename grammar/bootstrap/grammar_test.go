package bootstrap

import (
	"fmt"
	"log"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
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

	r := parse.NewRange("1+2*3")
	var v interface{}
	assert.True(t, g[expr].Parse(&r, &v))
	assert.Equal(t, `[expr 1 + [expr 2 * 3]]`, fmt.Sprintf("%v", v))

	r = parse.NewRange("1+(2-3/4)")
	assert.True(t, g[expr].Parse(&r, &v))
	assert.Equal(t, `[expr 1 + [expr ( [expr 2 - [expr 3 / 4]] )]]`, fmt.Sprintf("%v", v))
}

func TestGrammarGrammar(t *testing.T) {
	src := `
		// Simple expression grammar
		expr -> expr:/([-+])/
		      ^ expr:/[\/*]/
		      ^ expr | "-" expr
		      ^ /(\d+)/ | expr
		      ^ "(" expr ")";
	`
	g := GrammarGrammar.Compile()
	r := parse.NewRange(src)
	var v interface{}
	assert.True(t, g["grammar"].Parse(&r, &v), "r=%v\nv=%v", r.Context(), v)
	assert.Equal(t, len(src), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	assert.Equal(t,
		`[grammar `+
			`[stmt [comment // Simple expression grammar]] `+
			`[stmt [prod [ident expr] -> [expr [expr [atom [ident expr]] [quant [: [atom [re ([-+])]]]]] `+
			`^ [expr [atom [ident expr]] [quant [: [atom [re [\/*]]]]]] `+
			`^ [expr [expr [atom [ident expr]] []] | [expr [expr [atom [str -]] []] [expr [atom [ident expr]] []]]] `+
			`^ [expr [expr [atom [re (\d+)]] []] | [expr [atom [ident expr]] []]] `+
			`^ [expr [expr [atom [str (]] []] [expr [atom [ident expr]] []] [expr [atom [str )]] []]]] ;]]]`,
		fmt.Sprintf("%v", v))

	log.Print(v)
}

func TestGrammarGrammarGrammar(t *testing.T) {
	grammarGrammarSrc := `
		// Non-terminals
		grammar -> prod+;
		stmt    -> comment | prod;
		comment -> /(//.*$)/;
		prod    -> ident "->" expr+ ";";
		expr    -> expr:"^";
		         ^ expr:"|";
		         ^ expr+;
		         ^ expr | ("<" expr ">")? expr;
		         ^ atom quant?;
		atom    -> ident | str | re | "(" expr ")";
		quant   -> /([?*+])/ | "{" int? "," int? "}" | ":" atom;

		// Terminals
		ident   -> /([A-Za-z_\.]\w*)/;
		str     -> /"([^"\\]|\\.)*"/;
		i       -> /(\d+)/;
		re      -> /\/([^\/\\]|\\.)\//;
		.wrapRE -> /\s*()\s* /
	`

	g := GrammarGrammar.Compile()
	r := parse.NewRange(grammarGrammarSrc)
	var v interface{}
	assert.True(t, g["grammar"].Parse(&r, &v))
	log.Print(v)
}
