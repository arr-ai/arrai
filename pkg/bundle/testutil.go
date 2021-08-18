package bundle

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/arr-ai/arrai/syntax"
)

// ConfigFile is a test utility to create the content of a bundle config file.
func ConfigFile(mainRoot, mainFile string) string {
	return fmt.Sprintf("(main_root: %q, main_file: %q)", mainRoot, mainFile)
}

// ModuleFile is a test utility to create a path for an arrai script in a
// bundled script.
func ModuleFile(file string) string {
	return path.Join(syntax.ModuleDir, file)
}

// NoModuleFile is a test utility to create a path for an arrai script that is
// not part of a module.
func NoModuleFile(file string) string {
	return path.Join(syntax.NoModuleDir, file)
}

// SentinelFile is a test utility to create a path for a sentinel file in a
// module within a bundle.
func SentinelFile(file string) string {
	return ModuleFile(SentinelPath(file))
}

// SentinelPath is a test utility to create a path for a sentinel file without
// the sentinel module prefix.
func SentinelPath(file string) string {
	return path.Join(file, syntax.ModuleRootSentinel)
}

// MustCreateTestBundleFromMap takes in a map of files representing a memory
// filesystem and then path to main script and returns the bundled script.
func MustCreateTestBundleFromMap(t *testing.T, files map[string]string, script string) []byte {
	fs := ctxfs.CreateTestMemMapFs(t, files)
	return MustCreateTestBundleFromFs(t, fs, script)
}

// MustCreateTestBundleFromFs takes in a filesystem and path to main script
// and returns the bundled script.
func MustCreateTestBundleFromFs(t *testing.T, fs afero.Fs, script string) []byte {
	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
	ctx = ctxrootcache.WithRootCache(ctx)
	buf := &bytes.Buffer{}
	require.NoError(t, BundledScripts(ctx, syntax.MustAbs(t, script), buf))
	return buf.Bytes()
}
