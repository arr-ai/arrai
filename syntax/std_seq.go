package syntax

import (
	"fmt"
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
				return ArrayContains(args[1].(rel.Array), args[0])
			case rel.Bytes:
				return BytesContain(args[1].(rel.Bytes), args[0].(rel.Bytes))
			}

			return rel.NewBool(false)
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
				return nil
			case rel.Bytes:
				return nil
			}

			return nil
		}),
		createNestedFuncAttr("sub", 3, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
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
				return ArraySub(args[2].(rel.Array), args[0], args[1])
			case rel.Bytes:
				return BytesSub(args[2].(rel.Bytes), args[0].(rel.Bytes), args[1].(rel.Bytes))
			}

			return nil
		}),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.String:
				return rel.NewBool(strings.HasPrefix(mustAsString(args[1]), mustAsString(args[0])))
			case rel.Array:
				return ArrayPrefix(args[1].(rel.Array), args[0])
			case rel.Bytes:
				return rel.NewBool(strings.HasPrefix(args[1].String(), args[0].String()))
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.String:
				return rel.NewBool(strings.HasSuffix(mustAsString(args[1]), mustAsString(args[0])))
			case rel.Array:
				return ArraySuffix(args[1].(rel.Array), args[0])
			case rel.Bytes:
				return rel.NewBool(strings.HasSuffix(args[1].String(), args[0].String()))
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			switch args[1].(type) {
			case rel.Set:
				strs := args[1].(rel.Set)
				toJoin := make([]string, 0, strs.Count())
				for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
					toJoin = append(toJoin, mustAsString(i.Current()))
				}
				return rel.NewString([]rune(strings.Join(toJoin, mustAsString(args[0]))))
			case rel.Array:
				return ArrayContains(args[1].(rel.Array), args[0])
			case rel.Bytes:
				return nil
			}

			panic("couldn't find hanlder for subject sequence, the supported subject sequence are string, array and byte array.")
		}),
	)
}
