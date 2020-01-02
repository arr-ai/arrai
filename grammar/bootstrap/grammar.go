package bootstrap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

var (
	comment = Rule("comment")
	prod    = Rule("prod")
	stmt    = Rule("stmt")
	expr    = Rule("expr")
	choice  = Rule("choice")
	seq     = Rule("seq")
	tag     = Rule("tag")
	term    = Rule("term")
	atom    = Rule("atom")
	quant   = Rule("quant")
	ident   = Rule("ident")
	str     = Rule("str")
	i       = Rule("i")
	re      = Rule("re")

	// WrapRE is a special rule to indicate a wrapper around all regexps and
	// strings. When supplied in the form "pre()post", then all regexes will be
	// wrapped in "pre(?:" and ")post" and all strings will be escaped using
	// regexp.QuoteMeta then likewise wrapped.
	WrapRE = Rule(".wrapRE")
)

var GrammarGrammar = Grammar{
	"grammar": Some(stmt),
	stmt:      Choice{comment, prod},
	comment:   RE(`(//.*)$`),
	prod:      Seq{ident, S("->"), Some(expr), S(";")},
	expr:      choice,
	choice:    Delim{Term: seq, Sep: S("|")},
	seq:       Some(tag),
	tag:       Seq{Opt(Seq{S("<"), ident, S(">")}), term},
	term:      Seq{atom, Opt(quant)},
	atom:      Choice{ident, str, re, Seq{S("("), expr, S(")")}},
	quant: Choice{
		RE(`([?*+])`),
		Seq{S("{"), Opt(i), S(","), Opt(i), S("}")},
		Seq{S(":"), atom},
	},

	ident:  RE(`([A-Za-z_\.]\w*)`),
	str:    RE(`"((?:[^"\\]|\\.)*)"`),
	i:      RE(`(\d+)`),
	re:     RE(`/((?:[^/\\]|\\.)*)/`),
	WrapRE: RE(`\s*()\s*`),
}

type Grammar map[Rule]Term

type Term interface {
	fmt.Stringer
	IsTerm()
	Parser(name Rule, c cache) parse.Parser
}

func (Rule) IsTerm()   {}
func (RE) IsTerm()     {}
func (Seq) IsTerm()    {}
func (Delim) IsTerm()  {}
func (Quant) IsTerm()  {}
func (Choice) IsTerm() {}

type Rule string

type RE string

type Seq []Term

type Delim struct {
	Term Term
	Sep  Term
}

type Quant struct {
	Term Term
	Min  int
	Max  int // 0 = infinity
}

func S(s string) RE { return RE("(" + regexp.QuoteMeta(s) + ")") }

func Opt(term Term) Quant  { return Quant{Term: term, Max: 1} }
func Any(term Term) Quant  { return Quant{Term: term} }
func Some(term Term) Quant { return Quant{Term: term, Min: 1} }

type Choice []Term

func join(terms []Term, sep string) string {
	s := []string{}
	for _, t := range terms {
		s = append(s, t.String())
	}
	return strings.Join(s, sep)
}

func (g Rule) String() string   { return string(g) }
func (g RE) String() string     { return fmt.Sprintf("/%v/", string(g)) }
func (g Seq) String() string    { return join(g, " ") }
func (g Delim) String() string  { return fmt.Sprintf("%v:%v", g.Term, g.Sep) }
func (g Quant) String() string  { return fmt.Sprintf("%v{%d,%d}", g.Term, g.Min, g.Max) }
func (g Choice) String() string { return join(g, " | ") }

type cache struct {
	parsers map[Rule]parse.Parser
	grammar Grammar
}

func (c cache) MakeParsers(terms []Term) []parse.Parser {
	parsers := make([]parse.Parser, 0, len(terms))
	for _, t := range terms {
		parsers = append(parsers, t.Parser("", c))
	}
	return parsers
}

func nameOr(name Rule, descr string) string {
	if name != "" {
		return string(name)
	}
	return descr
}

func nameTag(name Rule, term Term) func(v ...interface{}) []interface{} {
	descr := nameOr(name, term.String())
	return func(v ...interface{}) []interface{} { return append([]interface{}{descr}, v...) }
}

func (g Grammar) Parsers() func(rule Rule) parse.Parser {
	c := cache{parsers: map[Rule]parse.Parser{}, grammar: g}
	for rule, term := range g {
		c.parsers[rule] = term.Parser(rule, c)
	}
	return func(rule Rule) parse.Parser {
		return c.parsers[rule]
	}
}

func captureForDebugging(interface{}) {}

func (g Rule) Parser(name Rule, c cache) parse.Parser {
	var parser parse.Parser
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		if parser == nil {
			var ok bool
			if parser, ok = c.parsers[g]; !ok {
				panic("missing parser: " + g)
			}
		}
		return parser.Parse(input)
	})
}

func (g RE) Parser(name Rule, c cache) parse.Parser {
	tag := nameTag(name, g)
	s := string(g)
	if wrap, has := c.grammar[WrapRE]; has {
		s = strings.Replace(string(wrap.(RE)), "()", "(?:"+s+")", 1)
	}
	parser := parse.Regexp(s)
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		if v, ok := parser.Parse(input); ok {
			return tag(v), ok
		}
		return nil, false
	})
}

func (g Seq) Parser(name Rule, c cache) parse.Parser {
	tag := nameTag(name, g)
	parsers := c.MakeParsers(g)
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		result := make([]interface{}, 0, len(parsers))
		for _, parser := range parsers {
			v, ok := parser.Parse(input)
			if !ok {
				return nil, false
			}
			result = append(result, v)
		}
		return tag(result...), true
	})
}

func (g Delim) Parser(name Rule, c cache) parse.Parser {
	tag := nameTag(name, g)
	term := g.Term.Parser("", c)
	sep := Seq{g.Sep, g.Term}.Parser("", c)
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		if v, ok := term.Parse(input); ok {
			result := []interface{}{v}
			for {
				v, ok := sep.Parse(input)
				if !ok {
					break
				}
				result = append(result, v.([]interface{})[1:]...)
			}
			return tag(result...), true
		}
		return nil, false
	})
}

func (g Quant) Parser(name Rule, c cache) parse.Parser {
	tag := nameTag(name, g)
	term := g.Term.Parser("", c)
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		result := make([]interface{}, 0, g.Min)
		i := 0
		max := g.Max
		if max == 0 {
			max = int(uint(^max) >> 1)
		}
		for ; i < max; i++ {
			v, ok := term.Parse(input)
			if !ok {
				break
			}
			result = append(result, v)
		}
		if i >= g.Min {
			return tag(result...), true
		}
		return nil, false
	})
}

func (g Choice) Parser(name Rule, c cache) parse.Parser {
	tag := nameTag(name, g)
	parsers := c.MakeParsers(g)
	return parse.Func(func(input *parse.Scanner) (interface{}, bool) {
		captureForDebugging(g)
		for _, parser := range parsers {
			if v, ok := parser.Parse(input); ok {
				return tag(v), ok
			}
		}
		return nil, false
	})
}
