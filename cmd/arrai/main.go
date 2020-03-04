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

	if path.Base(os.Args[0]) == "ax" {
		app.Name = "ax"
		app.Usage = "the ultimate data transformer"
		app.Action = transform
	} else {
		app.Name = "arrai"
		app.Usage = "the ultimate data engine"

		app.Commands = []*cli.Command{
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
