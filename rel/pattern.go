package rel

import (
	"context"
	"fmt"
)

// Pattern can be inside an Expr, Expr can be a Pattern.
type Pattern interface {
	// Require a String() method.
	fmt.Stringer

	// Bind matches value against the pattern, adding bindings to b (which
	// is shared by every part of an enclosing pattern: one flat frame per
	// match, each binding in a fixed slot). A mismatch is reported as an
	// error, which every implementation constructs lazily.
	Bind(ctx context.Context, scope Scope, value Value, b *scopeBuilder) (context.Context, error)

	// Bindings returns all the names a pattern expects to bind
	Bindings() []string
}
