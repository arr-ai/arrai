package rel

import (
	"context"
	"errors"
	"fmt"
)

// errPatternMismatch is every structural mismatch on the fast path: one
// shared value, so trying and discarding a match — cond's bread and butter —
// allocates nothing. Bind again with scopeBuilder.explain for the real
// message.
var errPatternMismatch = errors.New("pattern did not match")

// explainBind re-runs a failed bind to produce its user-facing error.
// b must be a fresh builder with explain set.
func explainBind(ctx context.Context, p Pattern, scope Scope, value Value) error {
	b := scopeBuilder{explain: true}
	_, err := p.Bind(ctx, scope, value, &b)
	if err == nil {
		// A pure re-match cannot succeed after failing; guard anyway.
		return errPatternMismatch
	}
	return err
}

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
