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
		args: ['cat', 'nonexistent'],
		exitCode: 1,
		stdout: {},
		stderr: <<'cat: nonexistent: No such file or directory\n'>>,
	)`, `//deprecated.exec(['cat', 'nonexistent'])`)
}

func TestDeprecatedExecFail(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, ``, `//deprecated.exec('')`)
	AssertCodeErrors(t, ``, `//deprecated.exec(['ls std_deprecated_test.go'])`)
}
