package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// RecursionExpr evaluates `let rec name = fn; ...` by tying the knot directly:
// name is bound to a recCell in the scope that fn is evaluated in, and the
// cell is filled with fn's value once it exists. Closures created while
// evaluating fn capture that scope, so recursive references resolve through
// the cell with no combinator machinery.
type RecursionExpr struct {
	ExprScanner
	name string
	fn   Expr
}

// NewRecursionExpr returns a new RecursionExpr.
func NewRecursionExpr(scanner parser.Scanner, name string, fn Expr) Expr {
	return RecursionExpr{ExprScanner{scanner}, name, fn}
}

// recCell is the scope entry for a recursive binding. Scope maps names to
// Exprs, so the cell is an Expr whose evaluation yields the bound value. It
// is written exactly once, after fn has been evaluated; every copy of the
// scope shares the pointer, so the surrounding immutable structures are never
// mutated. The cell is deliberately opaque to String, Hash and friends:
// printing it as its name is what prevents self-referential closures from
// recursing forever when they are formatted, hashed or compared.
type recCell struct {
	name string
	val  Value
	set  bool
}

func (c *recCell) Eval(context.Context, Scope) (Value, error) {
	if !c.set {
		return nil, errors.Errorf("recursive binding %q used before it is defined", c.name)
	}
	return c.val, nil
}

func (c *recCell) Source() parser.Scanner {
	return *parser.NewScanner(c.name)
}

func (c *recCell) String() string {
	return c.name
}

// Eval evaluates the recursive binding and returns its value.
func (r RecursionExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	cell := &recCell{name: r.name}
	val, err := r.fn.Eval(ctx, local.With(r.name, cell))
	if err != nil {
		return nil, WrapContextErr(err, r, local)
	}
	switch v := val.(type) {
	case Closure:
	case Tuple:
		for e := v.Enumerator(); e.MoveNext(); {
			_, attr := e.Current()
			if _, isFunction := attr.(Closure); !isFunction {
				return nil, WrapContextErr(
					errors.Errorf("Recursion requires a tuple of functions: %v", v.String()), r, local)
			}
		}
	default:
		return nil, WrapContextErr(errors.Errorf("Recursion does not support %s", ValueTypeAsString(val)), r, local)
	}
	cell.val = val
	cell.set = true
	return val, nil
}

// String returns a string representation of the expression.
func (r RecursionExpr) String() string {
	return fmt.Sprintf("\\%s %s", r.name, r.fn)
}
