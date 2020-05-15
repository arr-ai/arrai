package syntax

import (
	"testing"
)

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
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("", "this is a test")             `)
	AssertCodesEvalToSameValue(t, `true `, `//seq.contains("is a test", "this is a test")    `)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("is not a test", "this is a test")`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains("a is", "this is a test")`)
	assertExprPanics(t, `//seq.contains(124, 123)`)
}

func TestArrayContains(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(1, [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(3, [1,2,3,4,5])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(5, [1,2,3,4,5])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('A',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('E',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains('C',['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B','C'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['C','D','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C','D','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['B','C','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','E'],['A','B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `false`, `//seq.contains(['A','B','C','D','E','F'],['A','B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['A','B','C'], ['A', 'A', 'B','C','D','E'])`)
	AssertCodesEvalToSameValue(t, `true`, `//seq.contains(['B','C'],['A', 'A', 'B','C','D','E'])`)

	AssertCodesEvalToSameValue(t, `true`, `//seq.contains([['B','C']],[['A', 'B'], ['B','C'],['D','E']])`)
}

func TestBytesContains(t *testing.T) {
	t.Parallel()
	// hello bytes - 104 101 108 108 111
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 104)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 111)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains({ |@, @byte| (0, 108),(0, 108)},//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `true`,
		`//seq.contains(//unicode.utf8.encode('h'),//unicode.utf8.encode('hello'))`)
}

///////////////////
func TestStrSub(t *testing.T) {
	t.Parallel()
	// string
	AssertCodesEvalToSameValue(t,
		`"this is not a test"`,
		`//seq.sub("aaa", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("is not", "is", "this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is a test"`,
		`//seq.sub("not ", "","this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"t1his is not1 a t1est1"`,
		`//seq.sub("t", "t1","this is not a test")`)
	AssertCodesEvalToSameValue(t,
		`"this is still a test"`,
		`//seq.sub( "doesn't matter", "hello there","this is still a test")`)
	assertExprPanics(t, `//seq.sub("hello there", "test", 1)`)
}

func TestArraySub(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`['T', 'B', 'T', 'C', 'D', 'E']`,
		`//seq.sub('A', 'T', ['A', 'B', 'A', 'C', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t,
		`['T', 'B', 'T', 'C', 'D', 'E']`,
		`//seq.sub(['A'], ['T'], ['A', 'B', 'A', 'C', 'D', 'E'])`)
	AssertCodesEvalToSameValue(t,
		`[['A', 'B'], ['T','C'],['A','D']]`,
		`//seq.sub([['A','C']], [['T','C']], [['A', 'B'], ['A','C'],['A','D']])`)
	AssertCodesEvalToSameValue(t,
		`[2, 2, 3]`,
		`//seq.sub(1, 2, [1, 2, 3])`)
	AssertCodesEvalToSameValue(t,
		`[2, 2, 3]`,
		`//seq.sub([1], [2], [1, 2, 3])`)
	AssertCodesEvalToSameValue(t,
		`[[1,1], [4,4], [3,3]]`,
		`//seq.sub([[2,2]], [[4,4]], [[1,1], [2,2], [3,3]])`)
}

func TestBytesSub(t *testing.T) {
	t.Parallel()
	///
}

func TestStrSplit(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//seq.split('',"this is a test")`)
	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//seq.split(" ","this is a test") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//seq.split(",","this is a test") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//seq.split("is","this is a test")`)
	assertExprPanics(t, `//seq.split(1, "this is a test")`)
}

func TestArraySplit(t *testing.T) {
	t.Parallel()
	// TODO
	AssertCodesEvalToSameValue(t,
		`['A', 'B', 'A', 'C', 'D', 'E']`,
		`//seq.sub(['A', 'B', 'A', 'C', 'D', 'E'], 'T')`)
	// AssertCodesEvalToSameValue(t,
	// 	`[['B'],['C', 'D', 'E']]`,
	// 	`//seq.sub(['A', 'B', 'A', 'C', 'D', 'E'], 'A')`)
}

func TestBytesSplit(t *testing.T) {
	t.Parallel()
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
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `["You", "me"]`, `//seq.join("",["You", "me"])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join("",[1,2])`)

	AssertCodesEvalToSameValue(t, `["A","B"]`, `//seq.join([],["A","B"])`)
	AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// if joinee is empty, the final value will be empty
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	AssertCodesEvalToSameValue(t, `["A",",","B"]`, `//seq.join([","],["A","B"])`)
	AssertCodesEvalToSameValue(t, `[1,0,2,0,3,0,4,0,5]`, `//seq.join([0], [1,2,3,4,5])`)
	// TODO
	//AssertCodesEvalToSameValue(t, `[1, 2, 0, 3, 4, 0, 5, 6]`, `//seq.join([0], [[1, 2], [3, 4], [5, 6]])`)
	AssertCodesEvalToSameValue(t, `['A','A','B','A','C','A','D']`, `//seq.join(['A'], ['A','B','C','D'])`)
}

func TestBytesJoin(t *testing.T) {
	t.Parallel()
	// joiner "" is translated to rel.GenericSet
	AssertCodesEvalToSameValue(t, `{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) }`,
		`//seq.join("",//unicode.utf8.encode('hello'))`)
	AssertCodesEvalToSameValue(t, `{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) }`,
		`//seq.join([],{ |@, @byte| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111) })`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	// AssertCodesEvalToSameValue(t, `["A","B"]`, `//seq.join([],["A","B"])`)
	// AssertCodesEvalToSameValue(t, `[1,2]`, `//seq.join([],[1,2])`)
	// // if joinee is empty, the final value will be empty
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join([1],[])`)
	// AssertCodesEvalToSameValue(t, `[]`, `//seq.join(['A'],[])`)

	// AssertCodesEvalToSameValue(t, `true`, `//seq.has_prefix(//unicode.utf8.encode('hello'),//unicode.utf8.encode('h'))`)
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
