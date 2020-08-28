package syntax

import (
	"context"
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

// TODO: Make this more robust.
func formatValue(ctx context.Context, format string, value rel.Value) string {
	var v interface{}
	switch set := value.(type) {
	case rel.Set:
		if s, is := tools.ValueAsString(set); is {
			v = s
		} else if s, is := tools.ValueAsBytes(set); is {
			v = string(s)
		} else {
			v = rel.Repr(set)
		}
	case nil:
		panic(fmt.Errorf("unable to format nil value"))
	default:
		v = value.Export(ctx)
	}
	switch format[len(format)-1] {
	case 't':
		v = value.IsTrue()
	case 'c', 'd', 'o', 'O', 'x', 'X', 'U':
		v = int(value.Export(ctx).(float64))
	case 'q':
		if f, ok := v.(float64); ok {
			v = int(f)
		}
	}
	return fmt.Sprintf(format, v)
}

var (
	stdStrExpand = mustCreateNestedFunc("expand", 4, func(ctx context.Context, args ...rel.Value) (rel.Value, error) {
		format, is := tools.ValueAsString(args[0])
		if !is {
			return nil, fmt.Errorf("//str.expand: format not a string: %v", args[0])
		}
		if format != "" {
			format = "%" + format
		} else {
			format = "%v"
		}

		var s string
		delim, is := tools.ValueAsString(args[2])
		if !is {
			return nil, fmt.Errorf("//str.expand: delim not a string: %v", args[2])
		}
		if strings.HasPrefix(delim, ":") {
			if array, is := rel.AsArray(args[1].(rel.Set)); is {
				var sb strings.Builder
				for i, value := range array.Values() {
					if i > 0 {
						sb.WriteString(delim[1:])
					}
					if value != nil {
						sb.WriteString(formatValue(ctx, format, value))
					}
				}
				s = sb.String()
			} else {
				return nil, fmt.Errorf("//str..expand: arg not an array in ${arg::}: %v", args[1])
			}
		} else {
			s = formatValue(ctx, format, args[1])
		}
		if s != "" {
			tail, is := tools.ValueAsString(args[3])
			if !is {
				return nil, fmt.Errorf("//str.expand: tail not a string: %v", args[3])
			}
			s += tail
		}
		return rel.NewString([]rune(s)), nil
	})

	stdStrRepr = rel.NewNativeFunction("repr", func(_ context.Context, value rel.Value) (rel.Value, error) {
		return rel.NewString([]rune(rel.Repr(value))), nil
	})
)

func stdStr() rel.Attr {
	return rel.NewTupleAttr("str",
		rel.NewAttr("expand", stdStrExpand),
		createNestedFuncAttr("lower", 1, func(_ context.Context, args ...rel.Value) (rel.Value, error) {
			if s, is := tools.ValueAsString(args[0]); is {
				return rel.NewString([]rune(strings.ToLower(s))), nil
			}
			return nil, fmt.Errorf("//str.lower: arg not a string: %v", args[0])
		}),
		rel.NewAttr("repr", stdStrRepr),
		createNestedFuncAttr("title", 1, func(_ context.Context, args ...rel.Value) (rel.Value, error) {
			if s, is := tools.ValueAsString(args[0]); is {
				return rel.NewString([]rune(strings.Title(s))), nil
			}
			return nil, fmt.Errorf("//str.title: arg not a string: %v", args[0])
		}),
		createNestedFuncAttr("upper", 1, func(_ context.Context, args ...rel.Value) (rel.Value, error) {
			if s, is := tools.ValueAsString(args[0]); is {
				return rel.NewString([]rune(strings.ToUpper(s))), nil
			}
			return nil, fmt.Errorf("//str.upper: arg not a string: %v", args[0])
		}),
	)
}
