package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

type ExprScanner struct {
	Src parser.Scanner
}

// Source returns a scanner locating the expression's source code.
func (e ExprScanner) Source() parser.Scanner {
	return e.Src
}

type ContextErr struct {
	err    error
	source parser.Scanner
	scope  Scope
}

func (c ContextErr) Error() string {
	if cerr, is := c.err.(ContextErr); is {
		errString := cerr.Error()
		if c.source.Contains(cerr.source) {
			return errString
		}
		return fmt.Sprintf("%s\n%s", errString, c.source.Context(parser.DefaultLimit))
	}
	return fmt.Sprintf("%s\n%s", c.err.Error(), c.source.Context(parser.DefaultLimit))
}

// NextErr returns the error contained in ContextErr
func (c ContextErr) NextErr() error {
	return c.err
}

// GetLastScope gets the scope nearest to the error
func (c ContextErr) GetLastScope() Scope {
	ctxErr := c
	for {
		var isContextErr bool
		var currentErr ContextErr
		if currentErr, isContextErr = ctxErr.err.(ContextErr); !isContextErr {
			return ctxErr.scope
		}
		ctxErr = currentErr
	}
}

func wrapContext(err error, expr Expr, scope Scope) error {
	return ContextErr{err, expr.Source(), scope}
}

func EvalExpr(expr Expr, local Scope) (_ Value, err error) {
	//TODO: this is only the initial scope, how to get the last scope?
	defer wrapPanic(expr, &err, local)
	return expr.Eval(local)
}

func wrapPanic(expr Expr, err *error, scope Scope) {
	switch r := recover().(type) {
	case nil:
	case error:
		*err = wrapContext(r, expr, scope)
	default:
		*err = wrapContext(fmt.Errorf("unexpected panic: %v", r), expr, scope)
	}
}
