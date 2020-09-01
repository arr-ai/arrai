package syntax

import (
	"context"
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
func (pc ParseContext) MustParseString(ctx context.Context, s string) ast.Branch {
	return pc.MustParse(ctx, parser.NewScanner(s))
}

// MustParse parses input and returns the parsed Expr or panics.
func (pc ParseContext) MustParse(ctx context.Context, s *parser.Scanner) ast.Branch {
	ast, err := pc.Parse(ctx, s)
	if err != nil {
		panic(err)
	}
	return ast
}

// ParseString parses input string and returns the parsed Expr or an error.
func (pc ParseContext) ParseString(ctx context.Context, s string) (ast.Branch, error) {
	return pc.Parse(ctx, parser.NewScanner(s))
}

// Parse parses input and returns the parsed Expr or an error.
func (pc ParseContext) Parse(ctx context.Context, s *parser.Scanner) (ast.Branch, error) {
	rscopes := []rel.Scope{{}}
	v, err := arraiParsers.ParseWithExternals(parser.Rule("expr"), s, parser.ExternalRefs{
		"bind": func(pscope parser.Scope, _ *parser.Scanner) (parser.TreeElement, error) {
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
			expr, err := pc.CompileExpr(ctx, exprNode)
			if err != nil {
				return nil, err
			}
			exprClosure := rel.NewExprClosure(rscopes[len(rscopes)-1], expr)

			identStr := "."
			if _, ident, has := pscope.GetVal("IDENT"); has {
				identStr = ident.(parser.Scanner).String()
			}
			rscopes = append(rscopes, rscopes[len(rscopes)-1].With(identStr, exprClosure))

			var patNode ast.Branch
			var node parser.TreeElement
			if _, node, has = pscope.GetVal("pat"); has && len(node.(parser.Node).Children) > 0 {
				patNode = ast.FromParserNode(arraiParsers.Grammar(), node.(parser.Node).Children[0].(parser.Node).Children[1])
			} else if _, node, has = pscope.GetVal("pattern"); has {
				patNode = ast.FromParserNode(arraiParsers.Grammar(), node)
			}
			if has {
				source := expr.Source()
				pat, err := pc.compilePattern(ctx, patNode)
				if err != nil {
					return nil, err
				}
				bindings := pat.Bindings()
				for _, b := range bindings {
					rhs := rel.NewFunction(*parser.NewScanner(fmt.Sprintf("let %s = %s; %s", pat, source, b)),
						pat, rel.NewIdentExpr(*parser.NewScanner(b), b))
					rscopes = append(rscopes, rscopes[len(rscopes)-1].With(b, binops["->"](source, expr, rhs)))
				}
			}

			return nil, nil
		},
		"ast": func(scope parser.Scope, input *parser.Scanner) (parser.TreeElement, error) {
			_, elt, ok := scope.GetVal("macro")
			if !ok {
				panic("wat?")
			}
			_, ruleElt, ok := scope.GetVal("rule")
			if !ok {
				panic("wat?")
			}
			relScope := rscopes[len(rscopes)-1]
			macro, err := pc.unpackMacro(ctx, elt.(parser.Node), ruleElt.(parser.Node), relScope)
			if err != nil {
				return nil, err
			}

			subg := wbnf.NewFromAst(rel.ASTNodeFromValue(macro.grammar))
			rule := parser.Rule(macro.ruleName)
			parsers := subg.Compile(subg)

			childast, err := parsers.ParseWithExternals(rule, input, parser.ExternalRefs{
				"*:{()}:": func(scope parser.Scope, _ *parser.Scanner) (parser.TreeElement, error) {
					childast, err := pc.Parse(ctx, input)
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
			bodyValue, err := rel.SetCall(ctx, macro.transform, childastValue)
			if err != nil {
				return nil, err
			}

			return parser.Node{Tag: "extref", Children: nil, Extra: ast.Branch{
				"@rule": ast.One{Node: ast.Extra{Data: parser.Rule(macro.ruleName)}},
				"value": ast.One{Node: NewMacroValue(bodyValue, childastNode.Scanner())},
			}}, nil
		},
	})
	if err != nil {
		return nil, err
	}
	result := ast.FromParserNode(arraiParsers.Grammar(), v)
	if s.String() != "" {
		return result, parser.UnconsumedInput(*s, v)
	}
	return result, nil
}

func parseNest(lhs rel.Expr, branch ast.Branch) rel.Expr {
	attr := branch.One("IDENT").One("").Scanner()
	namesBranch, exist := branch["names"]
	if !exist {
		return rel.NewSingleNestExpr(attr, lhs, attr.String())
	}
	names := namesBranch.(ast.One).Node.(ast.Branch)["IDENT"].(ast.Many)
	namestrings := make([]string, len(names))
	for i, name := range names {
		namestrings[i] = name.One("").Scanner().String()
	}
	_, isInverse := branch["inv"]
	return rel.NewNestExpr(attr, isInverse, lhs, rel.NewNames(namestrings...), attr.String())
}
