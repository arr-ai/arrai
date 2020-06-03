package main

import (
	"os"
	"path"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var cmds = []*cli.Command{
	shellCommand,
	runCommand,
	evalCommand,
	jsonCommand,
	observeCommand,
	serveCommand,
	syncCommand,
	transformCommand,
	updateCommand,
}

func main() {
	app := cli.NewApp()
	// logrus.SetLevel(logrus.InfoLevel)

	app.EnableBashCompletion = true
	args := os.Args
	switch path.Base(args[0]) {
	case "ai":
		app.Name = "ai"
		app.Usage = "arr.ai interactive shell"
		app.Action = iShell
	case "ax":
		app.Name = "ax"
		app.Usage = "the ultimate data transformer"
		app.Action = transform
	default:
		app.Name = "arrai"
		app.Usage = "the ultimate data engine"
		app.Commands = cmds
		if len(os.Args) > 1 {
			if execCmd := fetchCommand(args[1:]); execCmd != "" && !isExecCommand(execCmd, app.Commands) {
				tmpArgs := append(make([]string, 0, 1+len(args)), args[0], "run")
				args = append(tmpArgs, args[1:]...)
				syntax.RunOmitted = true
			}
		}
		//nolint:lll
		cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

   or

   {{.HelpName}} [global options] filepath [arguments...]
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	}

	err := app.Run(args)
	if err != nil {
		logrus.Info(err)
		if isTerminal() {
			if _, isContextErr := err.(rel.ContextErr); isContextErr {
				if err = createDebuggerShell(err); err != nil {
					logrus.Info(err)
				}
			}
		} else {
			logrus.Info("unable to start debug shell: standard input is not a terminal")
		}
	}
}

func isTerminal() bool {
	return (isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())) ||
		(isatty.IsCygwinTerminal(os.Stdin.Fd()) && isatty.IsCygwinTerminal(os.Stdout.Fd()))
}
