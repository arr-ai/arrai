package tools

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
)

func IsTerminal() bool {
	return (isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())) ||
		(isatty.IsCygwinTerminal(os.Stdin.Fd()) && isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

//SetArgs set context arguments, for example command line `arrai -d r file.arrai arg1 arg2 arg3` will
//save `file.arrai arg1 arg2 arg3` as a slice to Arguments.
func SetArgs(c *cli.Context) {
	Arguments = c.Args().Slice()
}

var Arguments []string
