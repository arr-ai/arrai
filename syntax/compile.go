package syntax

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
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

const exprTag = "expr"

var loggingOnce sync.Once

// Compile compiles source string.
func Compile(filepath, source string) (_ rel.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("error compiling %q: %v", filepath, r)
			}
		}
	}()
	return MustCompile(filepath, source), nil
}

func MustCompile(filePath, source string) rel.Expr {
	dirpath := "."
	if filePath != "" {
		if filePath == NoPath {
			dirpath = NoPath
		} else {
			dirpath = path.Dir(filePath)
		}
	}
	pc := ParseContext{SourceDir: dirpath}
	if !filepath.IsAbs(filePath) {
		var err error
		filePath, err = filepath.Rel(".", filePath)
		if err != nil {
			panic(err)
		}
	}
	ast, err := pc.Parse(parser.NewScannerWithFilename(source, filePath))
	if err != nil {
		panic(err)
	}
	return pc.CompileExpr(ast)
}

func (pc ParseContext) CompileExpr(b ast.Branch) rel.Expr {
	// Note: please make sure if it is necessary to add new syntax name before `expr`.
	name, c := which(b,
		"amp", "arrow", "let", "unop", "binop", "compare", "rbinop", "if", "get",
		"tail_op", "postfix", "touch", "get", "rel", "set", "dict", "array", "bytes",
		"embed", "op", "fn", "pkg", "tuple", "xstr", "IDENT", "STR", "NUM", "CHAR",
		"cond", exprTag,
	)
	if c == nil {
		panic(fmt.Errorf("misshapen node AST: %v", b))
	}
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
		return pc.compileCond(c)
	case "postfix", "touch":
		return pc.compilePostfixAndTouch(b, c)
	case "get", "tail_op":
		return pc.compileCallGet(b)
	case "rel":
		return pc.compileRelation(c)
	case "set":
		return pc.compileSet(c)
	case "dict":
		return pc.compileDict(c)
	case "array":
		return pc.compileArray(c)
	case "bytes":
		return pc.compileBytes(c)
	case "embed":
		return pc.compileMacro(b)
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
	case "CHAR":
		return pc.compileChar(c)
	case "extrefmacro":
		return c.(rel.Value)
	case exprTag:
		if result := pc.compileExpr(c); result != nil {
			return result
		}
	}
	panic(fmt.Errorf("unhandled node: %v", b))
}

func (pc ParseContext) compilePattern(b ast.Branch) rel.Pattern {
	if ptn := b.One("pattern"); ptn != nil {
		return pc.compilePattern(ptn.(ast.Branch))
	}
	if arr := b.One("array"); arr != nil {
		return pc.compileArrayPattern(arr.(ast.Branch))
	}
	if tuple := b.One("tuple"); tuple != nil {
		return pc.compileTuplePattern(tuple.(ast.Branch))
	}
	if dict := b.One("dict"); dict != nil {
		return pc.compileDictPattern(dict.(ast.Branch))
	}
	if set := b.One("set"); set != nil {
		return pc.compileSetPattern(set.(ast.Branch))
	}
	if extra := b.One("extra"); extra != nil {
		return pc.compileExtraElementPattern(extra.(ast.Branch))
	}
	if ident := b.One("identpattern"); ident != nil {
		return rel.NewIdentPattern(ident.Scanner().String())
	}
	if expr := b.Many("exprpattern"); expr != nil {
		var elements []rel.Expr
		for _, e := range expr {
			expr := pc.CompileExpr(e.(ast.Branch))
			elements = append(elements, expr)
		}
		return rel.NewExprsPattern(elements...)
	}

	return rel.NewExprPattern(pc.CompileExpr(b))
}

func (pc ParseContext) compileExtraElementPattern(b ast.Branch) rel.Pattern {
	var ident string
	if id := b.One("ident"); id != nil {
		ident = id.Scanner().String()
	}
	return rel.NewExtraElementPattern(ident)
}

func (pc ParseContext) compilePatterns(exprs ...ast.Node) []rel.Pattern {
	result := make([]rel.Pattern, 0, len(exprs))
	for _, expr := range exprs {
		result = append(result, pc.compilePattern(expr.(ast.Branch)))
	}
	return result
}

func (pc ParseContext) compileSparsePatterns(b ast.Branch) []rel.Pattern {
	var nodes []ast.Node
	if firstItem, exists := b["first_item"]; exists {
		nodes = []ast.Node{firstItem.(ast.One).Node}
		if items, exists := b["item"]; exists {
			for _, i := range items.(ast.Many) {
				nodes = append(nodes, i)
			}
		}
	}
	result := make([]rel.Pattern, 0, len(nodes))
	for _, expr := range nodes {
		if expr.One("empty") != nil {
			result = append(result, nil)
			continue
		}
		result = append(result, pc.compilePattern(expr.(ast.Branch)))
	}
	return result
}

func (pc ParseContext) compileArrayPattern(b ast.Branch) rel.Pattern {
	return rel.NewArrayPattern(pc.compileSparsePatterns(b)...)
}

func (pc ParseContext) compileTuplePattern(b ast.Branch) rel.Pattern {
	if pairs := b.Many("pairs"); pairs != nil {
		attrs := make([]rel.TuplePatternAttr, 0, len(pairs))
		for _, pair := range pairs {
			var k string
			var v rel.Pattern

			if extra := pair.One("extra"); extra != nil {
				v = pc.compilePattern(pair.(ast.Branch))
			} else {
				v = pc.compilePattern(pair.One("v").(ast.Branch))
				if name := pair.One("name"); name != nil {
					k = parseName(name.(ast.Branch))
				} else {
					k = v.String()
				}
			}

			attr := rel.NewTuplePatternAttr(k, v)
			attrs = append(attrs, attr)
		}
		return rel.NewTuplePattern(attrs...)
	}
	return rel.NewTuplePattern()
}

func (pc ParseContext) compileDictPattern(b ast.Branch) rel.Pattern {
	keys := b["key"]
	values := b["value"]
	if (keys != nil) != (values != nil) {
		panic("mismatch between dict keys and values")
	}
	if (keys != nil) && (values != nil) {
		keyExprs := pc.compileExprs(keys.(ast.Many)...)
		valuePtns := pc.compilePatterns(values.(ast.Many)...)
		if len(keyExprs) == len(valuePtns) {
			entryPtns := make([]rel.DictPatternEntry, 0, len(keyExprs))
			for i, keyExpr := range keyExprs {
				entryPtns = append(entryPtns, rel.NewDictPatternEntry(keyExpr, valuePtns[i]))
			}
			if extra := b["ext"]; extra != nil {
				entryPtns = append(
					entryPtns,
					rel.NewDictPatternEntry(
						rel.DotIdent,
						pc.compileExtraElementPattern(extra.(ast.Many)[0].One("extra").(ast.Branch)),
					),
				)
			}
			return rel.NewDictPattern(entryPtns...)
		}
		panic("mismatch between dict keys and values")
	}
	return rel.NewDictPattern()
}

func (pc ParseContext) compileSetPattern(b ast.Branch) rel.Pattern {
	if elts := b["elt"]; elts != nil {
		return rel.NewSetPattern(pc.compilePatterns(elts.(ast.Many)...)...)
	}
	return rel.NewSetPattern()
}

func (pc ParseContext) compileArrow(b ast.Branch, name string, c ast.Children) rel.Expr {
	expr := pc.CompileExpr(b[exprTag].(ast.One).Node.(ast.Branch))
	source := c.Scanner()
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
				op := d.(ast.One).Node.One("").(ast.Leaf).Scanner()
				f := binops[op.String()]
				expr = f(op, expr, pc.CompileExpr(arrow.(ast.Branch)[exprTag].(ast.One).Node.(ast.Branch)))
			case "binding":
				rhs := pc.CompileExpr(arrow.(ast.Branch)[exprTag].(ast.One).Node.(ast.Branch))
				if pattern := arrow.One("pattern"); pattern != nil {
					p := pc.compilePattern(pattern.(ast.Branch))
					rhs = rel.NewFunction(source, p, rhs)
				}
				expr = binops["->"](source, expr, rhs)
			}
		}
	}
	if name == "amp" {
		for range c.(ast.Many) {
			expr = rel.NewFunction(source, rel.NewIdentExpr(*parser.NewScanner("-"), "-"), expr)
		}
	}
	return expr
}

// let PATTERN                     = EXPR1; EXPR2
// let c.(ast.One).Node.One("...") = expr;  rhs
// EXPR1 -> \PATTERN EXPR2
func (pc ParseContext) compileLet(c ast.Children) rel.Expr {
	exprs := c.(ast.One).Node.Many(exprTag)
	expr := pc.CompileExpr(exprs[0].(ast.Branch))
	rhs := pc.CompileExpr(exprs[1].(ast.Branch))
	source := c.Scanner()

	p := pc.compilePattern(c.(ast.One).Node.(ast.Branch))
	rhs = rel.NewFunction(source, p, rhs)

	if c.(ast.One).Node.One("rec") != nil {
		fix, fixt := FixFuncs()
		expr = rel.NewRecursionExpr(c.Scanner(), p, expr, fix, fixt)
	}

	expr = binops["->"](source, expr, rhs)
	return expr
}

func (pc ParseContext) compileUnop(b ast.Branch, c ast.Children) rel.Expr {
	ops := c.(ast.Many)
	result := pc.CompileExpr(b.One(exprTag).(ast.Branch))
	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i].One("").(ast.Leaf).Scanner()
		f := unops[op.String()]
		source, err := parser.MergeScanners(op, result.Source())
		if err != nil {
			// TODO: Figure out why some exprs don't have usable sources (could be native funcs).
			source = op
		}
		result = f(source, result)
	}
	return result
}

func (pc ParseContext) compileBinop(b ast.Branch, c ast.Children) rel.Expr {
	ops := c.(ast.Many)
	args := b.Many(exprTag)
	result := pc.CompileExpr(args[0].(ast.Branch))
	for i, arg := range args[1:] {
		op := ops[i].One("").(ast.Leaf).Scanner()
		f := binops[op.String()]
		rhs := pc.CompileExpr(arg.(ast.Branch))
		source, err := parser.MergeScanners(op, result.Source(), rhs.Source())
		if err != nil {
			// TODO: Figure out why some exprs don't have usable sources (could be native funcs).
			source = op
		}
		result = f(source, result, rhs)
	}
	return result
}

func (pc ParseContext) compileCompare(b ast.Branch, c ast.Children) rel.Expr {
	args := b.Many(exprTag)
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
	return rel.NewCompareExpr(ops[0].One("").(ast.Leaf).Scanner(), argExprs, comps, opStrs)
}

func (pc ParseContext) compileRbinop(b ast.Branch, c ast.Children) rel.Expr {
	ops := c.(ast.Many)
	args := b[exprTag].(ast.Many)
	result := pc.CompileExpr(args[len(args)-1].(ast.Branch))
	for i := len(args) - 2; i >= 0; i-- {
		op := ops[i].One("").(ast.Leaf).Scanner()
		f, has := binops[op.String()]
		if !has {
			panic("rbinop %q not found")
		}
		result = f(op, pc.CompileExpr(args[i].(ast.Branch)), result)
	}
	return result
}

func (pc ParseContext) compileIf(b ast.Branch, c ast.Children) rel.Expr {
	loggingOnce.Do(func() {
		log.Error(context.Background(),
			errors.New("operator if is deprecated and will be removed soon, please use operator cond instead. "+
				"Operator cond sample: let a = cond ( 2 > 1 : 1, 2 > 3 :2, * : 3)"))
	})

	result := pc.CompileExpr(b.One(exprTag).(ast.Branch))
	source := result.Source()
	for _, ifelse := range c.(ast.Many) {
		t := pc.CompileExpr(ifelse.One("t").(ast.Branch))
		var f rel.Expr = rel.None
		if fNode := ifelse.One("f"); fNode != nil {
			f = pc.CompileExpr(fNode.(ast.Branch))
		}
		result = rel.NewIfElseExpr(source, result, t, f)
	}
	return result
}

func (pc ParseContext) compileCond(c ast.Children) rel.Expr {
	if controlVar := c.(ast.One).Node.(ast.Branch)["controlVar"]; controlVar != nil {
		return pc.compileCondWithControlVar(c)
	}
	return pc.compileCondWithoutControlVar(c)
}

func (pc ParseContext) compileCondWithControlVar(c ast.Children) rel.Expr {
	conditions := pc.compileCondElements(c.(ast.One).Node.(ast.Branch)["condition"].(ast.Many)...)
	values := pc.compileCondExprs(c.(ast.One).Node.(ast.Branch)["value"].(ast.Many)...)

	if len(conditions) != len(values) {
		panic("mismatch between conditions and values")
	}

	conditionPairs := []rel.PatternExprPair{}
	for i, condition := range conditions {
		conditionPairs = append(conditionPairs, rel.NewPatternExprPair(condition, values[i]))
	}

	controlVar := c.(ast.One).Node.(ast.Branch)["controlVar"]
	return rel.NewCondPatternControlVarExpr(c.(ast.One).Node.Scanner(),
		pc.CompileExpr(controlVar.(ast.One).Node.(ast.Branch)),
		conditionPairs...)
}

func (pc ParseContext) compileCondElements(elements ...ast.Node) []rel.Pattern {
	result := make([]rel.Pattern, 0, len(elements))
	for _, element := range elements {
		name, c := which(element.(ast.Branch), "pattern")
		if c == nil {
			panic(fmt.Errorf("misshapen node AST: %v", element.(ast.Branch)))
		}

		if name == "pattern" {
			pattern := pc.compilePattern(element.(ast.Branch))
			if pattern != nil {
				result = append(result, pattern)
			}
		}
	}

	return result
}

func (pc ParseContext) compileCondWithoutControlVar(c ast.Children) rel.Expr {
	keys := c.(ast.One).Node.(ast.Branch)["key"]
	values := c.(ast.One).Node.(ast.Branch)["value"]
	var result rel.Expr
	if keys != nil && values != nil {
		keyExprs := pc.compileCondExprs(keys.(ast.Many)...)
		valueExprs := pc.compileCondExprs(values.(ast.Many)...)
		entryExprs := pc.compileDictEntryExprs(c, keyExprs, valueExprs)
		if entryExprs != nil {
			// Generates type DictExpr always to make sure it is easy to do Eval, only process type DictExpr.
			result = rel.NewDictExpr(c.(ast.One).Node.Scanner(), false, true, entryExprs...)
		} else {
			result = rel.NewDict(false)
		}
	}

	// Note, the default case `_:expr` which can match anything is parsed to condition/value pairs by current syntax.
	return rel.NewCondExpr(c.(ast.One).Node.Scanner(), result)
}

func (pc ParseContext) compilePostfixAndTouch(b ast.Branch, c ast.Children) rel.Expr {
	if _, has := b["touch"]; has {
		panic("unfinished")
	}
	switch c.Scanner().String() {
	case "count":
		return rel.NewCountExpr(b.Scanner(), pc.CompileExpr(b.One(exprTag).(ast.Branch)))
	case "single":
		return rel.NewSingleExpr(b.Scanner(), pc.CompileExpr(b.One(exprTag).(ast.Branch)))
	default:
		panic("wat?")
	}

	// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
	// result := p.parseExpr(b.One(exprTag).(ast.Branch))
}

func (pc ParseContext) compileCallGet(b ast.Branch) rel.Expr {
	var result rel.Expr
	if expr := b.One(exprTag); expr != nil {
		result = pc.CompileExpr(expr.(ast.Branch))
	} else {
		result = pc.compileGet(rel.DotIdent, b.One("get"))
	}
	for _, part := range b.Many("tail_op") {
		if safe := part.One("safe_tail"); safe != nil {
			result = pc.compileSafeTails(result, part.One("safe_tail"))
		} else {
			result = pc.compileTail(result, part.One("tail"))
		}
	}
	return result
}

func (pc ParseContext) compileTail(base rel.Expr, tail ast.Node) rel.Expr {
	if tail != nil {
		if call := tail.One("call"); call != nil {
			args := call.Many("arg")
			exprs := make([]ast.Node, 0, len(args))
			for _, arg := range args {
				exprs = append(exprs, arg.One(exprTag))
			}
			for _, arg := range pc.compileExprs(exprs...) {
				base = rel.NewCallExpr(handleAccessScanners(base.Source(), call.Scanner()), base, arg)
			}
		}
		base = pc.compileGet(base, tail.One("get"))
	}
	return base
}

func (pc ParseContext) compileTailFunc(tail ast.Node) rel.SafeTailCallback {
	if tail != nil {
		if call := tail.One("call"); call != nil {
			args := call.Many("arg")
			exprs := make([]ast.Node, 0, len(args))
			for _, arg := range args {
				exprs = append(exprs, arg.One("expr"))
			}
			compiledExprs := pc.compileExprs(exprs...)
			return func(v rel.Value, local rel.Scope) (rel.Value, error) {
				for _, arg := range compiledExprs {
					a, err := arg.Eval(local)
					if err != nil {
						return nil, err
					}
					//TODO: scanner won't highlight calls properly in safe call
					v, err = rel.SafeSetCall(v.(rel.Set), a)
					if err != nil {
						return nil, err
					}
				}
				return v, nil
			}
		}
		if get := tail.One("get"); get != nil {
			var scanner parser.Scanner
			var attr string
			if ident := get.One("IDENT"); ident != nil {
				scanner = ident.One("").(ast.Leaf).Scanner()
				attr = scanner.String()
			}
			if str := get.One("STR"); str != nil {
				scanner = str.One("").Scanner()
				attr = parseArraiString(scanner.String())
			}
			return func(v rel.Value, local rel.Scope) (rel.Value, error) {
				return rel.NewDotExpr(handleAccessScanners(v.Source(), scanner), v, attr).Eval(local)
			}
		}
	}
	panic("no tail")
}

func (pc ParseContext) compileGet(base rel.Expr, get ast.Node) rel.Expr {
	if get != nil {
		var scanner parser.Scanner
		var attr string
		if ident := get.One("IDENT"); ident != nil {
			scanner = ident.One("").(ast.Leaf).Scanner()
			attr = scanner.String()
		}
		if str := get.One("STR"); str != nil {
			scanner = str.One("").Scanner()
			attr = parseArraiString(scanner.String())
		}
		base = rel.NewDotExpr(handleAccessScanners(base.Source(), scanner), base, attr)
	}
	return base
}

func (pc ParseContext) compileSafeTails(base rel.Expr, tail ast.Node) rel.Expr {
	if tail != nil {
		firstSafe := tail.One("first_safe").One("tail")
		safeCallback := func(tailFunc rel.SafeTailCallback) rel.SafeTailCallback {
			return func(v rel.Value, local rel.Scope) (rel.Value, error) {
				val, err := tailFunc(v, local)
				if err != nil {
					switch e := err.(type) {
					case rel.NoReturnError:
						return nil, nil
					case rel.ContextErr:
						if _, isMissingAttrError := e.NextErr().(rel.MissingAttrError); isMissingAttrError {
							return nil, nil
						}
					}
					return nil, err
				}
				return val, nil
			}
		}

		exprStates := []rel.SafeTailCallback{safeCallback(pc.compileTailFunc(firstSafe))}
		fallback := pc.CompileExpr(tail.One("fall").(ast.Branch))

		for _, o := range tail.Many("ops") {
			if safeTail := o.One("safe"); safeTail != nil {
				exprStates = append(exprStates, safeCallback(pc.compileTailFunc(safeTail.One("tail"))))
			} else if tail := o.One("tail"); tail != nil {
				exprStates = append(exprStates, pc.compileTailFunc(tail))
			} else {
				panic("wat")
			}
		}

		return rel.NewSafeTailExpr(tail.Scanner(), fallback, base, exprStates)
	}
	//TODO: panic?
	return base
}

func handleAccessScanners(base, access parser.Scanner) parser.Scanner {
	if len(base.String()) == 0 {
		return access
	}
	// handles .a
	if base.String() == "." {
		return *access.Skip(-1)
	}
	scanner, err := parser.MergeScanners(base, access)
	if err != nil {
		panic(err)
	}
	return scanner
}

func (pc ParseContext) compileRelation(c ast.Children) rel.Expr {
	names := parseNames(c.(ast.One).Node.(ast.Branch)["names"].(ast.One).Node.(ast.Branch))
	tuples := c.(ast.One).Node.(ast.Branch)["tuple"].(ast.Many)
	tupleExprs := make([][]rel.Expr, 0, len(tuples))
	for _, tuple := range tuples {
		tupleExprs = append(tupleExprs, pc.compileExprs(tuple.(ast.Branch)["v"].(ast.Many)...))
	}
	result, err := rel.NewRelationExpr(
		c.(ast.One).Node.(ast.Branch)["names"].(ast.One).Node.(ast.Branch).Scanner(),
		names,
		tupleExprs...,
	)
	if err != nil {
		panic(err)
	}
	return result
}

func (pc ParseContext) compileSet(c ast.Children) rel.Expr {
	if elts := c.(ast.One).Node.(ast.Branch)["elt"]; elts != nil {
		return rel.NewSetExpr(elts.(ast.Many).Scanner(), pc.compileExprs(elts.(ast.Many)...)...)
	}
	return rel.NewSetExpr(c.(ast.One).Node.Scanner())
}

func (pc ParseContext) compileDict(c ast.Children) rel.Expr {
	entryExprs := pc.compileDictEntryExprs(c, nil, nil)
	if entryExprs != nil {
		return rel.NewDictExpr(c.(ast.One).Node.Scanner(), false, false, entryExprs...)
	}

	return rel.NewDict(false)
}

func (pc ParseContext) compileDictEntryExprs(c ast.Children, keyExprs []rel.Expr,
	valueExprs []rel.Expr) []rel.DictEntryTupleExpr {
	// C* "{" C* dict=((key=@ ":" value=@):",",?) "}" C*
	keys := c.(ast.One).Node.(ast.Branch)["key"]
	values := c.(ast.One).Node.(ast.Branch)["value"]
	if (keys != nil) || (values != nil) {
		if (keys != nil) && (values != nil) {
			if keyExprs == nil {
				keyExprs = pc.compileExprs(keys.(ast.Many)...)
			}
			if valueExprs == nil {
				valueExprs = pc.compileExprs(values.(ast.Many)...)
			}
			if len(keyExprs) == len(valueExprs) {
				entryExprs := make([]rel.DictEntryTupleExpr, 0, len(keyExprs))
				for i, keyExpr := range keyExprs {
					valueExpr := valueExprs[i]
					entryExprs = append(entryExprs, rel.NewDictEntryTupleExpr(keys.Scanner(), keyExpr, valueExpr))
				}
				return entryExprs
			}
		}
		panic("mismatch between dict keys and values")
	}
	return nil
}

func (pc ParseContext) compileArray(c ast.Children) rel.Expr {
	if exprs := pc.compileSparseItems(c); len(exprs) > 0 {
		return rel.NewArrayExpr(c.Scanner(), exprs...)
	}
	return rel.NewArray()
}

func (pc ParseContext) compileBytes(c ast.Children) rel.Expr {
	if items := c.(ast.One).Node.(ast.Branch)["item"]; items != nil {
		//TODO: support sparse bytes
		return rel.NewBytesExpr(items.Scanner(), pc.compileExprs(items.(ast.Many)...)...)
	}
	return rel.NewBytes([]byte{})
}

func (pc ParseContext) compileExprs(exprs ...ast.Node) []rel.Expr {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		result = append(result, pc.CompileExpr(expr.(ast.Branch)))
	}
	return result
}

func (pc ParseContext) compileSparseItems(c ast.Children) []rel.Expr {
	var nodes []ast.Node
	if firstItem := c.(ast.One).Node.One("first_item"); firstItem != nil {
		nodes = []ast.Node{firstItem}
		if items := c.(ast.One).Node.Many("item"); items != nil {
			nodes = append(nodes, items...)
		}
	}
	result := make([]rel.Expr, 0, len(nodes))
	for _, expr := range nodes {
		if expr.One("empty") != nil {
			result = append(result, nil)
			continue
		}
		result = append(result, pc.CompileExpr(expr.(ast.Branch)))
	}
	return result
}

// compileCondExprs parses conditons/keys and values expressions for syntax `cond`.
func (pc ParseContext) compileCondExprs(exprs ...ast.Node) []rel.Expr {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		var exprResult rel.Expr

		name, c := which(expr.(ast.Branch), exprTag)
		if c == nil {
			panic(fmt.Errorf("misshapen node AST: %v", expr.(ast.Branch)))
		}

		if name == exprTag {
			switch c := c.(type) {
			case ast.One:
				exprResult = pc.CompileExpr(c.Node.(ast.Branch))
			case ast.Many:
				if len(c) == 1 {
					exprResult = pc.CompileExpr(c[0].(ast.Branch))
				} else {
					var elements []rel.Expr
					for _, e := range c {
						expr := pc.CompileExpr(e.(ast.Branch))
						elements = append(elements, expr)
					}
					exprResult = rel.NewArrayExpr(c.Scanner(), elements...)
				}
			}
		}

		if exprResult != nil {
			result = append(result, exprResult)
		}
	}
	return result
}

func (pc ParseContext) compileFunction(b ast.Branch) rel.Expr {
	ident := b.One("IDENT")
	expr := pc.CompileExpr(b.One(exprTag).(ast.Branch))
	source := ident.One("").Scanner()
	return rel.NewFunction(b.Scanner(), rel.NewIdentExpr(source, source.String()), expr)
}

func (pc ParseContext) compileMacro(b ast.Branch) rel.Expr {
	childast := b.One("embed").One("subgrammar").One("ast")
	if value := childast.One("value"); value != nil {
		return value.(ast.Extra).Data.(rel.Expr)
	} else {
		return rel.ASTNodeToValue(childast)
	}
}

func (pc ParseContext) compilePackage(c ast.Children) rel.Expr {
	pkg := c.(ast.One).Node.(ast.Branch)
	if std, has := pkg["std"]; has {
		ident := std.(ast.One).Node.One("IDENT").One("")
		pkgName := ident.(ast.Leaf).Scanner()
		return NewPackageExpr(pkgName, rel.NewDotExpr(pkgName, rel.DotIdent, pkgName.String()))
	}

	if str := pkg.One("PKGPATH"); str != nil {
		scanner := str.One("").(ast.Leaf).Scanner()
		name := scanner.String()
		if strings.HasPrefix(name, "/") {
			filepath := strings.Trim(name, "/")
			fromRoot := pkg["dot"] == nil
			if pc.SourceDir == "" {
				panic(fmt.Errorf("local import %q invalid; no local context", name))
			}
			return rel.NewCallExpr(scanner,
				NewPackageExpr(scanner, importLocalFile(fromRoot)),
				rel.NewString([]rune(path.Join(pc.SourceDir, filepath))),
			)
		}
		return rel.NewCallExpr(scanner, NewPackageExpr(scanner, importExternalContent()), rel.NewString([]rune(name)))
	}
	return NewPackageExpr(pkg.Scanner(), rel.DotIdent)
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
			attr, err := rel.NewAttrExpr(pair.One("v").(ast.Branch).Scanner(), k, v)
			if err != nil {
				panic(err)
			}
			attrs = append(attrs, attr)
		}
		return rel.NewTupleExpr(c.(ast.One).Node.Scanner(), attrs...)
	}
	return rel.EmptyTuple
}

func (pc ParseContext) compileIdent(c ast.Children) rel.Expr {
	s := c.(ast.One).Node.One("").Scanner()
	switch s.String() {
	case "true":
		return rel.True
	case "false":
		return rel.False
	}
	return rel.NewIdentExpr(s, s.String())
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

func (pc ParseContext) compileChar(c ast.Children) rel.Expr {
	char := c.(ast.One).Node.One("").Scanner().String()
	quote := "\""
	if char[0] == '\'' {
		quote = "'"
	}
	runes := []rune(parseArraiStringFragment(char, quote, ""))
	return rel.NewNumber(float64(runes[0]))
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

func dotUnary(f binOpFunc) unOpFunc {
	return func(scanner parser.Scanner, e rel.Expr) rel.Expr {
		return f(scanner, rel.DotIdent, e)
	}
}

type unOpFunc func(scanner parser.Scanner, e rel.Expr) rel.Expr

var unops = map[string]unOpFunc{
	"+":  rel.NewPosExpr,
	"-":  rel.NewNegExpr,
	"^":  rel.NewPowerSetExpr,
	"!":  rel.NewNotExpr,
	"*":  rel.NewEvalExpr,
	"//": NewPackageExpr,
	"=>": dotUnary(rel.NewDArrowExpr),
	">>": dotUnary(rel.NewSequenceMapExpr),
	":>": dotUnary(rel.NewTupleMapExpr),
}

type binOpFunc func(scanner parser.Scanner, a, b rel.Expr) rel.Expr

var binops = map[string]binOpFunc{
	"->":      rel.NewArrowExpr,
	"=>":      rel.NewDArrowExpr,
	">>":      rel.NewSequenceMapExpr,
	">>>":     rel.NewIndexedSequenceMapExpr,
	":>":      rel.NewTupleMapExpr,
	"orderby": rel.NewOrderByExpr,
	"order":   rel.NewOrderExpr,
	"rank":    rel.NewRankExpr,
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
	"<:":  func(a, b rel.Value) bool { return b.(rel.Set).Has(a) },
	"!<:": func(a, b rel.Value) bool { return !b.(rel.Set).Has(a) },
	"=":   func(a, b rel.Value) bool { return a.Equal(b) },
	"!=":  func(a, b rel.Value) bool { return !a.Equal(b) },
	"<":   func(a, b rel.Value) bool { return a.Less(b) },
	">":   func(a, b rel.Value) bool { return b.Less(a) },
	"<=":  func(a, b rel.Value) bool { return !b.Less(a) },
	">=":  func(a, b rel.Value) bool { return !a.Less(b) },

	"(<)":   func(a, b rel.Value) bool { return subset(a, b) },
	"(>)":   func(a, b rel.Value) bool { return subset(b, a) },
	"(<=)":  func(a, b rel.Value) bool { return subsetOrEqual(a, b) },
	"(>=)":  func(a, b rel.Value) bool { return subsetOrEqual(b, a) },
	"(<>)":  func(a, b rel.Value) bool { return subsetOrSuperset(a, b) },
	"(<>=)": func(a, b rel.Value) bool { return subsetSupersetOrEqual(b, a) },

	"!(<)":   func(a, b rel.Value) bool { return !subset(a, b) },
	"!(>)":   func(a, b rel.Value) bool { return !subset(b, a) },
	"!(<=)":  func(a, b rel.Value) bool { return !subsetOrEqual(a, b) },
	"!(>=)":  func(a, b rel.Value) bool { return !subsetOrEqual(b, a) },
	"!(<>)":  func(a, b rel.Value) bool { return !subsetOrSuperset(a, b) },
	"!(<>=)": func(a, b rel.Value) bool { return !subsetSupersetOrEqual(b, a) },
}
