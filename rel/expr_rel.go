package rel

import (
	"strings"

	"github.com/go-errors/errors"
)

// NewRelationExpr returns a new relation for the given data.
func NewRelationExpr(names []string, tuples ...[]Expr) (Expr, error) {
	elements := make([]Expr, len(tuples))
	stringCharTuples := 0
	arrayItemTuples := 0
	dictEntryTuples := 0
	for i, tuple := range tuples {
		if len(tuple) != len(names) {
			return nil, errors.Errorf(
				"heading-tuple mismatch: %v vs %v", names, tuple)
		}
		attrs := make([]AttrExpr, len(names))
		for i, name := range names {
			attrs[i] = AttrExpr{name, tuple[i]}
		}
		if len(attrs) == 2 {
			if attrs[1].name == "@" {
				attrs[0], attrs[1] = attrs[1], attrs[0]
			}
			if attrs[0].name == "@" && strings.HasPrefix(attrs[1].name, "@") {
				switch attrs[1].name {
				case StringCharAttr:
					elements[i] = NewStringCharTupleExpr(attrs[0].expr, attrs[1].expr)
					stringCharTuples++
				case ArrayItemAttr:
					elements[i] = NewArrayItemTupleExpr(attrs[0].expr, attrs[1].expr)
					arrayItemTuples++
				case DictValueAttr:
					elements[i] = NewDictEntryTupleExpr(attrs[0].expr, attrs[1].expr)
					dictEntryTuples++
				default:
					elements[i] = NewTupleExpr(attrs...)
				}
				continue
			}
		}
		elements[i] = NewTupleExpr(attrs...)
	}
	switch len(elements) {
	case stringCharTuples:
		charExprs := make([]Expr, 0, len(elements))
		for _, e := range elements {
			charExprs = append(charExprs, e.(StringCharTupleExpr))
		}
		// TODO: Implement NewStringCharSetExpr.
		return NewSetExpr(charExprs...), nil
	case arrayItemTuples:
		entryExprs := make([]Expr, 0, len(elements))
		for _, e := range elements {
			entryExprs = append(entryExprs, e.(ArrayItemTupleExpr))
		}
		// TODO: Implement NewArrayItemSetExpr.
		return NewSetExpr(entryExprs...), nil
	case dictEntryTuples:
		entryExprs := make([]DictEntryTupleExpr, 0, len(elements))
		for _, e := range elements {
			entryExprs = append(entryExprs, e.(DictEntryTupleExpr))
		}
		return NewDictExpr(true, entryExprs...), nil
	}
	return NewSetExpr(elements...), nil
}

// NewJoinExpr evaluates a <&> b.
func NewJoinExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "<&>", "(%s <&> %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Join(x, y), nil
				}
				return nil, errors.Errorf("<&> rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("<&> lhs must be a Set, not %T", a)
		})
}

// NewUnionExpr evaluates a | b.
func NewUnionExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "|", "(%s | %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Union(x, y), nil
				}
				return nil, errors.Errorf("| rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("| lhs must be a Set, not %T", a)
		})
}

// NewDiffExpr evaluates a &~ b.
func NewDiffExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "&~", "(%s &~ %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Difference(x, y), nil
				}
				return nil, errors.Errorf("&~ rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("&~ lhs must be a Set, not %T", a)
		})
}

// NewSymmDiffExpr evaluates a ~~ b.
func NewSymmDiffExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "~~", "(%s ~~ %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return SymmetricDifference(x, y), nil
				}
				return nil, errors.Errorf("~~ rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("~~ lhs must be a Set, not %T", a)
		})
}

// NewConcatExpr evaluates a ++ b.
func NewConcatExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "++", "(%s ++ %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Concatenate(x, y)
				}
				return nil, errors.Errorf("++ rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("++ lhs must be a Set, not %T", a)
		})
}
