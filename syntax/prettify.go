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
//					d: {3, 1, 2}
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
	case nil:
		return "", nil
	default:
		return t.String(), nil
	}
}

func prettifySet(arr rel.GenericSet, indentsNum int) (string, error) {
	content, err := prettifyItems(arr.OrderedValues(), indentsNum)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%v}", content), nil
}

func prettifyArray(arr rel.Array, indentsNum int) (string, error) {
	content, err := prettifyItems(arr.Values(), indentsNum)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("[%v]", content), nil
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
		format := getPrettyFormat(",\n%s%v: %v", index, tuple.Count())
		_, err = fmt.Fprintf(&sb, format, indentsStr, name, formattedStr)
		if err != nil {
			return "", nil
		}
	}

	sb.WriteString(fmt.Sprintf("%s)", getIndents(indentsNum-1)))
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
		format := getPrettyFormat(",\n%s%v: %v", index, dict.Count())
		_, err = fmt.Fprintf(&sb, format, indentsStr, prettyKey, prettyVal)
		if err != nil {
			return "", err
		}
	}

	sb.WriteString(fmt.Sprintf("%s}", getIndents(indentsNum-1)))
	return sb.String(), nil
}

func prettifyString(str rel.String) (string, error) {
	return rel.Repr(str), nil
}

type Enumerable interface {
	ArrayEnumerator() (rel.OffsetValueEnumerator, bool)
}

// prettifyItems returns a pretty string representation of the contents of a set or array.
func prettifyItems(vals []rel.Value, indentsNum int) (string, error) {
	var sb strings.Builder
	simple := isSimpleValues(vals)

	for index, item := range vals {
		str, err := prettifyItem(index, item, simple, indentsNum)
		if err != nil {
			return "", err
		}
		sb.WriteString(str)
	}
	if !simple {
		sb.WriteString("\n" + getIndents(indentsNum-1))
	}
	return sb.String(), nil
}

// prettifyItem returns the pretty string for an item at an index within a collection.
func prettifyItem(index int, item rel.Value, simple bool, indent int) (string, error) {
	var sb strings.Builder
	formattedStr, err := PrettifyString(item, indent)
	if err != nil {
		return "", err
	}
	if index > 0 {
		sb.WriteString(",")
		if simple {
			sb.WriteString(" ")
		}
	}
	if !simple {
		sb.WriteString("\n" + getIndents(indent))
	}

	fmt.Fprintf(&sb, "%v", formattedStr)
	return sb.String(), nil
}

// isSimple returns true if the value should be pretty printed on a single line.
func isSimple(val rel.Value) bool {
	switch t := val.(type) {
	case rel.Number, rel.String, nil:
		return true
	case rel.Array:
		return isSimpleValues(t.Values())
	case rel.GenericSet:
		return isSimpleValues(t.OrderedValues())
	case rel.Dict:
		return t.Count() == 0
	}
	return false
}

// isSimpleEnumerator returns true if all of the enumerated values should be pretty printed on a
// single line.
func isSimpleValues(vals []rel.Value) bool {
	for _, item := range vals {
		if !isSimple(item) {
			return false
		}
	}
	return true
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
