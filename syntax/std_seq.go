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
			return rel.NewBool(strings.Contains(mustAsString(args[0]), mustAsString(args[1])))
		}),
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
	)
}
