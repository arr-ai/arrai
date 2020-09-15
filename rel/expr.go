package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// ExprScanner represents the source location of an Expr.
type ExprScanner struct {
	Src parser.Scanner
}

// Source returns a scanner locating the expression's source code.
func (e ExprScanner) Source() parser.Scanner {
	return e.Src
}

// ContextErr represents the whole stack frame of an error from arrai script.
type ContextErr struct {
	err    error
	source parser.Scanner
	scope  Scope
}

// NewContextErr creates a ContextErr.
func NewContextErr(err error, source parser.Scanner, scope Scope) ContextErr {
	return ContextErr{err, source, scope}
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

// GetImportantFrames returns an array of important frames whose last element
// is the last frame near the point of failure.
// Important frames are frames that don't contain the frame under it.
func (c ContextErr) GetImportantFrames() []ContextErr {
	if cerr, is := c.err.(ContextErr); is {
		currScope := cerr.GetImportantFrames()
		if c.source.Contains(cerr.source) {
			return currScope
		}
		return append([]ContextErr{c}, currScope...)
	}
	return []ContextErr{c}
}

// GetScope gets the scope of a ContextErr.
func (c ContextErr) GetScope() Scope {
	return c.scope
}

// GetSource gets the source of a ContextErr.
func (c ContextErr) GetSource() parser.Scanner {
	return c.source
}

// WrapContextErr wraps an error in a ContextErr
func WrapContextErr(err error, expr Expr, scope Scope) error {
	return ContextErr{err, expr.Source(), scope}
}
