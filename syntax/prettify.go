package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

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
		return prettifyTupleString(t, indentsNum)
	case rel.Dict: // {'a': 1}
		return prettifyDictString(t, indentsNum)
	case rel.GenericSet: // {1, 2}
		return prettifySetString(t, indentsNum)
	default:
		return t.String(), nil
	}
}

func prettifyTupleString(tuple rel.Tuple, indentsNum int) (string, error) {
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
			getPrettyFormat(index, tuple.Count()), indentsStr, name,
			formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s)", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifyDictString(dict rel.Dict, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("{")

	for index, item := range dict.OrderedEntries() {
		key, found := item.Get("@")
		if !found {
			return "", fmt.Errorf("couldn't find @ in %s", item)
		}
		val, found := item.Get(rel.DictValueAttr)
		if !found {
			return "", fmt.Errorf("couldn't find value in %s", item)
		}
		formattedStr, err := PrettifyString(val, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, getPrettyFormat(index, dict.Count()), indentsStr, key, formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifySetString(set rel.GenericSet, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("{")

	for index, item := range set.OrderedValues() {
		format := ",\n%s%v"
		if index == 0 && set.Count() == 1 {
			format = format[1:] + "\n"
		} else if index == 0 && set.Count() > 1 {
			format = format[1:]
		} else if index == set.Count()-1 && set.Count() > 1 {
			format = format + "\n"
		}

		formattedStr, err := PrettifyString(item, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, format, indentsStr, formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func getPrettyFormat(index, length int) string {
	format := ",\n%s%v: %v"
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

const indentStr = "  "
