package bootstrap

import (
	"fmt"
	"log"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	g := Grammar{
		"expr":  Rule("add"),
		"add":   Delim{Term: Rule("mul"), Sep: RE(`([-+])`)},
		"mul":   Delim{Term: Rule("neg"), Sep: RE(`([*/])`)},
		"neg":   Seq{Opt(S("-")), Rule("atom")},
		"atom":  Choice{RE(`(\d+)`), Rule("paren")},
		"paren": Seq{S("("), Rule("expr"), S(")")},
	}.Parsers()

	r := parse.NewRange("42+54")
	v, ok := g("expr").Parse(&r)
	assert.True(t, ok)
	assert.Equal(t,
		`[add [mul [neg ["-"{0,1}] [atom [/(\d+)/ 42]]]] [/([-+])/ +] [mul [neg ["-"{0,1}] [atom [/(\d+)/ 54]]]]]`,
		fmt.Sprintf("%v", v),
	)

	r = parse.NewRange("1+(2-3/4)")
	v, ok = g("expr").Parse(&r)
	assert.True(t, ok)
	assert.Equal(t,
		`[add [mul [neg ["-"{0,1}] [atom [/(\d+)/ 1]]]] `+
			`[/([-+])/ +] `+
			`[mul [neg ["-"{0,1}] [atom [paren ["(" (] `+
			`[add [mul [neg ["-"{0,1}] [atom [/(\d+)/ 2]]]] `+
			`[/([-+])/ -] `+
			`[mul [neg ["-"{0,1}] [atom [/(\d+)/ 3]]] `+
			`[/([*/])/ /] `+
			`[neg ["-"{0,1}] [atom [/(\d+)/ 4]]]]] [")" )]]]]]]`,
		fmt.Sprintf("%v", v),
	)
}

func TestGrammarGrammar(t *testing.T) {
	src := `
		// Simple expression grammar
		expr  -> add;
		add   -> mul:/([-+])/;
		mul   -> neg:/[/*]/;
		neg   -> "-"? atom;
		atom  -> /(\d+)/ | paren;
		paren -> "(" expr ")";
	`
	gg := GrammarGrammar.Parsers()
	r := parse.NewRange(src)
	v, ok := gg("grammar").Parse(&r)
	assert.True(t, ok, "%v", r)
	log.Print(v)
}

func TestGrammarGrammarGrammar(t *testing.T) {
	grammarGrammarSrc := `
		// Non-terminals
		grammar -> prod+;
		stmt    -> comment | prod;
		comment -> /(//.*$)/;
		prod    -> ident "->" expr+ ";";
		expr    -> choice;
		choice  -> seq:"|";
		seq     -> tag+;
		tag     -> ("<" ident ">")? term;
		term    -> atom quant?;
		atom    -> ident | str | re | "(" expr ")";
		quant   -> /([?*+])/ | "{" int? "," int? "}" | ":" atom;

		// Terminals
		ident   -> /([A-Za-z_\.]\w*)/;
		str     -> /"([^"\\]|\\.)*"/;
		i       -> /(\d+)/;
		re      -> /\/([^\/\\]|\\.)\//;
		.wrapRE -> /\s*()\s* /
	`

	gg := GrammarGrammar.Parsers()
	r := parse.NewRange(grammarGrammarSrc)
	v, ok := gg("grammar").Parse(&r)
	assert.True(t, ok)
	log.Print(v)
}
