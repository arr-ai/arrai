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
	stdStrExpand = mustCreateNestedFunc("expand", 4, func(args ...rel.Value) (rel.Value, error) {
		format, is := valueAsString(args[0])
		if !is {
			return nil, fmt.Errorf("expand: format not a string: %v", args[0])
		}
		if format != "" {
			format = "%" + format
		} else {
			format = "%v"
		}

		var s string
		delim, is := valueAsString(args[2])
		if !is {
			return nil, fmt.Errorf("expand: delim not a string: %v", args[2])
		}
		if strings.HasPrefix(delim, ":") {
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
				return nil, fmt.Errorf("arg not an array in ${arg::}: %v", args[1])
			}
		} else {
			s = formatValue(format, args[1])
		}
		if s != "" {
			tail, is := valueAsString(args[3])
			if !is {
				return nil, fmt.Errorf("expand: tail not a string: %v", args[3])
			}
			s += tail
		}
		return rel.NewString([]rune(s)), nil
	})

	stdStrRepr = rel.NewNativeFunction("repr", func(value rel.Value) (rel.Value, error) {
		return rel.NewString([]rune(rel.Repr(value))), nil
	})
)

func stdStr() rel.Attr {
	return rel.NewTupleAttr("str",
		rel.NewAttr("expand", stdStrExpand),
		createNestedFuncAttr("lower", 1, func(args ...rel.Value) (rel.Value, error) {
			return rel.NewString([]rune(strings.ToLower(mustValueAsString(args[0])))), nil
		}),
		rel.NewAttr("repr", stdStrRepr),
		createNestedFuncAttr("title", 1, func(args ...rel.Value) (rel.Value, error) {
			return rel.NewString([]rune(strings.Title(mustValueAsString(args[0])))), nil
		}),
		createNestedFuncAttr("upper", 1, func(args ...rel.Value) (rel.Value, error) {
			return rel.NewString([]rune(strings.ToUpper(mustValueAsString(args[0])))), nil
		}),
	)
}

func valueAsString(v rel.Value) (string, bool) {
	switch v := v.(type) {
	case rel.String:
		return v.String(), true
	case rel.GenericSet:
		return "", !v.IsTrue()
	}
	return "", false
}

// TODO: Remove
func mustValueAsString(v rel.Value) string {
	if s, is := valueAsString(v); is {
		return s
	}
	panic(fmt.Errorf("value not a string: %v", v))
}

func valueAsBytes(v rel.Value) ([]byte, bool) {
	switch v := v.(type) {
	case rel.Bytes:
		return v.Bytes(), true
	case rel.GenericSet:
		return nil, !v.IsTrue()
	}
	return nil, false
}

// TODO: Remove
func mustValueAsBytes(v rel.Value) []byte {
	if b, is := valueAsBytes(v); is {
		return b
	}
	panic(fmt.Errorf("value not a byte array: %v", v))
}
