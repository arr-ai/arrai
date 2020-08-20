package main

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkRunDFA(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RunTestDFA(b)
	}
}

func RunTestDFA(b *testing.B) {
	var buf bytes.Buffer
	evalFile(context.Background(), "./test/data_dictionary.arrai", &buf, "")
}
