package bootstrap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

func parseString(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			i++
			switch s[i] {
			case 'x':
				n, err := strconv.ParseInt(s[i:i+2], 16, 8)
				if err != nil {
					panic(err)
				}
				sb.WriteByte(uint8(n))
				i++
			case 'u':
				n, err := strconv.ParseInt(s[i:i+4], 16, 16)
				if err != nil {
					panic(err)
				}
				sb.WriteByte(uint8(n))
				i += 2
			case 'U':
				n, err := strconv.ParseInt(s[i:i+8], 16, 32)
				if err != nil {
					panic(err)
				}
				sb.WriteByte(uint8(n))
				i += 4
			case '0', '1', '2', '3', '4', '5', '6', '7':
				n, err := strconv.ParseInt(s[i:i+3], 8, 8)
				if err != nil {
					panic(err)
				}
				sb.WriteByte(uint8(n))
				i++
			case 'a':
				sb.WriteByte('\a')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'v':
				sb.WriteByte('\v')
			case '\\':
				sb.WriteByte('\\')
			case '\'':
				sb.WriteByte('\'')
			case '"':
				sb.WriteByte('"')
			default:
				panic(fmt.Errorf("unrecognized \\-escape: %q", s[i]))
			}
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func compileAtomNode(node parse.Node) Term {
	switch node.Extra.(int) {
	case 0:
		return Rule(node.GetString(0))
	case 1:
		return S(parseString(node.GetString(0)))
	case 2:
		return RE(strings.ReplaceAll(node.GetString(0), `\/`, `/`))
	case 3:
		return compileTermNode(node.GetNode(0, 1))
	case 4:
		return Seq{}
	default:
		panic(BadInput)
	}
}

func compileTermNamedNode(node parse.Node) Term {
	term := compileAtomNode(node.GetNode(1))
	if quant := node.GetNode(0); quant.Count() == 1 {
		return Named{
			Name: quant.GetString(0, 0),
			Term: term,
		}
	}
	return term
}

func compileTermQuantNode(node parse.Node) Term {
	term := compileTermNamedNode(node.GetNode(0))
	opt := node.GetNode(1)
	if opt.Count() == 1 {
		quant := opt.GetNode(0)
		switch quant.Extra.(int) {
		case 0:
			switch quant.GetString(0) {
			case "?":
				term = Opt(term)
			case "*":
				term = Any(term)
			case "+":
				term = Some(term)
			default:
				panic(BadInput)
			}
		case 1:
			seq := quant.GetNode(0)
			min, max := 0, 0
			minOpt := seq.GetNode(1)
			if minOpt.Count() == 1 {
				var err error
				min, err = strconv.Atoi(minOpt.GetString(0, 0))
				if err != nil {
					panic(err)
				}
			}
			maxOpt := seq.GetNode(3)
			if maxOpt.Count() == 1 {
				var err error
				max, err = strconv.Atoi(maxOpt.GetString(0, 0))
				if err != nil {
					panic(err)
				}
			}
			term = Quant{Term: term, Min: min, Max: max}
		case 2:
			seq := quant.GetNode(0)
			term = Delim{
				Term:  term,
				Sep:   compileTermNamedNode(seq.GetNode(2)),
				Assoc: NewAssociativity(seq.GetString(0)),
			}
		default:
			panic(BadInput)
		}
	}
	return term
}

func compileTermSeqNode(node parse.Node) Term {
	n := node.Count()
	if n == 1 {
		return compileTermQuantNode(node.Children[0].(parse.Node))
	}
	seq := make(Seq, 0, node.Count())
	for _, child := range node.Children {
		seq = append(seq, compileTermQuantNode(child.(parse.Node)))
	}
	return seq
}

func compileTermOneofNode(node parse.Node) Term {
	n := node.Count()
	if n == 1 {
		return compileTermSeqNode(node.GetNode(0))
	}
	oneof := make(Oneof, 0, n/2+1)
	for i := 0; i < n; i += 2 {
		oneof = append(oneof, compileTermSeqNode(node.GetNode(i)))
	}
	return oneof
}

func compileTermStackNode(node parse.Node) Term {
	if node.Count() == 1 {
		return compileTermOneofNode(node.GetNode(0))
	}
	var stack Stack
	for i := 0; i < node.Count(); i += 2 {
		stack = append(stack, compileTermOneofNode(node.GetNode(i)))
	}
	return stack
}

func compileTermNode(node parse.Node) Term {
	return compileTermStackNode(node)
}

func compileProdNode(node parse.Node) Term {
	children := node.GetNode(2).Children
	if len(children) == 1 {
		return compileTermNode(children[0].(parse.Node))
	}
	seq := make(Seq, 0, node.Count())
	for _, child := range children {
		seq = append(seq, compileTermNode(child.(parse.Node)))
	}
	return seq
}

// NewFromNode converts the output from parsing an input via GrammarGrammar into
// a Grammar, which can then be used to generate parsers.
func NewFromNode(node parse.Node) Grammar {
	g := Grammar{}
	for _, v := range node.Children {
		stmt := v.(parse.Node)
		switch stmt.Extra.(int) {
		case 0:
		// 	comment := v.(parse.Node).GetString(0)
		case 1:
			prod := stmt.GetNode(0)
			g[Rule(prod.GetString(0))] = compileProdNode(prod)
		default:
			panic(BadInput)
		}
	}
	return g
}
