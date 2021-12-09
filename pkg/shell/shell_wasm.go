//go:build wasm
// +build wasm

package shell

import (
	"context"

	"github.com/arr-ai/arrai/rel"
)

func Shell(_ context.Context, _ []rel.ContextErr) error {
	panic("not implemented")
}
