package syntax

import (
	"testing"
)

//nolint:lll
const arraiData = `[(decl: (target: 'xml', text: 'version="1.0"')), (text: '\n'), (elem: (attrs: {(name: 'xmlns', text: 'doop')}, children: [(text: '\n   '), (elem: (attrs: {(name: 'xmlns', text: 'woop')}, children: [(text: '\n      '), (elem: (attrs: {(name: 'id', text: 'bk101')}, children: [(text: 'yesman')], name: 'author', ns: 'woop')), (text: '\n      '), (elem: (children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description', ns: 'woop')), (text: '\n   ')], name: 'book', ns: 'woop')), (text: '\n')], name: 'catalog', ns: 'doop')), (text: '\n')]`

//nolint:lll
const strippedArraiData = `[(decl: (target: 'xml', text: 'version="1.0"')), (elem: (attrs: {(name: 'xmlns', text: 'doop')}, children: [(elem: (attrs: {(name: 'xmlns', text: 'woop')}, children: [(elem: (attrs: {(name: 'id', text: 'bk101')}, children: [(text: 'yesman')], name: 'author', ns: 'woop')), (elem: (children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description', ns: 'woop'))], name: 'book', ns: 'woop'))], name: 'catalog', ns: 'doop'))]`
const xmlData = `<<'<?xml version="1.0"?>
<catalog xmlns="doop">
   <book xmlns="woop">
      <author id="bk101">yesman</author>
      <description>An in-depth look at creating applications 
      with XML.</description>
   </book>
</catalog>
'>>`

func TestXMLEncode_declaration(t *testing.T) {
	t.Parallel()

	expected := `<<'<?xml version="1.0"?>'>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode([(decl: (target: 'xml', text: 'version="1.0"'))])`)
}

func TestXMLEncode_element(t *testing.T) {
	t.Parallel()

	expected := `<<'<catalog>hello</catalog>'>>`

	data := `[(elem: (children: [(text: 'hello')], name: 'catalog'))]`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode(`+data+`)`)
}

func TestXMLEncode_text(t *testing.T) {
	t.Parallel()

	expected := `<<'hello world'>>`

	data := `[(text: 'hello world')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode(`+data+`)`)
}

func TestXMLEncode_comment(t *testing.T) {
	t.Parallel()

	expected := `<<'<!--hello world comment-->'>>`

	data := `[(comment: 'hello world comment')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode(`+data+`)`)
}

func TestXMLEncode_directive(t *testing.T) {
	t.Parallel()

	expected := `<<'<!ATTLIST foo a CDATA #IMPLIED>'>>`
	data := `[(directive: 'ATTLIST foo a CDATA #IMPLIED')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode(`+data+`)`)
}

func TestXMLEncodeLarge(t *testing.T) {
	t.Parallel()

	expected := xmlData
	data := arraiData

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.encode("+data+")")
}

func TestXMLEncode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.encode(`woop`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<ldkfjroot>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root></hi>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<?xx?>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<!xx!>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<<'hel'>>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<<'[]'>>`)")
}

func TestXMLDecode(t *testing.T) {
	t.Parallel()

	expected := `[(decl: (target: 'xml', text: 'version="1.0"'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode('<?xml version="1.0"?>')`)
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(<<'<?xml version="1.0"?>'>>)`)

	expected = `[(text: 'woop')]`
	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decode('woop')")
}

func TestXMLDecode_element(t *testing.T) {
	t.Parallel()

	expected := arraiData
	data := xmlData

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_text(t *testing.T) {
	t.Parallel()

	expected := `[(text: 'hello world')]`
	data := `<<'hello world'>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_comment(t *testing.T) {
	t.Parallel()

	expected := `[(comment: 'hello world comment')]`
	data := `<<'<!--hello world comment-->'>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_directive(t *testing.T) {
	t.Parallel()

	data := `<<'<!ATTLIST foo a CDATA #IMPLIED>'>>`
	expected := `[(directive: 'ATTLIST foo a CDATA #IMPLIED')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

//nolint:lll
// children node inherits parent's implicit namespace
func TestXMLDecode_implicitNamespace(t *testing.T) {
	t.Parallel()

	data := `<<'<catalog xmlns="doop"><book>harry potter</book></catalog>'>>`
	expected := `[(elem: (attrs: {(name: 'xmlns', text: 'doop')}, children: [(elem: (children: [(text: 'harry potter')], name: 'book', ns: 'doop'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

//nolint:lll
// parent node uses explicit namespace
func TestXMLDecode_explicitNamespace(t *testing.T) {
	t.Parallel()

	data := `<<'<hello:catalog xmlns:hello="doop"><book>harry potter</book></hello:catalog>'>>`
	expected := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop')}, children: [(elem: (children: [(text: 'harry potter')], name: 'book'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

//nolint:lll
// children node inherits parent's implicit namespace
// parent node uses explicit namespace
func TestXMLDecode_dualNamespace(t *testing.T) {
	t.Parallel()

	data := `<<'<hello:catalog xmlns:hello="doop" xmlns="maaw"><book>harry potter</book></hello:catalog>'>>`
	expected := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop'), (name: 'xmlns', text: 'maaw')}, children: [(elem: (children: [(text: 'harry potter')], name: 'book', ns: 'maaw'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decode(`<root>`)")
}

func TestXMLDecoder_strip(t *testing.T) {
	t.Parallel()

	xml := xmlData
	expected := strippedArraiData

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder((stripFormatting: true))("+xml+")")
}

func TestXMLDecoder_dontStrip(t *testing.T) {
	t.Parallel()

	xml := xmlData
	expected := arraiData

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder((stripFormatting: false))("+xml+")")
}

func TestXMLDecoder_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decoder((stripFormatting: false))(`<root>`)")
	AssertCodeErrors(t, "", "//encoding.xml.decoder((unknown: false))(`<root>`)")
}
