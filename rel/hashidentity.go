//go:build hashidentity

package rel

// hashIdentity replaces GenericSet.Equal with a 128-bit comparison (🎯T22).
const hashIdentity = true
