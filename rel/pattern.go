package rel

import (
	"fmt"
)

// Pattern can be inside an Expr, Expr can be a Pattern.
type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	// Bind binds pattern with value and add the binding pair to scope
	Bind(scope Scope, value Value) (Scope, error)

	// Bindings returns all the names a pattern expects to bind
	Bindings() []string
}
