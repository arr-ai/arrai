package test

import "time"

type TestFile struct {
	Path     string
	Source   string
	WallTime time.Duration
	Results  []TestResult
}

type TestResult struct {
	Name    string
	Outcome TestOutcome
	Message string
}

type TestOutcome int

const (
	Failed TestOutcome = iota
	Invalid
	Ignored
	Passed
)
