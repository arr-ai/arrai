package syntax

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/ctxrootcache"
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
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")
}

func TestRetrieveModuleHonoursReplaceDir(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0

replace github.com/pkg/errors => ./errors-fork
`)

	fork := filepath.Join(root, "errors-fork")
	require.NoError(t, os.MkdirAll(fork, 0o700))
	gomod := []byte("module github.com/pkg/errors\n\ngo 1.21\n")
	require.NoError(t, os.WriteFile(filepath.Join(fork, "go.mod"), gomod, 0o600))
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
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")
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

func TestRetrieveModuleSelfHealsMissingGoSum(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)

	// Deliberately skip `go mod download`: go.sum has no entry for the pinned
	// module. Under plain -mod=readonly, `go list -m -json all` would fail and
	// the pin would be missed entirely, falling back to @latest.
	require.NoFileExists(t, filepath.Join(root, "go.sum"))

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Equal(t, "github.com/pkg/errors", m.Name)
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")
}

func TestPrimaryRootMakesNestedImportsUseTheTopLevelPin(t *testing.T) {
	outer := withTempModule(t, `module example.com/outer

go 1.21

require github.com/pkg/errors v0.8.0
`)
	cmd := exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = outer
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// Stands in for a dependency's own copy of go.mod (e.g. inside the
	// read-only module cache) pinning a different, stale version -- this
	// must lose to the top-level project's pin once it's part of a larger
	// build graph, not win just because it's nearer the importing file.
	nested := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(nested, "go.mod"), []byte(`module example.com/nested

go 1.21

require github.com/pkg/errors v0.9.1
`), 0o600))
	cmd = exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = nested
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	resetRequiredModulesCache()

	ctx := ctxrootcache.WithPrimaryRootCache(context.Background())

	// First resolution, from the top-level project itself: establishes the
	// primary root.
	root := ctxrootcache.PrimaryRoot(ctx, func() string { return outer })
	require.Equal(t, outer, root)

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")

	// Second resolution simulates a nested import made from a file under
	// `nested`. PrimaryRoot must still return the outer root -- not
	// `nested`, even though a compute() naively rooted at the importing
	// file would find `nested`'s go.mod instead.
	nestedRoot := ctxrootcache.PrimaryRoot(ctx, func() string { return nested })
	require.Equal(t, outer, nestedRoot, "nested imports must resolve against the primary (top-level) module root")

	m, err = retrieveModule("github.com/pkg/errors", "", nestedRoot)
	require.NoError(t, err)
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0",
		"the outer project's pin must win, not the nested directory's own (different) pin")
}

func TestLoadRequiredModulesDoesNotLeaveGoSumChanged(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)
	// No go.sum at all: resolving requires `go list -mod=mod` to self-heal
	// it from scratch, the most extreme case of the self-heal touching the
	// file. It must still come out exactly as it went in afterwards.
	require.NoFileExists(t, filepath.Join(root, "go.sum"))

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")

	require.NoFileExists(t, filepath.Join(root, "go.sum"),
		"resolving an import must not leave go.sum changed in the caller's project")
}

func TestLoadRequiredModulesPreservesExistingGoModAndGoSum(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)
	cmd := exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	resetRequiredModulesCache()

	goModBefore, err := os.ReadFile(filepath.Join(root, "go.mod"))
	require.NoError(t, err)
	goSumBefore, err := os.ReadFile(filepath.Join(root, "go.sum"))
	require.NoError(t, err)

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)
	require.Contains(t, filepath.ToSlash(m.Dir), "github.com/pkg/errors@v0.8.0")

	goModAfter, err := os.ReadFile(filepath.Join(root, "go.mod"))
	require.NoError(t, err)
	goSumAfter, err := os.ReadFile(filepath.Join(root, "go.sum"))
	require.NoError(t, err)
	require.Equal(t, goModBefore, goModAfter, "resolving an import must not change an existing go.mod")
	require.Equal(t, goSumBefore, goSumAfter, "resolving an import must not change an existing go.sum")
}

func TestRetrieveModuleErrorsWhenPinnedVersionUnavailable(t *testing.T) {
	root := t.TempDir()
	resetRequiredModulesCache()
	t.Cleanup(resetRequiredModulesCache)
	seedRequiredModules(root, map[string]requiredModule{
		"github.com/pkg/errors": {Path: "github.com/pkg/errors", Version: "v99.99.99-does-not-exist"},
	})

	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.Error(t, err)
	require.Nil(t, m)
	require.Contains(t, err.Error(), "github.com/pkg/errors@v99.99.99-does-not-exist")
}
