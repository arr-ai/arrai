//nolint:dupl
package rel

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func TestBytesSetCallAll(t *testing.T) {
	t.Parallel()

	ctx := arraictx.InitRunCtx(context.Background())
	_, err := NewBytes([]byte("abc")).CallAll(ctx, NewString([]rune("0")))
	assert.Error(t, err, "CallAll arg must be a number, not rel.String")
}
