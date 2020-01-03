package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/rel"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var jsonCommand = cli.Command{
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
	val, err := jsonToArrai(data)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(val)
	return nil
}

func jsonToArrai(data interface{}) (rel.Value, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return jsonObjToArrai(v)
	case []interface{}:
		return jsonArrToArrai(v)
	case string: // rel.NewValue cannot produce strings
		return rel.NewString([]rune(v)), nil
	default:
		return rel.NewValue(v)
	}
}

func jsonObjToArrai(data map[string]interface{}) (rel.Value, error) {
	tuples := make([]rel.Value, len(data))
	i := 0
	for key, val := range data {
		item, err := jsonToArrai(val)
		if err != nil {
			return nil, err
		}
		tuples[i] = rel.NewTuple(
			rel.Attr{"@", rel.NewString([]rune(key))},
			rel.Attr{"@item", item},
		)
		i++
	}
	return rel.NewSet(tuples...), nil
}

func jsonArrToArrai(data []interface{}) (rel.Value, error) {
	elts := make([]rel.Value, len(data))
	for i, val := range data {
		elt, err := jsonToArrai(val)
		if err != nil {
			return nil, err
		}
		elts[i] = elt
	}
	return rel.NewArray(elts...), nil
}
