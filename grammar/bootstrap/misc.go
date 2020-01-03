package bootstrap

const (
	// A function that should never have been called was called.
	Inconceivable Panicker = "How did this happen!?"

	// A function that isn't ready for use yet was called.
	Unfinished Panicker = "not yet implemented"
)

type Panicker string

var _ error = Panicker("")

func (p Panicker) Error() string {
	return string(p)
}
