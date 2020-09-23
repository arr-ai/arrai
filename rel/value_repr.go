package rel

import (
	"fmt"
	"io"
	"strings"
)

var reprEscapes = escapesWithDelim()

func escapesWithDelim() []byte {
	escapes := make([]byte, 32)
	escapes['\a'] = 'a'
	escapes['\b'] = 'b'
	escapes['\x1b'] = 'e'
	escapes['\f'] = 'f'
	escapes['\n'] = 'n'
	escapes['\r'] = 'r'
	escapes['\t'] = 't'
	escapes['\v'] = 'v'
	return escapes
}

type reprCommaSep bool

func (s *reprCommaSep) Sep(w io.Writer) {
	if !*s {
		*s = true
	} else {
		fmt.Fprint(w, ", ")
	}
}

func reprOffset(offset int, w io.Writer) {
	if offset != 0 {
		fmt.Fprintf(w, `%d\`, offset)
	}
}

func reprEscape(s string, delim byte, w io.Writer) {
	fmt.Fprintf(w, "%c", delim)
	// TODO: optimise by handling non-special characters in groups.
	for _, c := range s {
		if c == '\\' || c == rune(delim) {
			fmt.Fprintf(w, `\%c`, c)
		} else if c >= 32 {
			fmt.Fprintf(w, "%c", c)
		} else if escape := reprEscapes[c]; escape != 0 {
			fmt.Fprintf(w, `\%c`, escape)
		} else {
			fmt.Fprintf(w, `\x%02x`, c)
		}
	}
	fmt.Fprintf(w, "%c", delim)
}

func reprString(str String, w io.Writer) {
	reprOffset(str.offset, w)
	s := string(str.s)
	switch {
	case !strings.Contains(s, "'"):
		reprEscape(s, '\'', w)
	default:
		reprEscape(s, '"', w)
	}
}

func reprBytes(b Bytes, w io.Writer) {
	reprOffset(b.offset, w)
	fmt.Fprint(w, "<<")
	s := string(b.b)
	switch {
	case !strings.Contains(s, "'"):
		reprEscape(s, '\'', w)
	default:
		reprEscape(s, '"', w)
	}
	fmt.Fprint(w, ">>")
}

func reprArray(a Array, w io.Writer) {
	reprOffset(a.offset, w)
	fmt.Fprint(w, "[")
	var sep reprCommaSep
	for _, v := range a.values {
		sep.Sep(w)
		if v != nil {
			reprValue(v, w)
		}
	}
	fmt.Fprint(w, "]")
}

func reprDict(d Dict, w io.Writer) {
	fmt.Fprint(w, "{")
	var sep reprCommaSep
	for _, e := range d.OrderedEntries() {
		sep.Sep(w)
		reprValue(e.at, w)
		fmt.Fprint(w, ": ")
		reprValue(e.value, w)
	}
	fmt.Fprint(w, "}")
}

func reprSet(s GenericSet, w io.Writer) {
	if s.Equal(True) {
		fmt.Fprintf(w, "true")
		return
	}
	fmt.Fprint(w, "{")
	var sep reprCommaSep
	for _, v := range s.OrderedValues() {
		sep.Sep(w)
		reprValue(v, w)
	}
	fmt.Fprint(w, "}")
}

func reprClosure(c Closure, w io.Writer) {
	fmt.Fprintf(w, "%s", c.String())
}

func reprStringCharTuple(t StringCharTuple, w io.Writer) {
	fmt.Fprintf(w, "(@: %d, %s: %d)", t.at, StringCharAttr, t.char)
}

func reprArrayItemTuple(t ArrayItemTuple, w io.Writer) {
	fmt.Fprintf(w, "(@: %d, %s: ", t.at, ArrayItemAttr)
	reprValue(t.item, w)
	fmt.Fprint(w, ")")
}

func reprDictEntryTuple(t DictEntryTuple, w io.Writer) {
	fmt.Fprint(w, "(@: ")
	reprValue(t.at, w)
	fmt.Fprintf(w, ", %s: ", DictValueAttr)
	reprValue(t.value, w)
	fmt.Fprint(w, ")")
}

func reprTuple(t *GenericTuple, w io.Writer) {
	fmt.Fprint(w, "(")
	var sep reprCommaSep
	for _, name := range TupleOrderedNames(t) {
		sep.Sep(w)
		fmt.Fprintf(w, "%s: ", TupleNameRepr(name))
		reprValue(t.MustGet(name), w)
	}
	fmt.Fprint(w, ")")
}

func reprNumber(n Number, w io.Writer) {
	fmt.Fprint(w, n.String())
}

func reprNativeFunction(v Value, w io.Writer) {
	fmt.Fprintf(w, "native function %s", v.String())
}

func reprValue(v Value, w io.Writer) {
	switch v := v.(type) {
	case String:
		reprString(v, w)
	case Bytes:
		reprBytes(v, w)
	case Array:
		reprArray(v, w)
	case Dict:
		reprDict(v, w)
	case GenericSet:
		reprSet(v, w)
	case Closure:
		reprClosure(v, w)
	case StringCharTuple:
		reprStringCharTuple(v, w)
	case ArrayItemTuple:
		reprArrayItemTuple(v, w)
	case DictEntryTuple:
		reprDictEntryTuple(v, w)
	case *GenericTuple:
		reprTuple(v, w)
	case Number:
		reprNumber(v, w)
	case *NativeFunction:
		reprNativeFunction(v, w)
	default:
		panic(fmt.Errorf("Repr(): unexpected Value type %s: %[1]v", ValueTypeAsString(v))) //nolint:golint
	}
}

func Repr(v Value) string {
	var sb strings.Builder
	reprValue(v, &sb)
	return sb.String()
}
