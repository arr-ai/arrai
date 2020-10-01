package syntax

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageE(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2.718281828459045`, `//math.e`)
}

func TestPackagePi(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//math.pi`)
}

// func TestPackageRelativeImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, ``, `//{./myutil/work}(42)`)
// }

// func TestPackageYaml(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, ``, `//{./'myutil/work.yaml}'`)
// }

// func TestPackageRootImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `///`)
// }

func TestPackageImport(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{./examples/simple/simple}        `)
	AssertCodesEvalToSameValue(t, `1            `, `//{/examples/import/relative_import}`)
	AssertCodesEvalToSameValue(t, `2            `, `//{/examples/import/comb_import}    `)
}

func TestPackageImportFromRoot(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{/examples/simple/simple}       `)
	AssertCodesEvalToSameValue(t, `1            `, `//{/examples/import/module_import}`)
}

func TestJsonPackageImportFromModuleRoot(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`, `//{/examples/json/foo.json}`)
}

func TestJsonPackageImport(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`, `//{./examples/json/foo.json}`)
}

func TestJsonPackageImportNotExists(t *testing.T) {
	t.Parallel()
	//FIXME: empty error message because windows and UNIX error messages are different
	AssertCodeErrors(t, "", `//{./examples/json/fooooo.json}`)
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3.141592653589793`, `//github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }

func TestPackageExternalImportModule(t *testing.T) {
	// DO NOT t.Parallel()

	tempdir, err := ioutil.TempDir("", "arrai-TestPackageExternalImportModule-")
	require.NoError(t, err)
	t.Logf("tempdir: %s", tempdir)
	defer func() { require.NoError(t, os.RemoveAll(tempdir)) }()

	wd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempdir)
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(wd)) }()

	t.Run("", func(t *testing.T) {
		run := func(name string, f func(t *testing.T)) {
			t.Run(name, func(t *testing.T) {
				f(t)
			})
		}

		repo := "github.com/ChloePlanet/arrai-examples"
		run("Module", func(t *testing.T) { AssertCodesEvalToSameValue(t, `3`, `//{`+repo+`/add}`) })
		run("ModuleExt", func(t *testing.T) { AssertCodesEvalToSameValue(t, `3`, `//{`+repo+`/add.arrai}`) })
		run("ModuleJson", func(t *testing.T) { AssertCodesEvalToSameValue(t, `{}`, `//{`+repo+`/empty.json}`) })

	})
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
