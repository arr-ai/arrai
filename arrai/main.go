package main

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var completionCommand = cli.Command{
	Name:    "complete",
	Aliases: []string{"c"},
	Usage:   "complete a task on the list",
	Action: func(c *cli.Context) error {
		fmt.Println("completed task: ", c.Args().First())
		return nil
	},
}

func main() {
	app := cli.NewApp()

	app.EnableBashCompletion = true

	if path.Base(os.Args[0]) == "ax" {
		app.Name = "ax"
		app.Usage = "the ultimate data transformer"
		app.Action = transform
	} else {
		app.Name = "arrai"
		app.Usage = "the ultimate data engine"

		app.Commands = []cli.Command{
			completionCommand,
			transformCommand,
			evalCommand,
			serveCommand,
			observeCommand,
			updateCommand,
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
