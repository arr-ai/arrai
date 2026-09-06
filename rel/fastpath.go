//go:build !slowpath

package rel

// fastPaths enables the representation-specific shortcuts that skip a
// general fallback: shape-matched row conversion, the identifier and
// attribute caches, direct calls, index-answered `where`, and so on.
//
// Build with -tags slowpath to force every one of them to its fallback. The
// two must produce identical results, so running the test corpus both ways
// is a differential oracle over every shortcut at once — including the ones
// whose fallback ordinary tests would otherwise never reach. It is a
// constant, so the guarded branch costs nothing in a normal build.
const fastPaths = true
