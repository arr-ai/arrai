package bootstrap

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

var (
	stmt    = Rule("stmt")
	comment = Rule("comment")
	prod    = Rule("prod")
	expr    = Rule("expr")
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
	stmt:      Oneof{comment, prod},
	comment:   RE(`(//.*)$`),
	prod:      Seq{ident, S("->"), Some(expr), S(";")},
	expr: Tower{
		Delim{Term: expr, Sep: S("^")},
		Delim{Term: expr, Sep: S("|")},
		Some(expr),
		Oneof{expr, Seq{Opt(Seq{S("<"), ident, S(">")}), expr}},
		Oneof{Seq{atom, Opt(quant)}, atom},
	},
	atom: Oneof{ident, str, re, Seq{S("("), expr, S(")")}},
	quant: Oneof{
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
type Parsers map[Rule]parse.Parser

type Term interface {
	fmt.Stringer
	Parser(name Rule, c cache) parse.Parser
	ValidateParse(v interface{}) error
	Unparse(v interface{}, w io.Writer) (n int, err error)
	Resolve(oldRule, newRule Rule) Term
}

type (
	Rule  string
	S     string
	RE    string
	Seq   []Term
	Oneof []Term
	Tower []Term
	Delim struct {
		Term Term
		Sep  Term
	}
	Quant struct {
		Term Term
		Min  int
		Max  int // 0 = infinity
	}
	NamedTerm struct {
		Name string
		Term Term
	}
)

func Opt(term Term) *Quant  { return &Quant{Term: term, Max: 1} }
func Any(term Term) *Quant  { return &Quant{Term: term} }
func Some(term Term) *Quant { return &Quant{Term: term, Min: 1} }

func Name(name string, term Term) NamedTerm {
	return NamedTerm{Name: name, Term: term}
}

func join(terms []Term, sep string) string {
	s := []string{}
	for _, t := range terms {
		s = append(s, t.String())
	}
	return strings.Join(s, sep)
}

func (g Grammar) String() string {
	keys := make([]string, 0, len(g))
	for key := range g {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	var sb strings.Builder
	count := 0
	for _, key := range keys {
		if count > 0 {
			sb.WriteString("; ")
		}
		fmt.Fprintf(&sb, "%s -> %v", key, g[Rule(key)])
		count++
	}
	return sb.String()
}

func (t Rule) String() string      { return string(t) }
func (t S) String() string         { return fmt.Sprintf("%q", string(t)) }
func (t RE) String() string        { return fmt.Sprintf("/%v/", string(t)) }
func (t Seq) String() string       { return join(t, " ") }
func (t Oneof) String() string     { return join(t, " | ") }
func (t Tower) String() string     { return join(t, " >> ") }
func (t Delim) String() string     { return fmt.Sprintf("%v:%v", t.Term, t.Sep) }
func (t Quant) String() string     { return fmt.Sprintf("%v{%d,%d}", t.Term, t.Min, t.Max) }
func (t NamedTerm) String() string { return fmt.Sprintf("<%s>%v", t.Name, t.Term) }
