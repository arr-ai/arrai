package syntax

import (
	"strings"
	"testing"
)

func TestXStringSimple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""               `, `$""`)
	AssertCodesEvalToSameValue(t, `"42"             `, `$"${6*7}"`)
	AssertCodesEvalToSameValue(t, `"a42z"           `, `$"a${6*7}z"`)
	AssertCodesEvalToSameValue(t, `"a00042z"        `, `$"a${6*7:05d}z"`)
	AssertCodesEvalToSameValue(t, `"a001, 002, 003z"`, `$"a${[1, 2, 3]:03d:, }z"`)
	AssertCodesEvalToSameValue(t, `"a42k3.142z"     `, `$"a${6*7}k${//math.pi:.3f}z"`)
}

func TestXStringBackquote(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `""      `, "$``")
	AssertCodesEvalToSameValue(t, `"a\\n42"`, "$`a\\n${6*7}`")
}

func TestXStringStrings(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"hello"`, `$"${'hello'}"`)
}

func TestXStringIndent(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"a\nb"`, "$`\n  a\n  b`")
	AssertCodesEvalToSameValue(t, `"a\nb\n  c\nd"`, "$'\n  a\n  b\n    c\n  d'")
	AssertCodesEvalToSameValue(t,
		`
"letter: a
letter: b
letter: c"
`,
		`
$"${['a', 'b', 'c'] >> $'letter: ${.}'::\i}"
        `,
	)

	AssertCodesEvalToSameValue(t,
		`
"stuff xyzabc stuff 123123123"
`,
		`
$"stuff xyz${"abc"} stuff ${"123123123"}"
        `,
	)

	AssertCodesEvalToSameValue(t,
		`
"abc:
    letter: a
    letter: b
    letter: c"
`,
		`$"
        abc:
            ${['a', 'b', 'c'] >> $'letter: ${.}'::\i}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"stuff
    123
    321
    456
    654
333"
`,
		`$"
        stuff
            ${123}
            ${321}
            456
            ${654}
        ${333}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"abc:
    letter:
        a
    letter:
        b
    letter:
        c"
`,
		`$"
        abc:
            ${['a', 'b', 'c'] >> $'
            letter:
                ${.}
            '::\i}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"letters:
    letter:
        d
    letter:
        e
    letter:
        f

    numbers:
        number:
            1
        number:
            2
        number:
            3"`,
		`$"
        letters:
            ${['d', 'e', 'f'] >> $"
                letter:
                    ${.}
            "::\i}

            numbers:
                ${[1, 2, 3] >> $"
                    number:
                        ${.}
                "::\i}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"abc:
    letter:
        a
    letter:
        b
    letter:
        c"`,
		`$"
        abc:
            ${['a', 'b', 'c'] >> $'
            letter:
                ${.}
            '::\i}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"letters:
    letter:
        d
    letter:
        e
    letter:
        f
numbers:
number:
    1
number:
    2
number:
    3"`,
		`$"
        letters:
            ${['d', 'e', 'f'] >> $"
                letter:
                    ${.}
            "::\i}
        numbers:
        ${[1, 2, 3] >> $"
            number:
                ${.}
        "::\i}
        "`,
	)

	AssertCodesEvalToSameValue(t,
		`
"letters:letter:
    d
letter:
    e
letter:
    f
numbers:
number:
    1
number:
    2
number:
    3"`,
		`$"
        letters:${['d', 'e', 'f'] >> $"
                letter:
                    ${.}
            "::\i}
        numbers:
        ${[1, 2, 3] >> $"
            number:
                ${.}
        "::\i}
        "`,
	)
}

func TestXStringWS(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"1 2"`, `$"${1} ${2}"`)
	AssertCodesEvalToSameValue(t, `"1\n2"`, "$'\n  ${1}\n  ${2}'")
}

func TestXStringSuppressEmptyComputedLines(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"x\ny\n2"`, "$'\n  x\n  ${'y'}\n  ${2}'")
	AssertCodesEvalToSameValue(t, `"x\n2"`, "$'\n  x\n  ${''}\n  ${2}'")
	AssertCodesEvalToSameValue(t, `"x\n2"`, "$'\n  x\n  ${''}\n  ${''}\n  ${2}'")
}

func TestXStringSuppressNewlinesAfterEmptyComputedLines(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`
"stuff:
    abc
    ghi"
		`,
		`
$"
stuff:
    abc
    ${''}
    ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
    abc
ghi"
		`,
		`
$"
stuff:
    abc
    ${''}
ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
    abc
    ghi"
		`,
		`
$"
stuff:
    abc
    ${''}ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
    abc
        ghi"
		`,
		`
$"
stuff:
    abc
    ${''}    ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
    abc

	ghi"
		`,
		`
$"
stuff:
    abc
	${''}

	ghi"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
	abc
	def
	ghi"
		`,
		`
$"
stuff:
	abc
	${''}
	${'def'}
	ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
	abc
	def
	ghi"
		`,
		`
$"
stuff:
	abc
	${''}${'def'}
	ghi
"
		`,
	)
	AssertCodesEvalToSameValue(t,
		`
"stuff:
	abc
	defghi"
		`,
		`
$"
stuff:
	abc
	${''}${'def'}ghi
"
		`,
	)
}

func TestXStringSuppressLastLineWS(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"xy"`, `$"x${ [] :::=}y"`)
	AssertCodesEvalToSameValue(t, `"x123=y"`, `$"x${ [1, 2, 3] :::=}y"`)
	AssertCodesEvalToSameValue(t, `"x\n1\n2\n3\n"`, "$'\n  x\n  ${''}\n  ${''}\n  ${[1, 2, 3]\n  ::\\i:\\n}\n  '")
	AssertCodesEvalToSameValue(t, `"x\n"`, "$'\n  x\n  ${''}\n  ${''}\n  ${[]\n  ::\\i:\\n}\n  '")

	AssertCodesEvalToSameValue(t, `"1\n.\n\n2\n.\n\n3\n.\n"`, `
		$"
			${[1, 2, 3] >> $"
				${.}
				.
				"::\n\n:\n}
		"`)
	AssertCodesEvalToSameValue(t, `"1\n.\n\n2\n.\n\n3\n.\n"`, `
		$"
			${[1, 2, 3] >> $"
				${.}
				.
			"::\n\n:\n}
		"`)
	AssertCodesEvalToSameValue(t, `"stuff:\n\tletter:\n\t\td\n\tletter:\n\t\te\n\tletter:\n\t\tf"`,
		`
$"
stuff:
	${['d', 'e', 'f'] >> $"
		letter:
			${.}
	"::\i}
"
		`,
	)
	AssertCodesEvalToSameValue(t, `"stuff:\n\tletter:\n\t\td\n\tletter:\n\t\te\n\tletter:\n\t\tf\n\t\t\t"`,
		`
$"
stuff:
	${['d', 'e', 'f'] >> $"
		letter:
			${.}
	"::\i}
${"\t\t\t"}
"
		`,
	)
}

func TestXStringArrays(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `"x\n  1\n  2\n  3\ny"`, "$'x\n  ${[1, 2, 3]::\\i}\ny'")
}

func TestXStringMap(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"Getcustid() int"`,
		`(name: "custid", type: "int") -> $"Get${.name}() ${.type}"`,
	)
}

func TestXStringMap2(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`"GetCustid()"`,
		`[(name: "custid", type: "int")] -> $"${. >> $"Get${//str.title(.name)}()"::}"`,
	)
}

func TestXStringNested(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		strings.ReplaceAll(
			`"type Customer interface {
				IsCustomer()
				GetCustid() int
				GetDob() date
				GetAlias() string
			}"`, "\n\t\t\t", "\n",
		),
		`(name: "Customer", fields: [
			(name: "custid", type: "int"   ),
			(name: "dob",    type: "date"  ),
			(name: "alias",  type: "string"),
		]) -> $"
			type ${.name} interface {
				Is${.name}()
				${.fields >> $"Get${//str.title(.name)}() ${.type}"::\i}
			}"`,
	)
}
