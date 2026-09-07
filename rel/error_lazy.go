package rel

import "fmt"

// lazyError formats its message only when asked for it. Pattern matching
// reports a failed match as an error, and `cond` tries arms until one binds,
// so most of these errors are created and discarded without ever being
// printed; formatting a whole tuple or set into the message eagerly was a
// measurable cost. Values are immutable, so deferring the formatting cannot
// change the message.
type lazyError struct {
	format string
	args   []interface{}
}

func lazyErrorf(format string, args ...interface{}) error {
	return &lazyError{format: format, args: args}
}

func (e *lazyError) Error() string {
	return fmt.Sprintf(e.format, e.args...)
}
