package rel

import (
	"fmt"
	"strings"
)

// FormatString returns string of `rel.Value` with more reabable format.
// For examples:
// Dict:
//
//
// Tuple:
//
func FormatString(val Value, identsNum int) string {
	identsNum = identsNum + 1
	switch t := val.(type) {
	case Tuple: // (a: 1)
		return formatTupleString(t, identsNum)
	case Dict: // {'a': 1}
		return formatDictString(t, identsNum)
	case Array:
		return t.String()
	case Set: // {1, 2}
		return formatSetString(t, identsNum)
	default:
		return t.String()
	}
}

func formatTupleString(tuple Tuple, identsNum int) string {
	var sb strings.Builder
	identsStr := getIdents(identsNum)
	sb.WriteString("(")
	for n, enum := 0, tuple.Enumerator(); enum.MoveNext(); n++ {
		name, val := enum.Current()
		fmt.Fprintf(&sb, getFormat(n, tuple.Count()), identsStr, name, FormatString(val, identsNum))
	}

	sb.WriteString(fmt.Sprintf("%s)", getIdents(identsNum-1)))
	return sb.String()
}

func formatDictString(dict Dict, identsNum int) string {
	var sb strings.Builder
	identsStr := getIdents(identsNum)
	sb.WriteString("{")
	for n, enum := 0, dict.DictEnumerator(); enum.MoveNext(); n++ {
		key, val := enum.Current()
		fmt.Fprintf(&sb, getFormat(n, dict.Count()), identsStr, key, FormatString(val, identsNum))
	}

	sb.WriteString(fmt.Sprintf("%s}", getIdents(identsNum-1)))
	return sb.String()
}

func formatSetString(set Set, identsNum int) string {
	var sb strings.Builder
	// identsStr := getIdents(identsNum)
	sb.WriteString("{")
	for n, enum := 0, set.Enumerator(); enum.MoveNext(); n++ {
		format := ",\n %v"
		if n == 0 { // first item
			format = format[1:]
		} else if n > 0 && n == set.Count()-1 { // last item
			format = format + "\n"
		}

		enum.Current()
		val := enum.Current()
		fmt.Fprintf(&sb, format, FormatString(val, identsNum))
	}

	sb.WriteString("}")
	return sb.String()
}

func getFormat(index, length int) string {
	format := ",\n%s%v: %v\n"
	if index == 0 && length > 1 {
		format = format[1 : len(format)-1]
	} else if index == 0 && length == 1 {
		format = format[1:]
	}

	return format
}

func getIdents(identsNum int) string {
	var sb strings.Builder
	for i := 0; i < identsNum; i++ {
		sb.WriteString("\t")
	}

	return sb.String()
}
