package rel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/bootstrap"
	"github.com/arr-ai/wbnf/parser"
	parse "github.com/arr-ai/wbnf/parser"
)

func nodeToValue(p bootstrap.Parsers, rule bootstrap.Rule, v interface{}) (out Value) {
	a := astChildren{".nodeToValue": rule}
	defer enterf("nodeToValue(singletons:%v, rule:%v, v:%v)", p.Singletons(), rule, v).exitf("a:%v, out:%v", &a, &out)
	ruleNodeToAttr(p, rule, "", v, a)
	return a.Value("", p.Singletons())
}

func termNodeToAttr(p bootstrap.Parsers, t bootstrap.Term, prefix string, v interface{}, a astChildren) {
	switch t := t.(type) {
	case bootstrap.S:
		sNodeToAttr(p, t, prefix, v, a)
	case bootstrap.RE:
		reNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Rule:
		ruleNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Named:
		namedNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Seq:
		seqNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Oneof:
		oneofNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Delim:
		delimNodeToAttr(p, t, prefix, v, a)
	case bootstrap.Quant:
		quantNodeToAttr(p, t, prefix, v, a)
	default:
		panic(fmt.Errorf("unexpected term type: %v %[1]T", t))
	}
}

var delayerRE = regexp.MustCompile(`#\d+`)

func delayer(path string) string {
	return delayerRE.ReplaceAllLiteralString(path, "")
}

type astChildren map[string]interface{}

func (c astChildren) Add(p bootstrap.Parsers, name string, value Value) {
	defer enterf("%s.astChildren.Add(name:%v, value:%v)", c.ShortString(), name, value).exitf("")
	if p.Singletons().Has(name) {
		if _, has := c[name]; has {
			panic("???")
		}
		c[name] = value
	} else {
		var values []Value
		if v, has := c[name]; has {
			values = v.([]Value)
		}
		c[name] = append(values, value)
	}
}

func (c astChildren) Value(prefix string, singletons bootstrap.PathSet) (out Value) {
	defer enterf("%s.astChildren.Value(prefix:%v)", c.ShortString(), prefix).exitf("a:%v", &out)
	// if len(c) == 1 {
	// 	for name, values := range c {
	// 		if name == "" && len(values) == 1 {
	// 			return values[0]
	// 		}
	// 	}
	// }

	var b TupleBuilder
	for name, v := range c {
		if !strings.HasPrefix(name, ".") {
			switch v := v.(type) {
			case Value:
				b.Put(name, v)
			case []Value:
				b.Put(name, NewArray(v...))
			default:
				panic("???")
			}
		}
	}
	value := b.Finish()
	indentf("c.Value(): %v", value)
	return value
}

func (c astChildren) ShortString() string {
	for name, value := range c {
		if strings.HasPrefix(name, ".") {
			return fmt.Sprintf("{%q:%v, ...}", name, value)
		}
	}
	panic("???")
}

func sNodeToAttr(p bootstrap.Parsers, _ bootstrap.S, prefix string, v interface{}, a astChildren) {
	a.Add(p, prefix, NewString([]rune(v.(parse.Scanner).String())))
}

func reNodeToAttr(p bootstrap.Parsers, _ bootstrap.RE, prefix string, v interface{}, a astChildren) {
	a.Add(p, prefix, NewString([]rune(v.(parse.Scanner).String())))
}

func ruleNodeToAttr(p bootstrap.Parsers, t bootstrap.Rule, prefix string, v interface{}, a astChildren) {
	b := astChildren{".ruleNodeToAttr": prefix}
	defer enterf("ruleNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v)", t, prefix, v, a).exitf("a:%v b:%v", &a, &b)
	n := string(t)
	termNodeToAttr(p, p.Grammar()[t], n+".", v, b)
	if len(b) == 2 {
		for name, v := range b {
			if strings.HasPrefix(name, ".") {
				continue
			}
			var value Value
			switch v := v.(type) {
			case Value:
				value = v
			case []Value:
				if len(v) == 1 {
					value = v[0]
				}
			}
			if value != nil {
				base := n
				if i := strings.Index(n, "#"); i != -1 {
					base = n[:i+1]
				} else {
					if dot := strings.Index(base, "."); dot != -1 {
						base = base[:dot] + "#"
					} else {
						base += "#"
					}
				}
				if strings.HasPrefix(name, base) {
					indentf("singleton layer: %v", t)
					a.Add(p, name, value)
					return
				}
			}
		}
	}
	indentf("ruleNodeToAttr(): b = %v", b)
	value := b.Value(n+".", p.Singletons())
	a.Add(p, n, value)
}

func namedNodeToAttr(p bootstrap.Parsers, t bootstrap.Named, prefix string, v interface{}, a astChildren) {
	singleton := p.Singletons().Has(prefix + t.Name)
	defer enterf("namedNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v) 1=%v", t, prefix, v, a, singleton).exitf("a:%v", &a)
	if singleton {
		termNodeToAttr(p, t.Term, t.Name, v, a)
	} else {
		b := astChildren{".namedNodeToAttr": prefix}
		defer indentpf("b:%v", &b)
		termNodeToAttr(p, t.Term, prefix+t.Name+".", v, b)
		a.Add(p, prefix+t.Name, b.Value(prefix, p.Singletons()))
	}
}

func seqNodeToAttr(p bootstrap.Parsers, t bootstrap.Seq, prefix string, v interface{}, a astChildren) {
	defer enterf("seqNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v)", t, prefix, v, a).exitf("a:%v", &a)
	for i, child := range v.(parser.Node).Children {
		termNodeToAttr(p, t[i], prefix, child, a)
	}
}

func oneofNodeToAttr(p bootstrap.Parsers, t bootstrap.Oneof, prefix string, v interface{}, a astChildren) {
	defer enterf("oneofNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v)", t, prefix, v, a).exitf("a:%v", &a)
	node := v.(parser.Node)
	termNodeToAttr(p, t[node.Extra.(int)], prefix, node.Children[0], a)
}

func delimNodeToAttr(p bootstrap.Parsers, t bootstrap.Delim, prefix string, v interface{}, a astChildren) {
	defer enterf("delimNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v)", t, prefix, v, a).exitf("a:%v", &a)
	node := v.(parser.Node)
	left, right := t.LRTerms(node)
	terms := [2]bootstrap.Term{left, t.Sep}
	for i, child := range node.Children {
		termNodeToAttr(p, terms[i%2], prefix, child, a)
		terms[0] = right // Only use left once.
	}
}

func quantNodeToAttr(p bootstrap.Parsers, t bootstrap.Quant, prefix string, v interface{}, a astChildren) {
	defer enterf("quantNodeToAttr(t:%v, prefix:%v, v:%#v, a:%v)", t, prefix, v, a).exitf("a:%v", &a)
	for _, child := range v.(parser.Node).Children {
		termNodeToAttr(p, t.Term, prefix, child, a)
	}
}
