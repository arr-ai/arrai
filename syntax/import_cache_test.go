package syntax

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
)

func TestImportCache(t *testing.T) {
	msgs := []string{}
	var m sync.Mutex
	start := time.Now()
	log := func(format string, args ...interface{}) {
		m.Lock()
		defer m.Unlock()
		msgs = append(msgs, fmt.Sprintf("%d ", int(time.Since(start).Milliseconds()))+fmt.Sprintf(format, args...))
	}

	cache := newCache()
	var wg sync.WaitGroup
	add := func(whenMs int, key string, value rel.Value, sleepMs int, format string, args ...interface{}) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Until(start.Add(time.Duration(whenMs) * time.Millisecond)))
			actual := cache.getOrAdd(key, func() rel.Value {
				log(">"+format, args...)
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
				log("<"+format, args...)
				return value
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

	assert.Equal(t,
		[]string{
			"10 >a 1",
			"30 >b",
			"40 >c",
			"50 >d 1",
			"80 <b",
			"90 <d 1",
			"90 >d 2",
			"110 <a 1",
			"130 <d 2",
			"140 <c",
		},
		msgs,
	)
}
