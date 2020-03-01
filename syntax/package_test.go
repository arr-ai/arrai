package syntax

import (
	"testing"
)

// func TestPackageStd(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.`)
// }

func TestPackagePi(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.math.pi`)
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
