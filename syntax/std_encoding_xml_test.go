package syntax

import (
	"testing"
)

func TestXMLEncode_declaration(t *testing.T) {
	t.Parallel()

	expected := `<<'<?xml version="1.0"?>'>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode([(xmldecl: 'version="1.0"')])`)
}

func TestXMLEncode_element(t *testing.T) {
	t.Parallel()

	expected := `<<'<catalog>hello</catalog>'>>`

	data := `[(elem: (attrs: {}, children: [(text: 'hello')], name: 'catalog'))]`

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

//nolint:lll
func TestXMLEncode_implicitNamespace(t *testing.T) {
	t.Parallel()

	data := `[(elem: (attrs: {(name: 'xmlns', text: 'doop')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book', ns: 'doop'))], name: 'catalog', ns: 'doop'))]`
	expected := `<<'<catalog xmlns="doop"><book>harry potter</book></catalog>'>>`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.encode("+data+")")
}

//nolint:lll
func TestXMLEncode_explicitNamespace(t *testing.T) {
	t.Parallel()
	// NOTE: skipped due to the current implementation's limitation on explicit namespaces. remove when updated
	t.SkipNow()

	data := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book'))], name: 'catalog', ns: 'doop'))]`
	expected := `<<'<hello:catalog xmlns:hello="doop"><book>harry potter</book></hello:catalog>'>>`
	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.encode("+data+")")
}

//nolint:lll
func TestXMLEncode_dualNamespace(t *testing.T) {
	t.Parallel()
	// NOTE: skipped due to the current implementation's limitation on explicit namespaces. remove when updated
	t.SkipNow()

	data := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop'), (name: 'xmlns', text: 'maaw')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book', ns: 'maaw'))], name: 'catalog', ns: 'doop'))]`
	expected := `<<'<hello:catalog xmlns:hello="doop" xmlns="maaw"><book>harry potter</book></hello:catalog>'>>`
	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.encode("+data+")")
}

//nolint:lll
func TestXMLEncode_missingExplictNS(t *testing.T) {
	t.Parallel()
	// NOTE: skipped due to the current implementation's limitation on explicit namespaces. remove when updated
	t.SkipNow()

	data := `[(elem: (attrs: {}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book'))], name: 'catalog', ns: 'doop'))]`
	AssertCodeErrors(t, "", "//encoding.xml.encode("+data+")")
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

	expected := `[(xmldecl:  'version="1.0"')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode('<?xml version="1.0"?>')`)
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(<<'<?xml version="1.0"?>'>>)`)

	expected = `[(text: 'woop')]`
	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decode('woop')")
}

func TestXMLDecode_element(t *testing.T) {
	t.Parallel()

	data := `<<'<catalog>hello</catalog>'>>`
	expected := `[(elem: (attrs: {}, children: [(text: 'hello')], name: 'catalog'))]`

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
	expected := `[(elem: (attrs: {(name: 'xmlns', text: 'doop')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book', ns: 'doop'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

//nolint:lll
// parent node uses explicit namespace
func TestXMLDecode_explicitNamespace(t *testing.T) {
	t.Parallel()

	data := `<<'<hello:catalog xmlns:hello="doop"><book>harry potter</book></hello:catalog>'>>`
	expected := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

//nolint:lll
// children node inherits parent's implicit namespace
// parent node uses explicit namespace
func TestXMLDecode_dualNamespace(t *testing.T) {
	t.Parallel()

	data := `<<'<hello:catalog xmlns:hello="doop" xmlns="maaw"><book>harry potter</book></hello:catalog>'>>`
	expected := `[(elem: (attrs: {(name: 'hello', ns: 'xmlns', text: 'doop'), (name: 'xmlns', text: 'maaw')}, children: [(elem: (attrs: {}, children: [(text: 'harry potter')], name: 'book', ns: 'maaw'))], name: 'catalog', ns: 'doop'))]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_missingExplictNS(t *testing.T) {
	t.Parallel()
	// NOTE: skipped due to the current implementation's limitation on explicit namespaces. remove when updated
	t.SkipNow()

	data := `<<'<here:catalog></here:catalog>'>>`
	AssertCodeErrors(t, "", "//encoding.xml.decode("+data+")")
}

func TestXMLDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decode(`<root>`)")
}

//nolint:lll
func TestXMLDecoder_strip(t *testing.T) {
	t.Parallel()

	xml := `<<'<catalog>\n\t<book>Harry\nPotter</book>\n</catalog>'>>`
	expected := `[(elem: (attrs: {}, children: [(elem: (attrs: {}, children: [(text: 'Harry\nPotter')], name: 'book'))], name: 'catalog'))]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder((trimSurroundingWhitespace: true)).decode("+xml+")")
}

//nolint:lll
func TestXMLDecoder_dontStrip(t *testing.T) {
	t.Parallel()

	xml := `<<'<catalog>\n\t<book>Harry\nPotter</book>\n</catalog>'>>`
	expected := `[(elem: (attrs: {}, children: [(text: '\n\t'), (elem: (attrs: {}, children: [(text: 'Harry\nPotter')], name: 'book')), (text: '\n')], name: 'catalog'))]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder((trimSurroundingWhitespace: false)).decode("+xml+")")
}

func TestXMLDecoder_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decoder((trimSurroundingWhitespace: false)).decode(`<root>`)")
	AssertCodeErrors(t, "", "//encoding.xml.decoder((unknown: false)).decode(`<root>`)")
}
