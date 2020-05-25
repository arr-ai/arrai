package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func stdSeqConcat(seq rel.Value) rel.Value {
	if set, is := seq.(rel.Set); is {
		if !set.IsTrue() {
			return rel.None
		}
	}
	values := seq.(rel.Array).Values()
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
		createNestedFuncAttr("contains", 2, func(args ...rel.Value) rel.Value { //nolint:dupl
			sub, subject := args[0], args[1]
			switch subject.(type) {
			case rel.String:
				return rel.NewBool(strings.Contains(mustAsString(subject), mustAsString(sub)))
			case rel.Array:
				return arrayContains(sub, subject.(rel.Array))
			case rel.Bytes:
				return rel.NewBool(strings.Contains(asString(subject), asString(sub)))
			case rel.GenericSet:
				if emptySet, isSet := sub.(rel.GenericSet); isSet && !emptySet.IsTrue() {
					return rel.NewBool(true)
				}
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("has_prefix", 2, func(args ...rel.Value) rel.Value { //nolint:dupl
			prefix, subject := args[0], args[1]
			switch subject.(type) {
			case rel.String:
				return rel.NewBool(strings.HasPrefix(mustAsString(subject), mustAsString(prefix)))
			case rel.Array:
				return arrayHasPrefix(prefix, subject.(rel.Array))
			case rel.Bytes:
				return rel.NewBool(strings.HasPrefix(asString(subject), asString(prefix)))
			case rel.GenericSet:
				if emptySet, isSet := prefix.(rel.GenericSet); isSet && !emptySet.IsTrue() {
					return rel.NewBool(true)
				}
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("has_suffix", 2, func(args ...rel.Value) rel.Value { //nolint:dupl
			suffix, subject := args[0], args[1]
			switch subject.(type) {
			case rel.String:
				return rel.NewBool(strings.HasSuffix(mustAsString(subject), mustAsString(suffix)))
			case rel.Array:
				return arrayHasSuffix(suffix, subject.(rel.Array))
			case rel.Bytes:
				return rel.NewBool(strings.HasSuffix(asString(subject), asString(suffix)))
			case rel.GenericSet:
				if emptySet, isSet := suffix.(rel.GenericSet); isSet && !emptySet.IsTrue() {
					return rel.NewBool(true)
				}
			}

			return rel.NewBool(false)
		}),
		createNestedFuncAttr("sub", 3, func(args ...rel.Value) rel.Value {
			old, new, subject := args[0], args[1], args[2]
			switch subject := subject.(type) {
			case rel.String:
				return rel.NewString(
					[]rune(
						strings.ReplaceAll(
							mustAsString(subject),
							mustAsString(old),
							mustAsString(new),
						),
					),
				)
			case rel.Array:
				return arraySub(old, new, subject)
			case rel.Bytes:
				_, oldIsSet := old.(rel.GenericSet)
				_, newIsSet := new.(rel.GenericSet)
				if !oldIsSet && newIsSet {
					return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(),
						old.String(), "")))
				} else if oldIsSet && !newIsSet {
					return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(),
						"", new.String())))
				}
				return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(),
					old.String(), new.String())))
			case rel.GenericSet:
				_, oldIsSet := old.(rel.GenericSet)
				_, newIsSet := new.(rel.GenericSet)
				if oldIsSet && newIsSet {
					return subject
				} else if oldIsSet && !newIsSet {
					return new
				} else if !oldIsSet && newIsSet {
					return subject
				}
			}

			panic(fmt.Errorf("sub: unsupported args: %s, %s, %s", old, new, subject))
		}),
		createNestedFuncAttr("split", 2, func(args ...rel.Value) rel.Value {
			delimiter, subject := args[0], args[1]
			switch subject := subject.(type) {
			case rel.String:
				splitted := strings.Split(mustAsString(subject), mustAsString(delimiter))
				vals := make([]rel.Value, 0, len(splitted))
				for _, s := range splitted {
					vals = append(vals, rel.NewString([]rune(s)))
				}
				return rel.NewArray(vals...)
			case rel.Array:
				return arraySplit(delimiter, subject)
			case rel.Bytes:
				return bytesSplit(delimiter, subject)
			case rel.GenericSet:
				switch delimiter.(type) {
				case rel.String:
					return rel.NewArray(subject)
				case rel.Array, rel.Bytes:
					return rel.NewArray(rel.NewArray())
				case rel.GenericSet:
					return subject
				}
			}

			panic(fmt.Errorf("split: unsupported args: %s, %s", delimiter, subject))
		}),
		createNestedFuncAttr("join", 2, func(args ...rel.Value) rel.Value {
			joiner, subject := args[0], args[1]
			switch subject := subject.(type) {
			case rel.Array:
				switch subject.Values()[0].(type) {
				case rel.String:
					// if subject is rel.String
					return strJoin(args...)
				case rel.Value:
					if _, isStr := joiner.(rel.String); isStr {
						return strJoin(args...)
					}
					return arrayJoin(joiner, subject)
				}
			case rel.Bytes:
				if _, isSet := joiner.(rel.GenericSet); isSet {
					return subject
				}
				return bytesJoin(joiner.(rel.Bytes), subject)
			case rel.GenericSet:
				switch joiner.(type) {
				case rel.String:
					// if joiner is rel.String
					return strJoin(args...)
				case rel.Array, rel.GenericSet, rel.Bytes:
					return subject
				}
			}

			panic(fmt.Errorf("join: unsupported args: %s, %s", joiner, subject))
		}),
	)
}

func strJoin(args ...rel.Value) rel.Value {
	joiner, subject := args[0], args[1]
	strs := subject.(rel.Set)
	toJoin := make([]string, 0, strs.Count())
	for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
		toJoin = append(toJoin, mustAsString(i.Current()))
	}
	return rel.NewString([]rune(strings.Join(toJoin, mustAsString(joiner))))
}

func asString(val rel.Value) string {
	switch val := val.(type) {
	case rel.Bytes:
		return val.String()
	case rel.Set:
		return mustAsString(val)
	}
	panic("value can't be converted to a string")
}
