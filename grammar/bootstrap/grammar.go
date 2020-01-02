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

func captureForDebugging(interface{}) {}

func nameTag(name Rule, term Term, parser parse.Parser) parse.Parser {
	descr := nameOr(name, term.String())
	return parse.Func(func(input *parse.Scanner, output interface{}) bool {
		captureForDebugging(term)
		var v interface{}
		if parser.Parse(input, &v) {
			a, ok := v.([]interface{})
			if !ok {
				a = []interface{}{v}
			}
			parse.Put(append([]interface{}{descr}, a...), output)
			return true
		}
		return false
	})
}

func (g Grammar) Compile() func(rule Rule) parse.Parser {
	c := cache{parsers: map[Rule]parse.Parser{}, grammar: g}
	for rule, term := range g {
		c.parsers[rule] = term.Parser(rule, c)
	}
	return func(rule Rule) parse.Parser {
		return c.parsers[rule]
	}
}

func (g Rule) Parser(name Rule, c cache) parse.Parser {
	var parser parse.Parser
	return parse.Func(func(input *parse.Scanner, output interface{}) bool {
		captureForDebugging(g)
		if parser == nil {
			var ok bool
			if parser, ok = c.parsers[g]; !ok {
				panic("missing parser: " + g)
			}
		}
		return parser.Parse(input, output)
	})
}

func (g RE) Parser(name Rule, c cache) parse.Parser {
	s := string(g)
	if wrap, has := c.grammar[WrapRE]; has {
		s = strings.Replace(string(wrap.(RE)), "()", "(?:"+s+")", 1)
	}
	parser := parse.Regexp(s)
	return nameTag(name, g, parser)
}

func (g Seq) Parser(name Rule, c cache) parse.Parser {
	parsers := c.MakeParsers(g)
	return nameTag(name, g, parse.Func(func(input *parse.Scanner, output interface{}) bool {
		result := make([]interface{}, 0, len(parsers))
		for _, parser := range parsers {
			var v interface{}
			if !parser.Parse(input, &v) {
				return false
			}
			result = append(result, v)
		}
		parse.Put(result, output)
		return true
	}))
}

func (g Delim) Parser(name Rule, c cache) parse.Parser {
	term := g.Term.Parser("", c)
	sep := Seq{g.Sep, g.Term}.Parser("", c)
	return nameTag(name, g, parse.Func(func(input *parse.Scanner, output interface{}) bool {
		var v interface{}
		if !term.Parse(input, &v) {
			return false
		}
		result := []interface{}{v}
		for sep.Parse(input, &v) {
			result = append(result, v.([]interface{})[1:]...)
		}
		parse.Put(result, output)
		return true
	}))
}

func (g Quant) Parser(name Rule, c cache) parse.Parser {
	term := g.Term.Parser("", c)
	return nameTag(name, g, parse.Func(func(input *parse.Scanner, output interface{}) bool {
		result := make([]interface{}, 0, g.Min)
		var v interface{}
		i := 0
		for ; (g.Max == 0 || i < g.Max) && term.Parse(input, &v); i++ {
			result = append(result, v)
		}
		if i < g.Min {
			return false
		}
		parse.Put(result, output)
		return true
	}))
}

func (g Choice) Parser(name Rule, c cache) parse.Parser {
	parsers := c.MakeParsers(g)
	return nameTag(name, g, parse.Func(func(input *parse.Scanner, output interface{}) bool {
		for _, parser := range parsers {
			var v interface{}
			if parser.Parse(input, &v) {
				parse.Put([]interface{}{v}, output)
				return true
			}
		}
		return false
	}))
}
