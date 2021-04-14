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
	AssertCodeEvalsToType(t, rel.String{}, `"abc" >> .`)
	AssertCodeEvalsToType(t, rel.String{}, `"abc" ++ "def"`)
}

func TestStringWhere(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"abc"`, `"abc" where .@ < 10`)
	AssertCodesEvalToSameValue(t, `"abc"`, `"abc" where .@char < 100`)
	AssertCodesEvalToSameValue(t, `"ab"`, `"abc" where .@ < 2`)
	AssertCodesEvalToSameValue(t, `"ab"`, `"abc" where .@char < 99`)
	// TODO: Test for offset strings and holey strings.
}

func TestStringManipulation(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"Foo"`, `(\s //str.upper(s where .@=0) | (s where .@>0))("foo")`)
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

func TestDictWhere(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{"a": "b"}`, `{"a": "b"} where .@ = "a"`)
	AssertCodesEvalToSameValue(t, `{}`, `{"a": "b"} where .@ = "b"`)
}

func TestDictUnion(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{"a": "a",    "b": "a"} | {"a": "b"}`,
		`{"a": "a"} | {"b": "a",    "a": "b"}`)
}

func TestDictExpand(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"{'a': 'a', 'a': 'b', 'b': 'a'}"`,
		`//str.expand("", {"b": "a", "a": "b"} | {"a": "a"}, "", "")`)
}

func TestDictEqual(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true `, `{'a': 1} = {'a': 1}                                `)
	AssertCodesEvalToSameValue(t, `true `, `{'a': 1, 'b': 2, 'c': 3} = {'a': 1, 'b': 2, 'c': 3}`)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1} = {1}                                     `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1} = {1}                                     `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1, 'b': 1, 'c': 1} = {1, 2, 3}               `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1} = [1]                                     `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1, 'b': 1, 'c': 1} = [1, 2, 3]               `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1} = 1                                       `)
	AssertCodesEvalToSameValue(t, `false`, `{'a': 1} = (a: 1)                                  `)
}
