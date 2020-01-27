package rel

import (
	"bytes"

	"github.com/go-errors/errors"
)

// SetExpr returns the tuple or set with a single field replaced by an
// expression.
type SetExpr struct {
	elements []Expr
}

// NewSetExpr returns a new TupleExpr.
func NewSetExpr(elements ...Expr) Expr {
	values := make([]Value, len(elements))
	for i, expr := range elements {
		value, ok := expr.(Value)
		if !ok {
			return &SetExpr{elements}
		}
		values[i] = value
	}
	return NewSet(values...)
}

// NewRelationExpr returns a new relation for the given data.
func NewRelationExpr(names []string, tuples ...[]Expr) (Expr, error) {
	elements := make([]Expr, len(tuples))
	for i, tuple := range tuples {
		if len(tuple) != len(names) {
			return nil, errors.Errorf(
				"heading-tuple mismatch: %v vs %v", names, tuple)
		}
		attrs := make([]AttrExpr, len(names))
		for i, name := range names {
			attrs[i] = AttrExpr{name, tuple[i]}
		}
		elements[i] = NewTupleExpr(attrs...)
	}
	return NewSetExpr(elements...), nil
}

// NewArrayExpr returns a new TupleExpr.
func NewArrayExpr(elements ...Expr) Expr {
	values := make([]Value, len(elements))
	for i, expr := range elements {
		value, ok := expr.(Value)
		if !ok {
			tuples := make([]Expr, len(elements))
			for i, elt := range elements {
				posAttr, err := NewAttrExpr("@", NewNumber(float64(i)))
				if err != nil {
					panic(err)
				}
				valAttr, err := NewAttrExpr(ArrayItemAttr, elt)
				if err != nil {
					panic(err)
				}
				tuples[i] = NewTupleExpr(posAttr, valAttr)
			}
			return NewSetExpr(tuples...)
		}
		values[i] = value
	}
	return NewArray(values...)
}

// Elements returns a Set's elements.
func (e *SetExpr) Elements() []Expr {
	return e.elements
}

// String returns a string representation of the expression.
func (e *SetExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	i := 0
	for _, expr := range e.elements {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(expr.String())
		i++
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the subject
func (e *SetExpr) Eval(local, global Scope) (Value, error) {
	set := NewSet()
	for _, expr := range e.elements {
		value, err := expr.Eval(local, global)
		if err != nil {
			return nil, err
		}
		set = set.With(value)
	}
	return set, nil
}
