package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
)

func strArrToRelArr(s []string) rel.Value {
	values := make([]rel.Value, 0, len(s))
	for _, a := range s {
		values = append(values, rel.NewString([]rune(a)))
	}
	return rel.NewArray(values...)
}

// hasRule traverses to each node in an ast and determine if rule exists or not in the ast
func hasRule(b ast.Branch, rule string) bool {
	if _, has := b[rule]; has {
		return true
	}
	for _, c := range b {
		switch n := c.(type) {
		case ast.One:
			if branch, isBranch := n.Node.(ast.Branch); isBranch && hasRule(branch, rule) {
				return true
			}
		case ast.Many:
			for _, i := range n {
				if branch, isBranch := i.(ast.Branch); isBranch && hasRule(branch, rule) {
					return true
				}
			}
		}
	}
	return false
}
