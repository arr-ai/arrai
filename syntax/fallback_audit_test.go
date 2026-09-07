package syntax

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// T16: silent-fallback audit. Pin resolution that cannot fetch the pinned
// version must fail, not fall through to @latest (#736/#742). The happy
// path (missing go.sum, pin still honoured) is TestRetrieveModuleSelfHealsMissingGoSum.
func TestPinnedModuleDownloadFailureIsLoud(t *testing.T) {
	root := withTempModule(t, `module example.com/pintest

go 1.21

require github.com/pkg/errors v99.99.99
`)
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.sum"), []byte(""), 0o600))
	m, err := retrieveModule("github.com/pkg/errors", "", root)
	require.Error(t, err, "a pin that cannot be fetched must not resolve to @latest")
	require.Nil(t, m)
}
