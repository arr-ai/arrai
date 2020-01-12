package rel

import (
	"log"
	"testing"

	"github.com/arr-ai/wbnf/bootstrap"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeToValueSimple(t *testing.T) {
	grammar := `expr -> "+"|"*";`

	core := bootstrap.Core()
	rule := bootstrap.Rule("grammar")
	expr, err := core.Parse(rule, parser.NewScanner(grammar))
	assert.NoError(t, err)
	value := nodeToValue(core, rule, expr)
	log.Print(value)
}

func TestNodeToValueExpr(t *testing.T) {
	grammar := `
		expr -> expr:op="+"
		      ^ expr:op="*"
		      ^ n=/{[0-9]};
	`

	exprP, err := bootstrap.Compile(grammar)
	require.NoError(t, err)
	exprR := bootstrap.Rule("expr")
	math, err := exprP.Parse(exprR, parser.NewScanner("1+2"))
	assert.NoError(t, err)
	value := nodeToValue(exprP, exprR, math)
	log.Print(math)
	log.Print(value)
}
