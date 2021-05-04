package syntax

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestUnionSetHas(t *testing.T) {
	t.Parallel()
	unionSetAllMembers := []string{
		"(@: 0, @value: 97)",
		"(@: 1, @item: 98)",
		"(@: 2, @byte: 99)",
		"(@: 3, @char: 99)",
		"1", "2", "3", "{}",
		"()", "(a: 1)",
	}
	unionSetAll := `{` + strings.Join(unionSetAllMembers, ", ") + `}`
	for _, i := range unionSetAllMembers {
		AssertCodesEvalToSameValue(t, "true", fmt.Sprintf("%s <: %s", i, unionSetAll))
	}
}

func TestUnionSetWith(t *testing.T) {
	t.Parallel()

	u := `{1, (@: 0, @value: 1)}`
	AssertCodeEvalsToType(t, rel.UnionSet{}, u)
	AssertCodesEvalToSameValue(t, `{1, 2, (@: 0, @value: 1)}`, fmt.Sprintf("%s with 2", u))
	AssertCodesEvalToSameValue(t, `{1, (@: 1, @value: 2), (@: 0, @value: 1)}`, fmt.Sprintf("%s with (@: 1, @value: 2)", u))
	AssertCodesEvalToSameValue(t, `{1, 'abc', (@: 0, @value: 1)}`, fmt.Sprintf("%s with 'abc'", u))
	AssertCodesEvalToSameValue(t, u, fmt.Sprintf("%s with 1", u))
}

func TestUnionSetWithout(t *testing.T) {
	t.Parallel()

	// 3 elements to avoid changing set type
	u := `{1, (@: 0, @value: 0), (@: 0, @item: 0)}`
	AssertCodeEvalsToType(t, rel.UnionSet{}, u)
	AssertCodesEvalToSameValue(t, `{(@: 0, @value: 0), (@: 0, @item: 0)}`, fmt.Sprintf("%s without 1", u))
	AssertCodesEvalToSameValue(t, `{1, (@: 0, @value: 0)}`, fmt.Sprintf("%s without (@: 0, @item: 0)", u))
	AssertCodesEvalToSameValue(t, u, fmt.Sprintf("%s without 2", u))

	// test transformation into other Set
	AssertCodeEvalsToType(t, rel.Array{}, `{(@: 0, @item: 1), 1} without 1`)
	AssertCodeEvalsToType(t, rel.Bytes{}, `{(@: 0, @byte: 1), 1} without 1`)
	AssertCodeEvalsToType(t, rel.Dict{}, `{(@: 0, @value: 1), 1} without 1`)
	AssertCodeEvalsToType(t, rel.String{}, `{(@: 0, @char: 1), 1} without 1`)
	AssertCodeEvalsToType(t, rel.GenericSet{}, `{(@: 0, @value: 0), 1} without (@: 0, @value: 0)`)
	AssertCodeEvalsToType(t, rel.GenericSet{}, `{'a', 1} without 1`)
	AssertCodeEvalsToType(t, rel.UnionSet{}, `{(@: 0, @value: 0), "a", 1} without 1`)

	// TODO: test this when Relation is implemented
	// AssertCodeEvalsToType(t, rel.GenericSet{}, `{(a: 1), 1} without 1`)
}

func TestUnionSetWhere(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1]`, `{(@: 0, @item: 1), (@: 1, @value: 2)} where .@ != 1`)
	AssertCodesEvalToSameValue(t, `{1: 2}`, `{(@: 0, @item: 1), (@: 1, @value: 2)} where .@`)
	AssertCodeEvalsToType(t, rel.UnionSet{}, `{(@: 0, @item: 1), (@: 1, @value: 2), (@: 2, @byte: 3)} where .@ < 2`)
}

func TestUnionSetMap(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `'abc'`, `
		{(@: 0, @item: 97), (@: 1, @value: 98), (@: 2, @byte: 99)} => cond . {
			(:@item, :@): (:@, @char: @item),
			(:@value, :@): (:@, @char: @value),
			(:@byte, :@): (:@, @char: @byte),
		}
	`)
	AssertCodesEvalToSameValue(t, `{(@: 1, @item: 97), (@: 2, @value: 98), (@: 3, @byte: 99)}`, `
		{(@: 0, @item: 97), (@: 1, @value: 98), (@: 2, @byte: 99)} => (. +> (@: .@ + 1))
	`)
	AssertCodeEvalsToType(t, rel.UnionSet{}, `
		{(@: 0, @item: 97), (@: 1, @value: 98), (@: 2, @byte: 99)} => (. +> (@: .@ + 1))
	`)
}

func testTransform(t *testing.T, expectedType interface{}, expected, expr string) {
	AssertCodeEvalsToType(t, expectedType, expr)
	AssertCodesEvalToSameValue(t, expected, expr)
}

func testTransformIntoUnionSet(t *testing.T, expected, expr string) {
	testTransform(t, rel.UnionSet{}, expected, expr)
}

func TestUnionSetUnion(t *testing.T) {
	t.Parallel()

	testTransformIntoUnionSet(t,
		`{(@: 0, @byte: 97), (@: 2, @char: 97), (@: 0, @value: 0), (@: 1, @item: 1)}`,
		`{0: 0} | 1\[1] | <<'a'>> | 2\'a'`,
	)
	testTransformIntoUnionSet(t,
		`{(@: 0, @byte: 97), (@: 2, @char: 97), (@: 1, @item: 1)}`,
		`{(@: 0, @byte: 97), (@: 2, @char: 97)} | 1\[1]`,
	)
	testTransformIntoUnionSet(t,
		`{(@: 0, @byte: 97), (@: 2, @char: 97), (@: 1, @item: 1)}`,
		`1\[1] | {(@: 0, @byte: 97), (@: 2, @char: 97)}`,
	)
	testTransformIntoUnionSet(t,
		`{(@: 0, @byte: 97), (@: 2, @char: 97), (@: 1, @item: 1)}`,
		`{(@: 0, @byte: 97), (@: 2, @char: 97)} | {(@: 2, @char: 97), (@: 1, @item: 1)}`,
	)

	// TODO: uncomment when str, bytes, array duplicates union is fixed
	// testTransformIntoUnionSet(t,
	// 	`{(@: 0, @byte: 97), (@: 2, @char: 97), (@: 2, @char: 98), (@: 1, @item: 1)}`,
	// 	`{(@: 0, @byte: 97), (@: 2, @char: 98)} | {(@: 2, @char: 97), (@: 1, @item: 1)}`,
	// )
}

func TestUnionSetIntersect(t *testing.T) {
	t.Parallel()

	testTransformIntoUnionSet(t,
		`{(@: 0, @value: 0), (@: 0, @item: 0)}`,
		`{'a', (@: 0, @value: 0), (@: 0, @item: 0)} & {(@: 0, @value: 0), (@: 0, @item: 0), 'b'}`,
	)

	testTransform(t, rel.Array{}, `[0]`, `{'a', (@: 0, @item: 0)} & {(@: 0, @item: 0), 'b'}`)
	testTransform(t, rel.Dict{}, `{0: 0}`, `{'a', (@: 0, @value: 0)} & {(@: 0, @value: 0), 'b'}`)
	testTransform(t, rel.Bytes{}, `<<'a'>>`, `{'a', (@: 0, @byte: 97)} & {(@: 0, @byte: 97), 'b'}`)
	testTransform(t, rel.String{}, `1\'a'`, `{'a', (@: 1, @char: 97)} & {(@: 1, @char: 97), 'b'}`)
	testTransform(t, rel.GenericSet{}, `{1}`, `{1, (@: 1, @char: 97)} & {(@: 1, @byte: 97), 1}`)
	testTransform(t, rel.EmptySet{}, `{}`, `{2, (@: 1, @char: 97)} & {(@: 1, @byte: 97), 1}`)

	testTransform(t, rel.Array{}, `[0]`, `{'a', (@: 0, @item: 0)} & [0, 1]`)
	testTransform(t, rel.Dict{}, `{0: 0}`, `{'a', (@: 0, @value: 0)} & {0: 0, 1: 1}`)
	testTransform(t, rel.Bytes{}, `<<'a'>>`, `{'a', (@: 0, @byte: 97)} & <<'ab'>>`)
	testTransform(t, rel.String{}, `1\'a'`, `{'a', (@: 1, @char: 97)} & 1\'ab'`)
	testTransform(t, rel.GenericSet{}, `{1}`, `{1, (@: 1, @char: 97)} & {1}`)
	testTransform(t, rel.EmptySet{}, `{}`, `{2, (@: 1, @char: 97)} & {}`)
}

func TestUnionSetDifference(t *testing.T) {
	t.Parallel()

	testTransformIntoUnionSet(t,
		`{(@: 0, @value: 0), (@: 0, @item: 0)}`,
		`
			{(@: 0, @value: 0), (@: 0, @item: 0), (@: 0, @byte: 97), (@: 0, @char: 97)}
			&~ {(@: 0, @byte: 97), (@: 0, @char: 97)}
		`,
	)
	testTransform(t,
		rel.Array{},
		`[0]`,
		`
			{(@: 0, @item: 0), (@: 0, @byte: 97), (@: 0, @char: 97)}
			&~ {(@: 0, @byte: 97), (@: 0, @char: 97)}
		`,
	)
	testTransform(t,
		rel.Array{},
		`[0]`,
		`{(@: 0, @item: 0)} &~ {(@: 0, @byte: 97), (@: 0, @char: 97)}`,
	)
	testTransform(t,
		rel.Dict{},
		`{0: 0}`,
		`
			{(@: 0, @value: 0), (@: 0, @byte: 97), (@: 0, @char: 97)}
			&~ {(@: 0, @byte: 97), (@: 0, @char: 97)}
		`,
	)
	testTransform(t,
		rel.Dict{},
		`{0: 0}`,
		`{(@: 0, @value: 0)} &~ {(@: 0, @byte: 97), (@: 0, @char: 97)}`,
	)
	testTransform(t,
		rel.String{},
		`'a'`,
		`
			{(@: 0, @char: 97), (@: 0, @item: 0), (@: 0, @byte: 97)}
			&~ {(@: 0, @byte: 97), (@: 0, @item: 0)}
		`,
	)
	testTransform(t,
		rel.String{},
		`'a'`,
		`{(@: 0, @char: 97)} &~ {(@: 0, @byte: 97), (@: 0, @item: 0)}`,
	)
	testTransform(t,
		rel.Bytes{},
		`<<'a'>>`,
		`
			{(@: 0, @byte: 97), (@: 0, @item: 0), (@: 0, @char: 97)}
			&~ {(@: 0, @char: 97), (@: 0, @item: 0)}
		`,
	)
	testTransform(t,
		rel.Bytes{},
		`<<'a'>>`,
		`{(@: 0, @byte: 97)} &~ {(@: 0, @char: 97), (@: 0, @item: 0)}`,
	)
	testTransform(t,
		rel.GenericSet{},
		`{1}`,
		`{1, (@: 0, @byte: 97), (@: 0, @char: 97)} &~ {2, (@: 0, @byte: 97), (@: 0, @char: 97)}`,
	)
	testTransform(t, rel.GenericSet{}, `{1}`, `{1} &~ {2, (@: 0, @byte: 97), (@: 0, @char: 97)}`)
	testTransform(t, rel.EmptySet{},
		`{}`,
		`{(@: 0, @item: 0), (@: 0, @byte: 97)} &~ {(@: 0, @item: 0), (@: 0, @byte: 97)}`,
	)

	testTransformIntoUnionSet(t,
		`{(@: 0, @item: 0), (@: 0, @value: 0), (@: 1, @byte: 97)}`,
		`{(@: 0, @item: 0), (@: 0, @value: 0), (@: 0, @byte: 97), (@: 1, @byte: 97)} &~ <<'a'>>`,
	)
	testTransform(t, rel.Array{}, `[0]`, `{(@: 0, @item: 0), (@: 0, @value: 0)} &~ {0: 0}`)
	testTransform(t, rel.Dict{}, `{0: 0}`, `{(@: 0, @item: 0), (@: 0, @value: 0)} &~ [0]`)
	testTransform(t, rel.String{}, `'a'`, `{(@: 0, @char: 97), (@: 0, @byte: 97)} &~ <<'a'>>`)
	testTransform(t, rel.Bytes{}, `<<'a'>>`, `{(@: 0, @char: 97), (@: 0, @byte: 97)} &~ 'a'`)
	testTransform(t, rel.GenericSet{}, `{1}`, `{1, (@: 0, @char: 97)} &~ 'a'`)

	testTransform(t, rel.Array{}, `1\[1]`, `[0, 1] &~ {(@: 0, @item: 0), (@: 0, @value: 0)}`)
	testTransform(t, rel.Dict{}, `{1: 1}`, `{0: 0, 1: 1} &~ {(@: 0, @item: 0), (@: 0, @value: 0)}`)
	testTransform(t, rel.String{}, `1\'bc'`, `'abc' &~ {(@: 0, @char: 97), (@: 0, @byte: 97)}`)
	testTransform(t, rel.Bytes{}, `1\<<'bc'>>`, `<<'abc'>> &~ {(@: 0, @char: 97), (@: 0, @byte: 97)}`)
	testTransform(t, rel.GenericSet{}, `{2}`, `{1, 2} &~ {1, (@: 0, @char: 97)}`)
	testTransform(t, rel.EmptySet{}, `{}`, `{1} &~ {1, (@: 0, @char: 97)}`)
}
