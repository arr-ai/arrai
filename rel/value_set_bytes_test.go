//nolint:dupl
package rel

import (
	"context"
	"fmt"
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

func TestBytesFormat(t *testing.T) {
	t.Parallel()

	f := func(s string) string {
		return fmt.Sprintf("%#v", NewBytes([]byte(s)))
	}

	assert.Equal(t, "", NewBytes([]byte("")).String())
	assert.Equal(t, `{}`, f(""))
	assert.Equal(t, `<<1>>`, f("\x01"))
	assert.Equal(t, `<<'\r\n'>>`, f("\r\n"))
	assert.Equal(t, `<<'\e[1m'>>`, f("\x1b[1m"))
	assert.Equal(t, `<<'hello\n'>>`, f("hello\n"))
	assert.Equal(t, `<<104, 101, 108, 108, 111, 0>>`, f("hello\000"))
}
