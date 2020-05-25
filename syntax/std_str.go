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
			if array, is := rel.AsArray(args[1].(rel.Set)); is {
				var sb strings.Builder
				for i, value := range array.Values() {
					if i > 0 {
						sb.WriteString(delim[1:])
					}
					sb.WriteString(formatValue(format, value))
				}
				s = sb.String()
			} else {
				panic(fmt.Errorf("arg not an array in ${arg::}: %v", args[1]))
			}
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
		rel.NewAttr("expand", stdStrExpand),
		createNestedFuncAttr("lower", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToLower(mustAsString(args[0]))))
		}),
		rel.NewAttr("repr", stdStrRepr),
		createNestedFuncAttr("title", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.Title(mustAsString(args[0]))))
		}),
		createNestedFuncAttr("upper", 1, func(args ...rel.Value) rel.Value {
			return rel.NewString([]rune(strings.ToUpper(mustAsString(args[0]))))
		}),
	)
}

func mustAsString(v rel.Value) string {
	if s, ok := rel.AsString(v.(rel.Set)); ok {
		return s.String()
	}
	panic("value is not a string")
}
