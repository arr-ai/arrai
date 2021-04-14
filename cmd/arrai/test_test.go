package main

import (
	"context"
	"fmt"
	"github.com/arr-ai/arrai/internal/test"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/require"
)

func TestNewTest(t *testing.T) {
	t.Parallel()

	var s strings.Builder
	err := test.RunTestsInPath(arraictx.InitRunCtx(context.Background()), &s, "./../../examples/test_new")
	require.NotNil(t, err)

	fmt.Println("\n\n*****\n" + s.String() + "*****\n\n")
}

func TestWalkLeaves(t *testing.T) {
	t.Parallel()

	bytes, err := ioutil.ReadFile("./../../examples/test_new/multiple_cases_test.arrai")
	require.NoError(t, err)
	val, err := syntax.EvaluateExpr(context.Background(), "", string(bytes))
	require.NoError(t, err)

	fmt.Print("\n\n*****\n")
	test.ForeachLeaf(val, "", func(val rel.Value, path string) {
		fmt.Println(path + ": " + val.String())
	})
	fmt.Print("*****\n\n")
}

func TestColor(t *testing.T) {
	const text = "%d: \033[38;5;255;%d;1mTEST\033[0m\n"

	for i := 0; i < 100; i++ {
		fmt.Printf(text, i, i)
	}
}
