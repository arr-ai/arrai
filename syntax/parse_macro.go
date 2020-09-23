package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

// Macro represents the metadata of a macro invocation: the grammar and rule to parse with, and the
// transform to apply to the parsed body.
type Macro struct {
	ruleName  string
	grammar   rel.Tuple
	transform rel.Set
}

// unpackMacro processes the current parser scope when a macro invocation is detected. It extracts
// the details of the macro to invoke, and returns them as a Macro.
func (pc ParseContext) unpackMacro(
	ctx context.Context,
	macroElt parser.Node,
	ruleElt parser.Node,
	relScope rel.Scope,
) (Macro, error) {
	macroNode := ast.FromParserNode(arraiParsers.Grammar(), macroElt)
	macroExpr, err := pc.CompileExpr(ctx, macroNode)
	if err != nil {
		return Macro{}, err
	}
	macroValue, err := macroExpr.Eval(ctx, relScope)
	if err != nil {
		return Macro{ruleName: ""}, err
	}
	macroTuple := macroValue.(rel.Tuple)

	grammar := macroTuple
	if macroTuple.HasName("@grammar") {
		grammar = macroTuple.MustGet("@grammar").(rel.Tuple)
	}

	var ruleName string
	if ruleElt.Count() > 0 {
		ruleNode := ast.FromParserNode(arraiParsers.Grammar(), ruleElt.Get(0))
		ruleName = ruleNode.One("name").Scanner().String()
	} else if ruleName, err = getFirstRuleName(grammar); err != nil {
		return Macro{ruleName: ruleName}, err
	}

	transform := rel.None
	if transforms, ok := macroTuple.Get("@transform"); ok {
		// If @transform is present but there is no named transform for ruleName, fail
		// loudly rather than falling back on the default rule or nothing. A macro's
		// rule transforms should be as well-specified as the grammar.
		if transformValue, ok := transforms.(rel.Tuple).Get(ruleName); ok {
			transform = transformValue.(rel.Set)
		} else {
			return Macro{ruleName: ruleName}, fmt.Errorf("transform for rule %q not found", ruleName)
		}
	}
	return Macro{ruleName, grammar, transform}, nil
}

// getFirstRuleName finds the first rule declared in grammar and returns its name.
func getFirstRuleName(grammar rel.Tuple) (string, error) {
	stmts := grammar.MustGet("stmt").(rel.Array).Values()
	for _, stmt := range stmts {
		if prod, ok := stmt.(rel.Tuple).Get("prod"); ok {
			return prod.(rel.Tuple).MustGet("IDENT").(rel.Tuple).MustGet("").String(), nil
		}
	}
	return "", fmt.Errorf("no prod rule found in grammar")
}

// MacroValue is an Extra node with an Expr value and a Scanner for the macro source.
type MacroValue struct {
	ast.Extra
	expr    rel.Expr
	scanner parser.Scanner
}

// NewMacroValue returns a MacroValue with a given Expr and Scanner.
func NewMacroValue(expr rel.Expr, scanner parser.Scanner) MacroValue {
	return MacroValue{expr: expr, scanner: scanner}
}

// Scanner returns a scanner of the source that was replaced by the macro.
func (m MacroValue) Scanner() parser.Scanner {
	return m.scanner
}

// SubExpr returns the Expr resulting from evaluating the macro.
func (m MacroValue) SubExpr() rel.Expr {
	return m.expr
}
