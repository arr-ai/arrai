package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arr-ai/arrai/examples/bundle/internal/arrai"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

func eval(args ...string) (rel.Value, error) {
	bundle, err := arrai.Asset("internal/arrai/echo.arraiz")
	if err != nil {
		return nil, err
	}

	return syntax.EvaluateBundle(bundle, args...)
}

func main() {
	val, err := eval(os.Args...)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(val.String())
	}
}
