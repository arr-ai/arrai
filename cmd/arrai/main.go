package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var cmds = []*cli.Command{
	shellCommand,
	runCommand,
	bundleCommand,
	evalCommand,
	compileCommand,
	jsonCommand,
	observeCommand,
	serveCommand,
	syncCommand,
	transformCommand,
	updateCommand,
	infoCommand,
	testCommand,
}

func main() {
	prepareProfilers()

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
		os.Exit(1)
	}
}

func setupVersion(app *cli.App) {
	app.Version = Version

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("arrai %s %s/%s\n", Version, BuildOS, BuildArch)
	}
	// Construct build information here as these data are stored in main package only.
	syntax.BuildInfo = syntax.GetBuildInfo(Version, BuildDate, GitFullCommit, GitTags, BuildOS, BuildArch, GoVersion)
}

func prepareProfilers() {
	args := os.Args[1:]
loop:
	for len(args) >= 1 {
		flag := strings.SplitN(args[0], "=", 2)
		switch flag[0] {
		case "-cpuprofile":
			cpuprofile := flag[1]
			f, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatalf("could not create cpu profile: %v", err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Printf("error closing cpu profile: %v", err)
				}
			}()
			if err := pprof.StartCPUProfile(f); err != nil {
				logrus.Fatal(err)
			}
			defer pprof.StopCPUProfile()
		case "-memprofile":
			memprofile := flag[1]
			defer func() {
				f, err := os.Create(memprofile)
				if err != nil {
					log.Fatalf("could not create memory profile: %v", err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						log.Printf("error closing memory profile: %v", err)
					}
				}()
				runtime.GC()
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}()
		default:
			break loop
		}
		args = args[1:]
	}
	os.Args = append(os.Args[:1], args...)
}
