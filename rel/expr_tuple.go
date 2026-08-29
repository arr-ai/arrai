package rel

import (
	"bytes"
	"context"
	"sort"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// AttrExpr represents a single name:expr in a TupleExpr.
type AttrExpr struct {
	ExprScanner
	name string
	expr Expr
}

// NewAttrExpr constructs a new AttrExpr from the given arguments.
func NewAttrExpr(scanner parser.Scanner, name string, expr Expr) (AttrExpr, error) {
	isWildcard := false
	if dot, ok := expr.(*DotExpr); ok {
		if dot.Attr() == "*" {
			isWildcard = true
			expr = dot.Subject()
		}
	}
	if isWildcard != (name == "*") {
		return AttrExpr{}, errors.Errorf("Wildcard attr cannot have a name")
	}
	return AttrExpr{ExprScanner{scanner}, name, expr}, nil
}

// NewWildcardExpr constructs a new wildcard AttrExpr.
func NewWildcardExpr(scanner parser.Scanner, lhs Expr) AttrExpr {
	return AttrExpr{ExprScanner{scanner}, "*", lhs}
}

// IsWildcard returns true iff the AttrExpr is a wildcard expression.
func (e *AttrExpr) IsWildcard() bool {
	return e.name == "*"
}

// Apply applies the AttrExpr to the Tuple.
func (e *AttrExpr) Apply(
	ctx context.Context, local Scope, tuple Tuple,
) (Tuple, error) {
	value, err := e.expr.Eval(ctx, local)
	if err != nil {
		return nil, err
	}
	if e.IsWildcard() {
		if t, ok := value.(Tuple); ok {
			for e := t.Enumerator(); e.MoveNext(); {
				tuple = tuple.With(e.Current())
			}
			return tuple, nil
		}
		return nil, errors.Errorf(
			"LHS of wildcard must be tuple, not %s", ValueTypeAsString(value))
	}
	tuple = tuple.With(e.name, value)
	return tuple, nil
}

// TupleExpr returns a set from a slice of Exprs.
type TupleExpr struct {
	ExprScanner
	attrs   []AttrExpr
	attrMap map[string]Expr

	// shape and slots are set when the literal's attribute names are static
	// (no wildcards, no view attributes, no repeats): the tuple can then be
	// built by storing each value into its slot.
	shape *Shape
	slots []int
}

// staticShape precomputes the shape of a literal with static attribute names.
func (e *TupleExpr) staticShape() {
	names := make([]string, 0, len(e.attrs))
	for _, attr := range e.attrs {
		if attr.IsWildcard() || strings.HasPrefix(attr.name, "&") {
			return
		}
		names = append(names, attr.name)
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1] {
			return
		}
	}
	shape := shapeOf(sorted)
	// The "@"-indexed two-attribute shapes canonicalise to specialised
	// kinds; leave those to the general path.
	if len(sorted) == 2 && sorted[0] == "@" && strings.HasPrefix(sorted[1], "@") {
		return
	}
	e.slots = make([]int, len(names))
	for i, name := range names {
		e.slots[i], _ = shape.Index(name)
	}
	e.shape = shape
}

// NewTupleExpr returns a new TupleExpr.
func NewTupleExpr(scanner parser.Scanner, attrs ...AttrExpr) Expr {
	attrValues := make([]Attr, len(attrs))
	for i, attr := range attrs {
		if value, is := exprIsValue(attr.expr); is {
			attrValues[i] = Attr{attr.name, value}
		} else {
			attrMap := make(map[string]Expr, len(attrs))
			for _, attr := range attrs {
				attrMap[attr.name] = attr.expr
			}
			e := &TupleExpr{ExprScanner: ExprScanner{scanner}, attrs: attrs, attrMap: attrMap}
			e.staticShape()
			return e
		}
	}
	return NewLiteralExpr(scanner, NewTuple(attrValues...))
}

// NewTupleExprFromMap returns a new TupleExpr from a map[string]Expr.
func NewTupleExprFromMap(scanner parser.Scanner, attrMap map[string]Expr) Expr {
	attrValues := make([]Attr, len(attrMap))
	i := 0
	for name, expr := range attrMap {
		if value, ok := expr.(Value); ok {
			attrValues[i] = Attr{name, value}
			i++
		} else {
			attrs := make([]AttrExpr, len(attrMap))
			i := 0
			for name, expr := range attrMap {
				attrs[i] = AttrExpr{ExprScanner{scanner}, name, expr}
				i++
			}
			e := &TupleExpr{ExprScanner: ExprScanner{scanner}, attrs: attrs, attrMap: attrMap}
			e.staticShape()
			return e
		}
	}
	return NewTuple(attrValues...)
}

// String returns a string representation of the expression.
func (e *TupleExpr) String() string { //nolint:dupl
	var b bytes.Buffer
	b.WriteByte('(')
	for i, attr := range e.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		if attr.IsWildcard() {
			if ident, is := attr.expr.(IdentExpr); !is || ident.Ident() != "." {
				b.WriteString(attr.expr.String())
			}
			b.WriteString(".*")
		} else {
			b.WriteString(attr.name)
			b.WriteString(": ")
			b.WriteString(attr.expr.String())
		}
	}
	b.WriteByte(')')
	return b.String()
}

// Eval returns the subject
func (e *TupleExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	if e.shape != nil {
		vals := make([]Value, len(e.slots))
		for i, attr := range e.attrs {
			value, err := attr.expr.Eval(ctx, local)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			vals[e.slots[i]] = value
		}
		return newShapedTuple(e.shape, vals), nil
	}
	tuple := EmptyTuple
	var err error
	for _, attr := range e.attrs {
		tuple, err = attr.Apply(ctx, local, tuple)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
	}
	// TODO: Construct new tuple directly
	return tuple.(*GenericTuple).Canonical(), nil
}
