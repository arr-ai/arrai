package rel

import (
	"github.com/go-errors/errors"
)

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

// NewUnionExpr evaluates a <&> b.
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

// NewDiffExpr evaluates a <&> b.
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

// NewSymmDiffExpr evaluates a <&> b.
func NewSymmDiffExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "(-)", "(%s (-) %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return SymmetricDifference(x, y), nil
				}
				return nil, errors.Errorf("(-) rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("(-) lhs must be a Set, not %T", a)
		})
}

// NewConcatExpr evaluates a <&> b.
func NewConcatExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "|", "(%s | %s)",
		func(a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Concatenate(x, y)
				}
				return nil, errors.Errorf("(-) rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("(-) lhs must be a Set, not %T", a)
		})
}
