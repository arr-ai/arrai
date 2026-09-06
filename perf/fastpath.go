//go:build !slowpath

package perf

// fastPathsEnabled mirrors rel's build-tag switch. Under -tags slowpath the
// evaluator takes every general fallback, so the scenario's allocation
// budget does not apply — but its output still must, which is what makes
// that run a correctness oracle over every shortcut at once.
const fastPathsEnabled = true
