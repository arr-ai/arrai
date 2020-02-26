package syntax

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

type noParseType struct{}

type parseFunc func(v interface{}) (rel.Expr, error)

func (*noParseType) Error() string {
	return "No parse"
}

var noParse = &noParseType{}

func unfakeBackquote(s string) string {
	return strings.ReplaceAll(s, "‵", "`")
}

var arraiParsers = wbnf.MustCompile(unfakeBackquote(`
expr   -> amp="&"* @ arrow=(nest | unnest | ARROW @ | binding="->" "\\" IDENT %%bind @ | binding="->" %%bind @)*
        > @:binop=("with" | "without")
        > @:binop="||"
        > @:binop="&&"
        > @:binop=/{!?(?:<>?=?|>=?|=)}
        > @ if=("if" t=expr "else" f=expr)*
        > @:binop=/{[+|]|-%?|\(\+\)}
        > @:binop=/{&|[-<][-&][->]}
        > @:binop=/{//|[*/%]}
        > @:rbinop="^"
        > unop=/{:>|=>|>>|[-+!*^]}* @
        > @ count="count"? touch?
        > @ call=("(" arg=expr:",", ")")*
        > get+ | @ get*
        > "{" rel=(names tuple=("(" v=@:",", ")"):",",?) "}"
        | "{" set=(elt=@:",",?) "}"
        | "[" array=(item=@:",",?) "]"
        | "{:" embed=(grammar=@ ":" subgrammar=%%ast) ":}"
        | op="\\\\" @
		| fn="\\" IDENT @
        | "//" pkg=( "." ("/" local=name)+
                   | "." std=IDENT?
                   | http="http://"? fqdn=name:"." ("/" path=name)*
                   )
        | "(" tuple=(pairs=(name ":" v=@ | ":" vk=(@ "." k=IDENT)):",",?) ")"
		| "(" @ ")"
		| xstr | IDENT | STR | NUM;
nest   -> "nest" names IDENT;
unnest -> "unnest" IDENT;
touch  -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
get    -> dot="." ("&"? IDENT | STR | "*");
names  -> "|" IDENT:"," "|";
name   -> IDENT | STR;
xstr   -> quote=/{\$"\s*} part=( sexpr | fragment=/{(?: \\. | :[^{"] | [^\\":] )+} )* '"'
        | quote=/{\$'\s*} part=( sexpr | fragment=/{(?: \\. | :[^{'] | [^\\':] )+} )* "'"
		| quote=/{\$‵\s*} part=( sexpr | fragment=/{(?: ‵‵  | :[^{‵] | [^‵  :] )+} )* "‵";
sexpr  -> ":{" format=/{(?:[-+#*.\_0-9a-z]+:)?} expr delim=/{(?:(?: \\. | [^\\}] )+)?} "}:";

ARROW  -> /{:>|=>|>>|order|where|sum|max|mean|median|min};
IDENT  -> /{ \. | [$@A-Za-z_][0-9$@A-Za-z_]* };
STR    -> /{ " (?: \\. | [^\\"] )* "
           | ' (?: \\. | [^\\'] )* '
           | ‵ (?: ‵‵  | [^‵  ] )* ‵
           };
NUM    -> /{ (?: \d+(?:\.\d*)? | \.\d+ ) (?: [Ee][-+]?\d+ )? };

.wrapRE -> /{\s*()\s*};
`))

type ParseContext struct {
	SourceDir string
}

func (pc ParseContext) CompileExpr(b wbnf.Branch) rel.Expr {
	// fmt.Println(b)
	name, c := which(b,
		"amp", "arrow", "unop", "binop", "rbinop",
		"if", "call", "count", "touch", "get",
		"rel", "set", "array", "embed", "op", "fn", "pkg", "tuple",
		"xstr", "IDENT", "STR", "NUM",
		"expr",
	)
	if c == nil {
		panic(fmt.Errorf("misshapen node AST: %v", b))
	}
	// log.Println(name, "\n", b)
	switch name {
	case "amp", "arrow":
		expr := pc.CompileExpr(b["expr"].(wbnf.One).Node.(wbnf.Branch))
		if arrows, has := b["arrow"]; has {
			for _, arrow := range arrows.(wbnf.Many) {
				branch := arrow.(wbnf.Branch)
				part, d := which(branch, "nest", "unnest", "ARROW", "binding")
				switch part {
				case "nest":
					expr = parseNest(expr, branch["nest"].(wbnf.One).Node.(wbnf.Branch))
				case "unnest":
					panic("unfinished")
				case "ARROW":
					f := binops[d.(wbnf.One).Node.One("").(wbnf.Leaf).Scanner().String()]
					expr = f(expr, pc.CompileExpr(arrow.(wbnf.Branch)["expr"].(wbnf.One).Node.(wbnf.Branch)))
				case "binding":
					f := binops["->"]
					rhs := pc.CompileExpr(arrow.(wbnf.Branch)["expr"].(wbnf.One).Node.(wbnf.Branch))
					if ident := arrow.One("IDENT"); ident != nil {
						rhs = rel.NewFunction(ident.Scanner().String(), rhs)
					}
					expr = f(expr, rhs)
				}
			}
		}
		if name == "amp" {
			for range c.(wbnf.Many) {
				expr = rel.NewFunction("-", expr)
			}
		}
		return expr
	case "unop":
		ops := c.(wbnf.Many)
		result := pc.CompileExpr(b.One("expr").(wbnf.Branch))
		for i := len(ops) - 1; i >= 0; i-- {
			op := ops[i].One("").(wbnf.Leaf).Scanner().String()
			f := unops[op]
			result = f(result)
		}
		return result
	case "binop":
		ops := c.(wbnf.Many)
		args := b["expr"].(wbnf.Many)
		result := pc.CompileExpr(args[0].(wbnf.Branch))
		for i, arg := range args[1:] {
			op := ops[i].One("").(wbnf.Leaf).Scanner().String()
			f := binops[op]
			result = f(result, pc.CompileExpr(arg.(wbnf.Branch)))
		}
		return result
	case "rbinop":
		ops := c.(wbnf.Many)
		args := b["expr"].(wbnf.Many)
		result := pc.CompileExpr(args[len(args)-1].(wbnf.Branch))
		for i := len(args) - 2; i >= 0; i-- {
			op := ops[i].One("").(wbnf.Leaf).Scanner().String()
			f, has := binops[op]
			if !has {
				panic("rbinop %q not found")
			}
			result = f(pc.CompileExpr(args[i].(wbnf.Branch)), result)
		}
		return result
	case "if":
		result := pc.CompileExpr(b.One("expr").(wbnf.Branch))
		for _, ifelse := range c.(wbnf.Many) {
			t := pc.CompileExpr(ifelse.One("t").(wbnf.Branch))
			f := pc.CompileExpr(ifelse.One("f").(wbnf.Branch))
			result = rel.NewIfElseExpr(result, t, f)
		}
		return result
	case "call":
		result := pc.CompileExpr(b.One("expr").(wbnf.Branch))
		for _, call := range c.(wbnf.Many) {
			for _, arg := range pc.parseExprs(call.Many("arg")...) {
				result = rel.NewCallExpr(result, arg)
			}
		}
		return result
	case "count", "touch":
		if _, has := b["touch"]; has {
			panic("unfinished")
		}
		return rel.NewCountExpr(pc.CompileExpr(b.One("expr").(wbnf.Branch)))

		// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
		// result := p.parseExpr(b.One("expr").(wbnf.Branch))
	case "get":
		var result rel.Expr
		if expr := b.One("expr"); expr != nil {
			result = pc.CompileExpr(expr.(wbnf.Branch))
		} else {
			result = rel.DotIdent
		}
		if result == nil {
			result = rel.DotIdent
		}
		for _, dot := range c.(wbnf.Many) {
			ident := dot.One("IDENT").One("").(wbnf.Leaf).Scanner().String()
			result = rel.NewDotExpr(result, ident)
		}
		return result
	case "rel":
		names := parseNames(c.(wbnf.One).Node.(wbnf.Branch)["names"].(wbnf.One).Node.(wbnf.Branch))
		tuples := c.(wbnf.One).Node.(wbnf.Branch)["tuple"].(wbnf.Many)
		tupleExprs := make([][]rel.Expr, 0, len(tuples))
		for _, tuple := range tuples {
			tupleExprs = append(tupleExprs, pc.parseExprs(tuple.(wbnf.Branch)["v"].(wbnf.Many)...))
		}
		result, err := rel.NewRelationExpr(names, tupleExprs...)
		if err != nil {
			panic(err)
		}
		return result
	case "set":
		if elts := c.(wbnf.One).Node.(wbnf.Branch)["elt"]; elts != nil {
			return rel.NewSetExpr(pc.parseExprs(elts.(wbnf.Many)...)...)
		}
		return rel.NewSetExpr()
	case "array":
		if items := c.(wbnf.One).Node.(wbnf.Branch)["item"]; items != nil {
			return rel.NewArrayExpr(pc.parseExprs(items.(wbnf.Many)...)...)
		}
		return rel.NewArray()
	case "embed":
		return rel.ASTNodeToValue(b.One("embed").One("subgrammar").One("").One("@node"))
	case "fn":
		ident := b.One("IDENT")
		expr := pc.CompileExpr(b.One("expr").(wbnf.Branch))
		return rel.NewFunction(ident.One("").Scanner().String(), expr)
	case "pkg":
		pkg := c.(wbnf.One).Node.(wbnf.Branch)
		if std, has := pkg["std"]; has {
			ident := std.(wbnf.One).Node.One("IDENT").One("")
			pkgName := ident.(wbnf.Leaf).Scanner().String()
			return NewPackageExpr(rel.NewDotExpr(rel.DotIdent, pkgName))
		} else if local := pkg["local"]; local != nil {
			var sb strings.Builder
			for i, part := range local.(wbnf.Many) {
				if i > 0 {
					sb.WriteRune('/')
				}
				sb.WriteString(strings.Trim(parseName(part.One("name").(wbnf.Branch)), "'"))
			}
			filepath := sb.String()
			if pc.SourceDir == "" {
				panic(fmt.Errorf("local import %q invalid; no local context", filepath))
			}
			return rel.NewCallExpr(
				NewPackageExpr(rel.NewIdentExpr("//./")),
				rel.NewString([]rune(path.Join(pc.SourceDir, filepath))))
		} else if fqdn := pkg["fqdn"]; fqdn != nil {
			var sb strings.Builder
			if http := pkg["http"]; http != nil {
				sb.WriteString(http.(wbnf.One).Node.(wbnf.Leaf).Scanner().String())
			}
			for i, part := range fqdn.(wbnf.Many) {
				if i > 0 {
					sb.WriteRune('.')
				}
				sb.WriteString(strings.Trim(parseName(part.One("name").(wbnf.Branch)), "'"))
			}
			if path := pkg["path"]; path != nil {
				for _, part := range path.(wbnf.Many) {
					sb.WriteRune('/')
					sb.WriteString(strings.Trim(parseName(part.One("name").(wbnf.Branch)), "'"))
				}
			}
			return rel.NewCallExpr(NewPackageExpr(rel.NewIdentExpr("//")), rel.NewString([]rune(sb.String())))
		} else {
			return NewPackageExpr(rel.DotIdent)
		}
	case "tuple":
		if entries := c.(wbnf.One).Node.Many("pairs"); entries != nil {
			attrs := make([]rel.AttrExpr, 0, len(entries))
			for _, entry := range entries {
				k := parseName(entry.One("name").(wbnf.Branch))
				v := pc.CompileExpr(entry.One("v").(wbnf.Branch))
				attr, err := rel.NewAttrExpr(k, v)
				if err != nil {
					panic(err)
				}
				attrs = append(attrs, attr)
			}
			return rel.NewTupleExpr(attrs...)
		}
		return rel.EmptyTuple
	case "IDENT":
		s := c.(wbnf.One).Node.One("").Scanner().String()
		switch s {
		case "true":
			return rel.True
		case "false":
			return rel.False
		}
		return rel.NewIdentExpr(s)
	case "STR":
		s := c.(wbnf.One).Node.One("").Scanner().String()
		return rel.NewString([]rune(parseArraiString(s)))
	case "xstr":
		quote := c.(wbnf.One).Node.One("quote").Scanner().String()[1]
		parts := c.(wbnf.One).Node.Many("part")
		exprs := make([]rel.Expr, 0, len(parts))
		for _, part := range parts {
			p, part := which(part.(wbnf.Branch), "sexpr", "fragment")
			switch p {
			case "sexpr":
				sexpr := part.(wbnf.One).Node.(wbnf.Branch)
				format := sexpr.One("format").One("").(wbnf.Leaf).Scanner().String()
				if format != "" {
					format = format[:len(format)-1]
				}
				expr := sexpr.One("expr").(wbnf.Branch)
				delim := sexpr.One("delim").One("").(wbnf.Leaf).Scanner().String()
				if delim != "" {
					delim = delim[1:]
				}
				exprs = append(exprs,
					rel.NewCallExprCurry(libStrExpand,
						rel.NewString([]rune(format)),
						pc.CompileExpr(expr),
						rel.NewString([]rune(delim)),
					))
			case "fragment":
				s := part.(wbnf.One).Node.One("").Scanner().String()
				exprs = append(exprs, rel.NewString([]rune(parseArraiStringFragment(s, quote, ":"))))
			}
		}
		return rel.NewCallExpr(libStrConcat, rel.NewArrayExpr(exprs...))
	case "NUM":
		s := c.(wbnf.One).Node.One("").Scanner().String()
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("Wat?")
		}
		return rel.NewNumber(n)
	case "expr":
		switch c := c.(type) {
		case wbnf.One:
			return pc.CompileExpr(c.Node.(wbnf.Branch))
		case wbnf.Many:
			if len(c) == 1 {
				return pc.CompileExpr(c[0].(wbnf.Branch))
			}
			panic("too many expr children")
		}
	}
	panic(fmt.Errorf("unhandled node: %v", b))
}

func (pc ParseContext) parseExprs(exprs ...wbnf.Node) []rel.Expr {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		result = append(result, pc.CompileExpr(expr.(wbnf.Branch)))
	}
	return result
}

func parseNames(names wbnf.Branch) []string {
	idents := names["IDENT"].(wbnf.Many)
	result := make([]string, 0, len(idents))
	for _, ident := range idents {
		result = append(result, ident.One("").(wbnf.Leaf).Scanner().String())
	}
	return result
}

func parseName(name wbnf.Branch) string {
	ktype, children := which(name, "IDENT", "STR")
	switch ktype {
	case "IDENT":
		return children.(wbnf.One).Node.One("").(wbnf.Leaf).Scanner().String()
	case "STR":
		s := children.(wbnf.One).Node.One("").(wbnf.Leaf).Scanner().String()
		return parseArraiString(s)
	default:
		panic("wat?")
	}
}

// MustParseString parses input string and returns the parsed Expr or panics.
func (pc ParseContext) MustParseString(s string) wbnf.Branch {
	return pc.MustParse(parser.NewScanner(s))
}

// MustParse parses input and returns the parsed Expr or panics.
func (pc ParseContext) MustParse(s *parser.Scanner) wbnf.Branch {
	ast, err := pc.Parse(s)
	if err != nil {
		panic(err)
	}
	return ast
}

// ParseString parses input string and returns the parsed Expr or an error.
func (pc ParseContext) ParseString(s string) (wbnf.Branch, error) {
	return pc.Parse(parser.NewScanner(s))
}

// Parse parses input and returns the parsed Expr or an error.
func (pc ParseContext) Parse(s *parser.Scanner) (wbnf.Branch, error) {
	rscopes := []rel.Scope{{}}
	v, err := arraiParsers.ParsePartial(parser.Rule("expr"), s, parser.Externals{
		"bind": func(pscope parser.Scope, _ *parser.Scanner, end bool) (parser.TreeElement, parser.Node, error) {
			if end {
				rscopes = rscopes[:len(rscopes)-1]
			}

			identStr := "."
			if _, ident, has := pscope.GetVal("IDENT"); has {
				identStr = ident.(parser.Scanner).String()
			}

			_, exprElt, has := pscope.GetVal("expr@1")
			if !has {
				log.Println(pscope.Keys())
				log.Println(pscope)
				panic("wat?")
			}

			exprNode := wbnf.FromParserNode(arraiParsers.Grammar(), exprElt)
			expr := pc.CompileExpr(exprNode)
			expr = rel.NewExprClosure(rscopes[len(rscopes)-1], expr)
			rscopes = append(rscopes, rscopes[len(rscopes)-1].With(identStr, expr))
			return nil, parser.Node{}, nil
		},
		"ast": func(scope parser.Scope, input *parser.Scanner, _ bool) (parser.TreeElement, parser.Node, error) {
			_, elt, ok := scope.GetVal("grammar")
			if !ok {
				panic("wat?")
			}
			astNode := wbnf.FromParserNode(arraiParsers.Grammar(), elt)
			dotExpr := pc.CompileExpr(astNode).(*rel.DotExpr)
			astExpr := dotExpr.Subject()
			astValue, err := astExpr.Eval(rscopes[len(rscopes)-1])
			if err != nil {
				return nil, parser.Node{}, err
			}
			astValueNode := rel.ASTNodeFromValue(astValue).(wbnf.Branch)
			subgrammar := wbnf.ToParserNode(wbnf.Core().Grammar(), astValueNode).(parser.Node)
			rule := parser.Rule(dotExpr.Attr())
			parsers := wbnf.NewFromNode(subgrammar).Compile(&subgrammar)
			ast, err := parsers.ParsePartial(rule, input, parser.Externals{
				"*:{()}:": func(
					pscope parser.Scope, input *parser.Scanner, _ bool,
				) (parser.TreeElement, parser.Node, error) {
					ast, err := pc.Parse(input)
					switch err.(type) {
					case nil, parser.UnconsumedInputError:
					default:
						return nil, parser.Node{}, err
					}
					log.Printf("ast: %v", ast)
					node := wbnf.ToParserNode(arraiParsers.Grammar(), ast)
					return node, parser.Node{}, nil
				},
			})
			if err != nil {
				return nil, parser.Node{}, err
			}
			return ast, subgrammar, nil
		},
	})
	// log.Printf("Parse: v = %v", v)
	if err != nil {
		return nil, err
	}
	result := wbnf.FromParserNode(arraiParsers.Grammar(), v)
	// log.Printf("Parse: result = %v", result)
	if s.String() != "" {
		return result, parser.UnconsumedInput(*s)
	}
	return result, nil
}

func which(b wbnf.Branch, names ...string) (string, wbnf.Children) {
	if len(names) == 0 {
		panic("wat?")
	}
	for _, name := range names {
		if children, has := b[name]; has {
			return name, children
		}
	}
	return "", nil
}

var unops = map[string]unOpFunc{
	"+":  rel.NewPosExpr,
	"-":  rel.NewNegExpr,
	"^":  rel.NewPowerSetExpr,
	"!":  rel.NewNotExpr,
	"*":  rel.NewEvalExpr,
	"//": NewPackageExpr,
}

var binops = map[string]binOpFunc{
	"->":      rel.NewApplyExpr,
	"=>":      rel.NewMapExpr,
	">>":      rel.NewSequenceMapExpr,
	":>":      rel.NewTupleMapExpr,
	"order":   rel.NewOrderExpr,
	"where":   rel.NewWhereExpr,
	"sum":     rel.NewSumExpr,
	"max":     rel.NewMaxExpr,
	"mean":    rel.NewMeanExpr,
	"median":  rel.NewMedianExpr,
	"min":     rel.NewMinExpr,
	"with":    rel.NewWithExpr,
	"without": rel.NewWithoutExpr,
	"&&":      rel.NewAndExpr,
	"||":      rel.NewOrExpr,
	"=":       rel.MakeEqExpr("=", func(a, b rel.Value) bool { return a.Equal(b) }),
	"<":       rel.MakeEqExpr("<", func(a, b rel.Value) bool { return a.Less(b) }),
	">":       rel.MakeEqExpr(">", func(a, b rel.Value) bool { return b.Less(a) }),
	"!=":      rel.MakeEqExpr("!=", func(a, b rel.Value) bool { return !a.Equal(b) }),
	"<=":      rel.MakeEqExpr("<=", func(a, b rel.Value) bool { return !b.Less(a) }),
	">=":      rel.MakeEqExpr(">=", func(a, b rel.Value) bool { return !a.Less(b) }),
	"+":       rel.NewAddExpr,
	"-":       rel.NewSubExpr,
	"|":       unimplementedBinOpFunc, // rel.NewUnionExpr,
	"(+)":     unimplementedBinOpFunc, // rel.NewXorExpr,
	"<&>":     rel.NewJoinExpr,
	"*":       rel.NewMulExpr,
	"/":       rel.NewDivExpr,
	"%":       rel.NewModExpr,
	"-%":      rel.NewSubModExpr,
	"//":      rel.NewIdivExpr,
	"^":       rel.NewPowExpr,
}

func parseNest(lhs rel.Expr, branch wbnf.Branch) rel.Expr {
	attr := branch.One("IDENT").One("").Scanner().String()
	names := branch["names"].(wbnf.One).Node.(wbnf.Branch)["IDENT"].(wbnf.Many)
	namestrings := make([]string, len(names))
	for i, name := range names {
		namestrings[i] = name.One("").Scanner().String()
	}
	return rel.NewNestExpr(lhs, rel.NewNames(namestrings...), attr)
}

type binOpFunc func(a, b rel.Expr) rel.Expr
type unOpFunc func(e rel.Expr) rel.Expr

func unimplementedBinOpFunc(_, _ rel.Expr) rel.Expr {
	panic("unimplemented")
}
