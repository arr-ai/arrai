package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLTrivial(t *testing.T) {
	assertParse(t, rel.NewXML([]rune("a"), []rel.Attr{}), "<a/>")
	assertParse(t, rel.NewXML([]rune("@b-c_1$.D"), []rel.Attr{}), "<@b-c_1$.D/>")
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLTrivialWithEndTag(t *testing.T) {
	assertParse(t, rel.NewXML([]rune("a"), []rel.Attr{}), "<a></a>")
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLTrivialWithMismatchedEndTags(t *testing.T) {
	value, err := syntax.Parse(syntax.NewStringLexer("<a></ab>"))
	assert.Error(t, err, "%s", value)
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLNested(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), []rel.Attr{},
			rel.NewXML([]rune("b"), []rel.Attr{})),
		"<a><b/></a>")
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXML1Attr(t *testing.T) {
	assertParse(t,
		rel.NewXML(
			[]rune("a"),
			[]rel.Attr{{Name: "x", Value: rel.NewNumber(1)}},
		),
		`<a x=1/>`)
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXML2Attrs(t *testing.T) {
	assertParse(t,
		rel.NewXML(
			[]rune("abc"),
			[]rel.Attr{
				{Name: "x", Value: rel.NewNumber(1)},
				{Name: "yz", Value: rel.NewString([]rune("hello"))},
			}),
		`<abc x=1 yz="hello"/>`)
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXML1Data(t *testing.T) {
	assertParse(t,
		rel.NewXML(
			[]rune("abc"),
			[]rel.Attr{
				{Name: "x", Value: rel.NewNumber(1)},
				{Name: "yz", Value: rel.NewString([]rune("hello"))},
			}),
		`<abc x=1 yz="hello"/>`)
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLHtmlEntities(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), nil, rel.NewString([]rune("&"))),
		`<a>&amp;</a>`)
}

// TestParseNumber tests Parse recognising numbers.
func TestParseXMLHtmlEntitiesEuroBug(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), nil, rel.NewString([]rune("€"))),
		`<a>&euro;</a>`)
}

var xmlSpacePreserve = rel.Attr{
	Name:  "{https://www.w3.org/XML/1998/namespace}space",
	Value: rel.NewString([]rune("preserve")),
}

// TestParseTrimSpace tests Parse trimming whitespace.
// TODO: More edge-case coverage.
func TestParseTrimSpace(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), nil, rel.NewString([]rune("foo"))),
		`<a>
  foo
</a>`)
}

// TestParseXMLSpaceBadValue tests Parse error on xml:space="wrong".
func TestParseXMLSpaceBadValue(t *testing.T) {
	assertParseError(t, `<a xml:space="wrong"/>`)
}

// TestParseSpacePreserve tests Parse handling xml:space="preserve".
func TestParseSpacePreserve(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), []rel.Attr{xmlSpacePreserve},
			rel.NewString([]rune("\n  foo\n"))),
		`<a xml:space="preserve">
  foo
</a>`)
}

func xmlnsDefault(ns string) rel.Attr {
	return rel.Attr{Name: "xmlns", Value: rel.NewString([]rune(ns))}
}

func xmlnsAlias(alias string, ns string) rel.Attr {
	return rel.Attr{
		Name:  "{http://www.w3.org/2000/xmlns/}" + alias,
		Value: rel.NewString([]rune(ns)),
	}
}

// TestParseXmlns tests Parse handling xmlns="...".
func TestParseXmlns(t *testing.T) {
	assertParse(t,
		rel.NewXML(
			[]rune("{my-ns}foobar"),
			[]rel.Attr{xmlnsAlias("me", "my-ns")},
		),
		`<me:foobar xmlns:me="my-ns"/>`)
}

// TestParseXmlnsDefault tests Parse handling xmlns="...".
func TestParseXmlnsDefault(t *testing.T) {
	assertParse(t,
		rel.NewXML(
			[]rune("{my-ns}foobar"),
			[]rel.Attr{xmlnsDefault("my-ns")},
		),
		`<foobar xmlns="my-ns"/>`)
}

// TestParseXmlnsDefaultInAlias tests Parse handling xmlns="...".
func TestParseXmlnsDefaultInAlias(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("{my-ns}foobar"),
			[]rel.Attr{
				xmlnsDefault("def-ns"),
				xmlnsAlias("me", "my-ns"),
			},
			rel.NewXML([]rune("{def-ns}baz"), nil),
		),
		`<me:foobar xmlns="def-ns" xmlns:me="my-ns"><baz/></me:foobar>`)
}

// TestParseXmlnsAliasInDefault tests Parse handling xmlns="...".
func TestParseXmlnsAliasInDefault(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("{def-ns}foobar"),
			[]rel.Attr{
				xmlnsDefault("def-ns"),
				xmlnsAlias("me", "my-ns"),
			},
			rel.NewXML([]rune("{my-ns}baz"), nil),
		),
		`<foobar xmlns="def-ns" xmlns:me="my-ns"><me:baz/></foobar>`)
}

// TestParseXmlExprInsideElt tests Parse handling <a>{1}</a>.
func TestParseXmlExprInsideElt(t *testing.T) {
	assertParse(t,
		rel.NewXML([]rune("a"), []rel.Attr{},
			rel.NewNumber(1)),
		`<a>{1}</a>`)
}

// TODO: Fix
// // TestParseXmlDotExprAttr tests attr=.foo.
// func TestParseXmlDotExprAttr(t *testing.T) {
// 	assertParse(t,
// 		rel.NewXML([]rune("a"),
// 			[]rel.Attr{
// 				{Name: "attr", Value: rel.NewNumber(42)},
// 			}),
// 		`<a attr=.foo/>`,
// 	)
// }
