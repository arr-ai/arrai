package syntax

import "testing"

func TestDeprecatedExec(t *testing.T) {
	AssertCodesEvalToSameValue(t, `(
		args: ['ls', 'std_deprecated_test.go'],
		exitCode: 0,
		stdout: <<'std_deprecated_test.go\n'>>,
		stderr: {},
	)`, `//deprecated.exec(['ls', 'std_deprecated_test.go'])`)
}

func TestDeprecatedExecError(t *testing.T) {
	// TODO: Not sure why the stderr string looks the way it does.
	AssertCodesEvalToSameValue(t, `(
		args: ['ls', 'nonexistent'],
		exitCode: 1,
		stdout: {},
		stderr: <<'ls: nonexistent: No such file or directory\n'>>,
	)`, `//deprecated.exec(['ls', 'nonexistent'])`)
}

func TestDeprecatedExecFail(t *testing.T) {
	AssertCodeErrors(t, ``, `//deprecated.exec('')`)
	AssertCodeErrors(t, ``, `//deprecated.exec(['ls std_deprecated_test.go'])`)
}
