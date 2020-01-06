package parse

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Node struct {
	Tag      string
	Extra    interface{}
	Children []interface{}
}

func NewNode(tag string, extra interface{}, children ...interface{}) *Node {
	return &Node{Tag: tag, Extra: extra, Children: children}
}

func (n Node) Count() int {
	return len(n.Children)
}

func (n Node) Get(path ...int) interface{} {
	var v interface{} = n
	for _, i := range path {
		v = v.(Node).Children[i]
	}
	return v
}

func (n Node) GetNode(path ...int) Node {
	return n.Get(path...).(Node)
}

func (n Node) GetScanner(path ...int) Scanner {
	return n.Get(path...).(Scanner)
}

func (n Node) GetString(path ...int) string {
	return n.GetScanner(path...).String()
}

func (n Node) String() string {
	return fmt.Sprintf("%s", n) //nolint:gosimple
}

func (n Node) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%s", n.Tag)
	format := "%" + string(c)
	if n.Extra != nil {
		fmt.Fprintf(state, "║"+format, n.Extra)
	}
	state.Write([]byte{'('})
	for i, child := range n.Children {
		if i > 0 {
			fmt.Fprint(state, ", ")
		}
		fmt.Fprintf(state, format, child)
	}
	fmt.Fprint(state, ")")
}

type Parser interface {
	Parse(input, furthest *Scanner, output interface{}) bool
}

func PtrAssign(output, input interface{}) {
	*output.(*interface{}) = input
}

type Func func(input, furthest *Scanner, output interface{}) bool

func (f Func) Parse(input, furthest *Scanner, output interface{}) bool {
	return f(input, furthest, output)
}

func Transform(parser Parser, transform func(Node) Node) Parser {
	return Func(func(input, furthest *Scanner, output interface{}) bool {
		var v Node
		if parser.Parse(input, furthest, &v) {
			PtrAssign(output, transform(v))
			return true
		}
		return false
	})
}

type NodeDiff struct {
	A, B     *Node
	Children map[int]NodeDiff
	Types    map[int][2]reflect.Type
}

func (d NodeDiff) String() string {
	var sb strings.Builder
	d.report(nil, &sb)
	return sb.String()
}

func (d NodeDiff) report(path []string, w io.Writer) {
	if d.Equal() {
		return
	}
	prefix := ""
	if len(path) > 0 {
		prefix = fmt.Sprintf("[%s] ", strings.Join(path, "."))
	}
	if d.A.Tag != d.B.Tag {
		fmt.Fprintf(w, "%sTag: %v != %v", prefix, d.A.Tag, d.B.Tag)
	}
	if d.A.Extra != d.B.Extra {
		fmt.Fprintf(w, "%sExtra: %v != %v", prefix, d.A.Extra, d.B.Extra)
	}
	if len(d.A.Children) != len(d.B.Children) {
		fmt.Fprintf(w, "%slen(Children): %v != %v", prefix, len(d.A.Children), len(d.B.Children))
	}
	for i, t := range d.Types {
		fmt.Fprintf(w, "%sChildren[%d].Type: %v != %v", prefix, i, t[0], t[1])
	}
	for i, d := range d.Children {
		d.report(append(append([]string{}, path...), fmt.Sprintf("%s[%d]", d.A.Tag, i)), w)
	}
}

func (d NodeDiff) Equal() bool {
	return len(d.A.Children) == len(d.B.Children) &&
		d.A.Tag == d.B.Tag &&
		d.A.Extra == d.B.Extra &&
		len(d.Children) == 0 &&
		len(d.Types) == 0
}

func NewNodeDiff(a, b *Node) NodeDiff {
	children := map[int]NodeDiff{}
	types := map[int][2]reflect.Type{}
	n := len(a.Children)
	if n > len(b.Children) {
		n = len(b.Children)
	}
	for i, x := range a.Children[:n] {
		aType := reflect.TypeOf(a)
		bType := reflect.TypeOf(a)
		if aType != bType {
			types[i] = [2]reflect.Type{aType, bType}
		} else if childNodeA, ok := x.(Node); ok {
			childNodeB := b.Children[i].(Node)
			if d := NewNodeDiff(&childNodeA, &childNodeB); !d.Equal() {
				children[i] = d
			}
		}
	}
	return NodeDiff{
		A:        a,
		B:        b,
		Children: children,
		Types:    types,
	}
}
