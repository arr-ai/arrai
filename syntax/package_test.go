//nolint:goconst
package syntax

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
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

func TestIllegalImport(t *testing.T) {
	t.Parallel()

	errMessage := func(x string) string {
		return fmt.Sprintf("import path can not be pointing outside of the script's module directory: %s", x)
	}

	AssertCodeErrors(t, errMessage("../test"), `//{./../test}`)
	AssertCodeErrors(t, errMessage("../../../../test"), `//{./../../../../test}`)

	// this is allowed because .. on absolute import does not affect anything
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{/../../../examples/simple/simple}`)

	// this is allowed because it does not import outside the parent directory
	AssertCodesEvalToSameValue(t, `{1, 4, 9, 16}`, `//{./examples/../examples/simple/simple}`)
}

func TestJsonPackageImport(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`, `//{./examples/json/foo.json}`)
}

func TestExplicitNonArraiImport(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{
			|@row, ignore, non_heading|
			(1, 'Account, Balance', 'Customer ID'),
			(2, '100', 'foo'),
			(4, '200', 'bar'),
			(5, '300', 'bar'),
		}`,
		`//[//encoding.xlsx]{./examples/xlsx/foo.xlsx}`,
	)
	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['x', 'y', 'z']]`,
		`//[//encoding.csv]{./examples/csv/foo.csv}`,
	)
	AssertCodesEvalToSameValue(t,
		`{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`,
		`//[//encoding.json]{./examples/json/foo.json}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: [{'foo': (s: 'bar')}])`,
		`//[//encoding.yaml]{./examples/yaml/foo.yml}`,
	)
	if runtime.GOOS == "windows" {
		AssertCodesEvalToSameValue(t,
			`<<'1\r\n'>>`,
			`//[//encoding.bytes]{./examples/import/bar.arrai}`,
		)
		//FIXME: the expected value kept changing for windows
		// AssertCodesEvalToSameValue(t,
		// 	`<<'2\x0ev'>>`,
		// 	`//[(decode: \b b >> . + 1)]{./examples/import/bar.arrai}`,
		// )
	} else {
		AssertCodesEvalToSameValue(t,
			`<<'1\n'>>`,
			`//[//encoding.bytes]{./examples/import/bar.arrai}`,
		)
		AssertCodesEvalToSameValue(t,
			`<<'2\v'>>`,
			`//[(decode: \b b >> . + 1)]{./examples/import/bar.arrai}`,
		)
	}
}

func TestExplicitDecoderPrecedence(t *testing.T) {
	t.Parallel()

	// FIXME: this is to avoid crlf
	if runtime.GOOS != "windows" {
		AssertCodesEvalToSameValue(t,
			`<<'{\n    "name": "foo",\n    "location": "Melbourne"\n}\n'>>`,
			`//[//encoding.bytes]{./examples/json/foo.json}`,
		)
	}
}

func TestCacheDecoderImport(t *testing.T) {
	t.Parallel()

	// ensure that non arr.ai imports are cached based on the decoder
	// FIXME: this is to avoid crlf
	if runtime.GOOS != "windows" {
		AssertCodesEvalToSameValue(t,
			`{'location': (s: 'Melbourne'), 'name': <<'{\n    "name": "foo",\n    "location": "Melbourne"\n}\n'>>}`,
			`//[//encoding.json]{./examples/json/foo.json} +> {"name": //[//encoding.bytes]{./examples/json/foo.json}}`,
		)
	}
}

func TestImplicitNonArraiImport(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{
			|@row, ignore, non_heading|
			(1, 'Account, Balance', 'Customer ID'),
			(2, '100', 'foo'),
			(4, '200', 'bar'),
			(5, '300', 'bar'),
		}`,
		`//{./examples/xlsx/foo.xlsx}`,
	)
	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['x', 'y', 'z']]`,
		`//{./examples/csv/foo.csv}`,
	)
	AssertCodesEvalToSameValue(t,
		`{'location': (s: 'Melbourne'), 'name': (s: 'foo')}`,
		`//{./examples/json/foo.json}`,
	)
	AssertCodesEvalToSameValue(t,
		`(a: [{'foo': (s: 'bar')}])`,
		`//{./examples/yaml/foo.yml}`,
	)
}

func TestErrorExplicitImport(t *testing.T) {
	t.Parallel()

	// no decoder tuple
	AssertCodeErrors(t,
		`does not evaluate to a decoder tuple: 1`,
		`//[1]{./examples/import/bar.arrai}`,
	)

	// fail to evaluate decoder
	AssertCodeErrors(t,
		`fail to evalute decoder: name "a" not found in {}`,
		`//[1 + a]{./examples/import/bar.arrai}`,
	)

	// import is compile time, cannot use runtime variables
	AssertCodeErrors(t,
		`name "a" not found in {b}`,
		`let a = 1; //[(decode: \b b + a)]{./examples/import/bar.arrai}`,
	)

	// decoder is not a function
	AssertCodeErrors(t,
		`does not evaluate to a decoder function: 1`,
		`let a = 1; //[(decode: 1)]{./examples/import/bar.arrai}`,
	)
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

	cmd := exec.Command("go", "mod", "init", "github.com/arr-ai/arrai/fake")
	require.NoError(t, cmd.Run())

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

		raw := "raw.githubusercontent.com/ChloePlanet/arrai-examples/master"
		run("URLArraiHttps", func(t *testing.T) { AssertCodesEvalToSameValue(t, `3`, `//{https://`+raw+`/add.arrai}`) })
		run("URLArraiNoHttps", func(t *testing.T) { AssertCodesEvalToSameValue(t, `3`, `//{`+raw+`/add.arrai}`) })
		run("URLJsonHttps", func(t *testing.T) { AssertCodesEvalToSameValue(t, `{}`, `//{https://`+raw+`/empty.json}`) })
		run("URLJsonNoHttps", func(t *testing.T) { AssertCodesEvalToSameValue(t, `{}`, `//{`+raw+`/empty.json}`) })
	})
}

// func TestPackageExternalImport(t *testing.T) {
// 	AssertCodesEvalToSameValue(t, `3`, `//http://github.com/'arr-ai'/arrai/examples/'xml.wbnf'`)
// }
