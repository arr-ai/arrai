package bootstrap

const (
	// Inconceivable indicates that a function should never have been called.
	Inconceivable Error = "How did this happen!?"

	// Unfinished indicates that a function isn't ready for use yet.
	Unfinished Error = "not yet implemented"

	// BadInput indicates that a function was given bad inputs.
	BadInput Error = "bad input"
)

type Error string

var _ error = Error("")

func (p Error) Error() string {
	return string(p)
}
