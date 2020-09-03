package syntax

import (
	"context"
	"path"
	"testing"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/stretchr/testify/assert"
)

type memFile struct{ name, content string }

type bundleCases struct {
	name               string
	mainFile, rootFile memFile
	expectedConfig     bundleConfig
	expectedErr        string
}

func TestSetupBundle(t *testing.T) {
	t.Parallel()

	cases := []bundleCases{
		{
			"bundle with root",
			memFile{"/another/test/root/test.arrai", "1"},
			memFile{path.Join("/another/test/root/", ModuleRootSentinel), "module test/root\n"},
			bundleConfig{
				"test/root",
				path.Join(ModuleDir, "/test/root/test.arrai"),
				"/another/test/root",
			}, "",
		},
		{
			"bundle without root",
			memFile{"/another/test/root/test.arrai", "1"},
			memFile{},
			bundleConfig{
				"",
				path.Join(NoModuleDir, "test.arrai"),
				"/another/test/root",
			}, "",
		},
		{
			"root has no sentinel",
			memFile{"/another/test/root/test.arrai", "1"},
			memFile{path.Join("/another/test/root/", ModuleRootSentinel), "random"},
			bundleConfig{
				"test/root",
				path.Join(ModuleDir, "/test/root/test.arrai"),
				"/another/test/root",
			}, errSentinelHasNoModule.Error(),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			files := make(map[string]string)
			files[c.mainFile.name] = c.mainFile.content
			files[c.rootFile.name] = c.rootFile.content
			fs := ctxfs.CreateTestMemMapFs(t, files)
			ctx := ctxfs.SourceFsOnto(context.Background(), fs)
			ctx, err := SetupBundle(ctx, c.mainFile.name, []byte(c.mainFile.content))
			if c.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedConfig, fromBundleConfig(ctx))
			} else {
				assert.EqualError(t, err, c.expectedErr)
			}
		})
	}
}
