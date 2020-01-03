package parse

import (
	"fmt"
)

type Node struct {
	Tag      string
	Extra    interface{}
	Children []interface{}
}

func (n Node) String() string {
	return fmt.Sprintf("%s", n) //nolint:gosimple
}

func (n Node) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%s(", n.Tag)
	format := "%" + string(c)
	for i, child := range n.Children {
		if i > 0 {
			fmt.Fprint(state, ", ")
		}
		fmt.Fprintf(state, format, child)
	}
	fmt.Fprint(state, ")")
}

type Parser interface {
	Parse(input *Scanner, output interface{}) bool
}

func PtrAssign(output, input interface{}) {
	*output.(*interface{}) = input
}

type Func func(input *Scanner, output interface{}) bool

func (f Func) Parse(input *Scanner, output interface{}) bool {
	return f(input, output)
}

func Transform(parser Parser, transform func(Node) Node) Parser {
	return Func(func(input *Scanner, output interface{}) bool {
		var v Node
		if parser.Parse(input, &v) {
			PtrAssign(output, transform(v))
			return true
		}
		return false
	})
}
