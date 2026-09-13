package rel

import (
	"testing"
	"unique"
)

// unique.Handle evaluation for 🎯T27. Measured on Apple M4 Max / go1.25,
// testing.B.Loop, after NewGoString interned short contents (≤64):
//
//	NewGoString (interned)    9.0 ns/op   0 B   0 allocs
//	unique.Make.Value        12.1 ns/op   0 B   0 allocs
//	newGoString (uninterned) 15.4 ns/op  32 B   1 alloc
//
// Rejected unique.Handle: the map intern is faster and returns a String.

const internBenchKey = "name-42"

func BenchmarkNewGoStringShort(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = NewGoString(internBenchKey)
	}
}

func BenchmarkNewGoStringShortUninterned(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = newGoString(internBenchKey)
	}
}

func BenchmarkUniqueHandleShort(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = unique.Make(internBenchKey).Value()
	}
}
