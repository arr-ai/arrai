package syntax

import "testing"

func TestDeprecatedExec(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(
		args: ['ls', 'std_deprecated_test.go'],
		exitCode: 0,
		stdout: <<'std_deprecated_test.go\n'>>,
		stderr: {},
	)`, `//deprecated.exec(['ls', 'std_deprecated_test.go'])`)
	AssertCodesEvalToSameValue(t, `(
		args: ['echo'],
		exitCode: 0,
		stdout: <<'\n'>>,
		stderr: {},
	)`, `//deprecated.exec(['echo'])`)
}

func TestDeprecatedExecError(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(
		args: ['ls', 'nonexistent'],
		exitCode: 1,
		stdout: {},
		stderr: <<'ls: nonexistent: No such file or directory\n'>>,
	)`, `//deprecated.exec(['ls', 'nonexistent'])`)
}

func TestDeprecatedExecFail(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, ``, `//deprecated.exec('')`)
	AssertCodeErrors(t, ``, `//deprecated.exec(['ls std_deprecated_test.go'])`)
}
