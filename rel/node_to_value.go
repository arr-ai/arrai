package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"
)

func nodeToValue(n ast.Node) Value {
	switch n := n.(type) {
	case ast.Branch:
		result := EmptyTuple

		for name, children := range n {
			switch c := children.(type) {
			case ast.Many:
				values := make([]Value, 0, len(c))
				for _, child := range c {
					values = append(values, nodeToValue(child))
				}
				result = result.With(name, NewArray(values...))
			case ast.One:
				result = result.With(name, nodeToValue(c.Node))
			}
		}
		return result
	case ast.Leaf:
		s := n.Scanner()
		return NewOffsetString([]rune(s.String()), s.Offset())
	case ast.Extra:
		switch e := n.Data.(type) {
		case int:
			return NewNumber(float64(e))
		case wbnf.Rule:
			return NewString([]rune(string(e)))
		}
	}
	panic(fmt.Errorf("unhandled node: %v %[1]T", n))
}
