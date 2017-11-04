package rel

import (
	"log"
	"strings"
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

var depth = 0

func enter(format string, args ...interface{}) struct{} {
	ilog("ENTER: "+format, args...)
	depth++
	if depth > 10 {
		panic("Depth limit reached")
	}
	return struct{}{}
}

func ilog(format string, args ...interface{}) {
	log.Printf(strings.Repeat("  ", depth)+format, args...)
}

func exit(_ struct{}, args ...interface{}) {
	depth--
	ilog("EXIT: ", args...)
}

// // ConstantFold returns the Value for the given Expr if it doesn't depend on any
// // scope variables. Otherwise, returns the Expr.
// func ConstantFold(e Expr) (Expr, error) {
// 	if _, ok := e.(*DynExpr); ok {
// 		return e, nil
// 	}
// 	value, err := e.Eval(EmptyScope)
// 	if err == nil {
// 		return value, nil
// 	}
// 	if _, ok := err.(*IdentLookupFailed); ok {
// 		return e, nil
// 	}
// 	return nil, err
// }
