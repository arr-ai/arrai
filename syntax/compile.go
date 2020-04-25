package syntax

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/wbnf/ast"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
)

// type noParseType struct{}

// type parseFunc func(v interface{}) (rel.Expr, error)

// func (*noParseType) Error() string {
// 	return "No parse"
// }

// var noParse = &noParseType{}

const NoPath = "\000"

var loggingOnce sync.Once

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
		"amp", "arrow", "let", "unop", "binop", "compare", "rbinop", "if", "get",
		"tail", "count", "touch", "get", "rel", "set", "dict", "array",
		"embed", "op", "fn", "pkg", "tuple", "xstr", "IDENT", "STR", "NUM",
		"expr", "cond",
	)
	if c == nil {
		panic(fmt.Errorf("misshapen node AST: %v", b))
	}
	// log.Println(name, "\n", b)
	switch name {
	case "amp", "arrow":
		return pc.compileArrow(b, name, c)
	case "let":
		return pc.compileLet(c)
	case "unop":
		return pc.compileUnop(b, c)
	case "binop":
		return pc.compileBinop(b, c)
	case "compare":
		return pc.compileCompare(b, c)
	case "rbinop":
		return pc.compileRbinop(b, c)
	case "if":
		return pc.compileIf(b, c)
	case "cond":
		return pc.compileCond(b, c)
	case "count", "touch":
		return pc.compileCountTouch(b)
	case "get", "tail":
		return pc.compileCallGet(b)
	case "rel":
		return pc.compileRelation(c)
	case "set":
		return pc.compileSet(c)
	case "dict":
		return pc.compileDict(c)
	case "array":
		return pc.compileArray(c)
	case "embed":
		return rel.ASTNodeToValue(b.One("embed").One("subgrammar").One("ast"))
	case "fn":
		return pc.compileFunction(b)
	case "pkg":
		return pc.compilePackage(c)
	case "tuple":
		return pc.compileTuple(c)
	case "IDENT":
		return pc.compileIdent(c)
	case "STR":
		return pc.compileString(c)
	case "xstr":
		return pc.compileExpandableString(c)
	case "NUM":
		return pc.compileNumber(c)
	case "expr":
		if result := pc.compileExpr(c); result != nil {
			return result
		}
	}
	panic(fmt.Errorf("unhandled node: %v", b))
}

func (pc ParseContext) compileArrow(b ast.Branch, name string, c ast.Children) rel.Expr {
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
}

func (pc ParseContext) compileLet(c ast.Children) rel.Expr {
	exprs := c.(ast.One).Node.Many("expr")
	expr := pc.CompileExpr(exprs[0].(ast.Branch))
	rhs := pc.CompileExpr(exprs[1].(ast.Branch))
	if ident := c.(ast.One).Node.One("IDENT"); ident != nil {
		rhs = rel.NewFunction(ident.Scanner().String(), rhs)
	}
	expr = binops["->"](expr, rhs)
	return expr
}

func (pc ParseContext) compileUnop(b ast.Branch, c ast.Children) rel.Expr {
	ops := c.(ast.Many)
	result := pc.CompileExpr(b.One("expr").(ast.Branch))
	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i].One("").(ast.Leaf).Scanner().String()
		f := unops[op]
		result = f(result)
	}
	return result
}

func (pc ParseContext) compileBinop(b ast.Branch, c ast.Children) rel.Expr {
	ops := c.(ast.Many)
	args := b.Many("expr")
	result := pc.CompileExpr(args[0].(ast.Branch))
	for i, arg := range args[1:] {
		op := ops[i].One("").(ast.Leaf).Scanner().String()
		f := binops[op]
		result = f(result, pc.CompileExpr(arg.(ast.Branch)))
	}
	return result
}

func (pc ParseContext) compileCompare(b ast.Branch, c ast.Children) rel.Expr {
	args := b.Many("expr")
	argExprs := make([]rel.Expr, 0, len(args))
	comps := make([]rel.CompareFunc, 0, len(args))

	ops := c.(ast.Many)
	opStrs := make([]string, 0, len(ops))

	argExprs = append(argExprs, pc.CompileExpr(args[0].(ast.Branch)))
	for i, arg := range args[1:] {
		op := ops[i].One("").(ast.Leaf).Scanner().String()

		argExprs = append(argExprs, pc.CompileExpr(arg.(ast.Branch)))
		comps = append(comps, compareOps[op])

		opStrs = append(opStrs, op)
	}
	return rel.NewCompareExpr(argExprs, comps, opStrs)
}

func (pc ParseContext) compileRbinop(b ast.Branch, c ast.Children) rel.Expr {
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
}

func (pc ParseContext) compileIf(b ast.Branch, c ast.Children) rel.Expr {
	loggingOnce.Do(func() {
		log.Error(context.Background(),
			errors.New("operator if is deprecated and will be removed soon, please use operator cond instead. "+
				"Operator cond sample: let a = cond ( 2 > 1 : 1, 2 > 3 :2, * : 3)"))
	})

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
}

func (pc ParseContext) compileCond(b ast.Branch, c ast.Children) rel.Expr {
	// arrai eval 'cond (1 > 0:1, 2 > 3:2, *:10)'
	result := pc.compileDict(c)
	// TODO: pass arrai src expression to NewCondExpr, and include it in error messages which can help end user more.
	var identExpr, fExpr rel.Expr

	if cNode := c.(ast.One).Node; cNode != nil {
		if children, has := cNode.(ast.Branch)["IDENT"]; has {
			identExpr = pc.compileIdent(children)
		}
	}

	if fNode := c.(ast.One).Node.One("f"); fNode != nil {
		fExpr = pc.CompileExpr(fNode.(ast.Branch))
	}

	if identExpr != nil {
		result = rel.NewCondControlVarExpr(identExpr, result, fExpr)
	} else {
		result = rel.NewCondExpr(result, fExpr)
	}

	return result
}

func (pc ParseContext) compileCountTouch(b ast.Branch) rel.Expr {
	if _, has := b["touch"]; has {
		panic("unfinished")
	}
	return rel.NewCountExpr(pc.CompileExpr(b.One("expr").(ast.Branch)))

	// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
	// result := p.parseExpr(b.One("expr").(ast.Branch))
}

func (pc ParseContext) compileCallGet(b ast.Branch) rel.Expr {
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
}

func (pc ParseContext) compileRelation(c ast.Children) rel.Expr {
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
}

func (pc ParseContext) compileSet(c ast.Children) rel.Expr {
	if elts := c.(ast.One).Node.(ast.Branch)["elt"]; elts != nil {
		return rel.NewSetExpr(pc.parseExprs(elts.(ast.Many)...)...)
	}
	return rel.NewSetExpr()
}

func (pc ParseContext) compileDict(c ast.Children) rel.Expr {
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
}

func (pc ParseContext) compileArray(c ast.Children) rel.Expr {
	if items := c.(ast.One).Node.(ast.Branch)["item"]; items != nil {
		return rel.NewArrayExpr(pc.parseExprs(items.(ast.Many)...)...)
	}
	return rel.NewArray()
}

func (pc ParseContext) compileFunction(b ast.Branch) rel.Expr {
	ident := b.One("IDENT")
	expr := pc.CompileExpr(b.One("expr").(ast.Branch))
	return rel.NewFunction(ident.One("").Scanner().String(), expr)
}

func (pc ParseContext) compilePackage(c ast.Children) rel.Expr {
	pkg := c.(ast.One).Node.(ast.Branch)
	if std, has := pkg["std"]; has {
		ident := std.(ast.One).Node.One("IDENT").One("")
		pkgName := ident.(ast.Leaf).Scanner().String()
		return NewPackageExpr(rel.NewDotExpr(rel.DotIdent, pkgName))
	}
	if names := pkg.Many("name"); len(names) > 0 {
		var sb strings.Builder

		for i, name := range names {
			if i > 0 {
				sb.WriteRune('/')
			}
			sb.WriteString(parseName(name.(ast.Branch)))
		}
		filepath := sb.String()
		if pc.SourceDir == "" {
			panic(fmt.Errorf("local import %q invalid; no local context", filepath))
		}
		return rel.NewCallExpr(
			NewPackageExpr(importLocalFile(pkg["dot"] == nil)),
			rel.NewString([]rune(path.Join(pc.SourceDir, filepath))),
		)
	}
	if fqdn := pkg["fqdn"]; fqdn != nil {
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
	}
	return NewPackageExpr(rel.DotIdent)
}

func (pc ParseContext) compileTuple(c ast.Children) rel.Expr {
	if pairs := c.(ast.One).Node.Many("pairs"); pairs != nil {
		attrs := make([]rel.AttrExpr, 0, len(pairs))
		for _, pair := range pairs {
			var k string
			v := pc.CompileExpr(pair.One("v").(ast.Branch))
			if name := pair.One("name"); name != nil {
				k = parseName(name.(ast.Branch))
			} else {
				switch v := v.(type) {
				case *rel.DotExpr:
					k = v.Attr()
				case rel.IdentExpr:
					k = v.Ident()
				default:
					panic(fmt.Errorf("unnamed attr expression must be name or end in .name: %T(%[1]v)", v))
				}
			}
			attr, err := rel.NewAttrExpr(k, v)
			if err != nil {
				panic(err)
			}
			attrs = append(attrs, attr)
		}
		return rel.NewTupleExpr(attrs...)
	}
	return rel.EmptyTuple
}

func (pc ParseContext) compileIdent(c ast.Children) rel.Expr {
	s := c.(ast.One).Node.One("").Scanner().String()
	switch s {
	case "true":
		return rel.True
	case "false":
		return rel.False
	}
	return rel.NewIdentExpr(s)
}

func (pc ParseContext) compileString(c ast.Children) rel.Expr {
	s := c.(ast.One).Node.One("").Scanner().String()
	return rel.NewString([]rune(parseArraiString(s)))
}

func (pc ParseContext) compileNumber(c ast.Children) rel.Expr {
	s := c.(ast.One).Node.One("").Scanner().String()
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic("Wat?")
	}
	return rel.NewNumber(n)
}

func (pc ParseContext) compileExpr(c ast.Children) rel.Expr {
	switch c := c.(type) {
	case ast.One:
		return pc.CompileExpr(c.Node.(ast.Branch))
	case ast.Many:
		if len(c) == 1 {
			return pc.CompileExpr(c[0].(ast.Branch))
		}
		panic("too many expr children")
	}
	return nil
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

type unOpFunc func(e rel.Expr) rel.Expr

var unops = map[string]unOpFunc{
	"+":  rel.NewPosExpr,
	"-":  rel.NewNegExpr,
	"^":  rel.NewPowerSetExpr,
	"!":  rel.NewNotExpr,
	"*":  rel.NewEvalExpr,
	"//": NewPackageExpr,
}

type binOpFunc func(a, b rel.Expr) rel.Expr

var binops = map[string]binOpFunc{
	"->":      rel.NewApplyExpr,
	"=>":      rel.NewMapExpr,
	">>":      rel.NewSequenceMapExpr,
	">>>":     rel.NewIndexedSequenceMapExpr,
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
	"\\":      rel.NewOffsetExpr,
}

var compareOps = map[string]rel.CompareFunc{
	"<:": func(a, b rel.Value) bool { return b.(rel.Set).Has(a) },
	"=":  func(a, b rel.Value) bool { return a.Equal(b) },
	"<":  func(a, b rel.Value) bool { return a.Less(b) },
	">":  func(a, b rel.Value) bool { return b.Less(a) },
	"!=": func(a, b rel.Value) bool { return !a.Equal(b) },
	"<=": func(a, b rel.Value) bool { return !b.Less(a) },
	">=": func(a, b rel.Value) bool { return !a.Less(b) },
}
