package syntax

import "testing"

// func TestPackageStd(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//.`)
// }

func TestPackageE(t *testing.T) {
	AssertCodesEvalToSameValue(t, `2.718281828459045`, `//math.e`)
}

func TestPackagePi(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//math.pi`)
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

func TestPackageImport(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{./examples/simple/simple}`)
}

func TestPackageImportFromRoot(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{/examples/simple/simple}`)
}

func TestJsonPackageImportFromModuleRoot(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`, `//{/examples/json/foo.json}`)
}

func TestJsonPackageImport(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`, `//{./examples/json/foo.json}`)
}

func TestJsonPackageImportNotExists(t *testing.T) {
	AssertCodeErrors(t, `//{./examples/json/fooooo.json}`, "")
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

func TestPackageExternalImportModule(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `//{github.com/ChloePlanet/arrai-examples/add}`)
}

func TestPackageExternalImportModuleWithExt(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `//{github.com/ChloePlanet/arrai-examples/add.arrai}`)
}

func TestPackageExternalImportModuleJson(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `//{github.com/ChloePlanet/arrai-examples/empty.json}`)
}

func TestPackageExternalImportURLArrai(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `//{https://raw.githubusercontent.com/ChloePlanet/arrai-examples/master/add.arrai}`) // nolint:lll
}

func TestPackageExternalImportURLArraiWithoutHTTPS(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `//{raw.githubusercontent.com/ChloePlanet/arrai-examples/master/add.arrai}`)
}

func TestPackageExternalImportURLJson(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `//{https://raw.githubusercontent.com/ChloePlanet/arrai-examples/master/empty.json}`) // nolint:lll
}

func TestPackageExternalImportURLJsonWithoutHTTPS(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `//{raw.githubusercontent.com/ChloePlanet/arrai-examples/master/empty.json}`)
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
