package bootstrap

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

const towerDelim = "•"

var towerRE = regexp.MustCompile(fmt.Sprintf(`%s\d+%[1]s`, towerDelim))

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

type putter func(output interface{}, values ...interface{}) bool

func nameTag(name Rule) putter {
	towered := strings.HasSuffix(string(name), towerDelim)
	if towered {
		name = Rule(towerRE.ReplaceAllLiteralString(string(name), ""))
	}

	var head []interface{}
	if name != "" {
		// Create rather than append so cap(head) == 1.
		head = []interface{}{name}
	}

	return func(output interface{}, values ...interface{}) bool {
		if len(values) == 1 && (len(head) == 0 || towered) {
			parse.Put(values[0], output)
			return true
		}
		parse.Put(append(head, values...), output)
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
	log.Print(g)
	for rule, term := range g {
		if tower, ok := term.(Tower); ok {
			oldRule := rule
			for i, layer := range tower {
				newRule := rule
				if j := (i + 1) % len(tower); j > 0 {
					newRule = Rule(fmt.Sprintf("%s%s%d%[2]s", rule, towerDelim, j))
				}
				g[oldRule] = layer.Resolve(rule, newRule)
				oldRule = newRule
			}
		}
	}
	log.Print(g)
}

func (g Grammar) Compile() map[Rule]parse.Parser {
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
	t Rule
}

func (p ruleParser) Parse(input *parse.Scanner, output interface{}) bool {
	panic("should never get here")
}

func (t Rule) Parser(name Rule, c cache) parse.Parser {
	return ruleParser{t: t}
}

//-----------------------------------------------------------------------------

type reParser struct {
	t      RE
	parser parse.Parser
	put    putter
}

func (p *reParser) Parse(input *parse.Scanner, output interface{}) bool {
	var v interface{}
	if p.parser.Parse(input, &v) {
		return p.put(output, v)
	}
	return false
}

func (t RE) Parser(name Rule, c cache) parse.Parser {
	s := string(t)
	if wrap, has := c.grammar[WrapRE]; has {
		s = strings.Replace(string(wrap.(RE)), "()", "(?:"+s+")", 1)
	}
	return &reParser{t: t, parser: parse.Regexp(s), put: nameTag(name)}
}

//-----------------------------------------------------------------------------

type seqParser struct {
	t       Seq
	parsers []parse.Parser
	put     putter
}

func (p *seqParser) Parse(input *parse.Scanner, output interface{}) bool {
	result := make([]interface{}, 0, len(p.parsers))
	for _, parser := range p.parsers {
		var v interface{}
		if !parser.Parse(input, &v) {
			return false
		}
		result = append(result, v)
	}
	return p.put(output, result...)
}

func (t Seq) Parser(name Rule, c cache) parse.Parser {
	return &seqParser{
		t:       t,
		parsers: c.makeParsers(t),
		put:     nameTag(name),
	}
}

//-----------------------------------------------------------------------------

type delimParser struct {
	t    Delim
	term parse.Parser
	sep  parse.Parser
	put  putter
}

func (p *delimParser) Parse(input *parse.Scanner, output interface{}) bool {
	var v interface{}
	if !p.term.Parse(input, &v) {
		return false
	}
	result := []interface{}{v}
	var d interface{}
	for p.sep.Parse(input, &d) && p.term.Parse(input, &v) {
		result = append(result, d, v)
	}
	return p.put(output, result...)
}

func (t Delim) Parser(name Rule, c cache) parse.Parser {
	p := &delimParser{
		t:    t,
		term: t.Term.Parser("", c),
		sep:  t.Sep.Parser("", c),
		put:  nameTag(name),
	}
	c.registerRule(&p.term)
	c.registerRule(&p.sep)
	return p
}

//-----------------------------------------------------------------------------

type quantParser struct {
	t    Quant
	term parse.Parser
	put  putter
}

func (p *quantParser) Parse(input *parse.Scanner, output interface{}) bool {
	result := make([]interface{}, 0, p.t.Min)
	var v interface{}

	for i := 0; (p.t.Max == 0 || i < p.t.Max) && p.term.Parse(input, &v); i++ {
		result = append(result, v)
	}
	if len(result) >= p.t.Min {
		return p.put(output, result...)
	}
	return false
}

func (t Quant) Parser(name Rule, c cache) parse.Parser {
	p := &quantParser{
		t:    t,
		term: t.Term.Parser("", c),
		put:  nameTag(name),
	}
	c.registerRule(&p.term)
	return p
}

//-----------------------------------------------------------------------------

type choiceParser struct {
	t       Oneof
	parsers []parse.Parser
	put     putter
}

func (p *choiceParser) Parse(input *parse.Scanner, output interface{}) bool {
	for _, parser := range p.parsers {
		var v interface{}
		if parser.Parse(input, &v) {
			return p.put(output, v)
		}
	}
	return false
}

func (t Oneof) Parser(name Rule, c cache) parse.Parser {
	return &choiceParser{
		t:       t,
		parsers: c.makeParsers(t),
		put:     nameTag(name),
	}
}

func (t Tower) Parser(_ Rule, _ cache) parse.Parser {
	panic("should never get here")
}

//-----------------------------------------------------------------------------

func (t NamedTerm) Parser(name Rule, c cache) parse.Parser {
	return t.Term.Parser(name, c)
}
