package syntax

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func stdSeqConcat(arg rel.Value) rel.Value {
	if set, is := arg.(rel.Set); is {
		if !set.IsTrue() {
			return rel.None
		}
	}
	values := arg.(rel.Array).Values()
	if len(values) == 0 {
		return rel.None
	}
	switch v0 := values[0].(type) {
	case rel.String:
		var sb strings.Builder
		for _, value := range values {
			sb.WriteString(mustAsString(value))
		}
		return rel.NewString([]rune(sb.String()))
	case rel.Set:
		result := v0
		for _, value := range values[1:] {
			var err error
			result, err = rel.Concatenate(result, value.(rel.Set))
			if err != nil {
				panic(err)
			}
		}
		return result
	}
	panic(fmt.Errorf("concat: incompatible value: %v", values[0]))
}

func stdSeqRepeat(arg rel.Value) rel.Value {
	n := int(arg.(rel.Number))
	return rel.NewNativeFunction("repeat(n)", func(arg rel.Value) rel.Value {
		switch seq := arg.(type) {
		case rel.String:
			return rel.NewString([]rune(strings.Repeat(seq.String(), n)))
		case rel.Array:
			values := []rel.Value{}
			seqValues := seq.Values()
			for i := 0; i < n; i++ {
				values = append(values, seqValues...)
			}
			return rel.NewArray(values...)
		case rel.Set:
			if !seq.IsTrue() {
				return rel.None
			}
		}
		panic(fmt.Errorf("repeat: unsupported value: %v", arg))
	})
}

func stdSeq() rel.Attr {
	return rel.NewTupleAttr("seq",
		rel.NewNativeFunctionAttr("concat", stdSeqConcat),
		rel.NewNativeFunctionAttr("repeat", stdSeqRepeat),
		createNestedFuncAttr("contains", 2, func(args ...rel.Value) rel.Value {
			return process("contains", args...)
		}),
		createNestedFuncAttr("split", 2, func(args ...rel.Value) rel.Value {
			return process("split", args...)
		}),
		createNestedFuncAttr("sub", 3, func(args ...rel.Value) rel.Value {
			return process("sub", args...)
		}),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value {
			return process("has_prefix", args...)
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value {
			return process("has_suffix", args...)
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			return process("join", args...)
		}),
	)
}

func process(apiName string, args ...rel.Value) rel.Value {
	handler := handlerMapping[typeMethod{reflect.TypeOf(args[0]), apiName}]
	return handler(args...)
}

// This seq API handlders mapping, the key is API name + '_' + data type.
var (
	handlerMapping = map[typeMethod]func(...rel.Value) rel.Value{
		// API contains
		typeMethod{reflect.TypeOf(rel.String{}), "contains"}: func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.Contains(mustAsString(args[0]), mustAsString(args[1])))
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "contains"}: func(args ...rel.Value) rel.Value {
			a := args[0].(rel.Array)
			switch b := args[1].(type) {
			case rel.Array:
				return ContainsArray(a, b)
			case rel.Value:
				arrayEnum, _ := a.ArrayEnumerator()
				if arrayEnum != nil {
					for arrayEnum.MoveNext() {
						if arrayEnum.Current().Equal(b) {
							return rel.NewBool(true)
						}
					}
				}

			}
			return rel.NewBool(false)
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "contains"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		// API sub
		typeMethod{reflect.TypeOf(rel.String{}), "sub"}: func(args ...rel.Value) rel.Value {
			return rel.NewString(
				[]rune(
					strings.ReplaceAll(
						mustAsString(args[0]),
						mustAsString(args[1]),
						mustAsString(args[2]),
					),
				),
			)
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "sub"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "sub"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		// API split
		typeMethod{reflect.TypeOf(rel.String{}), "split"}: func(args ...rel.Value) rel.Value {
			splitted := strings.Split(mustAsString(args[0]), mustAsString(args[1]))
			vals := make([]rel.Value, 0, len(splitted))
			for _, s := range splitted {
				vals = append(vals, rel.NewString([]rune(s)))
			}
			return rel.NewArray(vals...)
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "split"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "split"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		// API join
		typeMethod{reflect.TypeOf(rel.GenericSet{}), "join"}: func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
				toJoin = append(toJoin, mustAsString(i.Current()))
			}
			return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[1]))))
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "join"}: func(args ...rel.Value) rel.Value {
			strs := args[0].(rel.Set)
			toJoin := make([]string, 0, strs.Count())
			for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
				toJoin = append(toJoin, mustAsString(i.Current()))
			}
			return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[1]))))
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "join"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		// API has_prefix
		typeMethod{reflect.TypeOf(rel.String{}), "has_prefix"}: func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasPrefix(mustAsString(args[0]), mustAsString(args[1])))
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "has_prefix"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "has_prefix"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		// API has_suffix
		typeMethod{reflect.TypeOf(rel.String{}), "has_suffix"}: func(args ...rel.Value) rel.Value {
			return rel.NewBool(strings.HasSuffix(mustAsString(args[0]), mustAsString(args[1])))
		},
		typeMethod{reflect.TypeOf(rel.Array{}), "has_suffix"}: func(args ...rel.Value) rel.Value {
			return nil
		},
		typeMethod{reflect.TypeOf(rel.Bytes{}), "has_suffix"}: func(args ...rel.Value) rel.Value {
			return nil
		},
	}
)

type typeMethod struct {
	t       reflect.Type
	apiName string
}
