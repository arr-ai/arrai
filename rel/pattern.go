package rel

import (
	"fmt"
)

// Pattern can be inside an Expr, Expr can be a Pattern.
type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	Bind(scope Scope, value Value) (Scope, error)
}
