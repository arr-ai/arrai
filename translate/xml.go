package translate

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/pkg/errors"
)

const procInst = "decl"
const directive = "directive"
const charData = "text"
const comment = "comment"
const element = "elem"

const textKey = "text"
const nameKey = "name"
const targetKey = "target"
const attributesKey = "attrs"
const childrenKey = "children"

type XMLDecodeConfig struct {
	StripFormatting bool
}

func BytesXMLToArrai(bs []byte, config XMLDecodeConfig) (rel.Value, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(bs))

	return parseXMLDFS(decoder, config)
}

// NOTE: there are subtle differences in a full xml -> arrai -> xml cycle
// 1. xml.CharData when written has escaped strings (looks like for http safety)
// 2. Self-closing tags are automatically expanded
func BytesXMLFromArrai(v rel.Value) (rel.Value, error) {
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

		if tup.Names().Count() != 1 {
			return nil, errors.New("tuple has multiple attributes")
		}

		// assume there is only a single attribute in the set
		switch tup.Names().TheOne() {
		case procInst:
			mTup, ok := tup.MustGet(procInst).(rel.Tuple)
			if !ok {
				return nil, errors.New("value is not a tuple")
			}

			target, ok := mTup.Get(targetKey)
			if !ok {
				return nil, fmt.Errorf("attribute does not exist: %s", targetKey)
			}
			text, ok := mTup.Get(textKey)
			if !ok {
				return nil, fmt.Errorf("attribute does not exist: %s", textKey)
			}
			rawTarget, err := RawString(target)
			if err != nil {
				return nil, err
			}
			rawText, err := RawString(text)
			if err != nil {
				return nil, err
			}
			xmlTokens = append(xmlTokens, xml.ProcInst{Target: rawTarget, Inst: []byte(rawText)})
		case directive:
			text := tup.MustGet(directive)
			rawText, err := RawString(text)
			if err != nil {
				return nil, err
			}
			var directive xml.Directive = []byte(rawText)
			xmlTokens = append(xmlTokens, directive)
		case comment:
			text := tup.MustGet(comment)
			rawText, err := RawString(text)
			if err != nil {
				return nil, err
			}
			var comment xml.Comment = []byte(rawText)
			xmlTokens = append(xmlTokens, comment)
		case charData:
			// NOTE: for some reason the xml.Encoder escapes the CharData text
			// https://golang.org/src/encoding/xml/marshal.go?s=7625:7671#L192
			text := tup.MustGet(charData)
			rawText, err := RawString(text)
			if err != nil {
				return nil, err
			}
			var chardata xml.CharData = []byte(rawText)
			xmlTokens = append(xmlTokens, chardata)
		case element:
			tup, ok := tup.MustGet(element).(rel.Tuple)
			if !ok {
				return nil, errors.New("value is not a tuple")
			}

			name, ok := tup.Get(nameKey)
			if !ok {
				return nil, fmt.Errorf("attribute does not exist: %s", nameKey)
			}
			children, ok := tup.Get(childrenKey)
			if !ok {
				return nil, fmt.Errorf("attribute does not exist: %s", childrenKey)
			}
			// attributes are ommited if empty
			attrs, attrOk := tup.Get(attributesKey)
			xmlAttrs := []xml.Attr{}

			// load attributes
			if attrOk {
				tAttrs, ok := attrs.(rel.Dict)
				log.Print(tAttrs)
				if !ok {
					return nil, errors.New("attributes is not a dictionary")
				}
				enum := tAttrs.DictEnumerator()
				for enum.MoveNext() {
					key, value := enum.Current()
					rawKey, err := RawString(key)
					if err != nil {
						return nil, err
					}
					rawValue, err := RawString(value)
					if err != nil {
						return nil, err
					}
					xmlAttrs = append(xmlAttrs, xml.Attr{
						Name:  xmlNameFromArrai(rawKey),
						Value: rawValue,
					})
				}
			}

			rawName, err := RawString(name)
			if err != nil {
				return nil, err
			}
			xmlName := xmlNameFromArrai(rawName)

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
func RawString(v rel.Value) (string, error) {
	set, ok := v.(rel.Set)
	if !ok {
		return "", errors.New("value is not a set")
	}
	str, ok := rel.AsString(set)
	if !ok {
		return "", errors.New("set is not a string")
	}

	return str.String(), nil
}

// NOTE: encoding/xml only handles somewhat well-formed xml. it does not validate the xml structure.
func parseXMLDFS(decoder *xml.Decoder, config XMLDecodeConfig) (rel.Value, error) {
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
				rel.NewTupleAttr(procInst,
					rel.NewStringAttr(targetKey, []rune(t.Target)),
					rel.NewStringAttr(textKey, []rune(string(t.Inst))),
				),
			)
		case xml.Directive:
			tuple = rel.NewTuple(rel.NewStringAttr(directive, []rune(string(t))))
		case xml.CharData:
			// ignore formatting new lines
			if config.StripFormatting && strings.Trim(string(t), " ") == "\n" {
				continue
			}
			tuple = rel.NewTuple(rel.NewStringAttr(charData, []rune(string(t))))
		case xml.Comment:
			tuple = rel.NewTuple(rel.NewStringAttr(comment, []rune(string(t))))
		case xml.StartElement:
			// NOTE: xml.Token() automatically expands self-closing tags. According to:
			// https://stackoverflow.com/questions/57494936/is-there-a-semantical-difference-between-tag-and-tag-tag-in-xml
			// there is no semantic differnce between them

			// recurse for child nodes
			child, err := parseXMLDFS(decoder, config)
			if err != nil {
				return nil, err
			}

			// parse attributes
			xmlAttrs := []rel.DictEntryTuple{}
			for _, attr := range t.Attr {
				xmlAttrs = append(xmlAttrs, rel.NewDictEntryTuple(
					rel.NewString([]rune(xmlNameToArrai(&attr.Name))),
					rel.NewString([]rune(attr.Value)),
				))
			}
			xmlAttrDict, err := rel.NewDict(false, xmlAttrs...)
			if err != nil {
				return nil, err
			}

			// element tuple attributes
			attrList := []rel.Attr{}
			attrList = append(attrList, rel.NewStringAttr(nameKey, []rune(xmlNameToArrai(&t.Name))))
			attrList = append(attrList, rel.NewAttr(childrenKey, child))
			// add attributes if there are some
			if xmlAttrDict.IsTrue() {
				attrList = append(attrList, rel.NewAttr(attributesKey, xmlAttrDict))
			}

			tuple = rel.NewTuple(rel.NewTupleAttr(element, attrList...))
		case xml.EndElement:
			//  NOTE: xml.Token() guarantees matching Start and End elements (so this will not prematurely exit)
			return rel.NewArray(values...), nil
		}

		values = append(values, tuple)
	}

	return rel.NewArray(values...), nil
}

func xmlNameToArrai(name *xml.Name) string {
	if len(name.Space) > 0 {
		return name.Space + ":" + name.Local
	}

	return name.Local
}

// XML names should only contain : as a namespace character https://www.w3.org/TR/xml/#NT-S
// technically it can prepent a localname but i doubt authors would want to do that
// assume the string is in the format "namespace:localname" or "localname" (if there is no namespace)
func xmlNameFromArrai(name string) xml.Name {
	var s = strings.Split(name, ":")
	if len(s) == 1 {
		return xml.Name{Local: s[0], Space: ""}
	}

	return xml.Name{Local: s[1], Space: s[0]}
}
