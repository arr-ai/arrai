package syntax

import (
	"fmt"
	"log"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

type MacroResult struct {
	Data rel.Value
}

func (MacroResult) IsExtra() {}

// type noParseType struct{}

// type parseFunc func(v interface{}) (rel.Expr, error)

// func (*noParseType) Error() string {
// 	return "No parse"
// }

// var noParse = &noParseType{}

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
			ruleName := "default"
			if _, ruleElt, ok := scope.GetVal("rule"); ok && ruleElt.(parser.Node).Count() > 0 {
				ruleNode := ast.FromParserNode(arraiParsers.Grammar(), ruleElt.(parser.Node).Get(0))
				ruleName = ruleNode.One("rulename").Scanner().String()
			}

			_, elt, ok := scope.GetVal("macro")
			if !ok {
				panic("wat?")
			}

			astNode := ast.FromParserNode(arraiParsers.Grammar(), elt)
			astExpr := pc.CompileExpr(astNode).(rel.Expr)
			astValue, err := astExpr.Eval(rscopes[len(rscopes)-1])
			if err != nil {
				return nil, err
			}

			transform := rel.NewSet()
			if macroGrammar, ok := astValue.(*rel.GenericTuple).Get("@grammar"); ok {
				if transforms, ok := astValue.(*rel.GenericTuple).Get("@transform"); ok {
					transform = transforms.(*rel.GenericTuple).MustGet(ruleName).(rel.Set)
				}
				astValue = macroGrammar
			}

			astValueNode := rel.ASTNodeFromValue(astValue).(ast.Branch)
			subg := wbnf.NewFromAst(astValueNode)
			rule := parser.Rule(ruleName)
			parsers := subg.Compile(subg)

			childast, err := parsers.ParseWithExternals(rule, input, parser.ExternalRefs{
				"*:{()}:": func(scope parser.Scope, _ *parser.Scanner) (parser.TreeElement, error) {
					childast, err := pc.Parse(input)
					switch err.(type) {
					case nil, parser.UnconsumedInputError:
					default:
						return nil, err
					}
					// log.Printf("ast: %v", ast)
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
			if !transform.IsTrue() {
				return ast.NewExtRefTreeElement(parsers.Grammar(), childast), nil
			}

			childastValue := rel.ASTNodeToValue(ast.FromParserNode(subg, childast))
			bodyValue := rel.SetCall(transform, childastValue)

			result := ast.Branch{}
			result["@rule"] = ast.One{Node: ast.Extra{parser.Rule(ruleName)}}
			result["value"] = ast.One{Node: ast.Extra{bodyValue}}
			return parser.Node{
				Tag:      "extref",
				Extra:    result,
				Children: nil,
			}, nil
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

//func FromParserNode(g parser.Grammar, v parser.Extra) ast.Branch {
//	rule := parser.Rule("default")
//	term := g[rule]
//	result := ast.Branch{}
//	extra := v
//	result["@rule"] = ast.One{Node: extra}
//	ctrs := newCounters(term)
//	result.fromParserNode(g, term, ctrs, e)
//	return result.collapse(0).(ast.Branch)
//}

func parseNest(lhs rel.Expr, branch ast.Branch) rel.Expr {
	attr := branch.One("IDENT").One("").Scanner()
	names := branch["names"].(ast.One).Node.(ast.Branch)["IDENT"].(ast.Many)
	namestrings := make([]string, len(names))
	for i, name := range names {
		namestrings[i] = name.One("").Scanner().String()
	}
	return rel.NewNestExpr(attr, lhs, rel.NewNames(namestrings...), attr.String())
}
