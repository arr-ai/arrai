package rel

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/arr-ai/arrai/pkg/fu"
)

var reprEscapes = func() [][]byte {
	escapes := make([][]byte, 32)
	escapes['\a'] = []byte(`\a`)
	escapes['\b'] = []byte(`\b`)
	escapes['\x1b'] = []byte(`\e`)
	escapes['\f'] = []byte(`\f`)
	escapes['\n'] = []byte(`\n`)
	escapes['\r'] = []byte(`\r`)
	escapes['\t'] = []byte(`\t`)
	escapes['\v'] = []byte(`\v`)
	return escapes
}()

// Build a slice of byte reprs: [][]byte{[]byte('0'), ..., []byte('255')}.
var byteReprs = func() [][]byte {
	buf := bytes.NewBuffer(make([]byte, 0, 3*256))
	offsets := append(make([]int, 0, 257), 0)
	for i := 0; i < 256; i++ {
		fu.Fprint(buf, i)
		offsets = append(offsets, buf.Len())
	}
	data := buf.Bytes()
	ret := make([][]byte, 0, 256)
	for i, end := range offsets[1:] {
		begin := offsets[i]
		ret = append(ret, data[begin:end])
	}
	return ret
}()

func writeSep(w io.Writer, i int, sep string) { //nolint:unparam
	if i > 0 {
		fu.WriteString(w, sep)
	}
}

func reprOffset(offset int, w io.Writer) {
	if offset != 0 {
		fu.Fprintf(w, `%d\`, offset)
	}
}

func reprEscape(s string, delim byte, w io.Writer) {
	fu.Fprintf(w, "%c", delim)
	// TODO: optimise by handling non-special characters in groups.
	for _, c := range s {
		if c == '\\' || c == rune(delim) {
			fu.Fprintf(w, `\%c`, c)
		} else if c >= 32 {
			fu.Fprintf(w, "%c", c)
		} else if escape := reprEscapes[c]; escape != nil {
			fu.Write(w, escape)
		} else {
			fu.Fprintf(w, `\x%02x`, c)
		}
	}
	fu.Fprintf(w, "%c", delim)
}

func reprString(str String, w io.Writer) {
	reprOffset(str.offset, w)
	reprStr(string(str.s), w)
}

func reprStr(s string, w io.Writer) {
	switch {
	case !strings.Contains(s, "'"):
		reprEscape(s, '\'', w)
	default:
		reprEscape(s, '"', w)
	}
}

// OrderableSet is a type used to repr GenericSet and UnionSet to avoid duplications.
type OrderableSet interface {
	Set
	OrderedValues() ValueEnumerator
}

func reprOrderableSet(w io.Writer, s OrderableSet) {
	if s.Equal(True) {
		fmt.Fprint(w, sTrue)
		return
	}
	fmt.Fprint(w, "{")
	for i, o := s.OrderedValues(), 0; i.MoveNext(); o++ {
		writeSep(w, o, ", ")
		fu.FRepr(w, i.Current())
	}
	fmt.Fprint(w, "}")
}
