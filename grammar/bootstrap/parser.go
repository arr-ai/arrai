package bootstrap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

const (
	towerDelim = "#"

	seqTag   = "_"
	oneofTag = "|"
	delimTag = ":"
	quantTag = "?"
)

type cache struct {
	parsers    map[Rule]parse.Parser
	grammar    Grammar
	rulePtrses map[Rule][]*parse.Parser
}

func (c cache) registerRule(parser *parse.Parser) {
	if rule, ok := (*parser).(ruleParser); ok {
		c.rulePtrses[rule.t] = append(c.rulePtrses[rule.t], parser)
	}
}

func (c cache) registerRules(parsers []parse.Parser) {
	for i := range parsers {
		c.registerRule(&parsers[i])
	}
}

func (c cache) makeParsers(terms []Term) []parse.Parser {
	parsers := make([]parse.Parser, 0, len(terms))
	for _, t := range terms {
		parsers = append(parsers, t.Parser("", c))
	}
	c.registerRules(parsers)
	return parsers
}

type putter func(output interface{}, extra interface{}, children ...interface{}) bool

func tag(rule Rule, alt Rule) putter {
	if rule == "" {
		rule = alt
	}

	return func(output interface{}, extra interface{}, children ...interface{}) bool {
		parse.PtrAssign(output, parse.Node{
			Tag:      string(rule),
			Extra:    extra,
			Children: children,
		})
		return true
	}
}

func (g Grammar) clone() Grammar {
	clone := make(Grammar, len(g))
	for rule, term := range g {
		clone[rule] = term
	}
	return clone
}

func (g Grammar) resolveTowers() {
	for rule, term := range g {
		if tower, ok := term.(Tower); ok {
			oldRule := rule
			for i, layer := range tower {
				newRule := rule
				if j := (i + 1) % len(tower); j > 0 {
					newRule = Rule(fmt.Sprintf("%s%s%d", rule, towerDelim, j))
				}
				g[oldRule] = layer.Resolve(rule, newRule)
				oldRule = newRule
			}
		}
	}
}

func (g Grammar) Compile() Parsers {
	for _, term := range g {
		if _, ok := term.(Tower); ok {
			g = g.clone()
			g.resolveTowers()
			break
		}
	}

	c := cache{
		parsers:    map[Rule]parse.Parser{},
		grammar:    g,
		rulePtrses: map[Rule][]*parse.Parser{},
	}
	for rule, term := range g {
		c.parsers[rule] = term.Parser(rule, c)
	}

	for rule, rulePtrs := range c.rulePtrses {
		term := c.parsers[rule]
		for _, rulePtr := range rulePtrs {
			*rulePtr = term
		}
	}

	return c.parsers
}

//-----------------------------------------------------------------------------

type ruleParser struct {
	rule Rule
	t    Rule
}

func (p ruleParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	panic(Inconceivable)
}

func (t Rule) Parser(rule Rule, c cache) parse.Parser {
	return ruleParser{
		rule: rule,
		t:    t,
	}
}

//-----------------------------------------------------------------------------

type sParser struct {
	rule Rule
	t    S
	re   *regexp.Regexp
}

func (p *sParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	var eaten parse.Scanner
	if input.EatRegexp(p.re, &eaten) {
		parse.PtrAssign(output, eaten)
		return true
	}
	return false
}

func (t S) Parser(rule Rule, c cache) parse.Parser {
	re := "(" + regexp.QuoteMeta(string(t)) + ")"
	if wrap, has := c.grammar[WrapRE]; has {
		re = strings.Replace(string(wrap.(RE)), "()", "(?:"+re+")", 1)
	}
	return &sParser{
		rule: rule,
		t:    t,
		re:   regexp.MustCompile(`(?m)\A(?:` + re + `)`),
	}
}

//-----------------------------------------------------------------------------

type reParser struct {
	rule Rule
	t    RE
	re   *regexp.Regexp
}

func (p *reParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	var eaten parse.Scanner
	if input.EatRegexp(p.re, &eaten) {
		parse.PtrAssign(output, eaten)
		return true
	}
	return false
}

func (t RE) Parser(rule Rule, c cache) parse.Parser {
	re := string(t)
	if wrap, has := c.grammar[WrapRE]; has {
		re = strings.Replace(string(wrap.(RE)), "()", "(?:"+re+")", 1)
	}
	return &reParser{
		rule: rule,
		t:    t,
		re:   regexp.MustCompile(`(?m)\A(?:` + re + `)`),
	}
}

//-----------------------------------------------------------------------------

type seqParser struct {
	rule    Rule
	t       Seq
	parsers []parse.Parser
	put     putter
}

func (p *seqParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	result := make([]interface{}, 0, len(p.parsers))
	for _, parser := range p.parsers {
		var n interface{}
		if !parser.Parse(input, &n) {
			return false
		}
		result = append(result, n)
	}
	return p.put(output, nil, result...)
}

func (t Seq) Parser(rule Rule, c cache) parse.Parser {
	return &seqParser{
		rule:    rule,
		t:       t,
		parsers: c.makeParsers(t),
		put:     tag(rule, seqTag),
	}
}

//-----------------------------------------------------------------------------

type delimParser struct {
	rule Rule
	t    Delim
	term parse.Parser
	sep  parse.Parser
	put  putter
}

func (p *delimParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	var n interface{}
	if !p.term.Parse(input, &n) {
		return false
	}
	result := []interface{}{n}
	var d interface{}
	for p.sep.Parse(input, &d) && p.term.Parse(input, &n) {
		result = append(result, d, n)
	}
	return p.put(output, nil, result...)
}

func (t Delim) Parser(rule Rule, c cache) parse.Parser {
	p := &delimParser{
		rule: rule,
		t:    t,
		term: t.Term.Parser("", c),
		sep:  t.Sep.Parser("", c),
		put:  tag(rule, delimTag),
	}
	c.registerRule(&p.term)
	c.registerRule(&p.sep)
	return p
}

//-----------------------------------------------------------------------------

type quantParser struct {
	rule Rule
	t    Quant
	term parse.Parser
	put  putter
}

func (p *quantParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	result := make([]interface{}, 0, p.t.Min)
	var n interface{}
	for i := 0; (p.t.Max == 0 || i < p.t.Max) && p.term.Parse(input, &n); i++ {
		result = append(result, n)
	}
	if len(result) >= p.t.Min {
		return p.put(output, nil, result...)
	}
	return false
}

func (t Quant) Parser(rule Rule, c cache) parse.Parser {
	p := &quantParser{
		rule: rule,
		t:    t,
		term: t.Term.Parser("", c),
		put:  tag(rule, quantTag),
	}
	c.registerRule(&p.term)
	return p
}

//-----------------------------------------------------------------------------

type oneofParser struct {
	rule    Rule
	t       Oneof
	parsers []parse.Parser
	put     putter
}

func (p *oneofParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	for i, parser := range p.parsers {
		var n interface{}
		start := *input
		if parser.Parse(&start, &n) {
			*input = start
			return p.put(output, i, n)
		}
	}
	return false
}

func (t Oneof) Parser(rule Rule, c cache) parse.Parser {
	return &oneofParser{
		rule:    rule,
		t:       t,
		parsers: c.makeParsers(t),
		put:     tag(rule, oneofTag),
	}
}

//-----------------------------------------------------------------------------

func (t Tower) Parser(_ Rule, _ cache) parse.Parser {
	panic(Inconceivable)
}

//-----------------------------------------------------------------------------

func (t NamedTerm) Parser(rule Rule, c cache) parse.Parser {
	return t.Term.Parser(Rule(t.Name), c)
}
