package syntax

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

func strJoin(joiner, subject rel.Value) (rel.Value, error) {
	strs := subject.(rel.Set)
	toJoin := make([]string, 0, strs.Count())
	index := 0
	for i, ok := strs.(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); index++ {
		s, is := tools.ValueAsString(i.Current())
		if !is {
			return nil, fmt.Errorf("//str.join: array item %d not a string: %v", index, i.Current())
		}
		toJoin = append(toJoin, s)
	}
	if j, is := tools.ValueAsString(joiner); is {
		return rel.NewString([]rune(strings.Join(toJoin, j))), nil
	}
	return nil, fmt.Errorf("join: sep not a string: %v", joiner)
}

func stdSeqConcat(_ context.Context, seq rel.Value) (rel.Value, error) {
	if set, is := seq.(rel.Set); is {
		if !set.IsTrue() {
			return rel.None, nil
		}
	}
	array, is := seq.(rel.Array)
	if !is {
		return nil, fmt.Errorf("//seq.concat: seq not an array: %v", seq)
	}
	values := array.Values()
	if len(values) == 0 {
		return rel.None, nil
	}
	switch v0 := values[0].(type) {
	case rel.String:
		var sb strings.Builder
		for i, value := range values {
			s, is := tools.ValueAsString(value)
			if !is {
				return nil, fmt.Errorf("//str.concat: array item %d not a string: %v", i, value)
			}
			sb.WriteString(s)
		}
		return rel.NewString([]rune(sb.String())), nil
	case rel.Set:
		result := v0
		for _, value := range values[1:] {
			var err error
			result, err = rel.Concatenate(result, value.(rel.Set))
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("concat: incompatible value: %v", values[0])
}

func stdSeqContains(_ context.Context, sub, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		if subStr, is := tools.ValueAsString(sub); is {
			return rel.NewBool(strings.Contains(subject.String(), subStr)), nil
		}
		return nil, fmt.Errorf("//seq.contains: sub not a string")
	case rel.Array:
		return arrayContains(sub, subject)
	case rel.Bytes:
		if subStr, is := tools.ValueAsBytes(sub); is {
			return rel.NewBool(bytes.Contains(subject.Bytes(), subStr)), nil
		}
		return nil, fmt.Errorf("//seq.contains: sub not a byte array")
	case rel.GenericSet:
		emptySet, isSet := sub.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqJoin(_ context.Context, joiner, subject rel.Value) (rel.Value, error) {
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

	return nil, fmt.Errorf("join: unsupported args: %s, %s", joiner, subject)
}

func stdSeqHasPrefix(_ context.Context, prefix, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		if prefixStr, is := tools.ValueAsString(prefix); is {
			return rel.NewBool(strings.HasPrefix(subject.String(), prefixStr)), nil
		}
		return nil, fmt.Errorf("//seq.has_prefix: prefix not a string: %v", prefix)
	case rel.Array:
		return arrayHasPrefix(prefix, subject)
	case rel.Bytes:
		if prefixBytes, is := tools.ValueAsBytes(prefix); is {
			return rel.NewBool(bytes.HasPrefix(subject.Bytes(), prefixBytes)), nil
		}
		return nil, fmt.Errorf("//seq.has_prefix: prefix not a byte array: %v", prefix)
	case rel.GenericSet:
		emptySet, isSet := prefix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqHasSuffix(_ context.Context, suffix, subject rel.Value) (rel.Value, error) { //nolint:dupl
	switch subject := subject.(type) {
	case rel.String:
		if suffixStr, is := tools.ValueAsString(suffix); is {
			return rel.NewBool(strings.HasSuffix(subject.String(), suffixStr)), nil
		}
		return nil, fmt.Errorf("//seq.has_suffix: suffix not a string: %v", suffix)
	case rel.Array:
		return arrayHasSuffix(suffix, subject)
	case rel.Bytes:
		if suffixBytes, is := tools.ValueAsBytes(suffix); is {
			return rel.NewBool(bytes.HasSuffix(subject.Bytes(), suffixBytes)), nil
		}
		return nil, fmt.Errorf("//seq.has_suffix: suffix not a byte array: %v", suffix)
	case rel.GenericSet:
		emptySet, isSet := suffix.(rel.GenericSet)
		return rel.NewBool(isSet && !emptySet.IsTrue()), nil
	}

	return rel.NewBool(false), nil
}

func stdSeqRepeat(_ context.Context, arg rel.Value) (rel.Value, error) {
	n := int(arg.(rel.Number))
	return rel.NewNativeFunction("repeat(n)", func(_ context.Context, arg rel.Value) (rel.Value, error) {
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

func stdSeqSub(_ context.Context, old, new, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		subjectStr := subject.String()
		oldStr, is := tools.ValueAsString(old)
		if !is {
			return nil, fmt.Errorf("//seq.sub: old not a string: %v", old)
		}
		newStr, is := tools.ValueAsString(new)
		if !is {
			return nil, fmt.Errorf("//seq.sub: new not a string: %v", new)
		}
		return rel.NewString([]rune(strings.ReplaceAll(subjectStr, oldStr, newStr))), nil
	case rel.Array:
		return arraySub(old, new, subject)
	case rel.Bytes:
		subjectBytes := subject.Bytes()
		oldBytes, is := tools.ValueAsBytes(old)
		if !is {
			return nil, fmt.Errorf("//seq.sub: old not a byte array: %v", old)
		}
		newBytes, is := tools.ValueAsBytes(new)
		if !is {
			return nil, fmt.Errorf("//seq.sub: new not a byte array: %v", new)
		}
		// TODO: Use a byte-aware implementation, not strings.ReplaceAll.
		return rel.NewBytes(
			[]byte(strings.ReplaceAll(
				string(subjectBytes),
				string(oldBytes),
				string(newBytes),
			)),
		), nil
	case rel.GenericSet:
		_, oldIsSet := old.(rel.GenericSet)
		_, newIsSet := new.(rel.GenericSet)
		if oldIsSet && newIsSet {
			return subject, nil
		} else if oldIsSet && !newIsSet {
			return new, nil
		} else {
			return subject, nil
		}
	}

	return nil, fmt.Errorf("sub: unsupported args: %s, %s, %s", old, new, subject)
}

func stdSeqSplit(_ context.Context, delimiter, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		delimStr, is := tools.ValueAsString(delimiter)
		if !is {
			return nil, fmt.Errorf("//seq.split: delim not a string: %v", delimiter)
		}
		splitted := strings.Split(subject.String(), delimStr)
		vals := make([]rel.Value, 0, len(splitted))
		for _, s := range splitted {
			vals = append(vals, rel.NewString([]rune(s)))
		}
		return rel.NewArray(vals...), nil
	case rel.Array:
		return arraySplit(delimiter, subject)
	case rel.Bytes:
		return bytesSplit(delimiter, subject)
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

	return nil, fmt.Errorf("split: unsupported args: %s, %s", delimiter, subject)
}

func stdSeqTrimPrefix(ctx context.Context, prefix, subject rel.Value) (rel.Value, error) {
	hasPrefix, err := stdSeqHasPrefix(ctx, prefix, subject)
	if err != nil {
		return nil, err
	}
	if hasPrefix.IsTrue() {
		switch subject := subject.(type) {
		case rel.String:
			subjectStr := subject.String()
			if prefixStr, is := tools.ValueAsString(prefix); is {
				return rel.NewString([]rune(strings.TrimPrefix(subjectStr, prefixStr))), nil
			}
			return nil, fmt.Errorf("//seq.trim_prefix: prefix not a string: %v", prefix)
		case rel.Array:
			return arrayTrimPrefix(prefix, subject)
		case rel.Bytes:
			subjectBytes := subject.Bytes()
			if prefixBytes, is := tools.ValueAsBytes(prefix); is {
				if bytes.HasPrefix(subjectBytes, prefixBytes) {
					return rel.NewBytes(subjectBytes[len(prefixBytes):]), nil
				}
				return subject, nil
			}
			return nil, fmt.Errorf("//seq.trim_prefix: prefix not a byte array: %v", prefix)
		}
	}
	return subject, nil
}

func stdSeqTrimSuffix(_ context.Context, suffix, subject rel.Value) (rel.Value, error) {
	switch subject := subject.(type) {
	case rel.String:
		subjectStr := subject.String()
		if suffixStr, is := tools.ValueAsString(suffix); is {
			return rel.NewString([]rune(strings.TrimSuffix(subjectStr, suffixStr))), nil
		}
		return nil, fmt.Errorf("//seq.trim_suffix: suffix not a string: %v", suffix)
	case rel.Array:
		return arrayTrimSuffix(suffix, subject)
	case rel.Bytes:
		subjectBytes := subject.Bytes()
		if suffixBytes, is := tools.ValueAsBytes(suffix); is {
			if bytes.HasSuffix(subjectBytes, suffixBytes) {
				return rel.NewBytes(subjectBytes[:len(subjectBytes)-len(suffixBytes)]), nil
			}
			return subject, nil
		}
		return nil, fmt.Errorf("//seq.trim_suffix: suffix not a byte array: %v", suffix)
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
