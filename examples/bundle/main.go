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
	return syntax.EvaluateBundle(arrai.EchoArraiz, args...)
}

func main() {
	val, err := eval(os.Args...)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(val.String())
	}
}
