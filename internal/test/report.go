package test

import (
	"fmt"
	"io"
)

// Report prints the results of the testing to the IO writer.
func Report(w io.Writer, rs Results) error {
	_, err := fmt.Fprintf(w, "%d/%d tests passed", rs.CountPassed(), rs.Count())
	if err != nil {
		return err
	}
	if rs.AllPassed() {
		_, err := fmt.Fprintf(w, "all tests passed")
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(w, "%d/%d tests passed", rs.CountPassed(), rs.Count())
		if err != nil {
			return err
		}
	}
	return nil
}
