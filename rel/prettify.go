package rel

import (
	"fmt"
	"strings"
)

const indentStr = "  "

// PrettifyString returns a string which represents `Value` with more reabable format.
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
func PrettifyString(val Value, indentsNum int) (string, error) {
	indentsNum = indentsNum + 1
	switch t := val.(type) {
	case Tuple: // (a: 1)
		return prettifyTuple(t, indentsNum)
	case Dict: // {'a': 1}
		return prettifyDict(t, indentsNum)
	case GenericSet: // {1, 2}
		return prettifySet(t, indentsNum)
	case Array:
		return prettifyArray(t, indentsNum)
	case String:
		return prettifyString(t, indentsNum)
	default:
		return t.String(), nil
	}
}

func prettifyTuple(tuple Tuple, indentsNum int) (string, error) {
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

func prettifyDict(dict Dict, indentsNum int) (string, error) {
	var sb strings.Builder
	indentsStr := getIndents(indentsNum)
	sb.WriteString("{")

	for index, item := range dict.OrderedEntries() {
		key, found := item.Get("@")
		if !found {
			return "", fmt.Errorf("couldn't find @ in %s", item)
		}
		val, found := item.Get(DictValueAttr)
		if !found {
			return "", fmt.Errorf("couldn't find value in %s", item)
		}
		formattedStr, err := PrettifyString(val, indentsNum)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&sb, getPrettyFormat(",\n%s%v: %v", index, dict.Count()), indentsStr, key, formattedStr)
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifySet(set GenericSet, indentsNum int) (string, error) {
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

func prettifyString(str String, indentsNum int) (string, error) {
	return fmt.Sprintf("'%s'", str), nil
}

// TODO: will change for #171 still
func prettifyArray(array Array, indentsNum int) (string, error) {
	var sb strings.Builder
	fmt.Fprint(&sb, "[")
	var sep reprCommaSep
	for _, v := range array.values {
		sep.Sep(&sb)
		if !isTrivialType(v) {
			return array.String(), nil
		}
		if v != nil {
			reprValue(v, &sb)
		}
	}
	fmt.Fprint(&sb, "]")

	return sb.String(), nil
}

// TODO: will change for #171 still
func isTrivialType(value Value) bool {
	switch value.(type) {
	case Number, String:
		return true
	default:
		return false
	}
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
