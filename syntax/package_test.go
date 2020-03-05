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

func TestPackageExternalImportModule(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3`, `//github.com/ChloePlanet/'arrai-examples'/add`)
}

func TestPackageExternalImportModuleWithExt(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3`, `//github.com/ChloePlanet/'arrai-examples'/'add.arrai'`)
}

func TestPackageExternalImportModuleJson(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `//github.com/ChloePlanet/'arrai-examples'/'empty.json'`)
}

func TestPackageExternalImportURLArrai(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3`, `//https://raw.githubusercontent.com/ChloePlanet/'arrai-examples'/master/'add.arrai'`)
}

func TestPackageExternalImportURLArraiWithoutHTTPS(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3`, `//raw.githubusercontent.com/ChloePlanet/'arrai-examples'/master/'add.arrai'`)
}

func TestPackageExternalImportURLJson(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `//https://jsonplaceholder.typicode.com/todos/'1'/userId`)
}

func TestPackageExternalImportURLJsonWithoutHTTPS(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `//jsonplaceholder.typicode.com/todos/'1'/userId`)
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
