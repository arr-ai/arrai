package rel

import (
	"github.com/go-errors/errors"
)

// NewJoinExpr evaluates a <&> b.
func NewJoinExpr(a, b Expr) Expr {
	return newBinExpr(a, b, "<&>", "(%s <&> %s)",
		func(a, b Value, _, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if y, ok := b.(Set); ok {
					return Join(x, y), nil
				}
				return nil, errors.Errorf("<&> rhs must be a Set, not %T", b)
			}
			return nil, errors.Errorf("<&> lhs must be a Set, not %T", a)
		})
}
