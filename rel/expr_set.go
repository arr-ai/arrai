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
	values := make([]Value, 0, len(elements))
	for _, expr := range elements {
		if value, ok := expr.(Value); ok {
			values = append(values, value)
			continue
		}
		tuples := make([]Expr, 0, len(elements))
		for i, elt := range elements {
			posAttr, err := NewAttrExpr("@", NewNumber(float64(i)))
			if err != nil {
				return nil
			}
			valAttr, err := NewAttrExpr(ArrayItemAttr, elt)
			if err != nil {
				return nil
			}
			tuples = append(tuples, NewTupleExpr(posAttr, valAttr))
		}
		return NewSetExpr(tuples...)
	}
	return NewArray(values...)
}

// NewDictExpr returns a new MapExpr from pairs.
func NewDictExpr(keyvals ...[2]Expr) Expr {
	values := make([]Value, 0, len(keyvals))
	for _, kv := range keyvals {
		if key, ok := kv[0].(Value); ok {
			if value, ok := kv[1].(Value); ok {
				values = append(values, NewTuple(NewAttr("@", key), NewAttr(DictValueAttr, value)))
				continue
			}
		}
		tuples := make([]Expr, 0, len(keyvals))
		for _, elt := range keyvals {
			posAttr, err := NewAttrExpr("@", elt[0])
			if err != nil {
				return nil
			}
			valAttr, err := NewAttrExpr(DictValueAttr, elt[1])
			if err != nil {
				return nil
			}
			tuples = append(tuples, NewTupleExpr(posAttr, valAttr))
		}
		return NewSetExpr(tuples...)
	}
	return NewSet(values...)
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
func (e *SetExpr) Eval(local Scope) (Value, error) {
	set := NewSet()
	for _, expr := range e.elements {
		value, err := expr.Eval(local)
		if err != nil {
			return nil, err
		}
		set = set.With(value)
	}
	return set, nil
}
