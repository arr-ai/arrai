package syntax

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

var (
	stdScopeOnce, fixOnce sync.Once
	stdScopeVar           rel.Scope
	fix, fixt             rel.Value
)

func FixFuncs() (rel.Value, rel.Value) {
	fixOnce.Do(func() {
		fix = mustParseLit(`(\f f(f))(\f \g \n g(f(f)(g))(n))`)
		fixt = mustParseLit(`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))`)
	})
	return fix, fixt
}

func StdScope() rel.Scope {
	stdScopeOnce.Do(func() {
		fixFn, fixtFn := FixFuncs()
		stdScopeVar = rel.EmptyScope.
			With("//", rel.NewTuple(
				rel.NewNativeFunctionAttr("dict", func(value rel.Value) (rel.Value, error) {
					if t, ok := value.(rel.Tuple); ok {
						if !t.IsTrue() {
							return rel.None, nil
						}
						entries := make([]rel.DictEntryTuple, 0, t.Count())
						for e := t.Enumerator(); e.MoveNext(); {
							name, value := e.Current()
							entries = append(entries, rel.NewDictEntryTuple(rel.NewString([]rune(name)), value))
						}
						return rel.NewDict(false, entries...), nil
					}
					return nil, fmt.Errorf("dict(): not a tuple")
				}),
				rel.NewNativeFunctionAttr("tuple", func(value rel.Value) (rel.Value, error) {
					switch t := value.(type) {
					case rel.Dict:
						attrs := make([]rel.Attr, 0, t.Count())
						for e := t.DictEnumerator(); e.MoveNext(); {
							key, value := e.Current()
							attrs = append(attrs, rel.NewAttr(key.(rel.String).String(), value))
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
					rel.NewNativeFunctionAttr("print", func(value rel.Value) (rel.Value, error) {
						log.Print(value)
						return value, nil
					}),
					createNestedFuncAttr("printf", 2, func(args ...rel.Value) (rel.Value, error) {
						format := args[0].(rel.String).String()
						strs := make([]interface{}, 0, args[1].(rel.Set).Count())
						for i, ok := args[1].(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
							strs = append(strs, i.Current())
						}
						log.Printf(format, strs...)
						return args[1], nil
					}),
				),
				stdArchive(),
				stdEncoding(),
				stdEval(),
				stdOs(),
				stdNet(),
				stdRe(),
				stdReflect(),
				stdRel(),
				stdSeq(),
				stdStr(),
				stdTest(),
				stdBits(),
				stdFmt(),
			))
	})
	return stdScopeVar
}

func createNestedFunc(
	name string, nArgs int, f func(...rel.Value) (rel.Value, error), args ...rel.Value,
) (rel.Value, error) {
	if nArgs == 0 {
		return f(args...)
	}

	return rel.NewNativeFunction(name+strconv.Itoa(nArgs), func(parent rel.Value) (rel.Value, error) {
		return createNestedFunc(name, nArgs-1, f, append(args, parent)...)
	}), nil
}

func mustCreateNestedFunc(
	name string, nArgs int, f func(...rel.Value) (rel.Value, error), args ...rel.Value,
) rel.Value {
	g, err := createNestedFunc(name, nArgs, f, args...)
	if err != nil {
		panic(err)
	}
	return g
}

func createNestedFuncAttr(name string, nArgs int, f func(...rel.Value) (rel.Value, error)) rel.Attr {
	g, err := createNestedFunc(name, nArgs, f)
	if err != nil {
		panic(err)
	}
	return rel.NewAttr(name, g)
}

func createFunc2(name string, f func(a, b rel.Value) (rel.Value, error)) rel.Value {
	return rel.NewNativeFunction(name, func(a rel.Value) (rel.Value, error) {
		return rel.NewNativeFunction(name+"_2", func(b rel.Value) (rel.Value, error) {
			return f(a, b)
		}), nil
	})
}

func createFunc2Attr(name string, f func(a, b rel.Value) (rel.Value, error)) rel.Attr {
	return rel.NewAttr(name, createFunc2(name, f))
}

func createFunc3(name string, f func(a, b, c rel.Value) (rel.Value, error)) rel.Value {
	return rel.NewNativeFunction(name, func(a rel.Value) (rel.Value, error) {
		return rel.NewNativeFunction(name+"$2", func(b rel.Value) (rel.Value, error) {
			return rel.NewNativeFunction(name+"$3", func(c rel.Value) (rel.Value, error) {
				return f(a, b, c)
			}), nil
		}), nil
	})
}

func createFunc3Attr(name string, f func(a, b, c rel.Value) (rel.Value, error)) rel.Attr {
	return rel.NewAttr(name, createFunc3(name, f))
}

func mustParseLit(s string) rel.Value {
	lit, err := MustCompile(NoPath, s).Eval(rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	return lit
}

func newFloatFuncAttr(name string, f func(float64) float64) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) (rel.Value, error) {
		return rel.NewNumber(f(value.(rel.Number).Float64())), nil
	})
}

func parseGrammar(v rel.Value) (rel.Value, error) {
	astNode := rel.ASTNodeFromValue(v).(ast.Branch)
	g := wbnf.NewFromAst(astNode)
	parsers := g.Compile(astNode)
	return rel.NewNativeFunction("parse(<grammar>)", func(v rel.Value) (rel.Value, error) {
		rule := v.String()
		return rel.NewNativeFunction(fmt.Sprintf("parse(%s)", rule), func(v rel.Value) (rel.Value, error) {
			node, err := parsers.Parse(parser.Rule(rule), parser.NewScanner(v.String()))
			if err != nil {
				return nil, err
			}
			return rel.ASTNodeToValue(ast.FromParserNode(parsers.Grammar(), node)), nil
		}), nil
	}), nil
}
