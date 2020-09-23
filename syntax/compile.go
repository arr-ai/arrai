package syntax

//TODO: should I add context to compile functions?
import (
	"context"
	"errors"
	"fmt"
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
func Compile(ctx context.Context, filePath, source string) (rel.Expr, error) {
	dirpath := "."
	if filePath != "" {
		if filePath == NoPath {
			dirpath = NoPath
		} else {
			dirpath = filepath.Dir(filePath)
		}
	}
	pc := ParseContext{SourceDir: dirpath}
	// bundle run will always get absolute UNIX filePath. This needs to happen
	// with windows too.
	if !filepath.IsAbs(filePath) && !isRunningBundle(ctx) {
		var err error
		filePath, err = filepath.Rel(".", filePath)
		if err != nil {
			return nil, err
		}
	}
	ast, err := pc.Parse(ctx, parser.NewScannerWithFilename(source, filePath))
	if err != nil {
		return nil, err
	}

	return pc.CompileExpr(ctx, ast)
}

func MustCompile(ctx context.Context, filePath, source string) rel.Expr {
	expr, err := Compile(ctx, filePath, source)
	if err != nil {
		panic(err)
	}
	return expr
}

func (pc ParseContext) CompileExpr(ctx context.Context, b ast.Branch) (rel.Expr, error) {
	// Note: please make sure if it is necessary to add new syntax name before `expr`.
	name, c := which(b,
		"amp", "arrow", "let", "unop", "binop", "compare", "rbinop", "if", "get",
		"tail_op", "postfix", "touch", "get", "rel", "set", "dict", "array", "bytes",
		"embed", "op", "fn", "pkg", "tuple", "xstr", "IDENT", "STR", "NUM", "CHAR",
		"cond", exprTag,
	)
	if c == nil {
		return nil, fmt.Errorf("misshapen node AST: %v", b)
	}
	switch name {
	case "amp", "arrow":
		return pc.compileArrow(ctx, b, name, c)
	case "let":
		return pc.compileLet(ctx, c)
	case "unop":
		return pc.compileUnop(ctx, b, c)
	case "binop":
		return pc.compileBinop(ctx, b, c)
	case "compare":
		return pc.compileCompare(ctx, b, c)
	case "rbinop":
		return pc.compileRbinop(ctx, b, c)
	case "if":
		return pc.compileIf(ctx, b, c)
	case "cond":
		return pc.compileCond(ctx, c)
	case "postfix", "touch":
		return pc.compilePostfixAndTouch(ctx, b, c)
	case "get", "tail_op":
		return pc.compileCallGet(ctx, b)
	case "rel":
		return pc.compileRelation(ctx, b, c)
	case "set":
		return pc.compileSet(ctx, b, c)
	case "dict":
		return pc.compileDict(ctx, b, c)
	case "array":
		return pc.compileArray(ctx, b, c)
	case "bytes":
		return pc.compileBytes(ctx, b, c)
	case "embed":
		return pc.compileMacro(ctx, b), nil
	case "fn":
		return pc.compileFunction(ctx, b)
	case "pkg":
		return pc.compilePackage(ctx, b, c)
	case "tuple":
		return pc.compileTuple(ctx, b, c)
	case "IDENT":
		return pc.compileIdent(ctx, c), nil
	case "STR":
		return pc.compileString(ctx, c), nil
	case "xstr":
		return pc.compileExpandableString(ctx, b, c)
	case "NUM":
		return pc.compileNumber(ctx, c)
	case "CHAR":
		return pc.compileChar(ctx, c), nil
	case exprTag:
		result, err := pc.compileExpr(ctx, b, c)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, fmt.Errorf("unhandled node: %v", b)
}

func (pc ParseContext) compilePattern(ctx context.Context, b ast.Branch) (rel.Pattern, error) {
	if ptn := b.One("pattern"); ptn != nil {
		return pc.compilePattern(ctx, ptn.(ast.Branch))
	}
	if arr := b.One("array"); arr != nil {
		return pc.compileArrayPattern(ctx, arr.(ast.Branch))
	}
	if tuple := b.One("tuple"); tuple != nil {
		return pc.compileTuplePattern(ctx, tuple.(ast.Branch))
	}
	if dict := b.One("dict"); dict != nil {
		return pc.compileDictPattern(ctx, dict.(ast.Branch))
	}
	if set := b.One("set"); set != nil {
		return pc.compileSetPattern(ctx, set.(ast.Branch))
	}
	if extra := b.One("extra"); extra != nil {
		return pc.compileExtraElementPattern(ctx, extra.(ast.Branch)), nil
	}
	if expr := b.Many("exprpattern"); expr != nil {
		var elements []rel.Expr
		for _, e := range expr {
			expr, err := pc.CompileExpr(ctx, e.(ast.Branch))
			if err != nil {
				return nil, err
			}
			elements = append(elements, expr)
		}
		if len(elements) > 0 {
			return rel.NewExprsPattern(elements...), nil
		}
	}
	p, err := pc.CompileExpr(ctx, b)
	if err != nil {
		return nil, err
	}
	return rel.NewExprPattern(p), nil
}

func (pc ParseContext) compileExtraElementPattern(_ context.Context, b ast.Branch) rel.Pattern {
	var ident string
	if id := b.One("ident"); id != nil {
		ident = id.Scanner().String()
	}
	return rel.NewExtraElementPattern(ident)
}

func (pc ParseContext) compilePatterns(ctx context.Context, exprs ...ast.Node) ([]rel.Pattern, error) {
	result := make([]rel.Pattern, 0, len(exprs))
	for _, expr := range exprs {
		p, err := pc.compilePattern(ctx, expr.(ast.Branch))
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (pc ParseContext) compileSparsePatterns(ctx context.Context, b ast.Branch) ([]rel.FallbackPattern, error) {
	var nodes []ast.Node
	if firstItem, exists := b["first_item"]; exists {
		nodes = []ast.Node{firstItem.(ast.One).Node}
		if items, exists := b["item"]; exists {
			for _, i := range items.(ast.Many) {
				nodes = append(nodes, i)
			}
		}
	}
	result := make([]rel.FallbackPattern, 0, len(nodes))
	for _, expr := range nodes {
		if expr.One("empty") != nil {
			result = append(result, rel.NewFallbackPattern(nil, nil))
			continue
		}
		ptn, err := pc.compilePattern(ctx, expr.(ast.Branch))
		if err != nil {
			return nil, err
		}
		if fall := expr.One("fall"); fall != nil {
			fallback, err := pc.CompileExpr(ctx, fall.(ast.Branch))
			if err != nil {
				return nil, err
			}
			result = append(result, rel.NewFallbackPattern(ptn, fallback))
			continue
		}
		result = append(result, rel.NewFallbackPattern(ptn, nil))
	}
	return result, nil
}

func (pc ParseContext) compileArrayPattern(ctx context.Context, b ast.Branch) (rel.Pattern, error) {
	arrPatterns, err := pc.compileSparsePatterns(ctx, b)
	if err != nil {
		return nil, err
	}
	return rel.NewArrayPattern(arrPatterns...), nil
}

func (pc ParseContext) compileTuplePattern(ctx context.Context, b ast.Branch) (_ rel.Pattern, err error) {
	if pairs := b.Many("pairs"); pairs != nil {
		attrs := make([]rel.TuplePatternAttr, 0, len(pairs))
		for _, pair := range pairs {
			var k string
			var v rel.Pattern

			if extra := pair.One("extra"); extra != nil {
				v, err = pc.compilePattern(ctx, pair.(ast.Branch))
				if err != nil {
					return nil, err
				}
				attrs = append(attrs, rel.NewTuplePatternAttr(k, rel.NewFallbackPattern(v, nil)))
			} else {
				v, err = pc.compilePattern(ctx, pair.One("v").(ast.Branch))
				if err != nil {
					return nil, err
				}
				if name := pair.One("name"); name != nil {
					k = parseName(name.(ast.Branch))
				} else {
					k = v.String()
				}

				tail := pair.One("tail")
				fall := pair.One("v").One("fall")
				if fall == nil {
					attrs = append(attrs, rel.NewTuplePatternAttr(k, rel.NewFallbackPattern(v, nil)))
				} else if tail != nil && fall != nil {
					fallExpr, err := pc.CompileExpr(ctx, fall.(ast.Branch))
					if err != nil {
						return nil, err
					}
					attrs = append(attrs, rel.NewTuplePatternAttr(k, rel.NewFallbackPattern(v, fallExpr)))
				} else {
					return nil, errors.New("fallback item does not match")
				}
			}
		}
		return rel.NewTuplePattern(attrs...), nil
	}
	return rel.NewTuplePattern(), nil
}

func (pc ParseContext) compileDictPattern(ctx context.Context, b ast.Branch) (_ rel.Pattern, err error) {
	if pairs := b.Many("pairs"); pairs != nil {
		entryPtns := make([]rel.DictPatternEntry, 0, len(pairs))
		for _, pair := range pairs {
			if extra := pair.One("extra"); extra != nil {
				p := pc.compileExtraElementPattern(ctx, extra.(ast.Branch))
				entryPtns = append(entryPtns, rel.NewDictPatternEntry(nil, rel.NewFallbackPattern(p, nil)))
				continue
			}
			key := pair.One("key")
			value := pair.One("value")
			keyExpr, err := pc.CompileExpr(ctx, key.(ast.Branch))
			if err != nil {
				return nil, err
			}
			valuePtn, err := pc.compilePattern(ctx, value.(ast.Branch))
			if err != nil {
				return nil, err
			}

			tail := key.One("tail")
			fall := value.One("fall")
			if fall == nil {
				entryPtns = append(entryPtns, rel.NewDictPatternEntry(keyExpr, rel.NewFallbackPattern(valuePtn, nil)))
			} else if tail != nil && fall != nil {
				fallExpr, err := pc.CompileExpr(ctx, fall.(ast.Branch))
				if err != nil {
					return nil, err
				}
				entryPtns = append(entryPtns, rel.NewDictPatternEntry(keyExpr,
					rel.NewFallbackPattern(valuePtn, fallExpr)))
			} else {
				panic("fallback item does not match")
			}
		}
		return rel.NewDictPattern(entryPtns...), nil
	}
	return rel.NewDictPattern(), nil
}

func (pc ParseContext) compileSetPattern(ctx context.Context, b ast.Branch) (rel.Pattern, error) {
	if elts := b["elt"]; elts != nil {
		patterns, err := pc.compilePatterns(ctx, elts.(ast.Many)...)
		if err != nil {
			return nil, err
		}
		return rel.NewSetPattern(patterns...), nil
	}
	return rel.NewSetPattern(), nil
}

func (pc ParseContext) compileArrow(ctx context.Context, b ast.Branch, name string, c ast.Children) (rel.Expr, error) {
	expr, err := pc.CompileExpr(ctx, b[exprTag].(ast.One).Node.(ast.Branch))
	if err != nil {
		return nil, err
	}
	source := c.Scanner()
	if arrows, has := b["arrow"]; has {
		for _, arrow := range arrows.(ast.Many) {
			branch := arrow.(ast.Branch)
			part, d := which(branch, "nest", "unnest", "ARROW", "binding", "FILTER")
			switch part {
			case "nest":
				expr = parseNest(expr, branch["nest"].(ast.One).Node.(ast.Branch))
			case "unnest":
				panic("unfinished")
			case "ARROW":
				op := d.(ast.One).Node.One("").(ast.Leaf).Scanner()
				f := binops[op.String()]
				rhs, err := pc.CompileExpr(ctx, arrow.(ast.Branch)[exprTag].(ast.One).Node.(ast.Branch))
				if err != nil {
					return nil, err
				}
				expr = f(b.Scanner(), expr, rhs)
			case "binding":
				rhs, err := pc.CompileExpr(ctx, arrow.(ast.Branch)[exprTag].(ast.One).Node.(ast.Branch))
				if err != nil {
					return nil, err
				}
				if pattern := arrow.One("pattern"); pattern != nil {
					p, err := pc.compilePattern(ctx, pattern.(ast.Branch))
					if err != nil {
						return nil, err
					}
					rhs = rel.NewFunction(source, p, rhs)
				}
				expr = binops["->"](source, expr, rhs)
			case "FILTER":
				transform, err := pc.CompileExpr(ctx, arrow.(ast.Branch))
				if err != nil {
					return nil, err
				}
				t := transform.(rel.CondPatternControlVarExpr)
				conditions := t.Conditions()
				trueConds := make([]rel.PatternExprPair, 0, len(conditions))
				for _, c := range conditions {
					trueConds = append(trueConds, rel.NewPatternExprPair(c.Pattern(), rel.True))
				}
				pred := rel.NewCondPatternControlVarExpr(
					t.ExprScanner.Src,
					t.Control(),
					trueConds...,
				)
				lhs := rel.NewWhereExpr(source, expr, pred)
				expr = rel.NewDArrowExpr(source, lhs, transform)
			}
		}
	}
	if name == "amp" {
		for range c.(ast.Many) {
			expr = rel.NewFunction(source, rel.IdentPattern("-"), expr)
		}
	}
	return expr, nil
}

// let PATTERN                     = EXPR1;      EXPR2
// let c.(ast.One).Node.One("...") = expr(lhs);  rhs
// EXPR1 -> \PATTERN EXPR2
func (pc ParseContext) compileLet(ctx context.Context, c ast.Children) (rel.Expr, error) {
	exprs := c.(ast.One).Node.Many(exprTag)
	expr, err := pc.CompileExpr(ctx, exprs[0].(ast.Branch))
	if err != nil {
		return nil, err
	}
	rhs, err := pc.CompileExpr(ctx, exprs[1].(ast.Branch))
	if err != nil {
		return nil, err
	}
	source := c.Scanner()

	var p rel.Pattern
	if pat := c.(ast.One).Node.(ast.Branch).One("pat"); pat != nil {
		p, err = pc.compilePattern(ctx, pat.(ast.Branch))
		if err != nil {
			return nil, err
		}
	} else {
		return rel.NewDynLetExpr(c.Scanner(), expr, rhs), nil
	}
	rhs = rel.NewFunction(source, p, rhs)

	if c.(ast.One).Node.One("rec") != nil {
		fix, fixt := FixFuncs()
		identPattern, is := p.(rel.IdentPattern)
		if !is {
			return nil, fmt.Errorf("let rec parameter must be IDENT, not %v", p)
		}
		name := identPattern.Ident()
		expr = rel.NewRecursionExpr(c.Scanner(), name, expr, fix, fixt)
	}

	return binops["->"](source, expr, rhs), nil
}

func (pc ParseContext) compileUnop(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	ops := c.(ast.Many)
	result, err := pc.CompileExpr(ctx, b.One(exprTag).(ast.Branch))
	if err != nil {
		return nil, err
	}
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
	return result, nil
}

func (pc ParseContext) compileBinop(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	ops := c.(ast.Many)
	args := b.Many(exprTag)
	result, err := pc.CompileExpr(ctx, args[0].(ast.Branch))
	if err != nil {
		return nil, err
	}
	for i, arg := range args[1:] {
		op := ops[i].One("").(ast.Leaf).Scanner()
		f := binops[op.String()]
		rhs, err := pc.CompileExpr(ctx, arg.(ast.Branch))
		if err != nil {
			return nil, err
		}
		source, err := parser.MergeScanners(op, result.Source(), rhs.Source())
		if err != nil {
			// TODO: Figure out why some exprs don't have usable sources (could be native funcs).
			source = op
		}
		result = f(source, result, rhs)
	}
	return result, nil
}

func (pc ParseContext) compileCompare(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	args := b.Many(exprTag)
	argExprs := make([]rel.Expr, 0, len(args))
	comps := make([]rel.CompareFunc, 0, len(args))

	ops := c.(ast.Many)
	opStrs := make([]string, 0, len(ops))

	argExpr, err := pc.CompileExpr(ctx, args[0].(ast.Branch))
	if err != nil {
		return nil, err
	}
	argExprs = append(argExprs, argExpr)
	for i, arg := range args[1:] {
		op := ops[i].One("").(ast.Leaf).Scanner().String()

		argExpr, err := pc.CompileExpr(ctx, arg.(ast.Branch))
		if err != nil {
			return nil, err
		}
		argExprs = append(argExprs, argExpr)
		comps = append(comps, compareOps[op])

		opStrs = append(opStrs, op)
	}
	scanner, err := parser.MergeScanners(argExprs[0].Source(), argExprs[len(argExprs)-1].Source())
	if err != nil {
		return nil, err
	}
	return rel.NewCompareExpr(scanner, argExprs, comps, opStrs), nil
}

func (pc ParseContext) compileRbinop(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	ops := c.(ast.Many)
	args := b[exprTag].(ast.Many)
	result, err := pc.CompileExpr(ctx, args[len(args)-1].(ast.Branch))
	if err != nil {
		return nil, err
	}
	for i := len(args) - 2; i >= 0; i-- {
		op := ops[i].One("").(ast.Leaf).Scanner()
		f, has := binops[op.String()]
		if !has {
			panic("rbinop %q not found")
		}
		rBinOpArg, err := pc.CompileExpr(ctx, args[i].(ast.Branch))
		if err != nil {
			return nil, err
		}
		result = f(op, rBinOpArg, result)
	}
	return result, nil
}

func (pc ParseContext) compileIf(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	loggingOnce.Do(func() {
		log.Error(ctx,
			errors.New("operator if is deprecated and will be removed soon, please use operator cond instead. "+
				"Operator cond sample: let a = cond {2 > 1: 1, 2 > 3: 2, _: 3}"))
	})

	result, err := pc.CompileExpr(ctx, b.One(exprTag).(ast.Branch))
	if err != nil {
		return nil, err
	}
	source := result.Source()
	for _, ifelse := range c.(ast.Many) {
		t, err := pc.CompileExpr(ctx, ifelse.One("t").(ast.Branch))
		if err != nil {
			return nil, err
		}
		var f rel.Expr = rel.None
		if fNode := ifelse.One("f"); fNode != nil {
			f, err = pc.CompileExpr(ctx, fNode.(ast.Branch))
			if err != nil {
				return nil, err
			}
		}
		result = rel.NewIfElseExpr(source, result, t, f)
	}
	return result, nil
}

func (pc ParseContext) compileCond(ctx context.Context, c ast.Children) (rel.Expr, error) {
	if controlVar := c.(ast.One).Node.(ast.Branch)["controlVar"]; controlVar != nil {
		return pc.compileCondWithControlVar(ctx, c)
	}
	return pc.compileCondWithoutControlVar(ctx, c)
}

func (pc ParseContext) compileCondWithControlVar(ctx context.Context, c ast.Children) (rel.Expr, error) {
	conditions, err := pc.compileCondElements(ctx, c.(ast.One).Node.(ast.Branch)["condition"].(ast.Many)...)
	if err != nil {
		return nil, err
	}
	values, err := pc.compileCondExprs(ctx, c.(ast.One).Node.(ast.Branch)["value"].(ast.Many)...)
	if err != nil {
		return nil, err
	}

	if len(conditions) != len(values) {
		return nil, fmt.Errorf(
			"compileCondWithControlVar: mismatch between conditions and values: %s and %s",
			conditions, values,
		)
	}

	conditionPairs := []rel.PatternExprPair{}
	for i, condition := range conditions {
		conditionPairs = append(conditionPairs, rel.NewPatternExprPair(condition, values[i]))
	}

	controlVar := c.(ast.One).Node.(ast.Branch)["controlVar"]
	controlVarExpr, err := pc.CompileExpr(ctx, controlVar.(ast.One).Node.(ast.Branch))
	if err != nil {
		return nil, err
	}
	return rel.NewCondPatternControlVarExpr(c.(ast.One).Node.Scanner(), controlVarExpr, conditionPairs...), nil
}

func (pc ParseContext) compileCondElements(ctx context.Context, elements ...ast.Node) ([]rel.Pattern, error) {
	result := make([]rel.Pattern, 0, len(elements))
	for _, element := range elements {
		name, c := which(element.(ast.Branch), "pattern")
		if c == nil {
			return nil, fmt.Errorf("misshapen node AST: %v", element.(ast.Branch))
		}

		if name == "pattern" {
			pattern, err := pc.compilePattern(ctx, element.(ast.Branch))
			if err != nil {
				return nil, err
			}
			if pattern != nil {
				result = append(result, pattern)
			}
		}
	}

	return result, nil
}

func (pc ParseContext) compileCondWithoutControlVar(ctx context.Context, c ast.Children) (rel.Expr, error) {
	var result rel.Expr
	entryExprs, err := pc.compileDictEntryExprs(ctx, c.(ast.One).Node.(ast.Branch))
	if err != nil {
		return nil, err
	}
	if entryExprs != nil {
		// Generates type DictExpr always to make sure it is easy to do Eval, only process type DictExpr.
		result, err = rel.NewDictExpr(c.(ast.One).Node.Scanner(), false, true, entryExprs...)
		if err != nil {
			return nil, err
		}
	} else {
		result = rel.MustNewDict(false)
	}

	// Note, the default case `_:expr` which can match anything is parsed to condition/value pairs by current syntax.
	return rel.NewCondExpr(c.(ast.One).Node.Scanner(), result), nil
}

func (pc ParseContext) compilePostfixAndTouch(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	if _, has := b["touch"]; has {
		panic("unfinished")
	}
	switch c.Scanner().String() {
	case "count":
		countArg, err := pc.CompileExpr(ctx, b.One(exprTag).(ast.Branch))
		if err != nil {
			return nil, err
		}
		return rel.NewCountExpr(b.Scanner(), countArg), nil
	case "single":
		singleArg, err := pc.CompileExpr(ctx, b.One(exprTag).(ast.Branch))
		if err != nil {
			return nil, err
		}
		return rel.NewSingleExpr(b.Scanner(), singleArg), nil
	default:
		panic("wat?")
	}

	// touch -> ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")";
	// result := p.parseExpr(b.One(exprTag).(ast.Branch))
}

func (pc ParseContext) compileCallGet(ctx context.Context, b ast.Branch) (_ rel.Expr, err error) {
	var result rel.Expr
	if expr := b.One(exprTag); expr != nil {
		result, err = pc.CompileExpr(ctx, expr.(ast.Branch))
		if err != nil {
			return nil, err
		}
	} else {
		get := b.One("get")
		dot := get.One("dot")
		result = pc.compileGet(ctx, rel.NewDotIdent(dot.Scanner()), get)
	}
	for _, part := range b.Many("tail_op") {
		if safe := part.One("safe_tail"); safe != nil {
			result, err = pc.compileSafeTails(ctx, result, part.One("safe_tail"))
		} else {
			result, err = pc.compileTail(ctx, result, part.One("tail"))
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (pc ParseContext) compileTail(ctx context.Context, base rel.Expr, tail ast.Node) (rel.Expr, error) {
	if tail != nil {
		if call := tail.One("call"); call != nil {
			args := call.Many("arg")
			exprs := make([]ast.Node, 0, len(args))
			for _, arg := range args {
				exprs = append(exprs, arg.One(exprTag))
			}
			compiledExprs, err := pc.compileExprs(ctx, exprs...)
			if err != nil {
				return nil, err
			}
			for _, arg := range compiledExprs {
				base = rel.NewCallExpr(handleAccessScanners(base.Source(), call.Scanner()), base, arg)
			}
		}
		base = pc.compileGet(ctx, base, tail.One("get"))
	}
	return base, nil
}

func (pc ParseContext) compileTailFunc(ctx context.Context, tail ast.Node) (rel.SafeTailCallback, error) {
	if tail != nil {
		if call := tail.One("call"); call != nil {
			args := call.Many("arg")
			exprs := make([]ast.Node, 0, len(args))
			for _, arg := range args {
				exprs = append(exprs, arg.One("expr"))
			}

			compiledExprs, err := pc.compileExprs(ctx, exprs...)
			if err != nil {
				return nil, err
			}
			return func(ctx context.Context, v rel.Value, local rel.Scope) (rel.Value, error) {
				for _, arg := range compiledExprs {
					a, err := arg.Eval(ctx, local)
					if err != nil {
						return nil, err
					}
					//TODO: scanner won't highlight calls properly in safe call
					set, is := v.(rel.Set)
					if !is {
						return nil, fmt.Errorf("not a set: %v", v)
					}
					v, err = rel.SetCall(ctx, set, a)
					if err != nil {
						return nil, err
					}
				}
				return v, nil
			}, nil
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
			return func(ctx context.Context, v rel.Value, local rel.Scope) (rel.Value, error) {
				return rel.NewDotExpr(handleAccessScanners(v.Source(), scanner), v, attr).Eval(ctx, local)
			}, nil
		}
	}
	return nil, fmt.Errorf("compileTailFunc: tail AST malformed: %s", tail)
}

func (pc ParseContext) compileGet(_ context.Context, base rel.Expr, get ast.Node) rel.Expr {
	if get != nil {
		if names := get.One("names"); names != nil {
			inverse := get.One("") != nil
			attrs := parseNames(names.(ast.Branch))
			return rel.NewTupleProjectExpr(
				handleAccessScanners(base.Source(), names.Scanner()),
				base, inverse, attrs,
			)
		}

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

func (pc ParseContext) compileSafeTails(ctx context.Context, base rel.Expr, tail ast.Node) (rel.Expr, error) {
	if tail != nil {
		firstSafe := tail.One("first_safe").One("tail")
		safeCallback := func(tailFunc rel.SafeTailCallback) rel.SafeTailCallback {
			return func(ctx context.Context, v rel.Value, local rel.Scope) (rel.Value, error) {
				val, err := tailFunc(ctx, v, local)
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

		firstTailFn, err := pc.compileTailFunc(ctx, firstSafe)
		if err != nil {
			return nil, err
		}
		exprStates := []rel.SafeTailCallback{safeCallback(firstTailFn)}
		fallback, err := pc.CompileExpr(ctx, tail.One("fall").(ast.Branch))
		if err != nil {
			return nil, err
		}

		for _, o := range tail.Many("ops") {
			if safeTail := o.One("safe"); safeTail != nil {
				safeTailFn, err := pc.compileTailFunc(ctx, safeTail.One("tail"))
				if err != nil {
					return nil, err
				}
				exprStates = append(exprStates, safeCallback(safeTailFn))
			} else if tail := o.One("tail"); tail != nil {
				tailFn, err := pc.compileTailFunc(ctx, tail)
				if err != nil {
					return nil, err
				}
				exprStates = append(exprStates, tailFn)
			} else {
				panic("wat")
			}
		}

		return rel.NewSafeTailExpr(tail.Scanner(), fallback, base, exprStates), nil
	}
	//TODO: panic?
	return base, nil
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

func (pc ParseContext) compileRelation(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	names := parseNames(c.(ast.One).Node.(ast.Branch)["names"].(ast.One).Node.(ast.Branch))
	tuples := c.(ast.One).Node.(ast.Branch)["tuple"].(ast.Many)
	tupleExprs := make([][]rel.Expr, 0, len(tuples))
	for _, tuple := range tuples {
		exprs, err := pc.compileExprs(ctx, tuple.(ast.Branch)["v"].(ast.Many)...)
		if err != nil {
			return nil, err
		}
		tupleExprs = append(tupleExprs, exprs)
	}
	result, err := rel.NewRelationExpr(
		delimsScanner(b),
		names,
		tupleExprs...,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pc ParseContext) compileSet(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := delimsScanner(b)
	if elts := c.(ast.One).Node.(ast.Branch)["elt"]; elts != nil {
		exprs, err := pc.compileExprs(ctx, elts.(ast.Many)...)
		if err != nil {
			return nil, err
		}
		return rel.NewSetExpr(scanner, exprs...)
	}
	return rel.NewLiteralExpr(scanner, rel.None), nil
}

func (pc ParseContext) compileDict(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := delimsScanner(b)
	entryExprs, err := pc.compileDictEntryExprs(ctx, c.(ast.One).Node.(ast.Branch))
	if err != nil {
		return nil, err
	}
	if entryExprs != nil {
		return rel.NewDictExpr(scanner, false, false, entryExprs...)
	}

	d, err := rel.NewDict(false)
	if err != nil {
		return nil, err
	}
	return rel.NewLiteralExpr(scanner, d), nil
}

func (pc ParseContext) compileDictEntryExprs(ctx context.Context, b ast.Branch) ([]rel.DictEntryTupleExpr, error) {
	if pairs := b.Many("pairs"); pairs != nil {
		entryExprs := make([]rel.DictEntryTupleExpr, 0, len(pairs))
		for _, pair := range pairs {
			key := pair.One("key")
			value := pair.One("value")
			keyExpr, err := pc.CompileExpr(ctx, key.(ast.Branch))
			if err != nil {
				return nil, err
			}
			valueExpr, err := pc.CompileExpr(ctx, value.(ast.Branch))
			if err != nil {
				return nil, err
			}
			entryExprs = append(entryExprs, rel.NewDictEntryTupleExpr(pair.Scanner(), keyExpr, valueExpr))
		}
		return entryExprs, nil
	}
	return nil, nil
}

func (pc ParseContext) compileArray(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := delimsScanner(b)
	exprs, err := pc.compileSparseItems(ctx, c)
	if err != nil {
		return nil, err
	}
	if len(exprs) > 0 {
		return rel.NewArrayExpr(scanner, exprs...), nil
	}
	return rel.NewLiteralExpr(scanner, rel.NewArray()), nil
}

func (pc ParseContext) compileBytes(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	if items := c.(ast.One).Node.(ast.Branch)["item"]; items != nil {
		//TODO: support sparse bytes
		exprs, err := pc.compileExprs(ctx, items.(ast.Many)...)
		if err != nil {
			return nil, err
		}
		return rel.NewBytesExpr(delimsScanner(b), exprs...), nil
	}
	return rel.NewBytes([]byte{}), nil
}

func (pc ParseContext) compileExprs(ctx context.Context, exprs ...ast.Node) ([]rel.Expr, error) {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		e, err := pc.CompileExpr(ctx, expr.(ast.Branch))
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (pc ParseContext) compileSparseItems(ctx context.Context, c ast.Children) ([]rel.Expr, error) {
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
		compiled, err := pc.CompileExpr(ctx, expr.(ast.Branch))
		if err != nil {
			return nil, err
		}
		result = append(result, compiled)
	}
	return result, nil
}

// compileCondExprs parses conditons/keys and values expressions for syntax `cond`.
func (pc ParseContext) compileCondExprs(ctx context.Context, exprs ...ast.Node) (_ []rel.Expr, err error) {
	result := make([]rel.Expr, 0, len(exprs))
	for _, expr := range exprs {
		var exprResult rel.Expr

		name, c := which(expr.(ast.Branch), exprTag)
		if c == nil {
			return nil, fmt.Errorf("misshapen node AST: %v", expr.(ast.Branch))
		}

		if name == exprTag {
			switch c := c.(type) {
			case ast.One:
				exprResult, err = pc.CompileExpr(ctx, c.Node.(ast.Branch))
				if err != nil {
					return nil, err
				}
			case ast.Many:
				if len(c) == 1 {
					exprResult, err = pc.CompileExpr(ctx, c[0].(ast.Branch))
					if err != nil {
						return nil, err
					}
				} else {
					var elements []rel.Expr
					for _, e := range c {
						expr, err := pc.CompileExpr(ctx, e.(ast.Branch))
						if err != nil {
							return nil, err
						}
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
	return result, nil
}

func (pc ParseContext) compileFunction(ctx context.Context, b ast.Branch) (rel.Expr, error) {
	p, err := pc.compilePattern(ctx, b)
	if err != nil {
		return nil, err
	}
	expr, err := pc.CompileExpr(ctx, b.One(exprTag).(ast.Branch))
	if err != nil {
		return nil, err
	}
	return rel.NewFunction(b.Scanner(), p, expr), nil
}

func (pc ParseContext) compileMacro(_ context.Context, b ast.Branch) rel.Expr {
	childast := b.One("embed").One("subgrammar").One("ast")
	if value := childast.One("value"); value != nil {
		return value.(MacroValue).SubExpr()
	}
	return rel.ASTNodeToValue(childast)
}

func (pc ParseContext) compilePackage(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	imp := b.One("import").Scanner()
	pkg := c.(ast.One).Node.(ast.Branch)
	if std, has := pkg["std"]; has {
		ident := std.(ast.One).Node.One("IDENT").One("")
		pkgName := ident.(ast.Leaf).Scanner()
		scanner, err := parser.MergeScanners(imp, pkgName)
		if err != nil {
			return nil, err
		}
		return NewPackageExpr(
			pkgName,
			rel.NewDotExpr(scanner, rel.NewIdentExpr(imp, imp.String()), pkgName.String()),
		), nil
	}

	if pkgpath := pkg.One("PKGPATH"); pkgpath != nil {
		scanner := pkgpath.One("").(ast.Leaf).Scanner()
		name := scanner.String()
		if strings.HasPrefix(name, "/") {
			filePath := strings.Trim(name, "/")
			fromRoot := pkg["dot"] == nil
			if pc.SourceDir == "" {
				return nil, fmt.Errorf("local import %q invalid; no local context", name)
			}
			importPath := filepath.Clean(filePath)
			if !fromRoot {
				importPath = filepath.Join(pc.SourceDir, filePath)
			}
			expr, err := importLocalFile(ctx, fromRoot, importPath, pc.SourceDir)
			if err != nil {
				return nil, err
			}
			return NewImportExpr(scanner, expr, importPath), nil
		}
		expr, err := importExternalContent(ctx, name)
		if err != nil {
			return nil, err
		}
		return NewImportExpr(scanner, expr, name), nil
	}
	return nil, fmt.Errorf("compilePackage: malformed package AST %s", pkg)
}

func (pc ParseContext) compileTuple(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	scanner := delimsScanner(b)
	if pairs := c.(ast.One).Node.Many("pairs"); pairs != nil {
		attrs := make([]rel.AttrExpr, 0, len(pairs))
		for _, pair := range pairs {
			var k string
			v, err := pc.CompileExpr(ctx, pair.One("v").(ast.Branch))
			if err != nil {
				return nil, err
			}
			if name := pair.One("name"); name != nil {
				k = parseName(name.(ast.Branch))
			} else {
				switch v := v.(type) {
				case *rel.DotExpr:
					k = v.Attr()
				case rel.IdentExpr:
					k = v.Ident()
				default:
					return nil, fmt.Errorf("unnamed attr expression must be name or end in .name: %T(%[1]v)", v)
				}
			}
			scanner := pair.One("v").(ast.Branch).Scanner()
			if pair.One("rec") != nil {
				fix, fixt := FixFuncs()
				v = rel.NewRecursionExpr(scanner, k, v, fix, fixt)
			}
			attr, err := rel.NewAttrExpr(scanner, k, v)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)
		}
		return rel.NewTupleExpr(scanner, attrs...), nil
	}
	return rel.NewLiteralExpr(scanner, rel.EmptyTuple), nil
}

func delimsScanner(b ast.Branch) parser.Scanner {
	result, err := parser.MergeScanners(b.One("odelim").Scanner(), b.One("cdelim").Scanner())
	if err != nil {
		panic(err)
	}
	return result
}

func (pc ParseContext) compileIdent(_ context.Context, c ast.Children) rel.Expr {
	scanner := c.(ast.One).Node.One("").Scanner()
	var value rel.Value
	switch scanner.String() {
	case "true":
		value = rel.True
	case "false":
		value = rel.False
	default:
		return rel.NewIdentExpr(scanner, scanner.String())
	}
	return rel.NewLiteralExpr(scanner, value)
}

func (pc ParseContext) compileString(_ context.Context, c ast.Children) rel.Expr {
	scanner := c.(ast.One).Node.One("").Scanner()
	return rel.NewLiteralExpr(scanner, rel.NewString([]rune(parseArraiString(scanner.String()))))
}

func (pc ParseContext) compileNumber(_ context.Context, c ast.Children) (rel.Expr, error) {
	scanner := c.(ast.One).Node.One("").Scanner()
	n, err := strconv.ParseFloat(scanner.String(), 64)
	if err != nil {
		//TODO: return custom error instead of golang's ParseFloat error
		return nil, err
	}
	return rel.NewLiteralExpr(scanner, rel.NewNumber(n)), nil
}

func (pc ParseContext) compileChar(_ context.Context, c ast.Children) rel.Expr {
	scanner := c.(ast.One).Node.One("").Scanner()
	char := scanner.String()[1:]
	runes := []rune(parseArraiStringFragment(char, "\"", ""))
	return rel.NewLiteralExpr(scanner, rel.NewNumber(float64(runes[0])))
}

func (pc ParseContext) compileExpr(ctx context.Context, b ast.Branch, c ast.Children) (rel.Expr, error) {
	switch c := c.(type) {
	case ast.One:
		expr, err := pc.CompileExpr(ctx, c.Node.(ast.Branch))
		if err != nil {
			return nil, err
		}
		if b.One("odelim") != nil {
			return rel.NewExprExpr(delimsScanner(b), expr), nil
		}
		return expr, nil
	case ast.Many:
		if len(c) == 1 {
			return pc.CompileExpr(ctx, c[0].(ast.Branch))
		}
		panic("too many expr children")
	}
	return nil, nil
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
		// TODO: Is scanner a suitable argument for rel.NewIdentExpr?
		return f(scanner, rel.NewIdentExpr(scanner, "."), e)
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
	">>": dotUnary(rel.NewSeqArrowExpr(false)),
	// TODO: >>>
	":>": dotUnary(rel.NewTupleMapExpr),
}

type binOpFunc func(scanner parser.Scanner, a, b rel.Expr) rel.Expr

var binops = map[string]binOpFunc{
	"->":      rel.NewArrowExpr,
	"=>":      rel.NewDArrowExpr,
	">>":      rel.NewSeqArrowExpr(false),
	">>>":     rel.NewSeqArrowExpr(true),
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
	"<->":     rel.NewComposeExpr,
	"-&-":     rel.NewJoinCommonExpr,
	"---":     rel.NewJoinExistsExpr,
	"-&>":     rel.NewRightMatchExpr,
	"<&-":     rel.NewLeftMatchExpr,
	"-->":     rel.NewRightResidueExpr,
	"<--":     rel.NewLeftResidueExpr,
	"*":       rel.NewMulExpr,
	"/":       rel.NewDivExpr,
	"%":       rel.NewModExpr,
	"-%":      rel.NewSubModExpr,
	"//":      rel.NewIdivExpr,
	"^":       rel.NewPowExpr,
	"\\":      rel.NewOffsetExpr,
	"+>":      rel.NewAddArrowExpr,
}

var compareOps = map[string]rel.CompareFunc{
	"<:": func(a, b rel.Value) (bool, error) {
		set, is := b.(rel.Set)
		if !is {
			return false, fmt.Errorf("<: rhs not a set: %v", b)
		}
		return set.Has(a), nil
	},
	"!<:": func(a, b rel.Value) (bool, error) {
		set, is := b.(rel.Set)
		if !is {
			return false, fmt.Errorf("!<: rhs not a set: %v", b)
		}
		return !set.Has(a), nil
	},
	"=":  func(a, b rel.Value) (bool, error) { return a.Equal(b), nil },
	"!=": func(a, b rel.Value) (bool, error) { return !a.Equal(b), nil },
	"<":  func(a, b rel.Value) (bool, error) { return a.Less(b), nil },
	">":  func(a, b rel.Value) (bool, error) { return b.Less(a), nil },
	"<=": func(a, b rel.Value) (bool, error) { return !b.Less(a), nil },
	">=": func(a, b rel.Value) (bool, error) { return !a.Less(b), nil },

	"(<)":   func(a, b rel.Value) (bool, error) { return subset(a, b), nil },
	"(>)":   func(a, b rel.Value) (bool, error) { return subset(b, a), nil },
	"(<=)":  func(a, b rel.Value) (bool, error) { return subsetOrEqual(a, b), nil },
	"(>=)":  func(a, b rel.Value) (bool, error) { return subsetOrEqual(b, a), nil },
	"(<>)":  func(a, b rel.Value) (bool, error) { return subsetOrSuperset(a, b), nil },
	"(<>=)": func(a, b rel.Value) (bool, error) { return subsetSupersetOrEqual(b, a), nil },

	"!(<)":   func(a, b rel.Value) (bool, error) { return !subset(a, b), nil },
	"!(>)":   func(a, b rel.Value) (bool, error) { return !subset(b, a), nil },
	"!(<=)":  func(a, b rel.Value) (bool, error) { return !subsetOrEqual(a, b), nil },
	"!(>=)":  func(a, b rel.Value) (bool, error) { return !subsetOrEqual(b, a), nil },
	"!(<>)":  func(a, b rel.Value) (bool, error) { return !subsetOrSuperset(a, b), nil },
	"!(<>=)": func(a, b rel.Value) (bool, error) { return !subsetSupersetOrEqual(b, a), nil },
}
