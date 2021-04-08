package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestHasRule(t *testing.T) {
	mustAST := func(source string) ast.Branch {
		pc := ParseContext{SourceDir: "."}
		ast, err := pc.Parse(context.Background(), parser.NewScannerWithFilename(source, "test"))
		assert.NoError(t, err)
		return ast
	}
	assert.True(t, hasRule(mustAST(`a +> (b+>: (c+>: 1))`), "nested_op"))
	assert.True(t, hasRule(mustAST(`a +> (b: (c+>: 1))`), "nested_op"))
	assert.False(t, hasRule(mustAST(`a +> (b+>: (c+>: 1))`), "binop"))
}
