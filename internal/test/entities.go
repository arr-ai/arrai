package test

import "time"

type TestFile struct {
	path     string
	source   string
	wallTime time.Duration
	results  []TestResult
}

type TestResult struct {
	name    string
	outcome TestOutcome
	message string
}

type TestOutcome int

const (
	Failed TestOutcome = iota
	Invalid
	Ignored
	Passed
)
