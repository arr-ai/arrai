// +build timingsensitive

package syntax

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

	cache := newCache()
	var wg sync.WaitGroup
	add := func(whenMs int, key string, value rel.Value, sleepMs int, descr string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Until(start.Add(time.Duration(whenMs) * time.Millisecond)))
			actual, _ := cache.getOrAdd(key, func() (rel.Value, error) {
				log("> %s", descr)
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
				log("< %s", descr)
				return value, nil
			})
			rel.AssertEqualValues(t, actual, value)
		}()
	}

	add(10, "a", rel.None, 100, "a 1")
	add(20, "a", rel.None, 0, "a 2")
	add(30, "b", rel.None, 50, "b")
	add(40, "c", rel.None, 100, "c")
	add(50, "d", nil, 40, "d 1")
	add(60, "d", nil, 40, "d 2")

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
