package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/rel/syntax"
)

// TestNumberJSON tests JSON-encoding of Numbers.
func TestNumberJSON(t *testing.T) {
	assertJSON(t, `0`, 0)
	assertJSON(t, `42`, 42)
	assertJSON(t, `-1.5`, -1.5)
	assertJSON(t, `3.14e-45`, 3.14e-45)
}

// TestTupleJSON tests JSON-encoding of Tuples.
func TestTupleJSON(t *testing.T) {
	assertJSON(t, `{}`, map[string]interface{}{})
	assertJSON(t, `{"a":1}`, map[string]interface{}{"a": 1})
	assertJSON(t, []string{`{"a":1,"b":2}`, `{"b":2,"a":1}`},
		map[string]interface{}{"a": 1, "b": 2})
}

// TestSetJSON tests JSON-encoding of Sets.
func TestSetJSON(t *testing.T) {
	assertJSON(t, `false`, []interface{}{})
	assertJSON(t, `{"{||}":[1]}`, []interface{}{1})
	assertJSON(t, []string{`{"{||}":[1,2]}`, `{"{||}":[2,1]}`},
		[]interface{}{1, 2})
}

// TestMixedJSON tests JSON-encoding of mixed Values.
func TestMixedJSON(t *testing.T) {
	assertJSON(t, `false`, []interface{}{})
	assertJSON(t, `{"{||}":[1]}`, []interface{}{1})
	assertJSON(t,
		ucl(`{"{||}":[`, `]}`,
			`2`,
			ucl(`{`, `}`,
				`"a":1`,
				ucl(`"b":{"{||}":[`, `]}`,
					`3`, `4`,
				),
			),
		).permutations(),
		[]interface{}{
			2,
			map[string]interface{}{
				"a": 1,
				"b": []interface{}{3, 4},
			},
		})
}

// TestXMLChildrenJSON tests JSON-encoding of XML node with child.
func TestXMLChildJSON(t *testing.T) {
	value, err := syntax.Parse([]byte(`<abc><def/></abc>`))
	require.NoError(t, err)

	assertJSON(t,
		ucl(`{"@xml":{`, `}}`,
			`"tag":"abc"`,
			ucl(`"children":[{"@xml":{`, `}}]`,
				`"tag":"def"`,
			),
		).permutations(),
		value,
	)
}

func TestPermutations(t *testing.T) {
	assert.Equal(t,
		[][]interface{}{
			{1, 2, 3},
			{1, 3, 2},
			{2, 1, 3},
			{2, 3, 1},
			{3, 1, 2},
			{3, 2, 1},
		},
		permutations([]interface{}{1, 2, 3}))
}

func TestUnorderedCommaList(t *testing.T) {
	assert.Equal(t,
		[]string{"[1,2]", "[2,1]"},
		ucl("[", "]", "1", "2").permutations())
}

func TestUnorderedCommaList2(t *testing.T) {
	assert.Equal(t,
		[]string{"[{1,2},3]", "[{2,1},3]", "[3,{1,2}]", "[3,{2,1}]"},
		ucl("[", "]", ucl("{", "}", "1", "2"), "3").permutations())
}

func permutations(elts []interface{}) [][]interface{} {
	if len(elts) == 0 {
		return [][]interface{}{elts}
	}

	exclude := func(i int) []interface{} {
		result := make([]interface{}, len(elts)-1)
		copy(result[:i], elts[:i])
		copy(result[i:], elts[i+1:])
		return result
	}

	result := [][]interface{}{}
	for i, elt := range elts {
		for _, tail := range permutations(exclude(i)) {
			perm := make([]interface{}, len(elts))
			perm[0] = elt
			copy(perm[1:], tail)
			result = append(result, perm)
		}
	}
	return result
}

type unorderedCommaList struct {
	prefix string
	suffix string
	elts   []interface{}
}

func ucl(prefix, suffix string, elts ...interface{}) *unorderedCommaList {
	return &unorderedCommaList{prefix, suffix, elts}
}

func generate(prefix, suffix string, parts []interface{}, out *[]string) {
	if len(parts) == 0 {
		*out = append(*out, prefix+suffix)
	} else {
		for _, head := range parts[0].([]string) {
			p := prefix + head
			if len(parts) > 1 {
				p += ","
			}
			generate(p, suffix, parts[1:], out)
		}
	}
}

func (u *unorderedCommaList) permutations() []string {
	if len(u.elts) == 0 {
		return []string{""}
	}

	eltPerms := make([]interface{}, len(u.elts))
	for i, elt := range u.elts {
		if s, ok := elt.(string); ok {
			eltPerms[i] = []string{s}
		} else {
			eltPerms[i] = elt.(*unorderedCommaList).permutations()
		}
	}

	result := []string{}
	for _, perm := range permutations(eltPerms) {
		generate(u.prefix, u.suffix, perm, &result)
	}
	return result
}

func assertJSON(t *testing.T, expected interface{}, value interface{}) {
	var expecteds []string
	switch x := expected.(type) {
	case []string:
		expecteds = x
		require.NotEqual(t, 0, len(x))
	case string:
		expecteds = []string{x}
	default:
		require.Fail(
			t, "expected must be string or []string", "%v", value)
	}
	v, err := rel.NewValue(value)
	require.NoError(t, err)
	j := rel.MarshalToJSON(v)
	s := string(j)
	ok := false
	for _, e := range expecteds {
		if e == s {
			ok = true
			break
		}
	}
	// Iff none succed, fail them all.
	if !ok {
		assert.Fail(t, "no permutation matches JSON",
			"%v doesn't contain %v", expected, s)
		return
	}

	v2, err := rel.UnmarshalFromJSON(j)
	if assert.NoError(t, err) {
		assert.True(t, v.Equal(v2), "%s == %s", v, v2)
	}

	for _, e := range expecteds {
		v3, err := rel.UnmarshalFromJSON([]byte(e))
		if assert.NoError(t, err) {
			assert.True(t, v.Equal(v3), "%s == %s", v, v2)
		}
	}
}
