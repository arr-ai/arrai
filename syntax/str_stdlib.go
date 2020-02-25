package syntax

import (
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func loadStrLib() rel.Attr {
	return rel.NewAttr("str", rel.NewTuple(
		createFunc("sub", 3, func(args ...rel.Value) rel.Value {
			return rel.NewString(
				[]rune(
					strings.ReplaceAll(
						args[0].(rel.String).String(),
						args[1].(rel.String).String(),
						args[2].(rel.String).String(),
					),
				),
			)
		}),
		createFunc("split", 2, func(args ...rel.Value) rel.Value {
			splitted := strings.Split(args[0].(rel.String).String(), args[1].(rel.String).String())
			vals := make([]rel.Value, 0, len(splitted))
			for _, s := range splitted {
				vals = append(vals, rel.NewString([]rune(s)))
			}
			return rel.NewArray(vals...)
		}),
		createFunc("lower", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToLower(args[0].(rel.String).String())))
		}),
		createFunc("upper", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToUpper(args[0].(rel.String).String())))
		}),
		createFunc("title", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.Title(args[0].(rel.String).String())))
		}),
		createFunc("split", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.Contains(args[0].(rel.String).String(), args[1].(rel.String).String()))
		}),
		createFunc("concat", 1, func(args ...rel.Value) rel.Value {
			var sb strings.Builder
			for i := args[0].(rel.Array).ArrayEnumerator(); i.MoveNext(); {
				sb.WriteString(i.Current().(rel.String).String())
			}
			return rel.NewString([]rune(sb.String()))
		}),
		createFunc("join", 2, func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i := strs.(rel.Array).ArrayEnumerator(); i.MoveNext(); {
				toJoin = append(toJoin, i.Current().(rel.String).String())
			}
			return rel.NewString([]rune(strings.Join(toJoin, args[1].(rel.String).String())))
		}),
	))
}
