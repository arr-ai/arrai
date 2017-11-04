package syntax

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/mediocregopher/seq"

	"github.com/arr-ai/arrai/rel"
)

// XML tokens
const (
	xmlOpenElt = xmlTokenBase + iota
	xmlAttrName
	xmlCloseNullElt
	xmlEndElt
	xmlCharData
)

var wsRE = regexp.MustCompile(`\A(\s+)()`)

type xmlContext struct {
	keepWS     bool
	namespaces *seq.HashMap
}

var xmlNS = "{https://www.w3.org/XML/1998/namespace}"
var xmlnsNS = "{http://www.w3.org/2000/xmlns/}"
var xmlSpaceAttr = xmlNS + "space"

var rootNS = seq.NewHashMap(
	&seq.KV{"", ""},
	&seq.KV{"xml", xmlNS},
	&seq.KV{"xmlns", xmlNS},
	&seq.KV{"arr.ai", "arr.ai"},
)

func newXMLContext() xmlContext {
	return xmlContext{
		false,
		rootNS,
	}
}

func (xc xmlContext) withAttrExprs(attrs ...rel.AttrExpr) (xmlContext, error) {
	nses := xc.namespaces
	for _, attr := range attrs {
		name := attr.Name()
		var alias string
		if name == "xmlns" {
			alias = ""
		} else if strings.HasPrefix(name, xmlnsNS) {
			alias = name[len(xmlnsNS):]
		} else {
			continue
		}
		if s, ok := rel.GetStringValue(attr.Expr()); ok {
			s = "{" + s + "}"
			if alias == "xml" && s != xmlNS {
				return xmlContext{}, errors.Errorf(
					"xml namespace must be %s", xmlNS[1:len(xmlNS)-1])
			}
			nses, _ = nses.Set(alias, s)
		} else {
			return xmlContext{}, errors.Errorf(
				"xmlns:... attr must be a string literal, not %s", attr)
		}
	}
	return xmlContext{xc.keepWS, nses}, nil
}

func (xc xmlContext) apply(name string, useDefault bool) (string, error) {
	if name[:1] == "{" {
		return name, nil
	}
	if colon := strings.IndexByte(name, byte(':')); colon != -1 {
		nsPrefix := name[:colon]
		ident := name[colon+1:]
		if ns, found := xc.namespaces.Get(nsPrefix); found {
			return ns.(string) + ident, nil
		}
		return "", errors.Errorf("Namespace %s not found", nsPrefix)
	}
	if !useDefault {
		return name, nil
	}
	if defaultNS, found := xc.namespaces.Get(""); found {
		return defaultNS.(string) + name, nil
	}
	panic("Missing default namespace")
}

// parseXML encodes XML syntax as arr.ai values.
//
// Example:
//   <div style={font-weight: "bold"}><b>Hello</b> world!</div>
//   ==
//   {@xml: {
//       tag: "div",
//       attributes: {style: {font-weight: "bold"}},
//       children: [
//         {tag: "b", children: ["Hello"]}
//         " world!",
//       ],
//   }}
//
func parseXML(l *Lexer, xc xmlContext) (rel.Expr, error) {
	l.PushState(xmlLexerState)
	l.eatRE(wsRE)
	defer l.PopState()

	if !l.Scan(xmlOpenElt) {
		return nil, noParse
	}

	tag := l.Data().(string)

	attrs, err := parseXMLAttributes(l, tag)
	if err != nil {
		return nil, err
	}

	// Extract new namespaces.
	xc, err = xc.withAttrExprs(attrs...)
	if err != nil {
		return nil, err
	}

	// Apply namespaces to tag.
	tag, err = xc.apply(tag, true)
	if err != nil {
		return nil, err
	}

	// Apply namespaces to attributes.
	nsAttrs := make(map[string]rel.Expr, len(attrs))
	for _, attr := range attrs {
		name, err := xc.apply(attr.Name(), false)
		if err != nil {
			return nil, err
		}
		expr := attr.Expr()
		nsAttrs[name] = expr

		if name == xmlSpaceAttr {
			if s, ok := rel.GetStringValue(expr); ok {
				switch s {
				case "default":
					xc.keepWS = false
				case "preserve":
					xc.keepWS = true
				default:
					return nil, errors.Errorf("Bad value xml:space=%q", s)
				}
			} else {
				return nil, expecting(l, "after "+xmlSpaceAttr+"=",
					`"default"`, `"preserve"`)
			}
		}
	}

	children, err := parseXMLChildren(l, tag, xc)
	if err != nil {
		return nil, err
	}

	parts := make([]rel.AttrExpr, 0, 3)
	tagAttr, _ := rel.NewAttrExpr("tag", rel.NewString([]rune(tag)))
	parts = append(parts, tagAttr)
	if len(nsAttrs) != 0 {
		attrs := rel.NewTupleExprFromMap(nsAttrs)
		attr, _ := rel.NewAttrExpr("attributes", attrs)
		parts = append(parts, attr)
	}
	if len(children) != 0 {
		attr, _ := rel.NewAttrExpr("children", rel.NewArrayExpr(children...))
		parts = append(parts, attr)
	}

	xmlExpr, _ := rel.NewAttrExpr("@xml", rel.NewTupleExpr(parts...))

	return rel.NewTupleExpr(xmlExpr), nil
}

func parseXMLAttributes(l *Lexer, tag string) ([]rel.AttrExpr, error) {
	l.PushState(xmlLexerAttrState)
	defer l.PopState()

	attrs := []rel.AttrExpr{}

	for l.Scan(xmlAttrName) {
		name := l.Data().(string)
		if l.Scan(Token('=')) {
			l.PushState(LexerInitState)
			expr, err := parsePrefix(l)
			l.PopState()
			if err != nil {
				return nil, err
			}

			if strings.HasPrefix(name, "xmlns:") {
				name = xmlnsNS + name[6:]
			}

			attr, _ := rel.NewAttrExpr(name, expr)
			attrs = append(attrs, attr)
		} else {
			boolAttr, _ := rel.NewAttrExpr(name, rel.True)
			attrs = append(attrs, boolAttr)
		}
	}
	if !l.Scan(Token('>'), xmlCloseNullElt) {
		return nil, expecting(l, "after xml attrs", "'>'", "'/>'")
	}
	return attrs, nil
}

func parseXMLChildren(
	l *Lexer, tag string, xc xmlContext,
) ([]rel.Expr, error) {
	children := []rel.Expr{}
	if l.Token() == Token('>') {
		for !l.Scan(xmlEndElt) {
			switch l.Peek() {
			case xmlOpenElt:
				element, err := parseXML(l, xc)
				if err != nil {
					return nil, err
				}
				children = append(children, element)
			case xmlCharData:
				l.Scan()
				cdata := html.UnescapeString(l.Data().(string))
				var trimmed string
				if xc.keepWS {
					trimmed = cdata
				} else {
					trimmed = strings.Trim(cdata, "\t\n\f\r ")
				}
				if trimmed != "" {
					children = append(children, rel.NewString([]rune(trimmed)))
				}
			case Token('{'):
				l.Scan()
				expr, err := func() (rel.Expr, error) {
					l.PushState(LexerInitState)
					defer l.PopState()

					expr, err := parseExpr(l)
					if err != nil {
						return nil, err
					}
					if !l.Scan(Token('}')) {
						return nil, expecting(l, "after xml-nested expr", "'}'")
					}
					return expr, nil
				}()
				if err != nil {
					return nil, err
				}
				children = append(children, expr)
			default:
				return nil, expecting(l, "after attrs", "<tag", "text", "'{'")
			}
		}

		endTag, err := xc.apply(l.Data().(string), true)
		if err != nil {
			return nil, err
		}
		if endTag != tag {
			return nil, expecting(l,
				fmt.Sprintf("after <%s>", tag),
				fmt.Sprintf("matching </%s>, not </%s>", tag, endTag))
		}
	}
	return children, nil
}

func xmlNameFromMatch(tok Token, match [][]byte) (interface{}, Token) {
	return string(match[2]), tok
}

// Like lex.go:tokRE(), but with an empty capture in lieu of leading whitespace.
func xmlRE(re string) *regexp.Regexp {
	return regexp.MustCompile(`\A()` + re)
}

var xmlLexerOperatorRE = xmlRE(`({)`)

var xmlIdentPat = `[$@A-Za-z_][-.0-9$@A-Za-z_]*`
var xmlNamePat = `(?:` + xmlIdentPat + `:)?` + xmlIdentPat + ``

var xmlLexerSymbols = []LexerSymbol{
	{xmlOpenElt, "xmlOpenElt", xmlRE(`<(` + xmlNamePat + `)`),
		xmlNameFromMatch},
	{xmlEndElt, "xmlEndElt", xmlRE(`</(` + xmlNamePat + `)>`),
		xmlNameFromMatch},
	{xmlCharData, "xmlCharData", xmlRE(`([^<{]+)`),
		func(tok Token, match [][]byte) (interface{}, Token) {
			return string(match[0]), tok
		}},
}

func xmlLexerState(l *Lexer) (Token, interface{}) {
	return l.ScanOperatorOrSymbol(xmlLexerOperatorRE, xmlLexerSymbols)
}

var xmlLexerAttrOperatorRe = tokRE(`[=>]`)

var xmlLexerAttrSymbols = []LexerSymbol{
	{xmlAttrName, "xmlAttrAssign", tokRE(xmlNamePat),
		func(tok Token, match [][]byte) (interface{}, Token) {
			return string(match[2]), tok
		}},
	{xmlCloseNullElt, "xmlCloseNullElt", tokRE(`/>`), nil},
}

func xmlLexerAttrState(l *Lexer) (Token, interface{}) {
	return l.ScanOperatorOrSymbol(xmlLexerAttrOperatorRe, xmlLexerAttrSymbols)
}
