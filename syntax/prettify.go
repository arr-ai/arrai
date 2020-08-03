package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

const indentStr = "  "

// PrettifyString returns a string which represents `rel.Value` with more reabable format.
// For example, `{b: 2, a: 1, c: (a: 2, b: {aa: {bb: (a: 22, d: {3, 1, 2})}})}` is formatted to:
//{
//	b: 2,
//	a: 1,
//	c: (
//		a: 2,
//		b: {
//			aa: {
//				bb: (
//					a: 22,
//					d: {
//						3,
//						1,
//						2
//					}
//				)
//			}
//		}
//	)
//}
func PrettifyString(val rel.Value, indentsNum int) (string, error) {
	indentsNum = indentsNum + 1
	switch t := val.(type) {
	case rel.Tuple: // (a: 1)
		return prettifyTuple(t, indentsNum)
	case rel.Array: // [1, 2]
		return prettifyArray(t, indentsNum)
	case rel.Dict: // {'a': 1}
		return prettifyDict(t, indentsNum)
	case rel.GenericSet: // {1, 2}
		return prettifySet(t, indentsNum)
	case rel.String:
		return prettifyString(t)
	default:
		return t.String(), nil
	}
}

func prettifyTuple(tuple rel.Tuple, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("(")

	for index, name := range tuple.Names().OrderedNames() {
		value, found := tuple.Get(name)
		if !found {
			return "", fmt.Errorf("couldn't find %s", name)
		}
		formattedStr, err := PrettifyString(value, indentsNum)
		if err != nil {
			return "", nil
		}
		fmt.Fprintf(&sb,
			getPrettyFormat(",\n%s%v: %v", index, tuple.Count()), indentsStr, name,
			formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s)", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifyArray(arr rel.Array, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("[")

	for index, item := range arr.Values() {
		formattedStr, err := PrettifyString(item, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, getPrettyFormat(",\n%s%v", index, arr.Count()), indentsStr, formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s]", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifyDict(dict rel.Dict, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("{")

	for index, item := range dict.OrderedEntries() {
		key, found := item.Get("@")
		if !found {
			return "", fmt.Errorf("couldn't find @ in %s", item)
		}
		prettyKey, err := PrettifyString(key, indentsNum)
		if err != nil {
			return "", err
		}
		val, found := item.Get(rel.DictValueAttr)
		if !found {
			return "", fmt.Errorf("couldn't find value in %s", item)
		}
		prettyVal, err := PrettifyString(val, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, getPrettyFormat(",\n%s%v: %v", index, dict.Count()), indentsStr, prettyKey, prettyVal)
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifySet(set rel.GenericSet, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("{")

	for index, item := range set.OrderedValues() {
		formattedStr, err := PrettifyString(item, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, getPrettyFormat(",\n%s%v", index, set.Count()), indentsStr, formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func getPrettyFormat(format string, index, length int) string {
	if index == 0 && length == 1 {
		format = format[1:] + "\n"
	} else if index == 0 && length > 1 {
		format = format[1:]
	} else if index == length-1 && length > 1 {
		format = format + "\n"
	}

	return format
}

func getIndents(indentsNum int) string {
	var sb strings.Builder
	for i := 0; i < indentsNum; i++ {
		sb.WriteString(indentStr)
	}

	return sb.String()
}

func prettifyString(str rel.String) (string, error) {
	return fmt.Sprintf("'%s'", str), nil
}
