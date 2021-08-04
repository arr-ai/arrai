package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

var jsonCommand = &cli.Command{
	Name:      "json",
	Aliases:   []string{"jx"},
	Usage:     "Convert json to arrai",
	UsageText: "Takes json as input from stdin, prints equivalent arrai to stdout",
	Action:    fromJSON,
}

func fromJSON(cli *cli.Context) error {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	val, err := translate.StrictTranslator().ToArrai(data)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(rel.Repr(val))
	return nil
}
