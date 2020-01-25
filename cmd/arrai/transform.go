package main

import (
	"github.com/urfave/cli"
)

var transformCommand = cli.Command{
	Name:    "transform",
	Aliases: []string{"x"},
	Usage:   "transform a stream of input data with an expression",
	Action:  transform,
}

func transform(c *cli.Context) error {
	panic("unfinished")
	// source := c.Args().Get(0)

	// expr, err := syntax.Parse(parser.NewScanner(source))
	// if err != nil {
	// 	return err
	// }

	// stream := parser.NewScanner(os.Stdin)
	// for {
	// 	if stream.Peek() == syntax.EOF {
	// 		break
	// 	}
	// 	value, err := syntax.Parse(stream)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	global := rel.EmptyScope.With(".", value)
	// 	xvalue, err := expr.Eval(global, global)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	s := xvalue.String()
	// 	os.Stdout.WriteString(s)
	// 	if s[len(s)-1] != '\n' {
	// 		os.Stdout.Write([]byte{'\n'})
	// 	}
	// }

	return nil
}
