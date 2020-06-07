package syntax

import (
	"io/ioutil"
	"os"
	"os/exec"
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
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{./examples/simple/simple}`)
}

func TestPackageImportFromRoot(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{/examples/simple/simple}`)
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
	AssertCodeErrors(t, `//{./examples/json/fooooo.json}`, "")
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

	cmd := exec.Command("go", "mod", "init", "github.com/arr-ai/arrai/fake")
	require.NoError(t, cmd.Run())

	repo := "github.com/ChloePlanet/arrai-examples"

	t.Run("Module", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `3`, `//{`+repo+`/add}`)
	})
	t.Run("ModuleWithExt", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `3`, `//{`+repo+`/add.arrai}`)
	})
	t.Run("ModuleJson", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `{}`, `//{`+repo+`/empty.json}`)
	})

	raw := "raw.githubusercontent.com/ChloePlanet/arrai-examples/master"

	t.Run("URLArrai", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `3`, `//{https://`+raw+`/add.arrai}`) // nolint:lll
	})
	t.Run("URLArraiWithoutHTTPS", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `3`, `//{`+raw+`/add.arrai}`)
	})
	t.Run("URLJson", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `{}`, `//{https://`+raw+`/empty.json}`) // nolint:lll
	})
	t.Run("URLJsonWithoutHTTPS", func(t *testing.T) {
		t.Parallel()
		AssertCodesEvalToSameValue(t, `{}`, `//{`+raw+`/empty.json}`)
	})
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
