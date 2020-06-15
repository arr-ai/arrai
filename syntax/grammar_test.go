package syntax

import (
	"fmt"
	"testing"
)

const expected = `(
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

func TestGrammarToValueExprQualified(t *testing.T) {
	AssertCodesEvalToSameValue(t, expected, `//grammar.parse(//grammar.lang.wbnf, "grammar", "a -> '1';")`)
}

func TestGrammarToValueExprScoped(t *testing.T) {
	AssertCodesEvalToSameValue(t, expected, `//grammar -> .parse(.lang.wbnf, "grammar", "a -> '1';")`)
}

func TestGrammarToValueExprInline(t *testing.T) {
	AssertCodesEvalToSameValue(t, expected, `{://grammar.lang.wbnf[grammar]: a -> '1'; :}`)
}

func TestGrammarToValueExprInlineDefault(t *testing.T) {
	AssertCodesEvalToSameValue(t, expected, `{://grammar.lang.wbnf: a -> '1'; :}`)
}

func TestMacroToValueInline(t *testing.T) {
	AssertCodesEvalToSameValue(t, `(year: 2020, month: 06, day: 09)`, `
		let time = (
			@grammar: {://grammar.lang.wbnf: date -> year=\d{4} "-" month=\d{2} "-" day=\d{2};:},
			@transform: (date: \ast ast -> (year: .year, month: .month, day: .day) :> //eval.value(.''))
		);
		{:time:2020-06-09:}
	`)
}

func TestArraiGrammarMacroEquality(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`//grammar.parse(//grammar.lang.arrai)("expr", "1")`,
		`{://grammar.lang.arrai:1:}`,
	)
}

// TODO(ladeo): Figure out why this fails and fix it.
//func TestArraiGrammarGrammarEquality(t *testing.T) {
//	AssertCodeEvalsToGrammar(t, arraiParsers.Grammar(), `//grammar.lang.arrai`)
//}

// TODO(ladeo): Figure out why this fails and fix it.
//func TestMacroToArraiValueInline(t *testing.T) {
//	AssertCodesEvalToSameValue(t,
//		`1`,
//		`{:(@grammar://grammar.lang.arrai, @transform:(expr:\ast 1)):1:}`,
//	)
//}

func TestGrammarToValueExprScopedAndInline(t *testing.T) {
	exprs := []string{
		`a -> '1';`,
		`expr -> @:"+" > @:"*" > \d;`,
	}
	for _, expr := range exprs {
		expr := expr
		t.Run(expr, func(t *testing.T) {
			t.Parallel()
			AssertCodesEvalToSameValue(t,
				"//grammar -> .parse(.lang.wbnf, 'grammar', `"+expr+"`)",
				`{://grammar.lang.wbnf[grammar]:`+expr+`:}`,
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
				fmt.Sprintf(`{:{://grammar.lang.wbnf[grammar]:%s:}[%s]:%s:}`, s.grammar, s.rule, s.text))
		})
	}
}

func TestGrammarParseParseScopeVar(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`(
			@rule: "x",
			'': ["1", "2"],
		)`,
		`//grammar -> (.parse(.lang.wbnf, "grammar", "x -> '1' '2';") -> \x .parse(x, "x", "12"))`)

	AssertCodesEvalToSameValue(t,
		`(
			"": ["+"],
			@rule: "expr",
			expr: [
				(expr: [("": "1")]),
				(
					"": ["*"],
					expr: [("": "2"), ("": "3")]
				)
			]
		)`,
		`//grammar -> (.parse(.lang.wbnf, "grammar", "expr -> @:'+' > @:'*' > \\d+;") -> \x .parse(x, "expr", "1+2*3"))`)

	scenarios := []struct{ grammar, rule, text string }{
		{`a -> "1" "2";`, "a", `12`},
		{`expr -> @:"+" > @:"*" > \d;`, "expr", `1+2*3`},
	}
	bindForms := []string{
		`{://grammar.lang.wbnf[grammar]:%s:} -> {:.[%s]:%s:}`,
		`let g = {://grammar.lang.wbnf[grammar]:%s:}; {:g[%s]:%s:}`,
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
// 		`{://grammar.lang.wbnf[grammar]: expr -> @:'+' > @:'*' > \d+; :} -> {:.expr:1+:{'2'}:*3:}`,
// 	)
// }
