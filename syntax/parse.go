package syntax

import (
	"fmt"
	"log"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

type ParseContext struct {
	SourceDir string
}

func parseNames(names ast.Branch) []string {
	idents := names["IDENT"].(ast.Many)
	result := make([]string, 0, len(idents))
	for _, ident := range idents {
		result = append(result, ident.One("").(ast.Leaf).Scanner().String())
	}
	return result
}

func parseName(name ast.Branch) string {
	ktype, children := which(name, "IDENT", "STR")
	switch ktype {
	case "IDENT":
		return children.(ast.One).Node.One("").(ast.Leaf).Scanner().String()
	case "STR":
		s := children.(ast.One).Node.One("").(ast.Leaf).Scanner().String()
		return parseArraiString(s)
	default:
		panic(fmt.Errorf("unexpected name: %v", name))
	}
}

// MustParseString parses input string and returns the parsed Expr or panics.
func (pc ParseContext) MustParseString(s string) ast.Branch {
	return pc.MustParse(parser.NewScanner(s))
}

// MustParse parses input and returns the parsed Expr or panics.
func (pc ParseContext) MustParse(s *parser.Scanner) ast.Branch {
	ast, err := pc.Parse(s)
	if err != nil {
		panic(err)
	}
	return ast
}

// ParseString parses input string and returns the parsed Expr or an error.
func (pc ParseContext) ParseString(s string) (ast.Branch, error) {
	return pc.Parse(parser.NewScanner(s))
}

// Parse parses input and returns the parsed Expr or an error.
func (pc ParseContext) Parse(s *parser.Scanner) (ast.Branch, error) {
	rscopes := []rel.Scope{{}}
	v, err := arraiParsers.ParseWithExternals(parser.Rule("expr"), s, parser.ExternalRefs{
		"bind": func(pscope parser.Scope, _ *parser.Scanner) (parser.TreeElement, error) {
			identStr := "."
			if _, ident, has := pscope.GetVal("IDENT"); has {
				identStr = ident.(parser.Scanner).String()
			} else if _, pattern, has := pscope.GetVal("pattern"); has {
				patNode := ast.FromParserNode(arraiParsers.Grammar(), pattern)
				pat := pc.compilePattern(patNode)
				if epat, is := pat.(rel.ExprPattern); is {
					if identExpr, is := epat.Expr.(rel.IdentExpr); is {
						identStr = identExpr.Ident()
					}
				}
				if identStr == "" {
					return nil, nil
				}
			}

			_, exprElt, has := pscope.GetVal("expr@1")
			if !has {
				_, exprElt, has = pscope.GetVal("expr")
				if !has {
					log.Println(pscope.Keys())
					log.Println(pscope)
					panic("wat?")
				}
			}

			exprNode := ast.FromParserNode(arraiParsers.Grammar(), exprElt)
			expr := pc.CompileExpr(exprNode)
			expr = rel.NewExprClosure(rscopes[len(rscopes)-1], expr)
			rscopes = append(rscopes, rscopes[len(rscopes)-1].With(identStr, expr))
			return nil, nil
		},
		"ast": func(scope parser.Scope, input *parser.Scanner) (parser.TreeElement, error) {
			macro, err := pc.unpackMacro(scope, rscopes[len(rscopes)-1])
			if err != nil {
				return nil, err
			}

			subg := wbnf.NewFromAst(rel.ASTNodeFromValue(macro.grammar))
			rule := parser.Rule(macro.ruleName)
			parsers := subg.Compile(subg)

			childast, err := parsers.ParseWithExternals(rule, input, parser.ExternalRefs{
				"*:{()}:": func(scope parser.Scope, _ *parser.Scanner) (parser.TreeElement, error) {
					childast, err := pc.Parse(input)
					switch err.(type) {
					case nil, parser.UnconsumedInputError:
					default:
						return nil, err
					}
					node := ast.ToParserNode(arraiParsers.Grammar(), childast)
					return node, nil
				},
			})
			if err != nil {
				if unconsumed, ok := err.(parser.UnconsumedInputError); ok {
					childast = unconsumed.Result()
					input = unconsumed.Residue()
				} else {
					return nil, err
				}
			}
			if !macro.transform.IsTrue() {
				return ast.NewExtRefTreeElement(parsers.Grammar(), childast), nil
			}

			childastNode := ast.FromParserNode(subg, childast)
			childastValue := rel.ASTNodeToValue(childastNode)
			bodyValue := rel.SetCall(macro.transform, childastValue)

			return parser.Node{Tag: "extref", Children: nil, Extra: ast.Branch{
				"@rule": ast.One{Node: ast.Extra{Data: parser.Rule(macro.ruleName)}},
				"value": ast.One{Node: NewMacroValue(bodyValue, childastNode.Scanner())},
			}}, nil
		},
	})
	//log.Printf("Parse: v = %v", v)
	if err != nil {
		return nil, err
	}
	result := ast.FromParserNode(arraiParsers.Grammar(), v)
	//log.Printf("Parse: result = %v", result)
	if s.String() != "" {
		return result, parser.UnconsumedInput(*s, v)
	}
	return result, nil
}

func parseNest(lhs rel.Expr, branch ast.Branch) rel.Expr {
	attr := branch.One("IDENT").One("").Scanner()
	names := branch["names"].(ast.One).Node.(ast.Branch)["IDENT"].(ast.Many)
	namestrings := make([]string, len(names))
	for i, name := range names {
		namestrings[i] = name.One("").Scanner().String()
	}
	return rel.NewNestExpr(attr, lhs, rel.NewNames(namestrings...), attr.String())
}

// Macro represents the "header" of a macro invocation: the grammar and rule to parse with, and the
// transform to apply to the parsed body.
type Macro struct {
	ruleName  string
	grammar   rel.Tuple
	transform rel.Set
}

// unpackMacro processes the current parser scope when a macro invocation is detected. It extracts
// the details of the macro to invoke, and returns them as a Macro.
func (pc ParseContext) unpackMacro(parseScope parser.Scope, relScope rel.Scope) (Macro, error) {
	_, elt, ok := parseScope.GetVal("macro")
	if !ok {
		panic("wat?")
	}
	macroNode := ast.FromParserNode(arraiParsers.Grammar(), elt)
	macroExpr := pc.CompileExpr(macroNode).(rel.Expr)
	macroValue, err := macroExpr.Eval(relScope)
	if err != nil {
		return Macro{ruleName: ""}, err
	}
	macroTuple := macroValue.(rel.Tuple)

	grammar := macroTuple
	if macroTuple.HasName("@grammar") {
		grammar = macroTuple.MustGet("@grammar").(rel.Tuple)
	}

	var ruleName string
	_, ruleElt, ok := parseScope.GetVal("rule")
	ruleNode := ruleElt.(parser.Node)
	if ruleNode.Count() > 0 {
		ruleNode := ast.FromParserNode(arraiParsers.Grammar(), ruleElt.(parser.Node).Get(0))
		ruleName = ruleNode.One("name").Scanner().String()
	} else {
		ruleName = getFirstRuleName(grammar)
	}

	transform := rel.NewSet()
	if transforms, ok := macroTuple.Get("@transform"); ok {
		// If @transform is present but there is no named transform for ruleName, fail
		// loudly rather than falling back on the default rule or nothing. A macro's
		// rule transforms should be as well-specified as the grammar.
		if transformValue, ok := transforms.(rel.Tuple).Get(ruleName); ok {
			transform = transformValue.(rel.Set)
		} else {
			panic(fmt.Errorf("transform for rule %q not found", ruleName))
		}
	}
	return Macro{ruleName, grammar, transform}, nil
}

// getFirstRuleName finds the first rule declared in grammar and returns its name.
func getFirstRuleName(grammar rel.Tuple) string {
	stmts := grammar.MustGet("stmt").(rel.Array).Values()
	for _, stmt := range stmts {
		if prod, ok := stmt.(rel.Tuple).Get("prod"); ok {
			return prod.(rel.Tuple).MustGet("IDENT").(rel.Tuple).MustGet("").String()
		}
	}
	panic("no prod rule found in grammar")
}

// MacroValue is an Extra node with an Expr value and a Scanner for the macro source.
type MacroValue struct {
	ast.Extra
	expr    rel.Expr
	scanner parser.Scanner
}

// NewMacroValue returns a MacroValue with a given Expr and Scanner.
func NewMacroValue(expr rel.Expr, scanner parser.Scanner) MacroValue {
	return MacroValue{expr: expr, scanner: scanner}
}

// Scanner returns a scanner of the source that was replaced by the macro.
func (m MacroValue) Scanner() parser.Scanner {
	return m.scanner
}

// SubExpr returns the Expr resulting from evaluating the macro.
func (m MacroValue) SubExpr() rel.Expr {
	return m.expr
}
