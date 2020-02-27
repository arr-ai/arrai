package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// TODO: Make this more robust.
func formatValue(format string, value rel.Value) string {
	v := value.Export()
	switch format[len(format)-1] {
	case 't':
		v = value.Bool()
	case 'c', 'd', 'o', 'O', 'x', 'X', 'U':
		v = int(value.Export().(float64))
	case 'q':
		if f, ok := v.(float64); ok {
			v = int(f)
		}
	}
	return fmt.Sprintf(format, v)
}

var (
	libStrConcat = createNestedFunc("concat", 1, func(args ...rel.Value) rel.Value {
		var sb strings.Builder
		for i := args[0].(rel.Set).ArrayEnumerator(); i.MoveNext(); {
			sb.WriteString(mustAsString(i.Current()).String())
		}
		return rel.NewString([]rune(sb.String()))
	})

	libStrExpand = createNestedFunc("expand", 3, func(args ...rel.Value) rel.Value {
		format := mustAsString(args[0]).String()
		if format != "" {
			format = "%" + format
		} else {
			format = "%v"
		}

		if delim := mustAsString(args[2]).String(); strings.HasPrefix(delim, ":") {
			var sb strings.Builder
			for n, i := 0, args[1].(rel.Set).ArrayEnumerator(); i.MoveNext(); n++ {
				if n > 0 {
					sb.WriteString(delim[1:])
				}
				sb.WriteString(formatValue(format, i.Current()))
			}
			return rel.NewString([]rune(sb.String()))
		}
		return rel.NewString([]rune(formatValue(format, args[1])))
	})
)

func loadStrLib() rel.Attr {
	return rel.NewAttr("str", rel.NewTuple(
		createFunc("sub", 3, func(args ...rel.Value) rel.Value {
			return rel.NewString(
				[]rune(
					strings.ReplaceAll(
						mustAsString(args[0]).String(),
						mustAsString(args[1]).String(),
						mustAsString(args[2]).String(),
					),
				),
			)
		}),
		createFunc("split", 2, func(args ...rel.Value) rel.Value {
			splitted := strings.Split(mustAsString(args[0]).String(), mustAsString(args[1]).String())
			vals := make([]rel.Value, 0, len(splitted))
			for _, s := range splitted {
				vals = append(vals, rel.NewString([]rune(s)))
			}
			return rel.NewArray(vals...)
		}),
		createFunc("lower", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToLower(mustAsString(args[0]).String())))
		}),
		createFunc("upper", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToUpper(mustAsString(args[0]).String())))
		}),
		createFunc("title", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.Title(mustAsString(args[0]).String())))
		}),
		createFunc("contains", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.Contains(mustAsString(args[0]).String(), mustAsString(args[1]).String()))
		}),
		rel.NewAttr("concat", libStrConcat),
		createFunc("join", 2, func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i := strs.(rel.Set).ArrayEnumerator(); i.MoveNext(); {
				toJoin = append(toJoin, mustAsString(i.Current()).String())
			}
			return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[1]).String())))
		}),
		rel.NewAttr("expand", libStrExpand),
	))
}

func mustAsString(v rel.Value) rel.String {
	s, isString := v.(rel.Set).AsString()
	if isString {
		return s
	}
	panic("can not be a string")
}
