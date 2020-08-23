package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	t.Parallel()

	var s strings.Builder
	err := testPath("./../../examples/test", &s)
	require.Nil(t, err)
	require.Equal(t,
		`Tests:
../../examples/test/multiple_cases_test.arrai
../../examples/test/single_case_test.arrai
2/2 tests passed
all tests passed
`, s.String())
}
