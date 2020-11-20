package syntax

import (
	"testing"
)

func TestXMLEncode(t *testing.T) {
	t.Parallel()

}

func TestXMLDecode(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0"?>`

	expected := `[(target: 'xml', text: 'version="1.0"', type: 'decl')]`
	AssertCodesEvalToSameValue(t, expected, `//encoding.xml.decode(`+"`"+xml+"`"+`)`)
}

func TestXMLEncode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", "//encoding.xml.encode(`woop`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root>`)")
	AssertCodeErrors(t, "", "//encoding.xml.encode(`<root>`)")
}
