package main

import (
	"fmt"
	"os"
	"path"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
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
	infoCommand,
	yamlCommand,
}

func main() {
	app := cli.NewApp()
	// logrus.SetLevel(logrus.InfoLevel)

	var debug bool

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
		app.Flags = []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "When evaluation fails, arrai will drop into the interactive shell with the last scope",
				Destination: &debug,
				EnvVars:     []string{"ARRAI_DEBUG"},
			},
		}
		if len(os.Args) > 1 {
			if execCmd := fetchCommand(args[1:]); execCmd != "" && !isExecCommand(execCmd, app.Commands) {
				args = insertRunCommand(app.Flags, os.Args)
				syntax.RunOmitted = true
			}
		}

		setupVersion(app)

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
		if tools.IsTerminal() {
			if _, isContextErr := err.(rel.ContextErr); isContextErr && debug {
				if err = createDebuggerShell(err); err != nil {
					logrus.Info(err)
				}
			}
		} else {
			logrus.Info("unable to start debug shell: standard input is not a terminal")
		}
	}
}

func setupVersion(app *cli.App) {
	app.Version = Version

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("arrai %s %s\n", Version, BuildOS)
	}
}
