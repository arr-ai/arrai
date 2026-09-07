package rel

import "sync/atomic"

// whereErr collects the error from a Where predicate. frozen runs the
// predicate on several goroutines once a set is large enough for its
// parallelism to engage, so the error cannot simply be captured by the
// closure: that is a data race, and it was one until this type existed.
//
// The first writer wins. Which error that is, when several rows fail
// concurrently, depends on scheduling — but any of them is a faithful
// report that the predicate failed, and evaluation stops either way.
type whereErr struct {
	err atomic.Pointer[error]
}

// failed reports whether the predicate has already failed, so the remaining
// elements can be skipped.
func (w *whereErr) failed() bool {
	return w.err.Load() != nil
}

// set records err if no error has been recorded yet.
func (w *whereErr) set(err error) {
	w.err.CompareAndSwap(nil, &err)
}

// get returns the recorded error, or nil.
func (w *whereErr) get() error {
	if e := w.err.Load(); e != nil {
		return *e
	}
	return nil
}
