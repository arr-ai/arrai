package syntax

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// withTempModule creates a temporary Go module with the given go.mod content
// and returns its absolute path. Clears the memoized module-graph cache so
// each test resolves its own go.mod.
func withTempModule(t *testing.T, goModContent string) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goModContent), 0o600))
	resetRequiredModulesCache()
	t.Cleanup(resetRequiredModulesCache)
	return dir
}

func TestRetrieveModuleUsesGoModPinnedVersion(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)

	cmd := exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// Force a fresh go list against this module root (download above may not
	// populate Dir until list runs in-module).
	resetRequiredModulesCache()

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Equal(t, "github.com/pkg/errors", m.Name)
	require.Contains(t, m.Dir, "github.com/pkg/errors@v0.8.0")
}

func TestRetrieveModuleHonoursReplaceDir(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0

replace github.com/pkg/errors => ./errors-fork
`)

	fork := filepath.Join(root, "errors-fork")
	require.NoError(t, os.MkdirAll(fork, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(fork, "go.mod"), []byte("module github.com/pkg/errors\n\ngo 1.21\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(fork, "errors.go"), []byte("package errors\n"), 0o600))

	resetRequiredModulesCache()

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Equal(t, "github.com/pkg/errors", m.Name)
	require.Equal(t, fork, m.Dir)
}

func TestRetrieveModuleUsesModuleRootNotCwd(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)

	cmd := exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	resetRequiredModulesCache()

	// Resolve while cwd is elsewhere — must still see the pin under root.
	wd, err := os.Getwd()
	require.NoError(t, err)
	elsewhere := t.TempDir()
	require.NoError(t, os.Chdir(elsewhere))
	t.Cleanup(func() { require.NoError(t, os.Chdir(wd)) })

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Contains(t, m.Dir, "github.com/pkg/errors@v0.8.0")
}

func TestRetrieveModuleFallsBackWithoutPin(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21
`)

	m, err := retrieveModule("github.com/arr-ai/arrai-import-tests", "", root)
	require.NoError(t, err)
	require.Equal(t, "github.com/arr-ai/arrai-import-tests", m.Name)
	require.NotEmpty(t, m.Dir)
}

func TestRequiredModuleOfMatchesLongestPrefix(t *testing.T) {
	resetRequiredModulesCache()
	t.Cleanup(resetRequiredModulesCache)
	seedRequiredModules("", map[string]requiredModule{
		"github.com/org/repo":     {Path: "github.com/org/repo", Version: "v1.0.0", Dir: "/mod/repo@v1.0.0"},
		"github.com/org/repo/sub": {Path: "github.com/org/repo/sub", Version: "v2.0.0", Dir: "/mod/sub@v2.0.0"},
	})

	mod, ok := requiredModuleOf("", "github.com/org/repo/sub/file.arrai")
	require.True(t, ok)
	require.Equal(t, "github.com/org/repo/sub", mod.Path)
	require.Equal(t, "v2.0.0", mod.Version)
	require.Equal(t, "/mod/sub@v2.0.0", mod.Dir)

	mod, ok = requiredModuleOf("", "github.com/org/repo/file.arrai")
	require.True(t, ok)
	require.Equal(t, "github.com/org/repo", mod.Path)

	_, ok = requiredModuleOf("", "github.com/other/repo/file.arrai")
	require.False(t, ok)
}

func TestRetrieveModuleWithInvalidPath(t *testing.T) {
	resetRequiredModulesCache()
	t.Cleanup(resetRequiredModulesCache)
	seedRequiredModules("", map[string]requiredModule{})

	m, err := retrieveModule("wrong_file_path/deps.arrai", "", "")
	require.Error(t, err)
	require.Nil(t, m)
}
