package translate

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

const procInst = "decl"
const directive = "directive"
const charData = "text"
const comment = "comment"
const element = "elem"

const textKey = "text"
const nameKey = "name"
const nsKey = "ns"
const targetKey = "target"
const attributesKey = "attrs"
const childrenKey = "children"

// NOTE: Currently the XML transform does not support documents with explicit namespaces.
// NOTE: A full cycle from XML -> Arr.ai -> XML reproduces semantically similar documents
//       with possibly different content.
type XMLDecodeConfig struct {
	TrimSurroundingWhitespace bool
}

// BytesXMLToArrai converts a well formatted XML document in byte representation
// to a structured Arr.ai object
func BytesXMLToArrai(bs []byte, config XMLDecodeConfig) (rel.Value, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(bs))

	return parseXML(decoder, config)
}

// BytesXMLFromArrai converts an Arr.ai object to an XML document in byte resentation
func BytesXMLFromArrai(v rel.Value) (rel.Value, error) {
	var b bytes.Buffer
	encoder := xml.NewEncoder(&b)

	// generate tokens from arrai rel.Value
	xmlTokens, err := unparseXML(v)
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

	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	return rel.NewBytes(b.Bytes()), nil
}

func unparseXML(v rel.Value) ([]xml.Token, error) {
	var xmlTokens []xml.Token

	arr, ok := rel.AsArray(v)
	if !ok {
		return nil, fmt.Errorf("value must be an array, not %s: %v", rel.ValueTypeAsString(v), v)
	}

	for _, val := range arr.Values() {
		tup, ok := val.(rel.Tuple)
		if !ok {
			return nil, fmt.Errorf("value must be tuple, not %s: %v", rel.ValueTypeAsString(val), val)
		}

		if tup.Names().Count() != 1 {
			return nil, fmt.Errorf("tuple has multiple attributes: %v", tup)
		}

		// assume there is only a single attribute in the set
		switch tup.Names().TheOne() {
		case procInst:
			val := tup.MustGet(procInst)
			mTup, ok := val.(rel.Tuple)
			if !ok {
				return nil, fmt.Errorf("value must be tuple, not %s: %v", rel.ValueTypeAsString(val), val)
			}

			target, ok := mTup.Get(targetKey)
			if !ok {
				return nil, fmt.Errorf("tuple attribute missing: %s", targetKey)
			}
			text, ok := mTup.Get(textKey)
			if !ok {
				return nil, fmt.Errorf("tuple attribute missing: %s", textKey)
			}
			rawTarget, ok := tools.ValueAsString(target)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", target)
			}
			rawText, ok := tools.ValueAsString(text)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", text)
			}
			xmlTokens = append(xmlTokens, xml.ProcInst{Target: rawTarget, Inst: []byte(rawText)})
		case directive:
			text := tup.MustGet(directive)
			rawText, ok := tools.ValueAsString(text)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", text)
			}
			var directive xml.Directive = []byte(rawText)
			xmlTokens = append(xmlTokens, directive)
		case comment:
			text := tup.MustGet(comment)
			rawText, ok := tools.ValueAsString(text)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", text)
			}
			var comment xml.Comment = []byte(rawText)
			xmlTokens = append(xmlTokens, comment)
		case charData:
			// NOTE: for some reason the xml.Encoder escapes the CharData text
			// https://golang.org/src/encoding/xml/marshal.go?s=7625:7671#L192
			text := tup.MustGet(charData)
			rawText, ok := tools.ValueAsString(text)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", text)
			}
			var chardata xml.CharData = []byte(rawText)
			xmlTokens = append(xmlTokens, chardata)
		case element:
			val := tup.MustGet(element)
			tup, ok := val.(rel.Tuple)
			if !ok {
				return nil, fmt.Errorf("value must be tuple, not %s: %v", rel.ValueTypeAsString(val), val)
			}

			name, ok := tup.Get(nameKey)
			if !ok {
				return nil, fmt.Errorf("tuple attribute missing: %s", nameKey)
			}
			children, ok := tup.Get(childrenKey)
			if !ok {
				return nil, fmt.Errorf("tuple attribute missing: %s", childrenKey)
			}
			attrs, ok := tup.Get(attributesKey)
			if !ok {
				return nil, fmt.Errorf("tuple attribute missing: %s", attributesKey)
			}

			xmlAttrs := []xml.Attr{}
			tAttrs, ok := attrs.(rel.Set)
			if !ok {
				return nil, fmt.Errorf("value must be a set, not %s: %v", rel.ValueTypeAsString(attrs), attrs)
			}
			enum := tAttrs.Enumerator()
			for enum.MoveNext() {
				tup, ok := enum.Current().(rel.Tuple)
				if !ok {
					return nil, fmt.Errorf("value must be tuple, not %s: %v", rel.ValueTypeAsString(tup), tup)
				}
				tupName, ok := tup.Get(nameKey)
				if !ok {
					return nil, fmt.Errorf("tuple attribute missing: %s", nameKey)
				}
				tupValue, ok := tup.Get(textKey)
				if !ok {
					return nil, fmt.Errorf("tuple attribute missing: %s", textKey)
				}
				tupNS, ok := tup.Get(nsKey)
				if !ok {
					tupNS = rel.NewString([]rune(""))
				}

				xmlValue, ok := tools.ValueAsString(tupValue)
				if !ok {
					return nil, fmt.Errorf("value cannot be converted to string: %s", tupValue)
				}
				xmlName, ok := tools.ValueAsString(tupName)
				if !ok {
					return nil, fmt.Errorf("value cannot be converted to string: %s", tupName)
				}
				xmlNS, ok := tools.ValueAsString(tupNS)
				if !ok {
					return nil, fmt.Errorf("value cannot be converted to string: %s", tupNS)
				}

				xmlAttrs = append(xmlAttrs, xml.Attr{
					Name: xml.Name{
						Local: xmlName,
						Space: xmlNS,
					},
					Value: xmlValue,
				})
			}

			rawName, ok := tools.ValueAsString(name)
			if !ok {
				return nil, fmt.Errorf("value cannot be converted to string: %s", name)
			}
			// NOTE: namespace does not need to be populated because the encoding/xml does not handle xml prefixes correctly.
			// namespace attributes are parsed in the attributes tuple
			xmlName := xml.Name{
				Local: rawName,
				Space: "",
			}

			startelement := xml.StartElement{Name: xmlName, Attr: xmlAttrs}
			xmlTokens = append(xmlTokens, startelement)

			// parse child nodes
			childTokens, err := unparseXML(children)
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

// Parses XML via the Go standard library "encoding/xml" tokeniser into an arr.ai structure.
// NOTE: encoding/xml only handles well-formed XML. It does not validate the XML structure.
func parseXML(decoder *xml.Decoder, config XMLDecodeConfig) (rel.Value, error) {
	values := []rel.Value{}

	var token interface{}
	var err error

	for {
		token, err = decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var tuple rel.Tuple
		switch t := token.(type) {
		case xml.ProcInst:
			tuple = rel.NewTuple(
				rel.NewTupleAttr(procInst,
					rel.NewStringAttr(targetKey, []rune(t.Target)),
					rel.NewStringAttr(textKey, []rune(string(t.Inst))),
				),
			)
		case xml.Directive:
			tuple = rel.NewTuple(rel.NewStringAttr(directive, []rune(string(t))))
		case xml.CharData:
			// ignore formatting new lines, tabs and spaces
			if config.TrimSurroundingWhitespace && strings.TrimSpace(string(t)) == "" {
				continue
			}
			tuple = rel.NewTuple(rel.NewStringAttr(charData, []rune(string(t))))
		case xml.Comment:
			tuple = rel.NewTuple(rel.NewStringAttr(comment, []rune(string(t))))
		case xml.StartElement:
			// NOTE: xml.Token() automatically expands self-closing tags. According to:
			// https://stackoverflow.com/questions/57494936/is-there-a-semantical-difference-between-tag-and-tag-tag-in-xml
			// there is no semantic difference between them

			// parse attributes
			xmlAttrs := []rel.Value{}
			for _, attr := range t.Attr {
				tupList := []rel.Attr{}
				tupList = append(tupList, rel.NewStringAttr(nameKey, []rune(attr.Name.Local)))
				tupList = append(tupList, rel.NewStringAttr(textKey, []rune(attr.Value)))
				if len(attr.Name.Space) > 0 {
					tupList = append(tupList, rel.NewStringAttr(nsKey, []rune(attr.Name.Space)))
				}
				xmlAttrs = append(xmlAttrs, rel.NewTuple(tupList...))
			}

			// recurse for child nodes
			child, err := parseXML(decoder, config)
			if err != nil {
				return nil, err
			}

			xmlAttrSet, err := rel.NewSet(xmlAttrs...)
			if err != nil {
				return nil, err
			}

			// element tuple attributes
			attrList := []rel.Attr{}
			if len(t.Name.Space) > 0 {
				attrList = append(attrList, rel.NewStringAttr(nsKey, []rune(t.Name.Space)))
			}
			attrList = append(attrList, rel.NewStringAttr(nameKey, []rune(t.Name.Local)))
			attrList = append(attrList, rel.NewAttr(childrenKey, child))
			attrList = append(attrList, rel.NewAttr(attributesKey, xmlAttrSet))

			tuple = rel.NewTuple(rel.NewTupleAttr(element, attrList...))
		case xml.EndElement:
			//  NOTE: xml.Token() guarantees matching Start and End elements (so this will not prematurely exit)
			return rel.NewArray(values...), nil
		}

		values = append(values, tuple)
	}

	return rel.NewArray(values...), nil
}
