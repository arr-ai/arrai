package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "arrai"
	app.Usage = "the ultimate data engine"

	app.Commands = []cli.Command{
		serveCommand,
		evalCommand,
		updateCommand,
		observeCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
