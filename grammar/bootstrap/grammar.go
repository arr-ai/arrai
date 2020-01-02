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
