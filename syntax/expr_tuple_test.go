package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestTupleType(t *testing.T) {
	t.Parallel()
	AssertCodeEvalsToType(t, rel.StringCharTuple{}, `(@: 1, @char: 65)`)
	AssertCodeEvalsToType(t, rel.ArrayItemTuple{}, `(@: 1, @item: 2)`)
	AssertCodeEvalsToType(t, rel.DictEntryTuple{}, `(@: {1, 2}, @value: 2)`)
}
