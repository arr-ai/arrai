package rel

// import (
// 	"fmt"

// 	"github.com/arr-ai/wbnf/bootstrap"
// 	"github.com/arr-ai/wbnf/parser"
// 	parse "github.com/arr-ai/wbnf/parser"
// )

// // func withRule(rule bootstrap.Rule, v Value) Value {
// // 	if rule == "" {
// // 		return v
// // 	}
// // 	return NewTuple(NewAttr(string(rule), v))
// // }

// func termFromValue(t bootstrap.Term, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	switch t := t.(type) {
// 	case bootstrap.Rule:
// 		return ruleToValue(t, g, rule, v)
// 	case bootstrap.S:
// 		return sToValue(t, g, rule, v)
// 	case bootstrap.RE:
// 		return reToValue(t, g, rule, v)
// 	case bootstrap.Seq:
// 		return seqToValue(t, g, rule, v)
// 	case bootstrap.Oneof:
// 		return oneofToValue(t, g, rule, v)
// 	case bootstrap.Stack:
// 		return stackToValue(t, g, rule, v)
// 	case bootstrap.Delim:
// 		return delimToValue(t, g, rule, v)
// 	case bootstrap.Quant:
// 		return quantToValue(t, g, rule, v)
// 	case bootstrap.Named:
// 		return namedToValue(t, g, rule, v)
// 	default:
// 		panic(fmt.Errorf("unknown term type: %v %[1]T", t))
// 	}
// }

// func sToValue(t bootstrap.S, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	return withRule(rule, NewString([]rune(v.(parse.Scanner).String())))
// }

// func reToValue(t bootstrap.RE, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	return withRule(rule, NewString([]rune(v.(parse.Scanner).String())))
// }

// func seqToValue(t bootstrap.Seq, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	node := v.(parser.Node)
// 	values := make([]Value, 0, len(t))
// 	for i, term := range t {
// 		values = append(values, termFromValue(term, g, rule, node.Children[i]))
// 	}
// 	return withRule(rule, NewArray(values...))
// }

// func oneofToValue(t bootstrap.Oneof, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	node := v.(parser.Node)
// 	return withRule(rule, termFromValue(t[node.Extra.(int)], g, rule, node.Children[0]))
// }

// func delimToValue(t bootstrap.Delim, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	node := v.(parser.Node)
// 	n := node.Count()

// 	left, right := t.LRTerms(node)

// 	values := make([]Value, 0, n)
// 	values = append(values, termFromValue(left, g, "", node.Children[0]))
// 	for i := 1; i < n; i += 2 {
// 		values = append(values, termFromValue(t.Sep, g, "", node.Children[i]))
// 		values = append(values, termFromValue(right, g, "", node.Children[i+1]))
// 	}
// 	return withRule(rule, NewArray(values...))
// }

// func quantToValue(t bootstrap.Quant, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	node := v.(parser.Node)
// 	n := node.Count()
// 	values := make([]Value, 0, n)
// 	for _, child := range node.Children {
// 		values = append(values, termFromValue(t.Term, g, "", child))
// 	}
// 	return withRule(rule, NewArray(values...))
// }

// func ruleToValue(t bootstrap.Rule, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	return termFromValue(g[t], g, t, v)
// }

// //-----------------------------------------------------------------------------

// func stackToValue(t bootstrap.Stack, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	panic(bootstrap.Inconceivable)
// }

// //-----------------------------------------------------------------------------

// func namedToValue(t bootstrap.Named, g bootstrap.Grammar, rule bootstrap.Rule, v interface{}) Value {
// 	return termFromValue(t.Term, g, bootstrap.Rule(t.Name), v)
// }
