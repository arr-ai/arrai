package syntax

import (
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// withTempModule creates a temporary Go module with the given go.mod content,
// chdirs into it for the duration of the test, and resets the memoized
// requiredModuleVersions cache so each test resolves its own go.mod.
func withTempModule(t *testing.T, goModContent string) {
	t.Helper()

	dir := t.TempDir()
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(wd)) })

	require.NoError(t, os.WriteFile("go.mod", []byte(goModContent), 0o600))

	requiredModuleVersionsOnce = sync.Once{}
	t.Cleanup(func() { requiredModuleVersionsOnce = sync.Once{} })
}

func TestRetrieveModuleUsesGoModPinnedVersion(t *testing.T) {
	// DO NOT t.Parallel(): mutates the process working directory and the
	// package-level requiredModuleVersions cache.

	withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)

	out, err := exec.Command("go", "mod", "download", "github.com/pkg/errors").CombinedOutput()
	require.NoError(t, err, string(out))

	m, err := retrieveModule("github.com/pkg/errors", "")
	require.NoError(t, err)
	require.Equal(t, "github.com/pkg/errors", m.Name)
	require.Contains(t, m.Dir, "github.com/pkg/errors@v0.8.0")
}

func TestRetrieveModuleFallsBackWithoutPin(t *testing.T) {
	// DO NOT t.Parallel(): mutates the process working directory and the
	// package-level requiredModuleVersions cache.

	withTempModule(t, `module example.com/pintest

go 1.21
`)

	m, err := retrieveModule("github.com/arr-ai/arrai-import-tests", "")
	require.NoError(t, err)
	require.Equal(t, "github.com/arr-ai/arrai-import-tests", m.Name)
	require.NotEmpty(t, m.Dir)
}

// seedRequiredModuleVersions marks the memoized go.mod graph as already
// loaded with the given contents, so requiredModuleVersionOf uses it instead
// of shelling out to `go list -m -json all`.
func seedRequiredModuleVersions(t *testing.T, versions map[string]string) {
	t.Helper()

	requiredModuleVersions = versions
	requiredModuleVersionsOnce = sync.Once{}
	requiredModuleVersionsOnce.Do(func() {})
	t.Cleanup(func() {
		requiredModuleVersions = nil
		requiredModuleVersionsOnce = sync.Once{}
	})
}

func TestRequiredModuleVersionOfMatchesLongestPrefix(t *testing.T) {
	seedRequiredModuleVersions(t, map[string]string{
		"github.com/org/repo":     "v1.0.0",
		"github.com/org/repo/sub": "v2.0.0",
	})

	modPath, version, ok := requiredModuleVersionOf("github.com/org/repo/sub/file.arrai")
	require.True(t, ok)
	require.Equal(t, "github.com/org/repo/sub", modPath)
	require.Equal(t, "v2.0.0", version)

	modPath, version, ok = requiredModuleVersionOf("github.com/org/repo/file.arrai")
	require.True(t, ok)
	require.Equal(t, "github.com/org/repo", modPath)
	require.Equal(t, "v1.0.0", version)

	_, _, ok = requiredModuleVersionOf("github.com/other/repo/file.arrai")
	require.False(t, ok)
}

func TestRetrieveModuleWithInvalidPath(t *testing.T) {
	// No network access needed: an empty cache means there's nothing pinned,
	// and the path is too short to ever look like a valid module.
	seedRequiredModuleVersions(t, map[string]string{})

	m, err := retrieveModule("wrong_file_path/deps.arrai", "")
	require.Error(t, err)
	require.Nil(t, m)
}
