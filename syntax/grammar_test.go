package syntax

import (
	"fmt"
	"testing"
)

func TestGrammarToValueExpr(t *testing.T) {
	expected := `(
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
	)`
	AssertCodesEvalToSameValue(t, expected, `//grammar.parse(//grammar.lang.wbnf, "grammar", "a -> '1';")`)
	AssertCodesEvalToSameValue(t, expected, `//grammar -> .parse(.lang.wbnf, "grammar", "a -> '1';")`)
	AssertCodesEvalToSameValue(t, expected, `{://grammar.lang.wbnf.grammar: a -> '1'; :}`)

	exprs := []string{
		`a -> '1';`,
		`expr -> @:"+" > @:"*" > \d;`,
	}
	for _, expr := range exprs {
		expr := expr
		t.Run(expr, func(t *testing.T) {
			AssertCodesEvalToSameValue(t,
				"//grammar -> .parse(.lang.wbnf, 'grammar', `"+expr+"`)",
				`{://grammar.lang.wbnf.grammar:`+expr+`:}`,
			)
		})
	}
}

func TestGrammarParseParseLiteral(t *testing.T) {
	expected := `(
		"": ["+"],
		@rule: "expr",
		expr: [
			(expr: [("": "1")]),
			(
				"": ["*"],
				expr: [("": "2"), ("": "3")]
			)
		]
	)`
	AssertCodesEvalToSameValue(t,
		expected,
		`//grammar -> .parse(.parse(.lang.wbnf, "grammar", "expr -> @:'+' > @:'*' > \\d+;"), "expr", "1+2*3")`)

	scenarios := []struct{ grammar, rule, text string }{
		{`a -> '1' '2';`, "a", `12`},
		{`expr -> @:"+" > @:"*" > \d;`, "expr", `1+2*3`},
	}
	for _, s := range scenarios {
		s := s
		t.Run(s.text, func(t *testing.T) {
			parse := fmt.Sprintf(
				"//grammar -> .parse(.parse(.lang.wbnf, 'grammar', `%s`), '%s', `%s`)",
				s.grammar, s.rule, s.text)
			AssertCodesEvalToSameValue(t,
				parse,
				fmt.Sprintf(`{:{://grammar.lang.wbnf.grammar:%s:}.%s:%s:}`, s.grammar, s.rule, s.text))
		})
	}
}

func TestGrammarParseParseScopeVar(t *testing.T) {
	// AssertCodesEvalToSameValue(t,
	// 	`(
	// 		@rule: "x",
	// 		'': ["1", "2"],
	// 	)`,
	// 	`//grammar -> (.parse(.lang.wbnf, "grammar", "x -> '1' '2';") -> \x .parse(x, "x", "12"))`)

	// AssertCodesEvalToSameValue(t,
	// 	`(
	// 		"": ["+"],
	// 		@rule: "expr",
	// 		expr: [
	// 			(expr: [("": "1")]),
	// 			(
	// 				"": ["*"],
	// 				expr: [("": "2"), ("": "3")]
	// 			)
	// 		]
	// 	)`,
	// 	`//grammar -> (.parse(.lang.wbnf, "grammar", "expr -> @:'+' > @:'*' > \\d+;") -> \x .parse(x, "expr", "1+2*3"))`)

	scenarios := []struct{ grammar, rule, text string }{
		{`a -> "1" "2";`, "a", `12`},
		{`expr -> @:"+" > @:"*" > \d;`, "expr", `1+2*3`},
	}
	bindForms := []string{
		`{://grammar.lang.wbnf.grammar:%s:} -> {:.%s:%s:}`,
		`let g = {://grammar.lang.wbnf.grammar:%s:}; {:g.%s:%s:}`,
	}
	for i, s := range scenarios {
		s := s
		for j, form := range bindForms {
			form := form
			t.Run(fmt.Sprintf("%d.%d", i, j), func(t *testing.T) {
				parse := fmt.Sprintf(
					"//grammar -> (.parse(.lang.wbnf, 'grammar', `%s`) -> \\g .parse(g, '%s', `%s`))",
					s.grammar, s.rule, s.text)
				AssertCodesEvalToSameValue(t,
					parse,
					fmt.Sprintf(form, s.grammar, s.rule, s.text))
			})
		}
	}
}

// func TestGrammarParseWithEscape(t *testing.T) {
// 	AssertCodesEvalToSameValue(t,
// 		`(
// 			"": ["+"],
// 			@rule: "expr",
// 			expr: [
// 				(expr: [("": "1")]),
// 				(
// 					"": ["*"],
// 					expr: [("": "2"), ("": "3")]
// 				)
// 			]
// 		)`,
// 		`{://grammar.lang.wbnf.grammar: expr -> @:'+' > @:'*' > \d+; :} -> {:.expr:1+:{'2'}:*3:}`,
// 	)
// }
