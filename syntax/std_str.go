package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// TODO: Make this more robust.
func formatValue(format string, value rel.Value) string {
	var v interface{}
	if set, ok := value.(rel.Set); ok {
		if s, is := rel.AsString(set); is {
			v = s
		} else {
			v = rel.Repr(set)
		}
	} else {
		v = value.Export()
	}
	switch format[len(format)-1] {
	case 't':
		v = value.IsTrue()
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
	stdStrExpand = createNestedFunc("expand", 4, func(args ...rel.Value) rel.Value {
		format := mustAsString(args[0])
		if format != "" {
			format = "%" + format
		} else {
			format = "%v"
		}

		var s string
		if delim := mustAsString(args[2]); strings.HasPrefix(delim, ":") {
			var sb strings.Builder
			n := 0
			for i, ok := args[1].(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); n++ {
				if n > 0 {
					sb.WriteString(delim[1:])
				}
				sb.WriteString(formatValue(format, i.Current()))
			}
			s = sb.String()
		} else {
			s = formatValue(format, args[1])
		}
		if s != "" {
			s += mustAsString(args[3])
		}
		return rel.NewString([]rune(s))
	})

	stdStrRepr = rel.NewNativeFunction("repr", func(value rel.Value) rel.Value {
		return rel.NewString([]rune(rel.Repr(value)))
	})
)

func stdStr() rel.Attr {
	return rel.NewTupleAttr("str",
		createNestedFuncAttr("contains", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.Contains(mustAsString(args[0]), mustAsString(args[1])))
		}),
		rel.NewAttr("expand", stdStrExpand),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasPrefix(mustAsString(args[0]), mustAsString(args[1])))
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasPrefix(mustAsString(args[0]), mustAsString(args[1])))
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
				toJoin = append(toJoin, mustAsString(i.Current()))
			}
			return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[1]))))
		}),
		createNestedFuncAttr("lower", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToLower(mustAsString(args[0]))))
		}),
		rel.NewAttr("repr", stdStrRepr),
		createNestedFuncAttr("split", 2, func(args ...rel.Value) rel.Value {
			splitted := strings.Split(mustAsString(args[0]), mustAsString(args[1]))
			vals := make([]rel.Value, 0, len(splitted))
			for _, s := range splitted {
				vals = append(vals, rel.NewString([]rune(s)))
			}
			return rel.NewArray(vals...)
		}),
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
		createNestedFuncAttr("title", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.Title(mustAsString(args[0]))))
		}),
		createNestedFuncAttr("upper", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToUpper(mustAsString(args[0]))))
		}),
	)
}

func mustAsString(v rel.Value) string {
	// log.Print(v)
	if s, ok := rel.AsString(v.(rel.Set)); ok {
		return s.String()
	}
	panic("value is not a string")
}
