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
			sb.WriteString(mustAsString(i.Current()))
		}
		return rel.NewString([]rune(sb.String()))
	})

	libStrExpand = createNestedFunc("expand", 3, func(args ...rel.Value) rel.Value {
		format := mustAsString(args[0])
		if format != "" {
			format = "%" + format
		} else {
			format = "%v"
		}

		if delim := mustAsString(args[2]); strings.HasPrefix(delim, ":") {
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

func stdStr() rel.Attr {
	return rel.NewAttr("str", rel.NewTuple(
		createNestedFuncAttr("sub", 3, func(args ...rel.Value) rel.Value {
			return rel.NewString(
				[]rune(
					strings.ReplaceAll(
						mustAsString(args[0]),
						mustAsString(args[1]),
						mustAsString(args[2]),
					),
				),
			)
		}),
		createNestedFuncAttr("split", 2, func(args ...rel.Value) rel.Value {
			splitted := strings.Split(mustAsString(args[0]), mustAsString(args[1]))
			vals := make([]rel.Value, 0, len(splitted))
			for _, s := range splitted {
				vals = append(vals, rel.NewString([]rune(s)))
			}
			return rel.NewArray(vals...)
		}),
		createNestedFuncAttr("lower", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToLower(mustAsString(args[0]))))
		}),
		createNestedFuncAttr("upper", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToUpper(mustAsString(args[0]))))
		}),
		createNestedFuncAttr("title", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.Title(mustAsString(args[0]))))
		}),
		createNestedFuncAttr("contains", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.Contains(mustAsString(args[0]), mustAsString(args[1])))
		}),
		createNestedFuncAttr("hasPrefix", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasPrefix(mustAsString(args[0]), mustAsString(args[1])))
		}),
		createNestedFuncAttr("hasSuffix", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasPrefix(mustAsString(args[0]), mustAsString(args[1])))
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i := strs.(rel.Set).ArrayEnumerator(); i.MoveNext(); {
				toJoin = append(toJoin, mustAsString(i.Current()))
			}
			return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[1]))))
		}),
		rel.NewAttr("concat", libStrConcat),
		rel.NewAttr("expand", libStrExpand),
	))
}

func mustAsString(v rel.Value) string {
	if s, ok := v.(rel.Set).AsString(); ok {
		return s.String()
	}
	panic("can not be a string")
}
