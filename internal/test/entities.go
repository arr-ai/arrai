package test

import "time"

type testFile struct {
	path     string
	source   string
	wallTime time.Duration
	results  []testResult
}

type testResult struct {
	name    string
	outcome testOutcome
	message string
}

type testOutcome int

const (
	Failed testOutcome = iota
	Invalid
	Ignored
	Passed
)
