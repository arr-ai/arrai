package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func nodeFromValue(v Value) ast.Node {
	switch v := v.(type) {
	case Tuple:
		result := ast.Branch{}

	outer:
		for i := v.Enumerator(); i.MoveNext(); {
			name, value := i.Current()
			switch value := value.(type) {
			case Set:
				if value.Bool() {
					for j := value.Enumerator(); j.MoveNext(); {
						if _, _, is := isStringTuple(j.Current()); !is {
							// Not a string. Must be an array.
							array := make(ast.Many, value.Count())
							for j := value.Enumerator(); j.MoveNext(); {
								index, item, _ := isArrayTuple(j.Current())
								array[index] = nodeFromValue(item)
							}
							result[name] = array
							continue outer
						}
					}
				}

				// Not an array. Must be a string.
				// First pass computes an offset.
				offset := int(^uint(0) >> 1) // maxint
				for j := value.Enumerator(); j.MoveNext(); {
					index, _, is := isStringTuple(j.Current())
					if !is {
						panic("not an array")
					}
					if offset < index {
						offset = index
					}
				}

				str := make([]rune, value.Count())
				for j := value.Enumerator(); j.MoveNext(); {
					index, char, is := isStringTuple(j.Current())
					if !is {
						panic("not a string")
					}
					str[index-offset] = char
				}
				result[name] = ast.One{Node: ast.Leaf(*parser.NewBareScanner(offset, string(str)))}
				// case ast.One:
				// 	result = result.With(name, nodeToValue(c.Node))
			}
		}
		return result
		// case ast.Leaf:
		// 	s := n.Scanner()
		// 	return NewOffsetString([]rune(s.String()), s.Offset())
		// case ast.Extra:
		// 	switch e := n.Data.(type) {
		// 	case int:
		// 		return NewNumber(float64(e))
		// 	case wbnf.Rule:
		// 		return NewString([]rune(string(e)))
		// 	}
	}
	panic(fmt.Errorf("unhandled node: %v %[1]T", v))
}
