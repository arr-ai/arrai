package syntax

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

const indentStr = "  "

type Enumerable interface {
	ArrayEnumerator() (rel.OffsetValueEnumerator, bool)
}

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
func PrettifyString(val interface{}, indentsNum int) (string, error) {
	switch t := val.(type) {
	case rel.DictEntryTuple:
		key := t.MustGet("@")
		prettyKey, err := PrettifyString(key, indentsNum)
		if err != nil {
			return "", err
		}
		val := t.MustGet(rel.DictValueAttr)
		prettyVal, err := PrettifyString(val, indentsNum)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v: %v", prettyKey, prettyVal), nil
	case rel.Attr: // a: 1
		prettyVal, err := PrettifyString(t.Value, indentsNum)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v: %v", t.Name, prettyVal), nil
	case rel.Tuple: // (a: 1)
		return prettifyTuple(t, indentsNum+1)
	case rel.Array: // [1, 2]
		return prettifyArray(t, indentsNum+1)
	case rel.Dict: // {'a': 1}
		return prettifyDict(t, indentsNum+1)
	case rel.GenericSet: // {1, 2}
		return prettifySet(t, indentsNum+1)
	case rel.String:
		return prettifyString(t)
	case nil:
		return "", nil
	case fmt.Stringer:
		return t.String(), nil
	default:
		return "", fmt.Errorf("unknown type: %T", t)
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

func prettifyDict(dict rel.Dict, indentsNum int) (string, error) {
	vals := make([]rel.Value, dict.Count())
	for i, item := range dict.OrderedEntries() {
		vals[i] = item
	}
	content, err := prettifyItems(vals, indentsNum)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%v}", content), nil
}

func prettifyTuple(tuple rel.Tuple, indentsNum int) (string, error) {
	var sb strings.Builder
	for index, name := range tuple.Names().OrderedNames() {
		str, err := prettifyItem(index, rel.NewAttr(name, tuple.MustGet(name)), false, indentsNum)
		if err != nil {
			return "", err
		}
		sb.WriteString(str)
	}
	if tuple.Count() > 0 {
		sb.WriteString("\n" + getIndents(indentsNum-1))
	}
	return fmt.Sprintf("(%v)", sb.String()), nil
}

func prettifyString(str rel.String) (string, error) {
	return rel.Repr(str), nil
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
func prettifyItem(index int, item interface{}, simple bool, indent int) (string, error) {
	var sb strings.Builder
	formattedStr, err := PrettifyString(item, indent)
	if err != nil {
		return "", err
	}
	if simple && index > 0 {
		sb.WriteString(", ")
	}
	if !simple {
		sb.WriteString("\n" + getIndents(indent))
	}

	_, err = fmt.Fprintf(&sb, "%v", formattedStr)
	if !simple {
		sb.WriteString(",")
	}
	if err != nil {
		return "", err
	}
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

func getIndents(indentsNum int) string {
	var sb strings.Builder
	for i := 0; i < indentsNum; i++ {
		sb.WriteString(indentStr)
	}

	return sb.String()
}
