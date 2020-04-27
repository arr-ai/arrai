package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineCollectorAppendLine(t *testing.T) {
	t.Parallel()

	c := newLineCollector()

	// testing string with no opener
	c.appendLine("random string")
	assert.Equal(t, []string{"random string"}, c.lines)
	assert.True(t, c.isBalanced())
	c.reset()

	// basic testing string with opener
	c.appendLine("random {")
	assert.Equal(t, []string{"random {"}, c.lines)
	assert.Equal(t, 1, len(c.stack))
	assert.Equal(t, "}", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	assert.False(t, c.isBalanced())
	c.appendLine("}")
	assert.Equal(t, []string{"random {", "}"}, c.lines)
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing with multiple openers
	c.appendLine("{([")
	assert.Equal(t, 3, len(c.stack))
	assert.Equal(t, "}", c.stack[0].char)
	assert.Equal(t, ")", c.stack[1].char)
	assert.Equal(t, "]", c.stack[2].char)
	c.appendLine("])}")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing context based opener
	c.appendLine("nested closer $'{([")
	c.appendLine("${")
	assert.Equal(t, []string{"nested closer $'{([", "${"}, c.lines)
	assert.Equal(t, 2, len(c.stack))
	assert.Equal(t, "'", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	assert.Equal(t, "}", c.stack[1].char)
	assert.Equal(t, true, c.stack[1].recursive)
	c.appendLine("}")
	assert.Equal(t, 1, len(c.stack))
	assert.Equal(t, "'", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	c.appendLine("'")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing non recursive opener
	c.appendLine("'")
	c.appendLine("\"")
	assert.Equal(t, []string{"'", "\""}, c.lines)
	assert.Equal(t, 1, len(c.stack))
	c.appendLine("'")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing escape
	c.appendLine("`\"`")
	assert.True(t, c.isBalanced())
}

func TestIsBalanced(t *testing.T) {
	t.Parallel()

	c := newLineCollector()
	assert.True(t, c.isBalanced())

	c.appendLine("random;")
	assert.False(t, c.isBalanced())

	c.appendLine("\\a \\b \\c")
	assert.False(t, c.isBalanced())

	c.appendLine("c:")
	assert.False(t, c.isBalanced())

	c.appendLine("\"")
	assert.False(t, c.isBalanced())

	c.appendLine("\"")
	assert.True(t, c.isBalanced())
}

