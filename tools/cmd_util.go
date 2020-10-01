package tools

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsTerminal() bool {
	return (isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())) ||
		(isatty.IsCygwinTerminal(os.Stdin.Fd()) && isatty.IsCygwinTerminal(os.Stdout.Fd()))
}
