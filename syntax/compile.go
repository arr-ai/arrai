package syntax

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/arr-ai/wbnf/ast"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

var leadingWSRE = regexp.MustCompile(`\A[\t ]*`)
var trailingWSRE = regexp.MustCompile(`[\t ]*\z`)
var expansionRE = regexp.MustCompile(`(?::([-+#*\.\_0-9a-z]*))(:(?:\\.|[^\\:}])*)?(?::((?:\\.|[^\\:}])*))?`)

// type noParseType struct{}

// type parseFunc func(v interface{}) (rel.Expr, error)

// func (*noParseType) Error() string {
// 	return "No parse"
// }

// var noParse = &noParseType{}

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
	name, c := which(b,
		"amp", "arrow", "let", "unop", "binop", "rbinop",
		"if", "get", "tail", "count", "touch", "get",
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
	case "count", "touch":
		if _, has := b["touch"]; has {
			panic("unfinished")
		}
		return rel.NewCountExpr(pc.CompileExpr(b.One("expr").(ast.Branch)))

		// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
		// result := p.parseExpr(b.One("expr").(ast.Branch))
	case "get", "tail":
		var result rel.Expr

		get := func(get ast.Node) {
			if get != nil {
				if ident := get.One("IDENT"); ident != nil {
					result = rel.NewDotExpr(result, ident.One("").(ast.Leaf).Scanner().String())
				}
				if str := get.One("STR"); str != nil {
					s := str.One("").Scanner().String()
					result = rel.NewDotExpr(result, parseArraiString(s))
				}
			}
		}

		if expr := b.One("expr"); expr != nil {
			result = pc.CompileExpr(expr.(ast.Branch))
		} else {
			result = rel.DotIdent
			get(b.One("get"))
		}

		for _, part := range b.Many("tail") {
			if call := part.One("call"); call != nil {
				args := call.Many("arg")
				exprs := make([]ast.Node, 0, len(args))
				for _, arg := range args {
					exprs = append(exprs, arg.One("expr"))
				}
				for _, arg := range pc.parseExprs(exprs...) {
					result = rel.NewCallExpr(result, arg)
				}
			}
			get(part.One("get"))
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
					entryExprs := make([]rel.DictEntryTupleExpr, 0, len(keyExprs))
					for i, keyExpr := range keyExprs {
						valueExpr := valueExprs[i]
						entryExprs = append(entryExprs, rel.NewDictEntryTupleExpr(keyExpr, valueExpr))
					}
					return rel.NewDictExpr(false, entryExprs...)
				}
			}
			panic("mismatch between dict keys and values")
		}
		return rel.NewDict(false)
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
			var pkgPathList []ast.Node
			localTree := local.(ast.Many)[0].(ast.Node)
			pkgNode := ast.First(localTree, "pkgname")
			pkgBranch := ast.First(pkgNode, "PKG_PATH")
			if pkgBranch != nil {
				pkgPath := ast.First(pkgBranch, "PATH")
				pkgPathList = ast.All(pkgPath, "")
			} else {
				pkgSTR := ast.First(pkgNode, "STR")
				pkgPathList = ast.All(pkgSTR, "")
			}

			for i, p := range pkgPathList {
				if i > 0 {
					sb.WriteRune('/')
				}
				sb.WriteString(strings.Trim(p.Scanner().String(), "'"))
			}
			filepath := sb.String()
			if pc.SourceDir == "" {
				panic(fmt.Errorf("local import %q invalid; no local context", filepath))
			}
			return rel.NewCallExpr(
				NewPackageExpr(importLocalFile(pkg["dot"] == nil)),
				rel.NewString([]rune(path.Join(pc.SourceDir, filepath))),
			)
		} else if fqdn := pkg["fqdn"]; fqdn != nil {
			var sb strings.Builder
			if http := pkg["http"]; http != nil {
				sb.WriteString(http.(ast.One).Node.Scanner().String())
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
			return rel.NewCallExpr(NewPackageExpr(importExternalContent()), rel.NewString([]rune(sb.String())))
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
					ws = sexpr.One("close").One("").(ast.Leaf).Scanner().String()[1:]
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
				exprs[i] = rel.NewCallExprCurry(stdStrExpand,
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
		return rel.NewCallExpr(stdStrConcat, rel.NewArrayExpr(exprs...))
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
	"~~":      rel.NewSymmDiffExpr,
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

type binOpFunc func(a, b rel.Expr) rel.Expr
type unOpFunc func(e rel.Expr) rel.Expr
