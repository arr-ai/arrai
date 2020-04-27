package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	// logrus.SetLevel(logrus.InfoLevel)

	app.EnableBashCompletion = true

	switch path.Base(os.Args[0]) {
	case "ai":
		app.Name = "ai"
		app.Usage = "arr.ai interactive shell"
		app.Action = shell
	case "ax":
		app.Name = "ax"
		app.Usage = "the ultimate data transformer"
		app.Action = transform
	default:
		app.Name = "arrai"
		app.Usage = "the ultimate data engine"

		app.Commands = []*cli.Command{
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
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
