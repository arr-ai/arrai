package test

// Result describes the outcome of a particular test case.
type Result struct {
	file string
	pass bool
}

// Results is a collection of test results from a test suite.
type Results struct {
	results []Result
}

// Add adds a result to the set of results.
func (rs *Results) Add(r Result) {
	rs.results = append(rs.results, r)
}

// AllPassed returns true if all results were successful.
func (rs Results) AllPassed() bool {
	for _, r := range rs.results {
		if !r.pass {
			return false
		}
	}
	return true
}

// Count returns the number of results in the set.
func (rs Results) Count() int {
	return len(rs.results)
}

// CountFailed returns the number of tests that failed.
func (rs Results) CountFailed() int {
	count := 0
	for _, r := range rs.results {
		if !r.pass {
			count++
		}
	}
	return count
}

// CountPassed returns the number of tests that passed.
func (rs Results) CountPassed() int {
	count := 0
	for _, r := range rs.results {
		if r.pass {
			count++
		}
	}
	return count
}
