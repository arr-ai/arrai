package syntax

import "testing"

func TestDeprecatedExec(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `<<'std_deprecated_test.go\n'>>`, `//deprecated.exec('ls std_deprecated_test.go')`)
}

func TestDeprecatedExecError(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, ``, `//deprecated.exec('')`)
	AssertCodeErrors(t, ``, `//deprecated.exec('ls | wc -l')`)
}
