package rel

import (
	"bytes"

	"github.com/go-errors/errors"
)

// AttrExpr represents a single name:expr in a TupleExpr.
type AttrExpr struct {
	name string
	expr Expr
}

// NewAttrExpr constructs a new AttrExpr from the given arguments.
func NewAttrExpr(name string, expr Expr) (AttrExpr, error) {
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
	return AttrExpr{name, expr}, nil
}

func MustNewAttrExpr(name string, expr Expr) AttrExpr {
	attrExpr, err := NewAttrExpr(name, expr)
	if err != nil {
		panic(err)
	}
	return attrExpr
}

// NewWildcardExpr constructs a new wildcard AttrExpr.
func NewWildcardExpr(lhs Expr) AttrExpr {
	return AttrExpr{"*", lhs}
}

// IsWildcard returns true iff the AttrExpr is a wildcard expression.
func (e *AttrExpr) IsWildcard() bool {
	return e.name == "*"
}

// Name returns the AttrExpr's name.
func (e *AttrExpr) Name() string {
	return e.name
}

// Expr returns the AttrExpr's expr.
func (e *AttrExpr) Expr() Expr {
	return e.expr
}

// Apply applies the AttrExpr to the Tuple.
func (e *AttrExpr) Apply(
	local Scope, tuple Tuple,
) (Tuple, error) {
	value, err := e.expr.Eval(local)
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
			"LHS of wildcard must be tuple, not %T", value)
	}
	tuple = tuple.With(e.name, value)
	return tuple, nil
}

// TupleExpr returns a set from a slice of Exprs.
type TupleExpr struct {
	attrs   []AttrExpr
	attrMap map[string]Expr
}

// NewTupleExpr returns a new TupleExpr.
func NewTupleExpr(attrs ...AttrExpr) Expr {
	attrValues := make([]Attr, len(attrs))
	for i, attr := range attrs {
		if value, ok := attr.expr.(Value); ok {
			attrValues[i] = Attr{attr.name, value}
		} else {
			attrMap := make(map[string]Expr, len(attrs))
			for _, attr := range attrs {
				attrMap[attr.name] = attr.expr
			}
			return &TupleExpr{attrs, attrMap}
		}
	}
	return NewTuple(attrValues...)
}

// NewTupleExprFromMap returns a new TupleExpr from a map[string]Expr.
func NewTupleExprFromMap(attrMap map[string]Expr) Expr {
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
				attrs[i] = AttrExpr{name, expr}
				i++
			}
			return &TupleExpr{attrs, attrMap}
		}
	}
	return NewTuple(attrValues...)
}

// Attrs returns a Tuple's attrs.
func (e *TupleExpr) Attrs() []AttrExpr {
	return e.attrs
}

// String returns a string representation of the expression.
func (e *TupleExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	for i, attr := range e.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		if attr.IsWildcard() {
			if attr.expr != DotIdent {
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
func (e *TupleExpr) Eval(local Scope) (Value, error) {
	tuple := EmptyTuple
	for _, attr := range e.attrs {
		var err error
		tuple, err = attr.Apply(local, tuple)
		if err != nil {
			return nil, err
		}
	}
	return tuple, nil
}

// Get returns the Expr for the given name or nil if not found.
func (e *TupleExpr) Get(name string) Expr {
	return e.attrMap[name]
}
