package syntax

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

var stdScope = rel.EmptyScope.
	With(".", rel.NewTuple(
		rel.NewAttr("math", rel.NewTuple(
			rel.NewAttr("pi", rel.NewNumber(math.Pi)),
			newFloatFuncAttr("sin", math.Sin),
			newFloatFuncAttr("cos", math.Cos),
		)),
		rel.NewAttr("grammar", rel.NewTuple(
			rel.NewNativeFunctionAttr("parse", parseGrammar),
			rel.NewAttr("lang", rel.NewTuple(
				rel.NewAttr("arrai", rel.ASTNodeToValue(ast.ParserNodeToNode(
					wbnf.Core().Grammar(), *arraiParsers.Node()))),
				rel.NewAttr("wbnf", rel.ASTNodeToValue(ast.CoreNode())),
			)),
		)),
	)).
	With("//./", rel.NewNativeFunction("//./", importLocalFile)).
	With("//", rel.NewNativeFunction("//", importURL))

func newFloatFuncAttr(name string, f func(float64) float64) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) rel.Value {
		return rel.NewNumber(f(value.(rel.Number).Float64()))
	})
}

func parseGrammar(v rel.Value) rel.Value {
	astNode := rel.ASTNodeFromValue(v).(ast.Branch)
	parserNode := ast.NodeToParserNode(wbnf.Core().Grammar(), astNode).(parser.Node)
	parsers := wbnf.NewFromNode(parserNode).Compile(&parserNode)
	return rel.NewNativeFunction("parse(<grammar>)", func(v rel.Value) rel.Value {
		rule := v.String()
		return rel.NewNativeFunction(fmt.Sprintf("parse(%s)", rule), func(v rel.Value) rel.Value {
			node, err := parsers.Parse(wbnf.Rule(rule), parser.NewScanner(v.String()))
			if err != nil {
				panic(err)
			}
			return rel.ASTNodeToValue(ast.ParserNodeToNode(wbnf.Core().Grammar(), node))
		})
	})
}

func importLocalFile(v rel.Value) rel.Value {
	data, err := ioutil.ReadFile(v.String())
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(string(data)))
}

func importURL(v rel.Value) rel.Value {
	url := v.String()
	if !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(string(data)))
}
