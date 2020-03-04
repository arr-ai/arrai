package syntax

import (
	"fmt"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

var leadingWSRE = regexp.MustCompile(`\A[\t ]*`)
var trailingWSRE = regexp.MustCompile(`[\t ]*\z`)
var expansionRE = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)

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
        > C* @:binop=/{&~|&|~|[-<][-&][->]} C*
        > C* @:binop=/{//|[*/%]} C*
        > C* @:rbinop="^" C*
        > C* unop=/{:>|=>|>>|[-+!*^]}* @ C*
        > C* @ count="count"? C* touch? C*
        > C* @ call=("(" arg=expr:",", ")")* C*
        > C* get+ C* | C* @ get* C*
        > C* "{" C* rel=(names tuple=("(" v=@:",", ")"):",",?) "}" C*
        | C* "{" C* set=(elt=@:",",?) "}" C*
        | C* "{" C* dict=((key=@ ":" value=@):",",?) "}" C*
        | C* "[" C* array=(item=@:",",?) "]" C*
        | C* "{:" C* embed=(grammar=@ ":" subgrammar=%%ast) ":}" C*
        | C* op="\\\\" @ C*
        | C* fn="\\" IDENT @ C*
        | C* "//" pkg=( dot="." ("/" local=name)+
                   | "." std=IDENT?
                   | http=/{https?://}? fqdn=name:"." ("/" path=name)*
                   )
        | C* "(" tuple=(pairs=(name ":" v=@ | ":" vk=(@ "." k=IDENT)):",",?) ")" C*
        | C* "(" @ ")" C*
        | C* let=("let" C* IDENT C* "=" C* @ %%bind C* "in"? C* @) C*
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
xstr   -> C* quote=/{\$"\s*} part=( sexpr | fragment=/{(?: \\. | :[^{"] | [^\\":] )+} )* '"' C*
        | C* quote=/{\$'\s*} part=( sexpr | fragment=/{(?: \\. | :[^{'] | [^\\':] )+} )* "'" C*
        | C* quote=/{\$‵\s*} part=( sexpr | fragment=/{(?: ‵‵  | :[^{‵] | [^‵  :] )+} )* "‵" C*;
sexpr  -> ":{"
          C* expr C*
          control=/{ (?: : [-+#*\.\_0-9a-z]* (?: : (?: \\. | [^\\:}] )* ){0,2} )? }
          close=/{\}:\s*};

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

const NoPath = "\000"

func Compile(filepath, source string) (_ rel.Expr, err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("error compiling %q: %v", filepath, e)
			}
		}
	}()
	return MustCompile(filepath, source), nil
}

func MustCompile(filepath, source string) rel.Expr {
	dirpath := "."
	if filepath != "" {
		if filepath == NoPath {
			dirpath = NoPath
		} else {
			dirpath = path.Dir(filepath)
		}
	}
	pc := ParseContext{SourceDir: dirpath}
	ast, err := pc.Parse(parser.NewScanner(source))
	if err != nil {
		panic(err)
	}
	return pc.CompileExpr(ast)
}

func (pc ParseContext) CompileExpr(b ast.Branch) rel.Expr {
	// fmt.Println(b)
	name, c := which(b,
		"amp", "arrow", "let", "unop", "binop", "rbinop",
		"if", "call", "count", "touch", "get",
		"rel", "set", "dict", "array", "embed", "op", "fn", "pkg", "tuple",
		"xstr", "IDENT", "STR", "NUM",
		"expr",
	)
	if c == nil {
		panic(fmt.Errorf("misshapen node AST: %v", b))
	}
	// log.Println(name, "\n", b)
	switch name {
	case "amp", "arrow":
		expr := pc.CompileExpr(b["expr"].(ast.One).Node.(ast.Branch))
		if arrows, has := b["arrow"]; has {
			for _, arrow := range arrows.(ast.Many) {
				branch := arrow.(ast.Branch)
				part, d := which(branch, "nest", "unnest", "ARROW", "binding")
				switch part {
				case "nest":
					expr = parseNest(expr, branch["nest"].(ast.One).Node.(ast.Branch))
				case "unnest":
					panic("unfinished")
				case "ARROW":
					f := binops[d.(ast.One).Node.One("").(ast.Leaf).Scanner().String()]
					expr = f(expr, pc.CompileExpr(arrow.(ast.Branch)["expr"].(ast.One).Node.(ast.Branch)))
				case "binding":
					rhs := pc.CompileExpr(arrow.(ast.Branch)["expr"].(ast.One).Node.(ast.Branch))
					if ident := arrow.One("IDENT"); ident != nil {
						rhs = rel.NewFunction(ident.Scanner().String(), rhs)
					}
					expr = binops["->"](expr, rhs)
				}
			}
		}
		if name == "amp" {
			for range c.(ast.Many) {
				expr = rel.NewFunction("-", expr)
			}
		}
		return expr
	case "let":
		exprs := c.(ast.One).Node.Many("expr")
		expr := pc.CompileExpr(exprs[0].(ast.Branch))
		rhs := pc.CompileExpr(exprs[1].(ast.Branch))
		if ident := c.(ast.One).Node.One("IDENT"); ident != nil {
			rhs = rel.NewFunction(ident.Scanner().String(), rhs)
		}
		expr = binops["->"](expr, rhs)
		return expr
	case "unop":
		ops := c.(ast.Many)
		result := pc.CompileExpr(b.One("expr").(ast.Branch))
		for i := len(ops) - 1; i >= 0; i-- {
			op := ops[i].One("").(ast.Leaf).Scanner().String()
			f := unops[op]
			result = f(result)
		}
		return result
	case "binop":
		ops := c.(ast.Many)
		args := b["expr"].(ast.Many)
		result := pc.CompileExpr(args[0].(ast.Branch))
		for i, arg := range args[1:] {
			op := ops[i].One("").(ast.Leaf).Scanner().String()
			f := binops[op]
			result = f(result, pc.CompileExpr(arg.(ast.Branch)))
		}
		return result
	case "rbinop":
		ops := c.(ast.Many)
		args := b["expr"].(ast.Many)
		result := pc.CompileExpr(args[len(args)-1].(ast.Branch))
		for i := len(args) - 2; i >= 0; i-- {
			op := ops[i].One("").(ast.Leaf).Scanner().String()
			f, has := binops[op]
			if !has {
				panic("rbinop %q not found")
			}
			result = f(pc.CompileExpr(args[i].(ast.Branch)), result)
		}
		return result
	case "if":
		result := pc.CompileExpr(b.One("expr").(ast.Branch))
		for _, ifelse := range c.(ast.Many) {
			t := pc.CompileExpr(ifelse.One("t").(ast.Branch))
			var f rel.Expr = rel.None
			if fNode := ifelse.One("f"); fNode != nil {
				f = pc.CompileExpr(fNode.(ast.Branch))
			}
			result = rel.NewIfElseExpr(result, t, f)
		}
		return result
	case "call":
		result := pc.CompileExpr(b.One("expr").(ast.Branch))
		for _, call := range c.(ast.Many) {
			for _, arg := range pc.parseExprs(call.Many("arg")...) {
				result = rel.NewCallExpr(result, arg)
			}
		}
		return result
	case "count", "touch":
		if _, has := b["touch"]; has {
			panic("unfinished")
		}
		return rel.NewCountExpr(pc.CompileExpr(b.One("expr").(ast.Branch)))

		// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
		// result := p.parseExpr(b.One("expr").(ast.Branch))
	case "get":
		var result rel.Expr
		if expr := b.One("expr"); expr != nil {
			result = pc.CompileExpr(expr.(ast.Branch))
		} else {
			result = rel.DotIdent
		}
		if result == nil {
			result = rel.DotIdent
		}
		for _, dot := range c.(ast.Many) {
			ident := dot.One("IDENT").One("").(ast.Leaf).Scanner().String()
			result = rel.NewDotExpr(result, ident)
		}
		return result
	case "rel":
		names := parseNames(c.(ast.One).Node.(ast.Branch)["names"].(ast.One).Node.(ast.Branch))
		tuples := c.(ast.One).Node.(ast.Branch)["tuple"].(ast.Many)
		tupleExprs := make([][]rel.Expr, 0, len(tuples))
		for _, tuple := range tuples {
			tupleExprs = append(tupleExprs, pc.parseExprs(tuple.(ast.Branch)["v"].(ast.Many)...))
		}
		result, err := rel.NewRelationExpr(names, tupleExprs...)
		if err != nil {
			panic(err)
		}
		return result
	case "set":
		if elts := c.(ast.One).Node.(ast.Branch)["elt"]; elts != nil {
			return rel.NewSetExpr(pc.parseExprs(elts.(ast.Many)...)...)
		}
		return rel.NewSetExpr()
	case "dict":
		// C* "{" C* dict=((key=@ ":" value=@):",",?) "}" C*
		keys := c.(ast.One).Node.(ast.Branch)["key"]
		values := c.(ast.One).Node.(ast.Branch)["value"]
		if (keys != nil) || (values != nil) {
			if (keys != nil) && (values != nil) {
				keyExprs := pc.parseExprs(keys.(ast.Many)...)
				valueExprs := pc.parseExprs(values.(ast.Many)...)
				if len(keyExprs) == len(valueExprs) {
					pairs := make([][2]rel.Expr, 0, len(keyExprs))
					for i, keyExpr := range keyExprs {
						valueExpr := valueExprs[i]
						pairs = append(pairs, [2]rel.Expr{keyExpr, valueExpr})
					}
					return rel.NewDictExpr(pairs...)
				}
			}
			panic("mismatch between dict keys and values")
		}
		return rel.NewDict()
	case "array":
		if items := c.(ast.One).Node.(ast.Branch)["item"]; items != nil {
			return rel.NewArrayExpr(pc.parseExprs(items.(ast.Many)...)...)
		}
		return rel.NewArray()
	case "embed":
		return rel.ASTNodeToValue(b.One("embed").One("subgrammar").One("ast"))
	case "fn":
		ident := b.One("IDENT")
		expr := pc.CompileExpr(b.One("expr").(ast.Branch))
		return rel.NewFunction(ident.One("").Scanner().String(), expr)
	case "pkg":
		pkg := c.(ast.One).Node.(ast.Branch)
		if std, has := pkg["std"]; has {
			ident := std.(ast.One).Node.One("IDENT").One("")
			pkgName := ident.(ast.Leaf).Scanner().String()
			return NewPackageExpr(rel.NewDotExpr(rel.DotIdent, pkgName))
		} else if local := pkg["local"]; local != nil {
			var sb strings.Builder
			for i, part := range local.(ast.Many) {
				if i > 0 {
					sb.WriteRune('/')
				}
				sb.WriteString(strings.Trim(parseName(part.One("name").(ast.Branch)), "'"))
			}
			filepath := sb.String()
			if pc.SourceDir == "" {
				panic(fmt.Errorf("local import %q invalid; no local context", filepath))
			}
			return rel.NewCallExpr(
				NewPackageExpr(importLocalFile()),
				rel.NewString([]rune(path.Join(pc.SourceDir, filepath))),
			)
		} else if fqdn := pkg["fqdn"]; fqdn != nil {
			var sb strings.Builder
			if http := pkg["http"]; http != nil {
				sb.WriteString(http.(ast.One).Node.(ast.Leaf).Scanner().String())
			}
			for i, part := range fqdn.(ast.Many) {
				if i > 0 {
					sb.WriteRune('.')
				}
				sb.WriteString(strings.Trim(parseName(part.One("name").(ast.Branch)), "'"))
			}
			if path := pkg["path"]; path != nil {
				for _, part := range path.(ast.Many) {
					sb.WriteRune('/')
					sb.WriteString(strings.Trim(parseName(part.One("name").(ast.Branch)), "'"))
				}
			}
			return rel.NewCallExpr(NewPackageExpr(importURL), rel.NewString([]rune(sb.String())))
		} else {
			return NewPackageExpr(rel.DotIdent)
		}
	case "tuple":
		if entries := c.(ast.One).Node.Many("pairs"); entries != nil {
			attrs := make([]rel.AttrExpr, 0, len(entries))
			for _, entry := range entries {
				k := parseName(entry.One("name").(ast.Branch))
				v := pc.CompileExpr(entry.One("v").(ast.Branch))
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
		s := c.(ast.One).Node.One("").Scanner().String()
		switch s {
		case "true":
			return rel.True
		case "false":
			return rel.False
		}
		return rel.NewIdentExpr(s)
	case "STR":
		s := c.(ast.One).Node.One("").Scanner().String()
		return rel.NewString([]rune(parseArraiString(s)))
	case "xstr":
		quote := c.(ast.One).Node.One("quote").Scanner().String()
		parts := []interface{}{}
		{
			ws := quote[2:]
			trim := ""
			trimIndent := func(s string) {
				s = ws + s
				ws = ""
				if trim == "" {
					s = strings.TrimPrefix(s, "\n")
					i := leadingWSRE.FindStringIndex(s)
					trim = "\n" + s[:i[1]]
					s = s[i[1]:]
				}
				if trim != "\n" {
					s = strings.ReplaceAll(s, trim, "\n")
				}
				if s != "" {
					parts = append(parts, s)
				}
			}
			for i, part := range c.(ast.One).Node.Many("part") {
				p, part := which(part.(ast.Branch), "sexpr", "fragment")
				switch p {
				case "sexpr":
					if i == 0 || ws != "" {
						trimIndent("")
					}
					sexpr := part.(ast.One).Node.(ast.Branch)
					ws = sexpr.One("close").One("").(ast.Leaf).Scanner().String()[2:]
					parts = append(parts, sexpr)
				case "fragment":
					s := part.(ast.One).Node.One("").Scanner().String()
					s = parseArraiStringFragment(s, quote[1:2]+":", "")
					trimIndent(s)
				}
			}
		}
		next := ""
		exprs := make([]rel.Expr, len(parts))
		for i := len(parts) - 1; i >= 0; i-- {
			part := parts[i]
			switch part := part.(type) {
			case ast.Branch:
				indent := ""
				if i > 0 {
					if s, ok := parts[i-1].(string); ok {
						indent = trailingWSRE.FindString(s)
					}
				}

				format := ""
				delim := ""
				appendIfNotEmpty := ""
				if control := part.One("control").One("").(ast.Leaf).Scanner().String(); control != "" {
					m := expansionRE.FindStringSubmatchIndex(control)
					if m[2] >= 0 {
						format = control[m[2]:m[3]]
					}
					if m[4] >= 0 {
						delim = parseArraiStringFragment(control[m[4]:m[5]], ":}", "\n"+indent)
					}
					if m[6] >= 0 {
						appendIfNotEmpty = parseArraiStringFragment(control[m[6]:m[7]], ":}", "\n"+indent)
					}
				}
				expr := part.One("expr").(ast.Branch)
				if strings.HasPrefix(next, "\n") {
					if i > 0 {
						if s, ok := parts[i-1].(string); ok {
							if strings.HasSuffix(s, "\n") {
								appendIfNotEmpty += "\n"
								parts[i+1] = next[1:]
							}
						}
					} else {
						appendIfNotEmpty += "\n"
						parts[i+1] = next[1:]
					}
					next = ""
				}
				exprs[i] = rel.NewCallExprCurry(libStrExpand,
					rel.NewString([]rune(format)),
					pc.CompileExpr(expr),
					rel.NewString([]rune(delim)),
					rel.NewString([]rune(appendIfNotEmpty)),
				)
			case string:
				next = part
			}
		}
		for i, part := range parts {
			if s, ok := part.(string); ok {
				exprs[i] = rel.NewString([]rune(s))
			}
		}
		return rel.NewCallExpr(libStrConcat, rel.NewArrayExpr(exprs...))
	case "NUM":
		s := c.(ast.One).Node.One("").Scanner().String()
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("Wat?")
		}
		return rel.NewNumber(n)
	case "expr":
		switch c := c.(type) {
		case ast.One:
			return pc.CompileExpr(c.Node.(ast.Branch))
		case ast.Many:
			if len(c) == 1 {
				return pc.CompileExpr(c[0].(ast.Branch))
			}
			panic("too many expr children")
		}
	}
	panic(fmt.Errorf("unhandled node: %v", b))
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
		panic("wat?")
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

func which(b ast.Branch, names ...string) (string, ast.Children) {
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
	"orderby": rel.NewOrderByExpr,
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
	"++":      rel.NewConcatExpr,
	"&~":      rel.NewDiffExpr,
	"~":       rel.NewSymmDiffExpr,
	"&":       rel.NewIntersectExpr,
	"|":       rel.NewUnionExpr,
	"<&>":     rel.NewJoinExpr,
	"*":       rel.NewMulExpr,
	"/":       rel.NewDivExpr,
	"%":       rel.NewModExpr,
	"-%":      rel.NewSubModExpr,
	"//":      rel.NewIdivExpr,
	"^":       rel.NewPowExpr,
	"<:":      rel.NewMemberExpr,
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

type binOpFunc func(a, b rel.Expr) rel.Expr
type unOpFunc func(e rel.Expr) rel.Expr

func unimplementedBinOpFunc(_, _ rel.Expr) rel.Expr {
	panic("unimplemented")
}
