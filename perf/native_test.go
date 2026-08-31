package perf

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNativeReconstruct keeps the handcrafted C++ implementation of the
// reconstruct scenario honest: it must stay byte-identical to the same
// golden output the arr.ai pipeline is held to. It skips when no C++
// compiler is available rather than failing.
func TestNativeReconstruct(t *testing.T) {
	if testing.Short() {
		t.Skip("perf scenario: skipped under -short")
	}
	cxx, err := exec.LookPath("clang++")
	if err != nil {
		if cxx, err = exec.LookPath("g++"); err != nil {
			t.Skip("no C++ compiler on PATH")
		}
	}

	dir, err := filepath.Abs("reconstruct")
	require.NoError(t, err)
	bin := filepath.Join(t.TempDir(), "reconstruct-native")
	build := exec.Command(cxx, "-std=c++20", "-O2", "-o", bin,
		filepath.Join(dir, "native", "reconstruct.cpp"))
	out, err := build.CombinedOutput()
	require.NoError(t, err, "compile: %s", out)

	expected, err := os.ReadFile(filepath.Join(dir, "expected.arrai"))
	require.NoError(t, err)

	start := time.Now()
	got, err := exec.Command(bin, filepath.Join(dir, "model.sysl.pb")).Output()
	require.NoError(t, err)
	t.Logf("native reconstruct: %s", time.Since(start).Round(time.Microsecond))

	require.True(t, bytes.Equal(expected, got),
		"native output differs from the v0.321.0 reference")
}
