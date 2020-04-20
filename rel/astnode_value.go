package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"

	"github.com/arr-ai/wbnf/parser"
)

func ASTNodeToValue(n ast.Node) Value {
	switch n := n.(type) {
	case ast.Leaf:
		return ASTLeafToValue(n)
	case ast.Branch:
		return ASTBranchToValue(n)
	case wbnf.GrammarNode:
		return ASTNodeToValue(n.Node)
	default:
		panic(fmt.Errorf("unexpected: %v %[1]T", n))
	}
}

func ASTLeafToValue(l ast.Leaf) Value {
	s := l.Scanner()
	return NewOffsetString([]rune(s.String()), s.Offset())
}

func ASTBranchToValue(b ast.Branch) Tuple {
	result := EmptyTuple

	for name, children := range b {
		var value Value
		switch name {
		case "@choice":
			ints := children.(ast.Many)
			values := make([]Value, 0, len(ints))
			for _, i := range ints {
				values = append(values, NewNumber(float64(i.(ast.Extra).Data.(parser.Choice))))
			}
			value = NewArray(values...)
		case "@rule":
			value = NewString([]rune(string(children.(ast.One).Node.(ast.Extra).Data.(parser.Rule))))
		case "@skip":
			value = NewNumber(float64(children.(ast.One).Node.(ast.Extra).Data.(int)))
		default:
			switch c := children.(type) {
			case ast.One:
				value = ASTNodeToValue(c.Node)
			case ast.Many:
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

func ASTNodeFromValue(value Value) ast.Node {
	switch value := value.(type) {
	case String:
		return ASTLeafFromValue(value)
	case Tuple:
		return ASTBranchFromValue(value)
	default:
		panic(fmt.Errorf("unexpected: %v %[1]T", value))
	}
}

func ASTLeafFromValue(s String) ast.Leaf {
	return ast.Leaf(*parser.NewScannerAt(s.String(), s.offset, s.Count()))
}

func ASTBranchFromValue(b Tuple) ast.Branch {
	result := ast.Branch{}
	for i := b.Enumerator(); i.MoveNext(); {
		name, value := i.Current()
		var children ast.Children
		switch name {
		case "@choice":
			values := value.(Array).values
			ints := make(ast.Many, 0, len(values))
			for _, v := range values {
				ints = append(ints, ast.Extra{Data: parser.Choice(v.(Number).Float64())})
			}
			children = ints
		case "@rule":
			children = ast.One{Node: ast.Extra{Data: parser.Rule(value.(String).String())}}
		case "@skip":
			children = ast.One{Node: ast.Extra{Data: int(value.(Number).Float64())}}
		default:
			switch value := value.(type) {
			case Tuple:
				children = ast.One{Node: ASTBranchFromValue(value)}
			case String:
				children = ast.One{Node: ASTLeafFromValue(value)}
			case GenericSet:
				c := make(ast.Many, 0, value.Count())
				for _, v := range value.OrderedValues() {
					c = append(c, ASTNodeFromValue(v.(Tuple).MustGet(ArrayItemAttr)))
				}
				children = c
			case Array:
				c := make(ast.Many, 0, value.Count())
				for _, v := range value.values {
					c = append(c, ASTNodeFromValue(v))
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
