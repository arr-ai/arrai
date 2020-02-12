package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

func ASTNodeToValue(n wbnf.Node) Value {
	switch n := n.(type) {
	case wbnf.Leaf:
		return ASTLeafToValue(n)
	case wbnf.Branch:
		return ASTBranchToValue(n)
	default:
		panic(fmt.Errorf("unexpected: %v %[1]T", n))
	}
}

func ASTLeafToValue(l wbnf.Leaf) Value {
	s := l.Scanner()
	return NewOffsetString([]rune(s.String()), s.Offset())
}

func ASTBranchToValue(b wbnf.Branch) Tuple {
	result := EmptyTuple

	for name, children := range b {
		var value Value
		switch name {
		case "@choice":
			ints := children.(wbnf.Many)
			values := make([]Value, 0, len(ints))
			for _, i := range ints {
				values = append(values, NewNumber(float64(i.(wbnf.Extra).Data.(parser.Choice))))
			}
			value = NewArray(values...)
		case "@rule":
			value = NewString([]rune(string(children.(wbnf.One).Node.(wbnf.Extra).Data.(parser.Rule))))
		case "@skip":
			value = NewNumber(float64(children.(wbnf.One).Node.(wbnf.Extra).Data.(int)))
		default:
			switch c := children.(type) {
			case wbnf.One:
				value = ASTNodeToValue(c.Node)
			case wbnf.Many:
				values := make([]Value, 0, len(c))
				for _, child := range c {
					values = append(values, ASTNodeToValue(child))
				}
				value = NewArray(values...)
			default:
				panic(fmt.Errorf("unexpected: %v %[1]T", value))
			}
		}
		result = result.With(name, value)
	}

	return result
}

func ASTNodeFromValue(value Value) wbnf.Node {
	switch value := value.(type) {
	case String:
		return ASTLeafFromValue(value)
	case Tuple:
		return ASTBranchFromValue(value)
	default:
		panic(fmt.Errorf("unexpected: %v %[1]T", value))
	}
}

func ASTLeafFromValue(s String) wbnf.Leaf {
	return wbnf.Leaf(*parser.NewBareScanner(s.offset, s.String()))
}

func ASTBranchFromValue(b Tuple) wbnf.Branch {
	result := wbnf.Branch{}
	for i := b.Enumerator(); i.MoveNext(); {
		name, value := i.Current()
		var children wbnf.Children
		switch name {
		case "@choice":
			values := value.(*genericSet).OrderedValues()
			ints := make(wbnf.Many, 0, len(values))
			for _, v := range values {
				ints = append(ints, wbnf.Extra{Data: parser.Choice(v.(Tuple).MustGet(ArrayItemAttr).(Number).Float64())})
			}
			children = ints
		case "@rule":
			children = wbnf.One{Node: wbnf.Extra{Data: parser.Rule(value.(String).String())}}
		case "@skip":
			children = wbnf.One{Node: wbnf.Extra{Data: int(value.(Number).Float64())}}
		default:
			switch value := value.(type) {
			case Tuple:
				children = wbnf.One{Node: ASTBranchFromValue(value)}
			case String:
				children = wbnf.One{Node: ASTLeafFromValue(value)}
			case *genericSet:
				c := make(wbnf.Many, 0, value.Count())
				for _, v := range value.OrderedValues() {
					c = append(c, ASTNodeFromValue(v.(Tuple).MustGet(ArrayItemAttr)))
				}
				children = c
			default:
				panic(fmt.Errorf("unexpected: %v %[1]T", value))
			}
		}
		result[name] = children
	}
	return result
}
