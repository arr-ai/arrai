package syntax

import "testing"

// Compiled grammars are memoised by grammar value; parsing must be
// unaffected, and distinct grammars must not share an entry.
func TestGrammarParseMemoised(t *testing.T) {
	t.Parallel()

	// The same grammar value parsed repeatedly (served from the cache after
	// the first compile) gives identical results.
	AssertCodesEvalToSameValue(t, `true`, `
		let g = {://grammar.lang.wbnf:
			a -> "x"+;
		:};
		let parse = \s //grammar.parse(g, 'a', s);
		parse('xxx') = parse('xxx') && parse('xx') != parse('xxx')`)
	// Distinct grammars never share a cache entry.
	AssertCodesEvalToSameValue(t, `true`, `
		let g1 = {://grammar.lang.wbnf: a -> "x"+; :};
		let g2 = {://grammar.lang.wbnf: a -> "y"+; :};
		let _ = //grammar.parse(g1, 'a', 'xx');
		//grammar.parse(g2, 'a', 'yy') != //grammar.parse(g1, 'a', 'xx')`)
}

// Compiled regexps are memoised by pattern.
func TestReCompileMemoised(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[[['ab']], [['ab']], []]`, `
		let f = \s //re.compile('a.').match(s);
		[f('ab'), f('ab'), f('x')]`)
	AssertCodesEvalToSameValue(t, `[[['1']], [['a']]]`,
		"[//re.compile(`\\d`).match('1'), //re.compile('[a-z]').match('a')]")
	AssertCodeErrors(t, ``, `//re.compile('(')`)
	AssertCodeErrors(t, ``, `let _ = //re.compile('a'); //re.compile('(')`)
}
