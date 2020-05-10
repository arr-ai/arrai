package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// LHSExpr represents any Expr that has a LHS component.
type LHSExpr interface {
	LHS() Expr
}

// GetStringValue returns the string value for expr or false if not a string.
func GetStringValue(expr Expr) (string, bool) {
	if set, ok := expr.(Set); ok {
		s := set.String()
		if s[:1] == `"` {
			// TODO: Fix this dirty hack. Maybe enhance Set.Export.
			return s[1 : len(s)-1], true
		}
	}
	return "", false
}

type ExprScanner struct {
	Src parser.Scanner
}

// Scanner returns the scanner
func (e ExprScanner) Scanner() parser.Scanner {
	return e.Src
}

func wrapContext(err error, expr Expr) error {
	return fmt.Errorf("%s\n%s", err.Error(), expr.Scanner().Context(parser.DefaultLimit))
}

func EvalExpr(expr Expr, local Scope) (_ Value, err error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case error:
			err = wrapContext(r, expr)
		default:
			err = wrapContext(fmt.Errorf("unexpected panic: %v", r), expr)
		}
	}()
	return expr.Eval(local)
}
