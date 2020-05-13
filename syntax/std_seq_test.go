package syntax

import "testing"

func TestSeqConcat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""              `, `//seq.concat([])                            `)
	AssertCodesEvalToSameValue(t, `""              `, `//seq.concat(["", "", "", ""])              `)
	AssertCodesEvalToSameValue(t, `"hello"         `, `//seq.concat(["", "", "", "", "hello", ""]) `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//seq.concat(["this", " is", " a", " test"])`)
	AssertCodesEvalToSameValue(t, `"this"          `, `//seq.concat(["this"])                      `)
	assertExprPanics(t, `//seq.concat("this")`)

	AssertCodesEvalToSameValue(t, `[]`, `//seq.concat([[]])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.concat([[], []])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `//seq.concat([[1, 2, 3]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `//seq.concat([[1, 2, 3], []])`)
	AssertCodesEvalToSameValue(t, `[4, 5, 6]`, `//seq.concat([[], [4, 5, 6]])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 4, 5, 6]`, `//seq.concat([[1, 2, 3], [4, 5, 6]])`)
	AssertCodePanics(t, `//seq.concat(42)`)
}

func TestSeqRepeat(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `""      `, `//seq.repeat(0, "ab")`)
	AssertCodesEvalToSameValue(t, `""      `, `//seq.repeat(3, "")`)
	AssertCodesEvalToSameValue(t, `"ab"    `, `//seq.repeat(1, "ab")`)
	AssertCodesEvalToSameValue(t, `"ababab"`, `//seq.repeat(3, "ab")`)

	AssertCodesEvalToSameValue(t, `[]                `, `//seq.repeat(0, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[]                `, `//seq.repeat(2, [])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]         `, `//seq.repeat(1, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, 1, 2, 3]`, `//seq.repeat(2, [1, 2, 3])`)

	AssertCodePanics(t, `//seq.repeat(2, 3.4)`)
}

func TestStrContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("this is a test", "")             `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("this is a test", "is a test")    `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("this is a test", "is not a test")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("this is a test", "a is")`)
	assertExprPanics(t, `//seq.contains(123, 124)`)
}

func TestArrayContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([1,2,3,4,5],1)`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([1,2,3,4,5],3)`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([1,2,3,4,5],5)`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],'A')`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],'E')`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],'C')`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['A','B','C'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['B','C'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E'],['B','C','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E'],['A','B','C','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E'],['A','B','C','D','E','F'])`)

	// TODO, it requires API contains code change
	// AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A', 'A', 'B','C','D','E'],['A','B','C'])`)
}
func TestBytesContains(t *testing.T) {
	t.Parallel()
}

///////////////////
func TestStrSub(t *testing.T) {
	t.Parallel()
	// string
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("this is not a test", "is not", "is")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("this is not a test", "not ", "")`)
	AssertCodesEvalToSameValue(t,
		`"this is still a test"`,
		`//seq.sub("this is still a test", "doesn't matter", "hello there")`)
	assertExprPanics(t, `//seq.sub("hello there", "test", 1)`)
}

func TestArraySub(t *testing.T) {
	t.Parallel()
	// AssertCodesEvalToSameValue(t,
	// 	`"this is a test"`,
	// 	`//seq.sub(["this", "is", "not", "a", "test"], ["is", "not"], "is")`)
}

func TestBytesSub(t *testing.T) {
	t.Parallel()
}

func TestStrSplit(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//seq.split("this is a test", "")`)
	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//seq.split("this is a test", " ") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//seq.split("this is a test", ",") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//seq.split("this is a test", "is")`)
	assertExprPanics(t, `//seq.split("this is a test", 1)`)
}

// TestStrJoin, joiner is string.
func TestStrJoin(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""                `, `//seq.join(",",[])                         `)
	AssertCodesEvalToSameValue(t, `",,"              `, `//seq.join(",",["", "", ""])               `)
	AssertCodesEvalToSameValue(t, `"this is a test"  `, `//seq.join(" ",["this", "is", "a", "test"])`)
	AssertCodesEvalToSameValue(t, `"this"            `, `//seq.join(",",["this"])                   `)
	AssertCodesEvalToSameValue(t, `"You and me"`, `//seq.join(" and ",["You", "me"])`)
	assertExprPanics(t, `//seq.join("this", 2)`)
}

func TestArrayJoin(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `["A","B"]`, `//seq.join([],["A","B"])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `["You", "me"]`, `//seq.join("",["You", "me"])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join("",[1,2])`)
	// AssertCodesEvalToSameValue(t, `["A",",","B"]`, `//seq.join([","],["A","B"])`)
	// AssertCodesEvalToSameValue(t, `[1,1,2]`, `//seq.join([1],[1,2])`)
	// AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// AssertCodesEvalToSameValue(t, `[1]`, `//seq.join([1],[])`)
}

func TestBytesJoin(t *testing.T) {
	t.Parallel()
}

func TestStrPrefix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("ABCDE", "A")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix("ABCDE", "AB")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("ABCDE", "BCD")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix("ABCDE", "CD")`)
}

func TestArrayPrefix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D','E'], 'A')`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D','E'], ['A'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D','E'], ['A','B'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D','E'], ['A','B','C'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(['A','B','C','D','E'], ['A','B','C','D'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E'], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E'], 'B')`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E'], ['B'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E'], ['B','C'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(['A','B','C','D','E'], ['A','B','C','D','E','F'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix([], ['A','B','C','D','E','F'])`)
}

func TestBytesPrefix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('h'))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('he'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('e'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('l'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('o'))`)
}

func TestStrSuffix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("ABCDE", "E")`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix("ABCDE", "DE")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix("ABCDE", "CD")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix("ABCDE", "D")`)
}

func TestArraySuffix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['A','B','C','D','E'], 'E')`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['A','B','C','D','E'], ['E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(['A','B','C','D','E'], ['D','E'])`)

	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E'], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix([], [])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix([], ['D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E'], 'D')`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E'], ['D'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E'], ['C','D'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(['A','B','C','D','E'], ['A','B','C','D','E','F'])`)
}

func TestBytesSuffix(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('o'))`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('lo'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('e'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('ell'))`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.has_suffix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('h'))`)
}
