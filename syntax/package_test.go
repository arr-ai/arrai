package syntax

import (
	"testing"
)

// func TestPackageStd(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.`)
// }

func TestPackageE(t *testing.T) {
	AssertCodesEvalToSameValue(t, `2.718281828459045`, `//.math.e`)
}

func TestPackagePi(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.math.pi`)
}

// func TestPackageRelativeImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, ``, `//./myutil/work(42)`)
// }

// func TestPackageYaml(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, ``, `//./'myutil/work.yaml'`)
// }

// func TestPackageRootImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `///`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
