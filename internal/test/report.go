package test

import (
	"fmt"
	"io"
)

// Report prints the results of the testing to the IO writer.
func Report(w io.Writer, rs Results) error {
	if rs.AllPassed() {
		_, err := fmt.Fprintln(w, "all tests passed")
		if err != nil {
			return err
		}
	} else {
		for _, r := range rs.results {
			if !r.pass {
				_, err := fmt.Fprintf(w, "%s failed\n", r.file)
				if err != nil {
					return err
				}
			}
		}
		_, err := fmt.Fprintf(w, "%d/%d tests passed\n", rs.CountPassed(), rs.Count())
		if err != nil {
			return err
		}
	}
	return nil
}
