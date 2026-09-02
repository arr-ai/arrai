package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// T16: signal/os.Exit cleanup (#737). The process-level Notify path is
// flushProfilersOnSignal; this tests the handler body CI can actually run.
func TestWaitAndFlushCallsStop(t *testing.T) {
	t.Parallel()
	sig := make(chan os.Signal, 1)
	stopped := make(chan struct{})
	exited := make(chan int, 1)
	go waitAndFlush(sig, func() { close(stopped) }, func(code int) { exited <- code })
	sig <- os.Interrupt
	<-stopped
	assert.Equal(t, 1, <-exited)
}
