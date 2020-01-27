package syntax

import "testing"

func TestGrammarToValueExprStd(t *testing.T) {
	AssertCodesEvalToSameValue(t, `1`, `//.grammar.parse(//.grammar.lang.wbnf, "grammar", "a-> '1';")`)
}
