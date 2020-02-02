package syntax

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestPackageStd(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.`)
// }

func TestPackagePi(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.math.pi`)
}

func TestPackageLocalImport(t *testing.T) {
	filename := "rel/value_set_rel.go"
	data, err := ioutil.ReadFile("../" + filename)
	require.NoError(t, err)
	AssertCodesEvalToSameValue(t,
		"'"+strings.ReplaceAll(string(data), "\n", `\n`)+"'",
		"//./'"+filename+"'")
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
