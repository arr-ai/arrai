// +build timingsensitive

package importcache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
)

func TestImportCache(t *testing.T) {
	t.Parallel()

	msgs := []string{}
	var m sync.Mutex
	start := time.Now()
	log := func(format string, args ...interface{}) {
		m.Lock()
		defer m.Unlock()
		msgs = append(msgs, fmt.Sprintf(format, args...))
	}

	cache := newImportCache()
	var wg sync.WaitGroup
	add := func(ctx context.Context, whenMs int, key string, value rel.Value, sleepMs int, descr string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Until(start.Add(time.Duration(whenMs) * time.Millisecond)))
			actual, _ := cache.getOrAdd(key, func() (rel.Expr, error) {
				log("> %s", descr)
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
				log("< %s", descr)
				return value, nil
			})
			// to avoid evaluating nil expression
			if value == nil {
				assert.Equal(t, nil, actual)
				return
			}
			actualValue, err := actual.Eval(ctx, rel.EmptyScope)
			require.NoError(t, err)
			rel.AssertEqualValues(t, actualValue, value)
		}()
	}

	ctx := arraictx.InitRunCtx(context.Background())
	add(ctx, 10, "a", rel.None, 100, "a 1")
	add(ctx, 20, "a", rel.None, 0, "a 2")
	add(ctx, 30, "b", rel.None, 50, "b")
	add(ctx, 40, "c", rel.None, 100, "c")
	add(ctx, 50, "d", nil, 40, "d 1")
	add(ctx, 60, "d", nil, 40, "d 2")

	wg.Wait()

	// TODO: Reinstate after figuring out how to make this more stable.
	// assert.Equal(t,
	// 	[]string{
	// 		"> a 1",
	// 		"> b",
	// 		"> c",
	// 		"> d 1",
	// 		"< b",
	// 		"< d 1",
	// 		"> d 2",
	// 		"< a 1",
	// 		"< d 2",
	// 		"< c",
	// 	},
	// 	msgs,
	// )
}
