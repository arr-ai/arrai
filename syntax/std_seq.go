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
			return includingProcess(func(args ...rel.Value) rel.Value {
				return rel.NewBool(strings.Contains(mustAsString(args[1]), mustAsString(args[0])))
			},
				func(args ...rel.Value) rel.Value {
					return arrayContains(args[0], args[1].(rel.Array))
				},
				func(args ...rel.Value) rel.Value {
					return rel.NewBool(strings.Contains(args[1].String(), args[0].String()))
				},
				args...)
		}),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value {
			return includingProcess(func(args ...rel.Value) rel.Value {
				return rel.NewBool(strings.HasPrefix(mustAsString(args[1]), mustAsString(args[0])))
			},
				func(args ...rel.Value) rel.Value {
					return arrayHasPrefix(args[0], args[1].(rel.Array))
				},
				func(args ...rel.Value) rel.Value {
					return rel.NewBool(strings.HasPrefix(args[1].String(), args[0].String()))
				},
				args...)
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value {
			return includingProcess(func(args ...rel.Value) rel.Value {
				return rel.NewBool(strings.HasSuffix(mustAsString(args[1]), mustAsString(args[0])))
			},
				func(args ...rel.Value) rel.Value {
					return arrayHasSuffix(args[0], args[1].(rel.Array))
				},
				func(args ...rel.Value) rel.Value {
					return rel.NewBool(strings.HasSuffix(args[1].String(), args[0].String()))
				},
				args...)
		}),
		createNestedFuncAttr("sub", 3, func(args ...rel.Value) rel.Value {
			switch args[2].(type) {
			case rel.String:
				return rel.NewString(
					[]rune(
						strings.ReplaceAll(
							mustAsString(args[2]),
							mustAsString(args[0]),
							mustAsString(args[1]),
						),
					),
				)
			case rel.Array:
				return arraySub(args[0], args[1], args[2].(rel.Array))
			case rel.Bytes:
				_, arg0IsSet := args[0].(rel.GenericSet)
				_, arg1IsSet := args[1].(rel.GenericSet)
				if !arg0IsSet && arg1IsSet {
					return rel.NewBytes([]byte(strings.ReplaceAll(args[2].String(),
						args[0].String(), "")))
				} else if arg0IsSet && !arg1IsSet {
					return rel.NewBytes([]byte(strings.ReplaceAll(args[2].String(),
						"", args[1].String())))
				}
				return rel.NewBytes([]byte(strings.ReplaceAll(args[2].String(),
					args[0].String(), args[1].String())))
			case rel.GenericSet:
				_, arg0IsSet := args[0].(rel.GenericSet)
				_, arg1IsSet := args[1].(rel.GenericSet)
				if arg0IsSet && arg1IsSet {
					return args[2]
				} else if arg0IsSet && !arg1IsSet {
					return args[1]
				} else if !arg0IsSet && arg1IsSet {
					return args[2]
				}
			}

			panic(fmt.Errorf(sharedError, reflect.TypeOf(args[2])))
		}),
		createNestedFuncAttr("split", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.String:
				splitted := strings.Split(mustAsString(args[1]), mustAsString(args[0]))
				vals := make([]rel.Value, 0, len(splitted))
				for _, s := range splitted {
					vals = append(vals, rel.NewString([]rune(s)))
				}
				return rel.NewArray(vals...)
			case rel.Array:
				return arraySplit(args[0], args[1].(rel.Array))
			case rel.Bytes:
				return bytesSplit(args[0], args[1].(rel.Bytes))
			case rel.GenericSet:
				switch args[0].(type) {
				case rel.String:
					return rel.NewArray(args[1])
				case rel.Array, rel.Bytes:
					return rel.NewArray(rel.NewArray())
				case rel.GenericSet:
					return args[1]
				}
			}

			panic(fmt.Errorf("expected subject sequence types are %s, %s and %s, but the actual type is %s",
				reflect.TypeOf(rel.String{}), reflect.TypeOf(rel.Array{}), reflect.TypeOf(rel.Bytes{}),
				reflect.TypeOf(args[2])))
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			switch a1 := args[1].(type) {
			case rel.Array:
				switch a1.Values()[0].(type) {
				case rel.String:
					// if subject is rel.String
					return strJoin(args...)
				case rel.Value:
					if _, isStr := args[0].(rel.String); isStr {
						return strJoin(args...)
					}
					return arrayJoin(args[0], a1)
				}
			case rel.Bytes:
				if _, isSet := args[0].(rel.GenericSet); isSet {
					return args[1]
				}
				return bytesJoin(args[0].(rel.Bytes), args[1].(rel.Bytes))
			case rel.GenericSet:
				switch args[0].(type) {
				case rel.String:
					// if joiner is rel.String
					return strJoin(args...)
				case rel.Array, rel.GenericSet, rel.Bytes:
					return args[1]
				}
			}

			panic(fmt.Errorf(sharedError, reflect.TypeOf(args[2])))
		}),
	)
}

// Shared method for contains, hasPrefix and hasSuffix
func includingProcess(
	strHandler,
	arrayHandler,
	bytesHandler func(...rel.Value) rel.Value,
	args ...rel.Value) rel.Value {

	switch args[1].(type) {
	case rel.String:
		return strHandler(args...)
	case rel.Array:
		return arrayHandler(args...)
	case rel.Bytes:
		if _, isSet := args[0].(rel.GenericSet); isSet {
			if len(args[1].String()) > 0 {
				return rel.NewBool(true)
			}
		}
		return bytesHandler(args...)
	case rel.GenericSet:
		if _, isSet := args[0].(rel.GenericSet); isSet {
			return rel.NewBool(true)
		}
	}

	return rel.NewBool(false)
}

func strJoin(args ...rel.Value) rel.Value {
	strs := args[1].(rel.Set)
	toJoin := make([]string, 0, strs.Count())
	for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
		toJoin = append(toJoin, mustAsString(i.Current()))
	}
	return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[0]))))
}

var sharedError = "expected subject sequence types are rel.String, rel.Array and rel.Bytes, but the actual type is %s"
