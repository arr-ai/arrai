package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestArrayToString(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"hello"`, `[104, 101, 108, 108, 111] => (@:.@, @char:.@item)`)
	AssertCodeEvalsToType(t, rel.String{}, `[104, 101, 108, 108, 111] => (@:.@, @char:.@item)`)
}

func TestArray(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{|@,@item| (0, 1), (1, 2), (2, 3)}`, `[1, 2, 3]`)
	AssertCodesEvalToSameValue(t, `{1: 2, 3: 4}`, `{(@: 1, @value: 2), (@: 3, @value: 4)}`)
}

func TestArrayType(t *testing.T) {
	t.Parallel()
	AssertCodeEvalsToType(t, rel.Array{}, `[1, 2, 3]`)
	AssertCodeEvalsToType(t, rel.Array{}, `{|@,@item| (0, 1)}`)
	AssertCodeEvalsToType(t, rel.Array{}, `{(@: 0, @item: 1)}`)
	AssertCodeEvalsToType(t, rel.Array{}, `{(@: 0, @item: 1), (@: 1, @item: 2)}`)
	AssertCodeEvalsToType(t, rel.Array{}, `{(@: 0, @item: 1), (@: 0, @item: 2)}`)
	AssertCodeEvalsToType(t, rel.Array{}, `[1, 2, 3] >> .`)
}

func TestArrayWhere(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[1, 2]`, `[1, 2] where .@ < 10`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `[1, 2] where .@item < 10`)
	AssertCodesEvalToSameValue(t, `[1]`, `[1, 2] where .@ < 1`)
	AssertCodesEvalToSameValue(t, `[1]`, `[1, 2] where .@item < 2`)
	AssertCodesEvalToSameValue(t, `1\[2]`, `[1, 2] where .@ > 0`)
	AssertCodesEvalToSameValue(t, `1\[2]`, `[1, 2] where .@item > 1`)
}

func TestArrayWithHoles(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[10] | 2\[12]`, `[10, , 12]`)
	AssertCodesEvalToSameValue(t, `[10, 11, 12]`, `[10, , 12] with (@: 1, @item: 11)`)
	AssertCodesEvalToSameValue(t, `[10, 11, 12]`, `[10, , 12, , , ] with (@: 1, @item: 11)`)
}

func TestSetOfArrays(t *testing.T) {
	t.Parallel()

	// This tests an error when stringifying sets of arrays.
	AssertCodesEvalToSameValue(t, `'{[1], [1, 1]}'`, `$'${{[1, 1], [1]}}'`)
}
