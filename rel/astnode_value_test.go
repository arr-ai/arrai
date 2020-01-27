package rel

import (
	"fmt"
	"testing"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stripASTNodeSrc(node ast.Node) ast.Node {
	switch node := node.(type) {
	case ast.Leaf:
		return ast.Leaf(node.Scanner().StripSource())
	case ast.Branch:
		result := make(ast.Branch, len(node))
		for name, children := range node {
			switch children := children.(type) {
			case ast.One:
				result[name] = ast.One{Node: stripASTNodeSrc(children.Node)}
			case ast.Many:
				many := make(ast.Many, 0, len(children))
				for _, node := range children {
					many = append(many, stripASTNodeSrc(node))
				}
				result[name] = many
			}
		}
		return result
	case ast.Extra:
		return node
	default:
		panic(fmt.Errorf("unexpected: %v %[1]T", node))
	}
}

func assertASTNodeToValueToNode(t *testing.T, p wbnf.Parsers, rule, src string) bool { //nolint:unparam
	v, err := p.Parse(wbnf.Rule(rule), parser.NewScanner(src))
	assert.NoError(t, err)
	ast1 := ast.ParserNodeToNode(p.Grammar(), v)
	value := ASTBranchToValue(ast1)
	ast2 := ASTBranchFromValue(value)
	return assert.EqualValues(t, stripASTNodeSrc(ast1), ast2)
}

func TestNodeToValueSimple(t *testing.T) {
	assertASTNodeToValueToNode(t, wbnf.Core(), "grammar", `expr -> "+"|"*";`)
}

func TestGrammarToValueExpr(t *testing.T) {
	assertASTNodeToValueToNode(t, wbnf.Core(), "grammar", `x->@:"+" > @:"*" > "1";`)
}

func TestGrammarToValueCore(t *testing.T) {
	assertASTNodeToValueToNode(t, wbnf.Core(), "grammar", wbnf.GrammarGrammar())
}

func TestNodeToValueExpr(t *testing.T) {
	grammar := `expr -> @:op="+" > @:op="*" > n=/{[0-9]};`

	exprP, err := wbnf.Compile(grammar)
	require.NoError(t, err)
	assertASTNodeToValueToNode(t, exprP, "expr", `1+2*3`)
}
