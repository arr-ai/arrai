package syntax

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"

	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

var (
	stdSafeScopeOnce, stdUnsafeScopeOnce, fixOnce sync.Once
	stdSafeScopeVar, stdUnsafeScopeVar            rel.Scope
	fix, fixt                                     rel.Value
)

func FixFuncs() (rel.Value, rel.Value) {
	fixOnce.Do(func() {
		fix = mustParseLit(`(\f f(f))(\f \g \n g(f(f)(g))(n))`)
		fixt = mustParseLit(`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))`)
	})
	return fix, fixt
}

func StdScope() rel.Scope {
	stdUnsafeScopeOnce.Do(func() {
		goStdlib := rel.MergeTuples(SafeStdScopeTuple(), rel.NewTuple(
			stdOsUnsafe(),
			stdNet(),
		))
		arraiUnsafeStdlib := mustParseBundle(stdlibUnsafeArraiz())
		stdlibVal, err := rel.NewCallExprCurry(*parser.NewScanner("stdlib"), arraiUnsafeStdlib, goStdlib).
			Eval(context.Background(), rel.EmptyScope)
		if err != nil {
			panic(err)
		}
		t, isTuple := stdlibVal.(rel.Tuple)
		if !isTuple {
			panic("standard library does not evaluate to a tuple")
		}
		stdUnsafeScopeVar = rel.EmptyScope.With("//", t)
	})
	return stdUnsafeScopeVar
}
func SafeStdScope() rel.Scope {
	stdSafeScopeOnce.Do(func() {
		safeStdlib := SafeStdScopeTuple()
		stdSafeScopeVar = rel.EmptyScope.With("//", safeStdlib)
	})
	return stdSafeScopeVar
}

func SafeStdScopeTuple() rel.Tuple {
	fixFn, fixtFn := FixFuncs()
	goStdlib := rel.NewTuple(
		rel.NewNativeFunctionAttr("dict", func(_ context.Context, value rel.Value) (rel.Value, error) {
			if t, ok := value.(rel.Tuple); ok {
				if !t.IsTrue() {
					return rel.None, nil
				}
				entries := make([]rel.DictEntryTuple, 0, t.Count())
				for e := t.Enumerator(); e.MoveNext(); {
					name, value := e.Current()
					entries = append(entries, rel.NewDictEntryTuple(rel.NewString([]rune(name)), value))
				}
				return rel.NewDict(false, entries...)
			}
			return nil, fmt.Errorf("dict(): not a tuple")
		}),
		rel.NewNativeFunctionAttr("tuple", func(_ context.Context, value rel.Value) (rel.Value, error) {
			switch t := value.(type) {
			case rel.Dict:
				attrs := make([]rel.Attr, 0, t.Count())
				for e := t.DictEnumerator(); e.MoveNext(); {
					keyVal, value := e.Current()
					var key string
					// keyVal won't be a rel.String if it's empty.
					if _, ok := keyVal.(rel.String); ok {
						key = keyVal.String()
					} else if _, ok := keyVal.(rel.EmptySet); ok {
						key = ""
					} else {
						return nil, fmt.Errorf(
							"all keys of arg to //tuple must be strings, not %s", rel.ValueTypeAsString(keyVal))
					}
					attrs = append(attrs, rel.NewAttr(key, value))
				}
				return rel.NewTuple(attrs...), nil
			case rel.Set:
				if !t.IsTrue() {
					return rel.NewTuple(), nil
				}
			}
			return nil, fmt.Errorf("tuple(): not a dict")
		}),
		rel.NewTupleAttr("math",
			rel.NewFloatAttr("pi", math.Pi),
			rel.NewFloatAttr("e", math.E),
			newFloatFuncAttr("sin", math.Sin),
			newFloatFuncAttr("cos", math.Cos),
		),
		rel.NewTupleAttr("grammar",
			rel.NewNativeFunctionAttr("parse", parseGrammar),
			rel.NewTupleAttr("lang",
				rel.NewAttr("arrai", rel.ASTNodeToValue(arraiParsers.Node().(ast.Node))),
				rel.NewAttr("wbnf", rel.ASTNodeToValue(wbnf.Core().Node().(ast.Node))),
			),
		),
		rel.NewTupleAttr("fn",
			rel.NewAttr("fix", fixFn),
			rel.NewAttr("fixt", fixtFn),
		),
		rel.NewTupleAttr("log",
			rel.NewNativeFunctionAttr("print", func(_ context.Context, value rel.Value) (rel.Value, error) {
				log.Print(rel.Repr(value))
				return value, nil
			}),
			createFunc2Attr("printf", func(ctx context.Context, a, b rel.Value) (rel.Value, error) {
				format := a.(rel.String).String()
				strs := make([]interface{}, 0, b.(rel.Set).Count())
				for i := b.(rel.Set).ArrayEnumerator(); i.MoveNext(); {
					strs = append(strs, i.Current())
				}
				log.Printf(format, strs...)
				return b, nil
			}),
		),
		rel.NewNativeFunctionAttr("error", func(_ context.Context, value rel.Value) (rel.Value, error) {
			// FIXME: this is a temporary error handling
			return nil, errors.New(value.String())
		}),
		rel.NewTupleAttr(
			"@internal", rel.NewTupleAttr(
				"xml", createFunc2Attr(
					"decode", func(_ context.Context, xmlConfig, value rel.Value) (rel.Value, error) {
						config, err := parseXMLConfig(xmlConfig)
						if err != nil {
							return nil, err
						}
						return decodeXML(value, *config)
					}),
				rel.NewNativeFunctionAttr("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
					return translate.BytesXMLFromArrai(value)
				})), rel.NewTupleAttr(
				"eval", createFunc2Attr(
					"eval", func(ctx context.Context, evalConfig, value rel.Value) (rel.Value, error) {
						config, err := parseEvalConfig(evalConfig)
						if err != nil {
							return nil, err
						}
						return contextualEval(ctx, *config, value)
					}))),
		stdArchive(),
		stdEncoding(),
		stdEval(),
		stdOsSafe(),
		stdRe(),
		stdReflect(),
		stdRel(),
		stdSeq(),
		stdStr(),
		stdTest(),
		stdBits(),
		stdFmt(),
		stdRuntime(),
		stdDeprecated(),
	)
	arraiSafeStdlib := mustParseBundle(stdlibSafeArraiz())
	stdlibVal, err := rel.NewCallExprCurry(*parser.NewScanner("stdlib"), arraiSafeStdlib, goStdlib).
		Eval(context.Background(), rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	t, isTuple := stdlibVal.(rel.Tuple)
	t = rel.MergeTuples(t, rel.NewTuple(rel.NewTupleAttr("std", rel.NewAttr("safe", t))))
	if !isTuple {
		panic("standard library does not evaluate to a tuple")
	}
	return t
}

type fnArgs struct {
	args []rel.Value
	ctx  context.Context
}

func createNestedFunc(
	name string, nArgs int,
	f func(context.Context, ...rel.Value) (rel.Value, error), args fnArgs) (rel.Value, error) {
	if nArgs == 0 {
		return f(args.ctx, args.args...)
	}
	return rel.NewNativeFunction(
		name+strconv.Itoa(nArgs),
		func(ctx context.Context, parent rel.Value) (rel.Value, error) {
			return createNestedFunc(name, nArgs-1, f, fnArgs{args: append(args.args, parent), ctx: ctx})
		}), nil
}

func mustCreateNestedFunc(
	name string, nArgs int, f func(context.Context, ...rel.Value) (rel.Value, error),
) rel.Value {
	if nArgs == 0 {
		panic(errors.New("mustCreateNestedFunc: function cannot have 0 arguments"))
	}
	g, err := createNestedFunc(name, nArgs, f, fnArgs{args: make([]rel.Value, 0)})
	if err != nil {
		panic(err)
	}
	return g
}

func createNestedFuncAttr(name string, nArgs int, f func(context.Context, ...rel.Value) (rel.Value, error)) rel.Attr {
	if nArgs == 0 {
		panic(errors.New("createNestedFuncAttr: function cannot have 0 arguments"))
	}
	g, err := createNestedFunc(name, nArgs, f, fnArgs{args: make([]rel.Value, 0)})
	if err != nil {
		panic(err)
	}
	return rel.NewAttr(name, g)
}

func createFunc2(name string, f func(ctx context.Context, a, b rel.Value) (rel.Value, error)) rel.Value {
	return rel.NewNativeFunction(name, func(_ context.Context, a rel.Value) (rel.Value, error) {
		return rel.NewNativeFunction(name+"$2", func(ctx context.Context, b rel.Value) (rel.Value, error) {
			return f(ctx, a, b)
		}), nil
	})
}

func createFunc2Attr(name string, f func(ctx context.Context, a, b rel.Value) (rel.Value, error)) rel.Attr {
	return rel.NewAttr(name, createFunc2(name, f))
}

func createFunc3(name string, f func(ctx context.Context, a, b, c rel.Value) (rel.Value, error)) rel.Value {
	return rel.NewNativeFunction(name, func(_ context.Context, a rel.Value) (rel.Value, error) {
		return rel.NewNativeFunction(name+"$2", func(_ context.Context, b rel.Value) (rel.Value, error) {
			return rel.NewNativeFunction(name+"$3", func(ctx context.Context, c rel.Value) (rel.Value, error) {
				return f(ctx, a, b, c)
			}), nil
		}), nil
	})
}

func createFunc3Attr(name string, f func(ctx context.Context, a, b, c rel.Value) (rel.Value, error)) rel.Attr {
	return rel.NewAttr(name, createFunc3(name, f))
}

func mustParseExpr(s string) rel.Expr {
	return MustCompile(context.TODO(), NoPath, s)
}

func mustParseLit(s string) rel.Value {
	// this shouldn't require any special key-value pairs in the context
	ctx := context.TODO()
	lit, err := mustParseExpr(s).Eval(ctx, rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	return lit
}

func mustParseBundle(b []byte) rel.Value {
	ctx := importcache.WithNewImportCache(context.TODO())
	ctx, err := WithBundleRun(ctx, b)
	if err != nil {
		panic(err)
	}
	ctx, mainFileSource, path := GetMainBundleSource(ctx)
	lit, err := MustCompile(ctx, path, string(mainFileSource)).Eval(ctx, rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	return lit
}

func newFloatFuncAttr(name string, f func(float64) float64) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(_ context.Context, value rel.Value) (rel.Value, error) {
		return rel.NewNumber(f(value.(rel.Number).Float64())), nil
	})
}

func parseGrammar(_ context.Context, v rel.Value) (rel.Value, error) {
	astNode := rel.ASTNodeFromValue(v).(ast.Branch)
	g := wbnf.NewFromAst(astNode)
	parsers := g.Compile(astNode)
	return rel.NewNativeFunction("parse(<grammar>)", func(_ context.Context, v rel.Value) (rel.Value, error) {
		rule := v.String()
		return rel.NewNativeFunction(
			fmt.Sprintf("parse(%s)", rule),
			func(_ context.Context, v rel.Value) (rel.Value, error) {
				node, err := parsers.Parse(parser.Rule(rule), parser.NewScanner(v.String()))
				if err != nil {
					return nil, err
				}
				return rel.ASTNodeToValue(ast.FromParserNode(parsers.Grammar(), node)), nil
			}), nil
	}), nil
}
