package syntax

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

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

	// Resolving an otherwise-unpinned import adds it to go.mod (restoring
	// v0.321.0's behaviour, lost when anz-bank/pkg/mod was replaced by an
	// early, unpinned version of this file) so this and every later
	// resolution -- in this run, and any future one -- is a fast, pinned
	// lookup instead of hitting this slow path again.
	goMod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	require.NoError(t, err)
	require.Contains(t, string(goMod), "github.com/arr-ai/arrai-import-tests",
		"resolving an unpinned import must add it as a require in go.mod")

	pinned, ok := requiredModuleOf(root, "github.com/arr-ai/arrai-import-tests")
	require.True(t, ok, "the newly-added module must be visible to a later lookup in the same run")
	require.Equal(t, m.Dir, pinned.Dir)
}

// TestRetrieveModulePersistsAcrossRuns simulates a later, separate run (fresh
// process, so a fresh memoized-graph cache) finding a module that an earlier
// run's fallback resolution added to go.mod: it must resolve via the fast
// pinned path, not fall back to `go get` again.
func TestRetrieveModulePersistsAcrossRuns(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21
`)

	_, err := retrieveModule("github.com/arr-ai/arrai-import-tests", "", root)
	require.NoError(t, err)

	// Simulate a fresh process: forget the in-memory graph, but go.mod/go.sum
	// on disk are exactly as the earlier "run" left them.
	resetRequiredModulesCache()

	pinned, ok := requiredModuleOf(root, "github.com/arr-ai/arrai-import-tests")
	require.True(t, ok, "a later run must find the module via go.mod, not need to re-add it")
	require.NotEmpty(t, pinned.Dir)
}

// TestRetrieveModuleStopsAtModuleBoundaryOnMissingPackage guards against
// retrieveModule's fallback loop shortening past a module `go` already
// confirmed exists, just because the specific package within it wasn't
// found (e.g. a stale proxy serving an older version than intended, or --
// as here -- a package path that plain doesn't exist in any version).
// Continuing to shorten past a confirmed module risks landing on a
// different, wrong (too-shallow) module instead of reporting the real one.
func TestRetrieveModuleStopsAtModuleBoundaryOnMissingPackage(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21
`)

	// cloud.google.com/go is a real, public multi-module monorepo: the repo
	// root and cloud.google.com/go/pubsub are each their own Go module. Using
	// a public module here (resolved via the Go module proxy) keeps this
	// test independent of any private-repo network access.
	m, err := retrieveModule(
		"cloud.google.com/go/pubsub/does/not/exist", "", root)
	require.NoError(t, err)
	require.Equal(t, "cloud.google.com/go/pubsub", m.Name)
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

// TestLoadRequiredModulesPreservesGoSumMtime guards against a regression
// where restore() wrote go.sum's original content back with os.WriteFile,
// which updates the mtime even when the bytes are unchanged. A build tool
// comparing mtimes (e.g. make, deciding whether a .arraiz needs rebuilding
// because it depends on go.mod) would see an untouched go.sum as "just
// modified" on every single resolve.
func TestLoadRequiredModulesPreservesGoSumMtime(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v0.8.0
`)
	cmd := exec.Command("go", "mod", "download", "github.com/pkg/errors")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	resetRequiredModulesCache()

	goSumPath := filepath.Join(root, "go.sum")
	before, err := os.Stat(goSumPath)
	require.NoError(t, err)

	// mtimes can have coarser resolution than the time between statting and
	// resolving; back-date the file so a real (bugged) rewrite is detectable
	// regardless of clock granularity.
	backdated := before.ModTime().Add(-time.Hour)
	require.NoError(t, os.Chtimes(goSumPath, backdated, backdated))

	_, err = retrieveModule("github.com/pkg/errors", "", root)
	require.NoError(t, err)

	after, err := os.Stat(goSumPath)
	require.NoError(t, err)
	require.True(t, after.ModTime().Equal(backdated),
		"resolving an import must not touch go.sum's mtime when its content is unchanged")
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
