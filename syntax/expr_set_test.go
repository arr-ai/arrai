package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestString(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{|@,@char| (0, 97), (1, 98), (2, 99)}`, `"abc"`)
	AssertCodesEvalToSameValue(t, `{1: 2, 3: 4}`, `{(@: 1, @value: 2), (@: 3, @value: 4)}`)
}

func TestStringType(t *testing.T) {
	t.Parallel()
	AssertCodeEvalsToType(t, rel.String{}, `"abc"`)
	AssertCodeEvalsToType(t, rel.String{}, `{|@,@char| (0, 97)}`)
	AssertCodeEvalsToType(t, rel.String{}, `{(@: 0, @char: 97)}`)
	AssertCodeEvalsToType(t, rel.String{}, `{(@: 0, @char: 97), (@: 1, @char: 98)}`)
	AssertCodeEvalsToType(t, rel.String{}, `{(@: 0, @char: 97), (@: 0, @char: 98)}`)
	AssertCodeEvalsToType(t, rel.String{}, `"abc" >> .`)
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

func TestDict(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{|@,@value| ("x", "y")}`, `{"x": "y"}`)
	AssertCodesEvalToSameValue(t, `{1: 2, 3: 4}`, `{(@: 1, @value: 2), (@: 3, @value: 4)}`)
}

func TestDictType(t *testing.T) {
	t.Parallel()
	AssertCodeEvalsToType(t, rel.Dict{}, `{1:2}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{|@,@value| (1, 2)}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{(@: 1, @value: 2)}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{(@: 1, @value: 2), (@: 3, @value: 4)}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{(@: 1, @value: 2), (@: 1, @value: 3)}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{1:2} >> .`)
}
