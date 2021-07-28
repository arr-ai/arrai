package test

import "time"

type File struct {
	Path     string
	Source   string
	WallTime time.Duration
	Results  []Result
}

type Result struct {
	Name    string
	Outcome Outcome
	Message string
}

type Outcome int

const (
	Failed Outcome = iota
	Invalid
	Ignored
	Passed
)
