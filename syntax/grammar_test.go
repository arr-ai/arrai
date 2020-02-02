package syntax

import (
	"regexp"
	"testing"
)

var nlIndentRE = regexp.MustCompile(`\n\t*`)

func TestGrammarToValueExprStd(t *testing.T) {
	expected := nlIndentRE.ReplaceAllString(
		`(
			@rule: "grammar",
			stmt: [
				(
					@choice: [1],
					prod: (
						'': ["->", ";"],
						IDENT: ('': "a"),
						term: [(term: [(term: [(term: [(named: (atom: (
								@choice: [1],
								STR: ('': "'1'")
						)))])])])]
					)
				)
			]
		)`, "")
	AssertCodesEvalToSameValue(t, expected, `//.grammar.parse(//.grammar.lang.wbnf, "grammar", "a-> '1';")`)
	AssertCodesEvalToSameValue(t, expected, `//.grammar -> .parse(.lang.wbnf, "grammar", "a-> '1';")`)
}

func TestGrammarParseParse(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`("": ["+"], @rule: "expr", expr: [(expr: [("": "1")]), ("": ["*"], expr: [("": "2"), ("": "3")])])`,
		`//.grammar -> .parse(.parse(.lang.wbnf, "grammar", "expr -> @:'+' > @:'*' > /{\\d+};"), "expr", "1+2*3")`)
}
