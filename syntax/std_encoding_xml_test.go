package syntax

import (
	"testing"
)

//nolint:lll
const arraiData = `[(decl: (target: 'xml', text: 'version="1.0"')), (text: '\n'), (elem: (children: [(text: '\n   '), (elem: (attrs: {'id': 'bk101'}, children: [(text: '\n      '), (elem: (children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (children: [(text: 'XML Developers Guide')], name: 'title')), (text: '\n      '), (elem: (children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

//nolint:lll
const strippedArraiData = `[(decl: (target: 'xml', text: 'version="1.0"')), (elem: (children: [(elem: (attrs: {'id': 'bk101'}, children: [(elem: (children: [(text: 'Gambardella, Matthew')], name: 'author')), (elem: (children: [(text: 'XML Developers Guide')], name: 'title')), (elem: (children: [(text: 'Computer')], name: 'genre')), (elem: (children: [(text: '44.95')], name: 'price')), (elem: (children: [(text: '2000-10-01')], name: 'publish_date')), (elem: (children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description'))], name: 'book'))], name: 'catalog'))]`
const xmlData = `<<'<?xml version="1.0"?>
<catalog>
   <book id="bk101">
      <author>Gambardella, Matthew</author>
      <title>XML Developers Guide</title>
      <genre>Computer</genre>
      <price>44.95</price>
      <publish_date>2000-10-01</publish_date>
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

func TestXMLDecode_roundTrip(t *testing.T) {
	t.Parallel()

	xml := xmlData
	expected := arraiData

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decode("+xml+")")
	AssertCodesEvalToSameValue(t, xml, "//encoding.xml.encode("+expected+")")
}
