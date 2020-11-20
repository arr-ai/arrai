package translate

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/pkg/errors"
)

const ProcInst = "decl"
const Directive = "dir"
const CharData = "text"
const Comment = "cmt"
const Element = "elem"

const TextKey = "text"
const NamespaceKey = "ns"
const TargetKey = "target"
const NameKey = "name"
const AttributesKey = "attrs"
const ChildrenKey = "children"
const TypeKey = "type"

func BytesXMLToArrai(bs []byte) (rel.Value, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(bs))

	return parseXMLDFS(decoder)
}

// NOTE: there are subtle differences in a full xml -> arrai -> xml cycle
// 1. xml.CharData when written has escaped strings (looks like for http safety)
// 2. Self-closing tags are automatically expanded
func BytesXMLFromArrai(v rel.Value) (rel.Value, error) {
	var xmlTokens []xml.Token
	var b bytes.Buffer
	encoder := xml.NewEncoder(&b)

	// generate tokens from arrai rel.Value
	xmlTokens, err := unparseXMLDFS(v)
	if err != nil {
		return nil, err
	}

	// load tokens into the encoder
	for _, i := range xmlTokens {
		err := encoder.EncodeToken(i)
		if err != nil {
			return nil, err
		}
	}

	// flush everything into the buffer
	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	return rel.NewBytes(b.Bytes()), nil
}

// NOTE: kind of shabby, both panics and returns and error. inconsisistent error handling
// panic due to MustGet(). Used to reduce LOC
func unparseXMLDFS(v rel.Value) ([]xml.Token, error) {
	var xmlTokens []xml.Token

	arr, ok := rel.AsArray(v)
	if !ok {
		return nil, errors.New("node is not an array")
	}

	for _, val := range arr.Values() {
		tup, ok := val.(rel.Tuple)
		if !ok {
			return nil, errors.New("value is not a tuple")
		}
		vType := tup.MustGet(TypeKey)

		switch RawString(vType) {
		case ProcInst:
			target := tup.MustGet(TargetKey)
			text := tup.MustGet(TextKey)
			xmlTokens = append(xmlTokens, xml.ProcInst{Target: RawString(target), Inst: []byte(RawString(text))})
		case Directive:
			text := tup.MustGet(TextKey)
			var directive xml.Directive = []byte(RawString(text))
			xmlTokens = append(xmlTokens, directive)
		case Comment:
			text := tup.MustGet(TextKey)
			var comment xml.Comment = []byte(RawString(text))
			xmlTokens = append(xmlTokens, comment)
		case CharData:
			// NOTE: for some reason the xml.Encoder escapes the CharData text
			// https://golang.org/src/encoding/xml/marshal.go?s=7625:7671#L192
			text := tup.MustGet(TextKey)
			var chardata xml.CharData = []byte(RawString(text))
			xmlTokens = append(xmlTokens, chardata)
		case Element:
			ns := tup.MustGet(NamespaceKey)
			name := tup.MustGet(NameKey)
			attrs := tup.MustGet(AttributesKey)
			tAttrs, ok := rel.AsArray(attrs)
			if !ok {
				return nil, errors.New("attributes is not an array")
			}
			xmlAttrs := []xml.Attr{}
			for _, attr := range tAttrs.Values() {
				tup, ok := attr.(rel.Tuple)
				if !ok {
					return nil, errors.New("attribute is not a tuple")
				}
				ns := tup.MustGet(NamespaceKey)
				name := tup.MustGet(NameKey)
				value := tup.MustGet(TextKey)
				xmlAttrs = append(xmlAttrs, xml.Attr{
					Name:  xml.Name{Local: RawString(name), Space: RawString(ns)},
					Value: RawString(value),
				})
			}

			xmlName := xml.Name{Local: RawString(name), Space: RawString(ns)}

			// start element dir
			startelement := xml.StartElement{Name: xmlName, Attr: xmlAttrs}
			xmlTokens = append(xmlTokens, startelement)

			// parse child nodes
			children := tup.MustGet(ChildrenKey)
			childTokens, err := unparseXMLDFS(children)
			if err != nil {
				return nil, err
			}
			xmlTokens = append(xmlTokens, childTokens...)

			// on end of element parsing append EndElement
			xmlTokens = append(xmlTokens, xml.EndElement{Name: xmlName})
		}
	}

	return xmlTokens, nil
}

// Helper function for printing
// given value if {} -> "" or {ss} -> "ss"
func RawString(v rel.Value) string {
	set, ok := v.(rel.Set)
	if !ok {
		log.Fatal("value is not a set")
	}
	str, ok := rel.AsString(set)
	if !ok {
		log.Fatal("set is not a string")
	}

	return str.String()
}

// NOTE: encoding/xml only handles well-formed xml. it does not validate the xml structure.
func parseXMLDFS(decoder *xml.Decoder) (rel.Value, error) {
	values := []rel.Value{}

	var token interface{}
	var err error

	for {
		token, err = decoder.Token()
		if err == io.EOF {
			// end of file (break out of loop) (this is fine)
			break
		}

		if err != nil {
			// something fishy happened
			return nil, err
		}

		var tuple rel.Tuple
		// otherwise token should not be nil
		switch t := token.(type) {
		case xml.ProcInst:
			tuple = rel.NewTuple(rel.NewStringAttr(TypeKey, []rune(ProcInst)),
				rel.NewStringAttr(TargetKey, []rune(t.Target)),
				rel.NewStringAttr(TextKey, []rune(string(t.Inst))))
		case xml.Directive:
			tuple = rel.NewTuple(rel.NewStringAttr(TypeKey, []rune(Directive)),
				rel.NewStringAttr(TextKey, []rune(string(t))))
		case xml.CharData:
			// ignore formatting new lines
			if strings.Trim(string(t), " ") == "\n" {
				continue
			}
			tuple = rel.NewTuple(rel.NewStringAttr(TypeKey, []rune(CharData)),
				rel.NewStringAttr(TextKey, []rune(string(t))))
		case xml.Comment:
			tuple = rel.NewTuple(rel.NewStringAttr(TypeKey, []rune(Comment)),
				rel.NewStringAttr(TextKey, []rune(string(t))))
		case xml.StartElement:
			// NOTE: xml.Token() automatically expands self-closing tags. According to:
			// https://stackoverflow.com/questions/57494936/is-there-a-semantical-difference-between-tag-and-tag-tag-in-xml
			// there is no semantic differnce between them

			// recurse for child nodes
			child, err := parseXMLDFS(decoder)
			if err != nil {
				return nil, err
			}

			// parse attributes
			attrs := []rel.Value{}
			for _, attr := range t.Attr {
				attrs = append(attrs, rel.NewTuple(rel.NewStringAttr(TextKey, []rune(attr.Value)),
					rel.NewStringAttr(NamespaceKey, []rune(attr.Name.Space)),
					rel.NewStringAttr(NameKey, []rune(attr.Name.Local))))
			}

			tuple = rel.NewTuple(rel.NewStringAttr(TypeKey, []rune(Element)),
				rel.NewStringAttr(NamespaceKey, []rune(t.Name.Space)),
				rel.NewAttr(AttributesKey, rel.NewArray(attrs...)),
				rel.NewAttr(ChildrenKey, child),
				rel.NewStringAttr(NameKey, []rune(t.Name.Local)))
		case xml.EndElement:
			//  NOTE: xml.Token() guarantees matching Start and End elements (so this will not prematurely exit)
			return rel.NewArray(values...), nil
		}

		values = append(values, tuple)
	}

	return rel.NewArray(values...), nil
}
