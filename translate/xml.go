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
const Directive = "directive"
const CharData = "text"
const Comment = "comment"
const Element = "elem"

const TextKey = "text"
const NameKey = "name"
const TargetKey = "target"
const AttributesKey = "attrs"
const ChildrenKey = "children"

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

		// assume there is only a single attribute in the set
		switch tup.Names().TheOne() {
		case ProcInst:
			mTup, ok := tup.MustGet(ProcInst).(rel.Tuple)
			if !ok {
				return nil, errors.New("value is not a tuple")
			}

			target := mTup.MustGet(TargetKey)
			text := mTup.MustGet(TextKey)
			xmlTokens = append(xmlTokens, xml.ProcInst{Target: RawString(target), Inst: []byte(RawString(text))})
		case Directive:
			text := tup.MustGet(Directive)
			var directive xml.Directive = []byte(RawString(text))
			xmlTokens = append(xmlTokens, directive)
		case Comment:
			text := tup.MustGet(Comment)
			var comment xml.Comment = []byte(RawString(text))
			xmlTokens = append(xmlTokens, comment)
		case CharData:
			// NOTE: for some reason the xml.Encoder escapes the CharData text
			// https://golang.org/src/encoding/xml/marshal.go?s=7625:7671#L192
			text := tup.MustGet(CharData)
			var chardata xml.CharData = []byte(RawString(text))
			xmlTokens = append(xmlTokens, chardata)
		case Element:
			tup, ok := tup.MustGet(Element).(rel.Tuple)
			if !ok {
				return nil, errors.New("value is not a tuple")
			}

			name := tup.MustGet(NameKey)
			attrs := tup.MustGet(AttributesKey)
			children := tup.MustGet(ChildrenKey)

			// load attributes
			tAttrs, ok := attrs.(rel.Set)
			if !ok {
				return nil, errors.New("attributes is not a set")
			}
			xmlAttrs := []xml.Attr{}
			enum := tAttrs.Enumerator()
			for enum.MoveNext() {
				tup, ok := enum.Current().(rel.Tuple)
				if !ok {
					return nil, errors.New("attribute is not a tuple")
				}
				name := tup.Names().TheOne()
				value := tup.MustGet(name)
				xmlAttrs = append(xmlAttrs, xml.Attr{
					Name:  xmlNameFromArrai(name),
					Value: RawString(value),
				})
			}

			xmlName := xmlNameFromArrai(RawString(name))

			// start element dir
			startelement := xml.StartElement{Name: xmlName, Attr: xmlAttrs}
			xmlTokens = append(xmlTokens, startelement)

			// parse child nodes
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
			tuple = rel.NewTuple(
				rel.NewTupleAttr(ProcInst,
					rel.NewStringAttr(TargetKey, []rune(t.Target)),
					rel.NewStringAttr(TextKey, []rune(string(t.Inst))),
				),
			)
		case xml.Directive:
			tuple = rel.NewTuple(rel.NewStringAttr(Directive, []rune(string(t))))
		case xml.CharData:
			// ignore formatting new lines
			if strings.Trim(string(t), " ") == "\n" {
				continue
			}
			tuple = rel.NewTuple(rel.NewStringAttr(CharData, []rune(string(t))))
		case xml.Comment:
			tuple = rel.NewTuple(rel.NewStringAttr(Comment, []rune(string(t))))
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
				attrs = append(attrs, rel.NewTuple(
					rel.NewStringAttr(xmlNameToArrai(&attr.Name), []rune(attr.Value)),
				))
			}
			attrSet, err := rel.NewSet(attrs...)
			if err != nil {
				return nil, err
			}

			tuple = rel.NewTuple(rel.NewTupleAttr(Element,
				rel.NewStringAttr(NameKey, []rune(xmlNameToArrai(&t.Name))),
				rel.NewAttr(AttributesKey, attrSet),
				rel.NewAttr(ChildrenKey, child),
			))
		case xml.EndElement:
			//  NOTE: xml.Token() guarantees matching Start and End elements (so this will not prematurely exit)
			return rel.NewArray(values...), nil
		}

		values = append(values, tuple)
	}

	return rel.NewArray(values...), nil
}

func xmlNameToArrai(name *xml.Name) string {
	return name.Space + ":" + name.Local
}

// assume the string is in the format "namespace:localname"
func xmlNameFromArrai(name string) xml.Name {
	var s = strings.Split(name, ":")
	return xml.Name{Local: s[1], Space: s[0]}
}
