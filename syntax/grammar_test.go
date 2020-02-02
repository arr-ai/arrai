package syntax

import (
	"regexp"
	"testing"
)

var nlIndentRE = regexp.MustCompile(`\n\t*`)

func TestGrammarToValueExprStd(t *testing.T) {
	AssertCodesEvalToSameValue(t, nlIndentRE.ReplaceAllString(
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
		)`, ""),
		`//.grammar.parse(//.grammar.lang.wbnf, "grammar", "a-> '1';")`,
	)
}
