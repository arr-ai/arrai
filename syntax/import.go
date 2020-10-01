package syntax

import (
	"context"
	"errors"
	"fmt"
	"github.com/joshcarp/gop/gop/cli"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"github.com/arr-ai/arrai/translate"
	"github.com/joshcarp/gop/gop"
)

// ModuleRootSentinel is a file which marks the module root of a project.
const ModuleRootSentinel = "go.mod"

var (
	cache             = newCache()
	errModuleNotExist = errors.New("module root not found")
)

func importLocalFile(ctx context.Context, fromRoot bool, importPath, sourceDir string) (rel.Expr, error) {
	if fromRoot {
		rootPath, err := findRootFromModule(ctx, sourceDir)
		if err != nil {
			return nil, err
		}
		if err = addLocalRoot(ctx, rootPath); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(importPath, "/") {
			importPath = rootPath + "/" + strings.ReplaceAll(importPath, "../", "")
		}
	}

	if err := bundleLocalFile(ctx, importPath); err != nil {
		return nil, err
	}

	v, err := fileValue(ctx, importPath)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func importExternalContent(ctx context.Context, importPath string) (rel.Expr, error) {
	var moduleErr error
	if !strings.HasPrefix(importPath, "http://") && !strings.HasPrefix(importPath, "https://") {
		v, err := importModuleFile(ctx, importPath)
		if err == nil {
			return v, nil
		}
		moduleErr = err

		// Since an explicit schema is allowed, it's OK to assume https as the default.
		importPath = "https://" + importPath
	}

	v, err := importURL(ctx, importPath)
	if err != nil {
		if moduleErr != nil {
			return nil, fmt.Errorf("failed to import %s - %s, and %s", importPath, moduleErr.Error(), err.Error())
		}
		return nil, err
	}

	return v, nil
}

func importModuleFile(ctx context.Context, importPath string) (rel.Expr, error) {
	if isRunningBundle(ctx) {
		return fileValue(ctx, path.Join(ModuleDir, importPath))
	}
	repo, resource, ver, _ := gop.ProcessRequest(importPath)
	if path.Ext(resource) == ""{
		resource += ".arrai"
	}
	if ver == ""{
		ver = "HEAD"
	}
	if repo != ""{
		repo += "/"
	}
	var tokens map[string]string
	f, _ := os.Open("~/.git-credentials")
	if f != nil{
		a, _ := ioutil.ReadAll(f)
		tokens, _ = cli.TokensFromGitCredentialsFile(a)
	}

	r := cli.Moduler(ctxfs.SourceFsFrom(ctx), "arrai_modules.yaml","arrai_modules", os.Getenv("ARRAI_PROXY"), tokens)
	bytes, _, err := r.Retrieve(path.Join(repo,resource)+"@"+ver)
	if err != nil{
		return nil, err
	}
	switch filepath.Ext(resource) {
	case ".json":
		return bytesJSONToArrai(bytes)
	case ".yml", ".yaml":
		return translate.BytesYamlToArrai(bytes)
	}
	return bytesValue(ctx, resource, bytes)
}

func findRootFromModule(ctx context.Context, modulePath string) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	currentPath = bundleToValidPath(ctx, currentPath)

	if r, exists, err := ctxrootcache.LoadRoot(ctx, currentPath); err != nil {
		return "", err
	} else if exists {
		return r, nil
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	systemRoot = bundleToValidPath(ctx, systemRoot)

	// 16 is enough for pretty much all cases.
	paths := append(make([]string, 0, 16), currentPath)

	// Keep walking up the directories to find nearest root marker
	for {
		exists, err := tools.FileExists(ctx, filepath.Join(currentPath, ModuleRootSentinel))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			for _, p := range paths {
				if err := ctxrootcache.StoreRoot(ctx, p, currentPath); err != nil {
					return "", err
				}
			}
			return currentPath, nil
		case reachedRoot:
			return "", errModuleNotExist
		case err != nil:
			return "", err
		}
		currentPath = filepath.Dir(currentPath)
		paths = append(paths, currentPath)
	}
}

func importURL(ctx context.Context, url string) (rel.Expr, error) {
	if isRunningBundle(ctx) {
		url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")
		return fileValue(ctx, path.Join(ModuleDir, url))
	}
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err = bundleRemoteFile(ctx, url, data); err != nil {
			return nil, err
		}
		val, err := cache.getOrAdd(url, func() (rel.Expr, error) { return bytesValue(ctx, NoPath, data) })
		return val, err
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(ctx context.Context, filename string) (rel.Expr, error) {
	if filepath.Ext(filename) == "" {
		filename += ".arrai"
	}

	bytes, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), filename)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(filename) {
	case ".json":
		return bytesJSONToArrai(bytes)
	case ".yml", ".yaml":
		return translate.BytesYamlToArrai(bytes)
	}
	return bytesValue(ctx, filename, bytes)
}

func bytesValue(ctx context.Context, filename string, data []byte) (rel.Expr, error) {
	compile := func() (rel.Expr, error) {
		return Compile(ctx, filename, string(data))
	}
	if filename != NoPath {
		return cache.getOrAdd(filename, compile)
	}
	return compile()
}
