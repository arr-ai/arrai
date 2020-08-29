package main

import (
	"context"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	t.Parallel()

	var s strings.Builder
	err := testPath(arraictx.InitRunCtx(context.Background()), "./../../examples/test", &s)
	require.Nil(t, err)
	windowsOsStr := `Tests:
..\..\examples\test\multiple_cases_test.arrai
..\..\examples\test\single_case_test.arrai
all tests passed
`
	linuxOsStr := `Tests:
../../examples/test/multiple_cases_test.arrai
../../examples/test/single_case_test.arrai
all tests passed
`
	require.True(t, (windowsOsStr == s.String()) || (linuxOsStr == s.String()))
}
