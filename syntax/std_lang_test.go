package syntax

import (
	"testing"
)

func TestStdLangGoParseString(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(
			(@type: 'File', Comments: {}, Decls: {}, Doc: {}, Imports: {}, Package: 1, 
			Name: (@type: 'Ident', Name: 'foo', NamePos: 9, Obj: {}),
			Scope: (@type: 'Scope', Objects: {}, Outer: {}), Unresolved: {})
		)`, `//lang.go.parse("package foo")`)
}

func TestStdLangGoParseBytes(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(
			(@type: 'File', Comments: {}, Decls: {}, Doc: {}, Imports: {}, Package: 1, 
			Name: (@type: 'Ident', Name: 'foo', NamePos: 9, Obj: {}),
			Scope: (@type: 'Scope', Objects: {}, Outer: {}), Unresolved: {})
		)`, `//lang.go.parse(<<"package foo">>)`)
}

func TestStdLangGoParseEmpty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `//lang.go.parse("")`)
}

func TestStdLangGoParseTuple(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//lang.go.parse((src: "package foo"))`)
}
