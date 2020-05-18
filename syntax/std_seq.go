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
			switch args[1].(type) {
			case rel.String:
				return rel.NewBool(strings.Contains(mustAsString(args[1]), mustAsString(args[0])))
			case rel.Array:
				return arrayContain(args[1].(rel.Array), args[0])
			case rel.Bytes:
				switch args[0].(type) {
				case rel.GenericSet:
					if len(args[1].String()) > 0 {
						return rel.NewBool(true)
					}
				}
				return rel.NewBool(strings.Contains(args[1].String(), args[0].String()))
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.String:
				return rel.NewBool(strings.HasPrefix(mustAsString(args[1]), mustAsString(args[0])))
			case rel.Array:
				return arrayHasPrefix(args[1].(rel.Array), args[0])
			case rel.Bytes:
				switch args[0].(type) {
				case rel.GenericSet:
					if len(args[1].String()) > 0 {
						return rel.NewBool(true)
					}
				}
				return rel.NewBool(strings.HasPrefix(args[1].String(), args[0].String()))
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.String:
				return rel.NewBool(strings.HasSuffix(mustAsString(args[1]), mustAsString(args[0])))
			case rel.Array:
				return arrayHasSuffix(args[1].(rel.Array), args[0])
			case rel.Bytes:
				switch args[0].(type) {
				case rel.GenericSet:
					if len(args[1].String()) > 0 {
						return rel.NewBool(true)
					}
				}
				return rel.NewBool(strings.HasSuffix(args[1].String(), args[0].String()))
			}

			return rel.NewBool(false)
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
				return arraySub(args[2].(rel.Array), args[0], args[1])
			case rel.Bytes:
				return rel.NewBytes([]byte(strings.ReplaceAll(args[2].String(), args[0].String(), args[1].String())))
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
				return arraySplit(args[1].(rel.Array), args[0])
			case rel.Bytes:
				return bytesSplit(args[1].(rel.Bytes), args[0])
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
					return strJoin(args...)
				case rel.Value:
					return arrayJoin(a1, args[0])
				}
			case rel.Bytes:
				return bytesJoin(args[1].(rel.Bytes), args[0].(rel.Bytes))
			}

			panic(fmt.Errorf(sharedError, reflect.TypeOf(args[2])))
		}),
	)
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
