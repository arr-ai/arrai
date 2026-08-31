//go:build !race

package perf

// raceEnabled is true when built with -race, which allocates on its own
// account and so puts the scenario's allocation budget out of reach. The
// output check still applies.
const raceEnabled = false
