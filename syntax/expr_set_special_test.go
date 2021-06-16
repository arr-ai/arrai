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

func TestRelationMap(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{0: 0} `, `{|@, @value, z| (0, 0, 0) } => .~|z|`)
	AssertCodesEvalToSameValue(t, `'a'    `, `{|@, @char,  z| (0, 97, 0)} => .~|z|`)
	AssertCodesEvalToSameValue(t, `[0]    `, `{|@, @item,  z| (0, 0, 0) } => .~|z|`)
	AssertCodesEvalToSameValue(t, `<<'a'>>`, `{|@, @byte,  z| (0, 97, 0)} => .~|z|`)

	AssertCodesEvalToSameValue(t,
		`{|a, b, c| (2, 3, -1), (3, 4, -1), (4, 5, -1)}`,
		`{|a, b| (1, 2), (2, 3), (3, 4)} => (a: .a+1, b: .b+1, c: .a - .b)`,
	)
}

func TestRelationNest(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{|z, dict| (0, {0: 0})}`,
		`{|@, @value, z| (0, 0, 0) } nest ~|z|dict`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, dict| (0, {0: 0})}`,
		`{|@, @value, z| (0, 0, 0) } nest |@, @value|dict`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, str| (0, 'a')}`,
		`{|@, @char, z| (0, 97, 0) } nest ~|z|str`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, str| (0, 'a')}`,
		`{|@, @char, z| (0, 97, 0) } nest |@, @char|str`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, array| (0, [0])}`,
		`{|@, @item, z| (0, 0, 0) } nest ~|z|array`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, array| (0, [0])}`,
		`{|@, @item, z| (0, 0, 0) } nest |@, @item|array`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, byte| (0, <<'a'>>)}`,
		`{|@, @byte, z| (0, 97, 0) } nest ~|z|byte`,
	)
	AssertCodesEvalToSameValue(t,
		`{|z, byte| (0, <<'a'>>)}`,
		`{|@, @byte, z| (0, 97, 0) } nest |@, @byte|byte`,
	)

	AssertCodesEvalToSameValue(t,
		`{|a, nested| (1, {|b, c| (1, 2), (2, 3)}), (2, {|b, c| (2, 2)})}`,
		`{|a, b, c| (1, 1, 2), (1, 2, 3), (2, 2, 2)} nest |b, c|nested`,
	)
}

func TestRelationWhere(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{|a, b, c| (0, 1, 2), (2, 3, 4), (4, 5, 6)}`,
		`{|a, b, c| (0, 1, 2), (1, 2, 3), (2, 3, 4), (3, 4, 5), (4, 5, 6)} where .a % 2 = 0`,
	)
}

func TestSetBuilder(t *testing.T) {
	t.Parallel()

	AssertCodeEvalsToType(t, rel.GenericSet{}, `{1}`)
	AssertCodeEvalsToType(t, rel.GenericSet{}, `{1, 2, 3}`)
	AssertCodeEvalsToType(t, rel.GenericSet{}, `{"test", "another test", 123}`)
	AssertCodeEvalsToType(t, rel.GenericSet{}, `{{(a: 1), (b: 1)}, {(c: 1), 1}}`)
	AssertCodeEvalsToType(t, rel.True, `{()}`)
	AssertCodeEvalsToType(t, rel.None, `{}`)

	AssertCodeEvalsToType(t, rel.String{}, `"abc"`)
	AssertCodeEvalsToType(t, rel.String{}, `{(@: 0, @char: 97), (@: 1, @char: 98), (@: 2, @char: 99)}`)
	AssertCodeEvalsToType(t, rel.String{}, `{|@, @char| (0, 97), (1, 98), (2, 99)}`)

	AssertCodeEvalsToType(t, rel.Bytes{}, `<<"abc">>`)
	AssertCodeEvalsToType(t, rel.Bytes{}, `{(@: 0, @byte: 97), (@: 1, @byte: 98), (@: 2, @byte: 99)}`)
	AssertCodeEvalsToType(t, rel.Bytes{}, `{|@, @byte| (0, 97), (1, 98), (2, 99)}`)

	AssertCodeEvalsToType(t, rel.Array{}, `[1, 2, 3]`)
	AssertCodeEvalsToType(t, rel.Array{}, `{(@: 0, @item: 97), (@: 1, @item: 98), (@: 2, @item: 99)}`)
	AssertCodeEvalsToType(t, rel.Array{}, `{|@, @item| (0, 97), (1, 98), (2, 99)}`)

	AssertCodeEvalsToType(t, rel.Dict{}, `{0: 1, 1: 2, 2: 3}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{(@: 0, @value: 97), (@: 1, @value: 98), (@: 2, @value: 99)}`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{|@, @value| (0, 97), (1, 98), (2, 99)}`)

	AssertCodeEvalsToType(t, rel.UnionSet{},
		`{(@: 0, @value: 97), (@: 1, @item: 98), (@: 2, @byte: 99), (@: 3, @char: 99), 1}`)
	AssertCodeEvalsToType(t, rel.UnionSet{}, `{(a: 1), (b: 1), (a: 1, b: 1)}`)

	AssertCodeEvalsToType(t, rel.Relation{}, `{|a, b| (0, 0), (1, 1), (2, 2)}`)
	AssertCodeEvalsToType(t, rel.Relation{}, `{(a: 1)}`)
}
