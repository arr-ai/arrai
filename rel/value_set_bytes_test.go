//nolint:dupl
package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func TestBytesSetCallAll(t *testing.T) {
	t.Parallel()

	ctx := arraictx.InitRunCtx(context.Background())

	b := NewSetBuilder()
	err := NewBytes([]byte("abc")).CallAll(ctx, NewString([]rune("0")), b)
	if assert.NoError(t, err) {
		set, err := b.Finish()
		require.NoError(t, err)
		assert.False(t, set.IsTrue())
	}
}
