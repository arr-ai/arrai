package syntax

import (
	"testing"
)

func TestXMLEncode_declaration(t *testing.T) {
	t.Parallel()

	expected := `<<'<?xml version="1.0"?>'>>`

	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode([(decl: (target: 'xml', text: 'version="1.0"'))])`)
}

func TestXMLEncode_element(t *testing.T) {
	t.Parallel()

	expected := `<<'<catalog>
   <book id="bk101">
      <author>Gambardella, Matthew</author>
      <title>XML Developers Guide</title>
      <genre>Computer</genre>
      <price>44.95</price>
      <publish_date>2000-10-01</publish_date>
      <description>An in-depth look at creating applications with XML.</description>
   </book>
</catalog>
'>>`

	data := `[(elem: (attrs: {}, children: [(text: '\n   '), (elem: (attrs: {(id: 'bk101')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'XML Developers Guide')], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

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

	expected := `<<'<?xml version="1.0"?>'>>`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.encode([(decl: (target: 'xml', text: 'version="1.0"'))])`)
}

//nolint:lll
func TestXMLEncodeLarge(t *testing.T) {
	t.Parallel()

	expected := `<<'<?xml version="1.0"?>
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
   <book id="bk112">
      <author>Galos, Mike</author>
      <title>Visual Studio 7: A Comprehensive Guide</title>
      <genre>Computer</genre>
      <price>49.95</price>
      <publish_date>2001-04-16</publish_date>
      <description>Microsoft Visual Studio 7 is explored in depth,
      looking at how Visual Basic, Visual C++, C#, and ASP+ are 
      integrated into a comprehensive development 
      environment.</description>
   </book>
</catalog>
'>>`

	data := `[(decl: (target: 'xml', text: 'version="1.0"')), (text: '\n'), (elem: (attrs: {}, children: [(text: '\n   '), (elem: (attrs: {(id: 'bk101')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: "XML Developers Guide")], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n   '), (elem: (attrs: {(id: 'bk112')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Galos, Mike')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Visual Studio 7: A Comprehensive Guide')], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '49.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2001-04-16')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Microsoft Visual Studio 7 is explored in depth,\n      looking at how Visual Basic, Visual C++, C#, and ASP+ are \n      integrated into a comprehensive development \n      environment.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.encode("+data+")")
}

func TestXMLEncode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.encode(`woop`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<ldkfjroot>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`:`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root></hi>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<?xx?>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root></root>`)")
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

	expected := `[(elem: (attrs: {}, children: [(text: '\n   '), (elem: (attrs: {(id: 'bk101')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'XML Developers Guide')], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

	data := `<<'<catalog>
   <book id="bk101">
      <author>Gambardella, Matthew</author>
      <title>XML Developers Guide</title>
      <genre>Computer</genre>
      <price>44.95</price>
      <publish_date>2000-10-01</publish_date>
      <description>An in-depth look at creating applications with XML.</description>
   </book>
</catalog>
'>>`

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

	expected := `[(decl: (target: 'xml', text: 'version="1.0"'))]`
	data := `<<'<?xml version="1.0"?>'>>`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+data+`)`)
}

func TestXMLDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decode(`<root>`)")
}

//nolint:lll
func TestXMLDecoder_strip(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0"?>
<catalog>
   <book id="bk101">
      <author>Gambardella, Matthew</author>
      <title>XML Developer's Guide</title>
      <genre>Computer</genre>
      <price>44.95</price>
      <publish_date>2000-10-01</publish_date>
      <description>An in-depth look at creating applications 
      with XML.</description>
   </book>
   <book id="bk112">
      <author>Galos, Mike</author>
      <title>Visual Studio 7: A Comprehensive Guide</title>
      <genre>Computer</genre>
      <price>49.95</price>
      <publish_date>2001-04-16</publish_date>
      <description>Microsoft Visual Studio 7 is explored in depth,
      looking at how Visual Basic, Visual C++, C#, and ASP+ are 
      integrated into a comprehensive development 
      environment.</description>
   </book>
</catalog>
`

	expected := `[(decl: (target: 'xml', text: 'version="1.0"')), (elem: (attrs: {}, children: [(elem: (attrs: {(id: 'bk101')}, children: [(elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (elem: (attrs: {}, children: [(text: "XML Developer's Guide")], name: 'title')), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description'))], name: 'book')), (elem: (attrs: {(id: 'bk112')}, children: [(elem: (attrs: {}, children: [(text: 'Galos, Mike')], name: 'author')), (elem: (attrs: {}, children: [(text: 'Visual Studio 7: A Comprehensive Guide')], name: 'title')), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (elem: (attrs: {}, children: [(text: '49.95')], name: 'price')), (elem: (attrs: {}, children: [(text: '2001-04-16')], name: 'publish_date')), (elem: (attrs: {}, children: [(text: 'Microsoft Visual Studio 7 is explored in depth,\n      looking at how Visual Basic, Visual C++, C#, and ASP+ are \n      integrated into a comprehensive development \n      environment.')], name: 'description'))], name: 'book'))], name: 'catalog'))]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder(true)(`"+xml+"`)")
}

//nolint:lll
func TestXMLDecoder_dont_strip(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0"?>
<catalog>
   <book id="bk101">
      <author>Gambardella, Matthew</author>
      <title>XML Developer's Guide</title>
      <genre>Computer</genre>
      <price>44.95</price>
      <publish_date>2000-10-01</publish_date>
      <description>An in-depth look at creating applications 
      with XML.</description>
   </book>
   <book id="bk112">
      <author>Galos, Mike</author>
      <title>Visual Studio 7: A Comprehensive Guide</title>
      <genre>Computer</genre>
      <price>49.95</price>
      <publish_date>2001-04-16</publish_date>
      <description>Microsoft Visual Studio 7 is explored in depth,
      looking at how Visual Basic, Visual C++, C#, and ASP+ are 
      integrated into a comprehensive development 
      environment.</description>
   </book>
</catalog>
`

	expected := `[(decl: (target: 'xml', text: 'version="1.0"')), (text: '\n'), (elem: (attrs: {}, children: [(text: '\n   '), (elem: (attrs: {(id: 'bk101')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: "XML Developer's Guide")], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n   '), (elem: (attrs: {(id: 'bk112')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Galos, Mike')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Visual Studio 7: A Comprehensive Guide')], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '49.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2001-04-16')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Microsoft Visual Studio 7 is explored in depth,\n      looking at how Visual Basic, Visual C++, C#, and ASP+ are \n      integrated into a comprehensive development \n      environment.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decoder(false)(`"+xml+"`)")
}

func TestXMLDecoder_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.decoder(false)(`<root>`)")
}

//nolint:lll
func TestXMLDecode_round_trip(t *testing.T) {
	t.Parallel()

	xml := `<<'<?xml version="1.0"?>
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
   <book id="bk112">
      <author>Galos, Mike</author>
      <title>Visual Studio 7: A Comprehensive Guide</title>
      <genre>Computer</genre>
      <price>49.95</price>
      <publish_date>2001-04-16</publish_date>
      <description>Microsoft Visual Studio 7 is explored in depth,
      looking at how Visual Basic, Visual C++, C#, and ASP+ are 
      integrated into a comprehensive development 
      environment.</description>
   </book>
</catalog>
'>>
`

	expected := `[(decl: (target: 'xml', text: 'version="1.0"')), (text: '\n'), (elem: (attrs: {}, children: [(text: '\n   '), (elem: (attrs: {(id: 'bk101')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Gambardella, Matthew')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: "XML Developers Guide")], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '44.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2000-10-01')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'An in-depth look at creating applications \n      with XML.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n   '), (elem: (attrs: {(id: 'bk112')}, children: [(text: '\n      '), (elem: (attrs: {}, children: [(text: 'Galos, Mike')], name: 'author')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Visual Studio 7: A Comprehensive Guide')], name: 'title')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Computer')], name: 'genre')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '49.95')], name: 'price')), (text: '\n      '), (elem: (attrs: {}, children: [(text: '2001-04-16')], name: 'publish_date')), (text: '\n      '), (elem: (attrs: {}, children: [(text: 'Microsoft Visual Studio 7 is explored in depth,\n      looking at how Visual Basic, Visual C++, C#, and ASP+ are \n      integrated into a comprehensive development \n      environment.')], name: 'description')), (text: '\n   ')], name: 'book')), (text: '\n')], name: 'catalog')), (text: '\n')]`

	AssertCodesEvalToSameValue(t, expected, "//encoding.xml.decode("+xml+")")
	AssertCodesEvalToSameValue(t, xml, "//encoding.xml.encode("+expected+")")
}
