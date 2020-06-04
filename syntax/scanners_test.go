package syntax

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestHandleAccessScanners(t *testing.T) {
	base := parser.NewScanner(".")
	access := parser.NewScannerAt(".a.b", 1, 2)

	assert.Equal(t, 0, handleAccessScanners(*base, *access).Offset())

	assert.Equal(t, *access, handleAccessScanners(parser.Scanner{}, *access))
}
