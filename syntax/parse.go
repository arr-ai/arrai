package syntax

import (
	"fmt"
	"log"
	"strings"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

// type noParseType struct{}

// type parseFunc func(v interface{}) (rel.Expr, error)

// func (*noParseType) Error() string {
// 	return "No parse"
// }

// var noParse = &noParseType{}

func unfakeBackquote(s string) string {
	return strings.ReplaceAll(s, "‵", "`")
}

var arraiParsers = wbnf.MustCompile(unfakeBackquote(`
expr   -> C* amp="&"* @ C* arrow=(
              nest |
              unnest |
              ARROW @ |
              binding="->" C* "\\" C* IDENT C* %%bind C* @ |
              binding="->" C* %%bind @
          )* C*
        > C* @:binop=("with" | "without") C*
        > C* @:binop="||" C*
        > C* @:binop="&&" C*
        > C* @:binop=/{!?(?:<:|<>?=?|>=?|=)} C*
		> C* @ if=("if" t=expr ("else" f=expr)?)* C*
        > C* @:binop=/{\+\+|[+|]|-%?} C*
        > C* @:binop=/{&~|&|~~?|[-<][-&][->]} C*
        > C* @:binop=/{//|[*/%]|\\} C*
        > C* @:rbinop="^" C*
        > C* unop=/{:>|=>|>>|[-+!*^]}* @ C*
        > C* @ count="count"? C* touch? C*
        > C* (get | @) tail=(
              get
            | call=("("
                  arg=(
                      expr (":" end=expr? (":" step=expr)?)?
                      |     ":" end=expr  (":" step=expr)?
                  ):",",
              ")")
          )* C*
        > C* "{" C* rel=(names tuple=("(" v=@:",", ")"):",",?) "}" C*
        | C* "{" C* set=(elt=@:",",?) "}" C*
		| C* "{" C* dict=((key=@ ":" value=@):",",?) "}" C*
        | C* cond=("cond" "(" (key=@ ":" value=@):",",? ")") C*
        | C* "[" C* array=(item=@:",",?) "]" C*
        | C* "{:" C* embed=(grammar=@ ":" subgrammar=%%ast) ":}" C*
        | C* op="\\\\" @ C*
        | C* fn="\\" IDENT @ C*
        | C* "//" pkg=( dot="."? ("/" name)+
                      | "." std=IDENT?
                      | http=/{https?://}? fqdn=name:"." ("/" path=name)*
                      )
        | C* "(" tuple=(pairs=(name? ":" v=@):",",?) ")" C*
        | C* "(" @ ")" C*
        | C* let=("let" C* IDENT C* "=" C* @ %%bind C* ";" C* @) C*
        | C* xstr C*
        | C* IDENT C*
        | C* STR C*
		| C* NUM C*;
nest   -> C* "nest" names IDENT C*;
unnest -> C* "unnest" IDENT C*;
touch  -> C* ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")" C*;
get    -> C* dot="." ("&"? IDENT | STR | "*") C*;
names  -> C* "|" C* IDENT:"," C* "|" C*;
name   -> C* IDENT C* | C* STR C*;
xstr   -> C* quote=/{\$"\s*} part=( sexpr | fragment=/{(?: \\. | \$[^{"] | [^\\"$] )+} )* '"' C*
        | C* quote=/{\$'\s*} part=( sexpr | fragment=/{(?: \\. | \$[^{'] | [^\\'$] )+} )* "'" C*
        | C* quote=/{\$‵\s*} part=( sexpr | fragment=/{(?: ‵‵  | \$[^{‵] | [^‵  $] )+} )* "‵" C*;
sexpr  -> "${"
          C* expr C*
          control=/{ (?: : [-+#*\.\_0-9a-z]* (?: : (?: \\. | [^\\:}] )* ){0,2} )? }
		  close=/{\}\s*};
cond1   -> "cond" "(" kvs ("*" f=expr)? ")";
kvs    -> (key=@ ":" value=@):",",?;

ARROW  -> /{:>|=>|>>|orderby|order|where|sum|max|mean|median|min};
IDENT  -> /{ \. | [$@A-Za-z_][0-9$@A-Za-z_]* };
STR    -> /{ " (?: \\. | [^\\"] )* "
           | ' (?: \\. | [^\\'] )* '
           | ‵ (?: ‵‵  | [^‵  ] )* ‵
           };
NUM    -> /{ (?: \d+(?:\.\d*)? | \.\d+ ) (?: [Ee][-+]?\d+ )? };
C      -> /{ # .* $ };

.wrapRE -> /{\s*()\s*};
`), nil)

type ParseContext struct {
	SourceDir string
}

func (pc ParseContext) parseExprs(exprs ...ast.Node) []rel.Expr {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		result = append(result, pc.CompileExpr(expr.(ast.Branch)))
	}
	return result
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
			_, elt, ok := scope.GetVal("grammar")
			if !ok {
				panic("wat?")
			}
			astNode := ast.FromParserNode(arraiParsers.Grammar(), elt)
			dotExpr := pc.CompileExpr(astNode).(*rel.DotExpr)
			astExpr := dotExpr.Subject()
			astValue, err := astExpr.Eval(rscopes[len(rscopes)-1])
			if err != nil {
				return nil, err
			}
			astValueNode := rel.ASTNodeFromValue(astValue).(ast.Branch)
			subg := wbnf.NewFromAst(astValueNode)
			rule := parser.Rule(dotExpr.Attr())
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
			return ast.NewExtRefTreeElement(parsers.Grammar(), childast), nil
		},
	})
	// log.Printf("Parse: v = %v", v)
	if err != nil {
		return nil, err
	}
	result := ast.FromParserNode(arraiParsers.Grammar(), v)
	// log.Printf("Parse: result = %v", result)
	if s.String() != "" {
		return result, parser.UnconsumedInput(*s, v)
	}
	return result, nil
}

func parseNest(lhs rel.Expr, branch ast.Branch) rel.Expr {
	attr := branch.One("IDENT").One("").Scanner().String()
	names := branch["names"].(ast.One).Node.(ast.Branch)["IDENT"].(ast.Many)
	namestrings := make([]string, len(names))
	for i, name := range names {
		namestrings[i] = name.One("").Scanner().String()
	}
	return rel.NewNestExpr(lhs, rel.NewNames(namestrings...), attr)
}
