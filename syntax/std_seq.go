package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func strJoin(joiner, subject rel.Value) rel.Value {
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

func stdSeqContains(sub, subject rel.Value) rel.Value { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewBool(strings.Contains(mustAsString(subject), mustAsString(sub)))
	case rel.Array:
		return arrayContains(sub, subject)
	case rel.Bytes:
		return rel.NewBool(strings.Contains(asString(subject), asString(sub)))
	case rel.GenericSet:
		emptySet, isSet := sub.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue())
	}

	return rel.NewBool(false)
}

func stdSeqJoin(joiner, subject rel.Value) rel.Value {
	switch subject := subject.(type) {
	case rel.Array:
		switch subject.Values()[0].(type) {
		case rel.String:
			// if subject is rel.String
			return strJoin(joiner, subject)
		case rel.Value:
			if _, isStr := joiner.(rel.String); isStr {
				return strJoin(joiner, subject)
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
			return strJoin(joiner, subject)
		case rel.Array, rel.GenericSet, rel.Bytes:
			return subject
		}
	}

	panic(fmt.Errorf("join: unsupported args: %s, %s", joiner, subject))
}

func stdSeqHasPrefix(prefix, subject rel.Value) rel.Value {
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewBool(strings.HasPrefix(mustAsString(subject), mustAsString(prefix)))
	case rel.Array:
		return arrayHasPrefix(prefix, subject)
	case rel.Bytes:
		return rel.NewBool(strings.HasPrefix(asString(subject), asString(prefix)))
	case rel.GenericSet:
		emptySet, isSet := prefix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue())
	}

	return rel.NewBool(false)
}

func stdSeqHasSuffix(suffix, subject rel.Value) rel.Value {
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewBool(strings.HasSuffix(mustAsString(subject), mustAsString(suffix)))
	case rel.Array:
		return arrayHasSuffix(suffix, subject)
	case rel.Bytes:
		return rel.NewBool(strings.HasSuffix(asString(subject), asString(suffix)))
	case rel.GenericSet:
		emptySet, isSet := suffix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue())
	}

	return rel.NewBool(false)
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

func stdSeqSub(old, new, subject rel.Value) rel.Value {
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
}

func stdSeqSplit(delimiter, subject rel.Value) rel.Value {
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
}

func stdSeqTrimPrefix(prefix, subject rel.Value) rel.Value {
	if stdSeqHasPrefix(prefix, subject).IsTrue() {
		switch subject := subject.(type) {
		case rel.String:
			prefixStr := mustAsString(prefix)
			subjectStr := mustAsString(subject)
			if strings.HasPrefix(subjectStr, prefixStr) {
				return rel.NewString([]rune(subjectStr[len(prefixStr):]))
			}
		case rel.Array:
			return arrayTrimPrefix(prefix, subject)
		case rel.Bytes:
			prefixStr := mustAsBytes(prefix)
			subjectStr := mustAsBytes(subject)
			if strings.HasPrefix(subjectStr, prefixStr) {
				return rel.NewBytes([]byte(subjectStr[len(prefixStr):]))
			}
		}
	}
	return subject
}

func stdSeqTrimSuffix(suffix, subject rel.Value) rel.Value {
	switch subject := subject.(type) {
	case rel.String:
		suffixStr := mustAsString(suffix)
		subjectStr := mustAsString(subject)
		if strings.HasSuffix(subjectStr, suffixStr) {
			return rel.NewString([]rune(subjectStr[:len(subjectStr)-len(suffixStr)]))
		}
	case rel.Array:
		return arrayTrimSuffix(suffix, subject)
	case rel.Bytes:
		suffixStr := mustAsBytes(suffix)
		subjectStr := mustAsBytes(subject)
		if strings.HasSuffix(subjectStr, suffixStr) {
			return rel.NewBytes([]byte(subjectStr[:len(subjectStr)-len(suffixStr)]))
		}
	}
	return subject
}

func stdSeq() rel.Attr {
	return rel.NewTupleAttr("seq",
		rel.NewNativeFunctionAttr("concat", stdSeqConcat),
		createFunc2Attr("contains", stdSeqContains),
		createFunc2Attr("has_prefix", stdSeqHasPrefix),
		createFunc2Attr("has_suffix", stdSeqHasSuffix),
		createFunc2Attr("join", stdSeqJoin),
		rel.NewNativeFunctionAttr("repeat", stdSeqRepeat),
		createFunc3Attr("sub", stdSeqSub),
		createFunc2Attr("split", stdSeqSplit),
		createFunc2Attr("trim_prefix", stdSeqTrimPrefix),
		createFunc2Attr("trim_suffix", stdSeqTrimSuffix),
	)
}
