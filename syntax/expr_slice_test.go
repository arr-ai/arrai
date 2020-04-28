package syntax

import "testing"

//nolint:dupl
func TestArraySlice(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1, 2, 3]      `, `[0, 1, 2, 3, 4](1;4)    `)
	AssertCodesEvalToSameValue(t, `[2, 3, 4]      `, `[0, 1, 2, 3, 4](2;)     `)
	AssertCodesEvalToSameValue(t, `[0, 1, 2]      `, `[0, 1, 2, 3, 4](;3)     `)
	AssertCodesEvalToSameValue(t, `[1, 3]         `, `[0, 1, 2, 3, 4](1;;2)   `)
	AssertCodesEvalToSameValue(t, `[0, 2]         `, `[0, 1, 2, 3, 4](;4;2)   `)
	AssertCodesEvalToSameValue(t, `[4, 2]         `, `[0, 1, 2, 3, 4](4;1;-2) `)
	AssertCodesEvalToSameValue(t, `[4, 3, 2, 1, 0]`, `[0, 1, 2, 3, 4](;;-1)   `)
	AssertCodesEvalToSameValue(t, `[0]            `, `[0, 1, 2, 3, 4](0;;-1)  `)
	AssertCodesEvalToSameValue(t, `[4, 3]         `, `[0, 1, 2, 3, 4](;2;-1)  `)
	AssertCodesEvalToSameValue(t, `[0, 1, 2, 3, 4]`, `[0, 1, 2, 3, 4](;)      `)
	AssertCodesEvalToSameValue(t, `{}             `, `[0, 1, 2, 3, 4](1;3;-1) `)
	AssertCodesEvalToSameValue(t, `{}             `, `[0, 1, 2, 3, 4](1;1)    `)

	AssertCodeErrors(t, `[0, 1, 2, 3, 4](10;1;-1)`, `10 is out of range`)
	AssertCodeErrors(t, `[0, 1, 2, 3, 4](-10;1;-1)`, `-10 is out of range`)
	AssertCodeErrors(t, `[0, 1, 2, 3, 4](0;42;-1)`, `42 is out of range`)
	AssertCodeErrors(t, `[0, 1, 2, 3, 4](0;-42;-1)`, `-42 is out of range`)
	AssertCodeErrors(t, `[](0;0)`, "set is empty")
}

//nolint:dupl
func TestArrayString(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `"bcd"  `, `"abcde"(1;4)    `)
	AssertCodesEvalToSameValue(t, `"cde"  `, `"abcde"(2;)     `)
	AssertCodesEvalToSameValue(t, `"abc"  `, `"abcde"(;3)     `)
	AssertCodesEvalToSameValue(t, `"bd"   `, `"abcde"(1;;2)   `)
	AssertCodesEvalToSameValue(t, `"ac"   `, `"abcde"(;4;2)   `)
	AssertCodesEvalToSameValue(t, `"ec"   `, `"abcde"(4;1;-2) `)
	AssertCodesEvalToSameValue(t, `"edcba"`, `"abcde"(;;-1)   `)
	AssertCodesEvalToSameValue(t, `"a"    `, `"abcde"(0;;-1)  `)
	AssertCodesEvalToSameValue(t, `"ed"   `, `"abcde"(;2;-1)  `)
	AssertCodesEvalToSameValue(t, `"abcde"`, `"abcde"(;)      `)
	AssertCodesEvalToSameValue(t, `{}     `, `"abcde"(1;3;-1) `)
	AssertCodesEvalToSameValue(t, `{}     `, `"abcde"(1;1)    `)

	AssertCodeErrors(t, `"abcde"(10;1;-1)`, `10 is out of range`)
	AssertCodeErrors(t, `"abcde"(-10;1;-1)`, `-10 is out of range`)
	AssertCodeErrors(t, `"abcde"(0;42;-1)`, `42 is out of range`)
	AssertCodeErrors(t, `"abcde"(0;-42;-1)`, `-42 is out of range`)
	AssertCodeErrors(t, `""(0;0)`, "set is empty")
}

//nolint:dupl
func TestArrayBytes(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 98), (1, 99), (2, 100) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(1;4)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 99), (1, 100), (2, 101) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(2;)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(;3)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 98), (1, 100) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(1;;2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 97), (1, 99) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(;4;2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 101), (1, 99) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(4;1;-2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 101), (1, 100), (2, 99), (3, 98), (4, 97) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(;;-1)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 97) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(0;;-1)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 101), (1, 100) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(;2;-1)`,
	)
	AssertCodesEvalToSameValue(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(;)`,
	)
	AssertCodesEvalToSameValue(t,
		`{}`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(1;3;-1)`,
	)
	AssertCodesEvalToSameValue(t,
		`{}`,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(1;1)`,
	)

	AssertCodeErrors(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(10;1;-1)`,
		`10 is out of range`)
	AssertCodeErrors(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(-10;1;-1)`,
		`-10 is out of range`)
	AssertCodeErrors(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(0;42;-1)`,
		`42 is out of range`)
	AssertCodeErrors(t,
		`{ |@, @byte| (0, 97), (1, 98), (2, 99), (3, 100), (4, 101) }(0;-42;-1)`,
		`-42 is out of range`)
}

//nolint:dupl
func TestDictSlice(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[10, "abc", 30]`,
		`{1: 10, 2: "abc", 3: 30, 4: 40, 5: 50}(1;4)`,
	)
	AssertCodesEvalToSameValue(t,
		`["abc", 30, 40, 50]`,
		`{1: 10, 2: "abc", 3: 30, 4: 40, 5: 50}(2;)`,
	)
	AssertCodesEvalToSameValue(t,
		`[10, "abc", 30, 40]`,
		`{1: 10, 2: "abc", 3: 30, 4: 40, 5: 50}(;5)`,
	)
	AssertCodesEvalToSameValue(t,
		`[50, 40, 30, "abc"]`,
		`{1: 10, 2: "abc", 3: 30, 4: 40, 5: 50}(5;1;-1)`,
	)
	AssertCodesEvalToSameValue(t,
		`[50, 30]`,
		`{1: 10, 2: "abc", 3: 30, 4: 40, 5: 50}(5;1;-2)`,
	)
	AssertCodesEvalToSameValue(t,
		`{}`,
		`{1: 10, 2: "abc", "c": 30, 4: 40, 5: 50}(1;1)`,
	)
}

func TestGenericSetSlice(t *testing.T) {
	t.Parallel()

	AssertCodePanics(t, `{1, 2, 3}(1;3)`)
}

func TestSliceFails(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		`{"a": 1, "b": 2, "c": 3}("a";)`,
		`lower bound does not evaluate to a Number: a`,
	)
	AssertCodeErrors(t,
		`{"a": 1, "b": 2, "c": 3}(;"c")`,
		`upper bound does not evaluate to a Number: c`,
	)
	AssertCodeErrors(t,
		`{"a": 1, "b": 2, "c": 3}(;;"b")`,
		`step does not evaluate to a Number: b`,
	)
	AssertCodeErrors(t,
		`{"a": 1, "b": 2, "c": 3}(;;0)`,
		`step cannot be 0`,
	)
}
