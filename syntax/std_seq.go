package syntax

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func strJoin(joiner, subject rel.Value) (rel.Value, error) {
	strs := subject.(rel.Set)
	toJoin := make([]string, 0, strs.Count())
	for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
		toJoin = append(toJoin, mustValueAsString(i.Current()))
	}
	if j, is := valueAsString(joiner); is {
		return rel.NewString([]rune(strings.Join(toJoin, j))), nil
	}
	return nil, fmt.Errorf("join: sep not a string: %v", joiner)
}

func stdSeqConcat(seq rel.Value) (rel.Value, error) {
	if set, is := seq.(rel.Set); is {
		if !set.IsTrue() {
			return rel.None, nil
		}
	}
	values := seq.(rel.Array).Values()
	if len(values) == 0 {
		return rel.None, nil
	}
	switch v0 := values[0].(type) {
	case rel.String:
		var sb strings.Builder
		for _, value := range values {
			sb.WriteString(mustValueAsString(value))
		}
		return rel.NewString([]rune(sb.String())), nil
	case rel.Set:
		result := v0
		for _, value := range values[1:] {
			var err error
			result, err = rel.Concatenate(result, value.(rel.Set))
			if err != nil {
				panic(err)
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("concat: incompatible value: %v", values[0])
}

func stdSeqContains(sub, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		if subStr, is := valueAsString(sub); is {
			return rel.NewBool(strings.Contains(subject.String(), subStr)), nil
		}
		return nil, fmt.Errorf("//seq.contains: sub not a string")
	case rel.Array:
		return arrayContains(sub, subject)
	case rel.Bytes:
		if subStr, is := valueAsBytes(sub); is {
			return rel.NewBool(bytes.Contains(subject.Bytes(), subStr)), nil
		}
		return nil, fmt.Errorf("//seq.contains: sub not a byte array")
	case rel.GenericSet:
		emptySet, isSet := sub.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqJoin(joiner, subject rel.Value) (rel.Value, error) {
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
			return subject, nil
		}
		return bytesJoin(joiner.(rel.Bytes), subject), nil
	case rel.GenericSet:
		switch joiner.(type) {
		case rel.String:
			// if joiner is rel.String
			return strJoin(joiner, subject)
		case rel.Array, rel.GenericSet, rel.Bytes:
			return subject, nil
		}
	}

	panic(fmt.Errorf("join: unsupported args: %s, %s", joiner, subject))
}

func stdSeqHasPrefix(prefix, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewBool(strings.HasPrefix(mustValueAsString(subject), mustValueAsString(prefix))), nil
	case rel.Array:
		return arrayHasPrefix(prefix, subject)
	case rel.Bytes:
		return rel.NewBool(bytes.HasPrefix(mustValueAsBytes(subject), mustValueAsBytes(prefix))), nil
	case rel.GenericSet:
		emptySet, isSet := prefix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqHasSuffix(suffix, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewBool(strings.HasSuffix(mustValueAsString(subject), mustValueAsString(suffix))), nil
	case rel.Array:
		return arrayHasSuffix(suffix, subject)
	case rel.Bytes:
		return rel.NewBool(bytes.HasSuffix(mustValueAsBytes(subject), mustValueAsBytes(suffix))), nil
	case rel.GenericSet:
		emptySet, isSet := suffix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqRepeat(arg rel.Value) (rel.Value, error) {
	n := int(arg.(rel.Number))
	return rel.NewNativeFunction("repeat(n)", func(arg rel.Value) (rel.Value, error) {
		switch seq := arg.(type) {
		case rel.String:
			return rel.NewString([]rune(strings.Repeat(seq.String(), n))), nil
		case rel.Array:
			values := []rel.Value{}
			seqValues := seq.Values()
			for i := 0; i < n; i++ {
				values = append(values, seqValues...)
			}
			return rel.NewArray(values...), nil
		case rel.Set:
			if !seq.IsTrue() {
				return rel.None, nil
			}
		}
		return nil, fmt.Errorf("repeat: unsupported value: %v", arg)
	}), nil
}

func stdSeqSub(old, new, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		return rel.NewString(
			[]rune(
				strings.ReplaceAll(
					mustValueAsString(subject),
					mustValueAsString(old),
					mustValueAsString(new),
				),
			),
		), nil
	case rel.Array:
		return arraySub(old, new, subject)
	case rel.Bytes:
		_, oldIsSet := old.(rel.GenericSet)
		_, newIsSet := new.(rel.GenericSet)
		if !oldIsSet && newIsSet {
			return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(), old.String(), ""))), nil
		} else if oldIsSet && !newIsSet {
			return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(), "", new.String()))), nil
		}
		return rel.NewBytes([]byte(strings.ReplaceAll(subject.String(), old.String(), new.String()))), nil
	case rel.GenericSet:
		_, oldIsSet := old.(rel.GenericSet)
		_, newIsSet := new.(rel.GenericSet)
		if oldIsSet && newIsSet {
			return subject, nil
		} else if oldIsSet && !newIsSet {
			return new, nil
		} else if !oldIsSet && newIsSet {
			return subject, nil
		}
	}

	panic(fmt.Errorf("sub: unsupported args: %s, %s, %s", old, new, subject))
}

func stdSeqSplit(delimiter, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		splitted := strings.Split(mustValueAsString(subject), mustValueAsString(delimiter))
		vals := make([]rel.Value, 0, len(splitted))
		for _, s := range splitted {
			vals = append(vals, rel.NewString([]rune(s)))
		}
		return rel.NewArray(vals...), nil
	case rel.Array:
		return arraySplit(delimiter, subject)
	case rel.Bytes:
		return bytesSplit(delimiter, subject), nil
	case rel.GenericSet:
		switch delimiter.(type) {
		case rel.String:
			return rel.NewArray(subject), nil
		case rel.Array, rel.Bytes:
			return rel.NewArray(rel.NewArray()), nil
		case rel.GenericSet:
			return subject, nil
		}
	}

	panic(fmt.Errorf("split: unsupported args: %s, %s", delimiter, subject))
}

func stdSeqTrimPrefix(prefix, subject rel.Value) (rel.Value, error) {
	hasPrefix, err := stdSeqHasPrefix(prefix, subject)
	if err != nil {
		return nil, err
	}
	if hasPrefix.IsTrue() {
		switch subject := subject.(type) {
		case rel.String:
			prefixStr := mustValueAsString(prefix)
			subjectStr := mustValueAsString(subject)
			if strings.HasPrefix(subjectStr, prefixStr) {
				return rel.NewString([]rune(subjectStr[len(prefixStr):])), nil
			}
		case rel.Array:
			return arrayTrimPrefix(prefix, subject)
		case rel.Bytes:
			prefixStr := mustValueAsBytes(prefix)
			subjectStr := mustValueAsBytes(subject)
			if bytes.HasPrefix(subjectStr, prefixStr) {
				return rel.NewBytes(subjectStr[len(prefixStr):]), nil
			}
		}
	}
	return subject, nil
}

func stdSeqTrimSuffix(suffix, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		suffixStr := mustValueAsString(suffix)
		subjectStr := mustValueAsString(subject)
		if strings.HasSuffix(subjectStr, suffixStr) {
			return rel.NewString([]rune(subjectStr[:len(subjectStr)-len(suffixStr)])), nil
		}
	case rel.Array:
		return arrayTrimSuffix(suffix, subject)
	case rel.Bytes:
		suffixStr := mustValueAsBytes(suffix)
		subjectStr := mustValueAsBytes(subject)
		if bytes.HasSuffix(subjectStr, suffixStr) {
			return rel.NewBytes(subjectStr[:len(subjectStr)-len(suffixStr)]), nil
		}
	}
	return subject, nil
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
