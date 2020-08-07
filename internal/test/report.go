package test

import (
	"fmt"
	"io"
)

// Report prints the results of the testing to the IO writer.
func Report(rs Results, w io.Writer) error {
	if rs.AllPassed() {
		_, err := fmt.Fprintf(w, "all tests passed")
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(w, "%d tests failed", rs.CountFailed())
		if err != nil {
			return err
		}
	}
	return nil
}
