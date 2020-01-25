package rel

import (
	"log"
	"testing"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeToValueSimple(t *testing.T) {
	grammar := `expr -> "+"|"*";`

	core := wbnf.Core()
	rule := wbnf.Rule("grammar")
	expr, err := core.Parse(rule, parser.NewScanner(grammar))
	assert.NoError(t, err)
	value := nodeToValue(ast.ParserNodeToNode(core.Grammar(), expr))
	log.Print(value)
}

func TestNodeToValueExpr(t *testing.T) {
	grammar := `
		expr -> @:op="+"
		      > @:op="*"
		      > n=/{[0-9]};
	`

	exprP, err := wbnf.Compile(grammar)
	require.NoError(t, err)
	exprR := wbnf.Rule("expr")
	math, err := exprP.Parse(exprR, parser.NewScanner("1+2*3"))
	assert.NoError(t, err)
	value := nodeToValue(ast.ParserNodeToNode(exprP.Grammar(), math))
	log.Print(math)
	log.Print(value)
}

func TestGrammarToValueExpr(t *testing.T) {
	grammar := `x->@:"+" > @:"*" > "1";`

	expr, err := wbnf.Core().Parse(wbnf.Rule("grammar"), parser.NewScanner(grammar))
	require.NoError(t, err)
	value := nodeToValue(ast.ParserNodeToNode(wbnf.Core().Grammar(), expr))
	log.Print(value)
}
