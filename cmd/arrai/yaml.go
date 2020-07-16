package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
	"github.com/urfave/cli/v2"
)

var yamlCommand = &cli.Command{
	Name:      "yaml",
	Aliases:   []string{"yx"},
	Usage:     "Convert yaml to arrai",
	UsageText: "Takes yaml as input from stdin, prints equivalent arrai to stdout",
	Action:    fromYAML,
}

func fromYAML(cli *cli.Context) error {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	val, err := translate.BytesYamlToArrai(raw)
	if err != nil {
		return err
	}
	fmt.Println(rel.Repr(val))
	return nil
}
