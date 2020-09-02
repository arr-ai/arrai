package syntax

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/anz-bank/sysl/pkg/mod"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/spf13/afero/zipfs"
)

type bundleConfig struct {
	mainRoot, mainFile string

	// absRootPath is only meant for bundling, to trim absolute filepaths.
	// It is not stored to the config file or used during running bundled
	// scripts.
	absRootPath string
}

type bundleKey int

const (
	moduleDir   = "/module"
	noModuleDir = "/unnamed"
	configFile  = "/config.arrai"

	bundleFsKey bundleKey = iota
	bundleConfKey
	runBundleMode
)

var (
	rootModuleRE  = regexp.MustCompile("^module ([^\n]+)\n")
	bundledConfig bundleConfig
)

var compileBundledConfig sync.Once

func (b bundleConfig) String() string {
	return fmt.Sprintf("(main_root: %q, main_file: %q)", b.mainRoot, b.mainFile)
}

//TODO: use createConfig
func createConfig(ctx context.Context) error {
	return ctxfs.ZipCreate(
		ctx, bundleFsKey,
		configFile, []byte(fromBundleConfig(ctx).String()),
	)
}

// WithBundleRun adds necessary values to the context to allow running bundled arrai scripts.
func WithBundleRun(ctx context.Context, filePath string, buf []byte) (context.Context, error) {
	f, err := ctxfs.SourceFsFrom(ctx).Open(filePath)
	if err != nil {
		return ctx, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return ctx, err
	}

	z, err := zip.NewReader(bytes.NewReader(buf), fi.Size())
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctxfs.SourceFsOnto(ctx, zipfs.New(z)), runBundleMode, true), nil
}

// GetMainBundleSource gets the path for the main arrai script in the bundled arrai scripts.
func GetMainBundleSource(ctx context.Context) ([]byte, string) {
	if !isRunningBundle(ctx) {
		return nil, ""
	}
	mainFile := getBundledConfig(ctx).mainFile
	buf, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), mainFile)
	if err != nil {
		panic(fmt.Errorf("not bundled properly: %s", mainFile))
	}
	return buf, mainFile
}

// OutputArraiz writes the zip binary to the provided writer.
func OutputArraiz(ctx context.Context, w io.Writer) error {
	if !isBundling(ctx) {
		return errors.New("cannot output bundled arrai because it is not bundling")
	}

	return ctxfs.OutputZip(ctx, bundleFsKey, w)
}

func getBundledConfig(ctx context.Context) bundleConfig {
	if !isRunningBundle(ctx) {
		//FIXME: return error?
		return bundleConfig{}
	}
	//TODO: better error message
	compileBundledConfig.Do(func() {
		buf, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), configFile)
		if err != nil {
			// config file not generated
			panic(err)
		}

		expr, err := Compile(ctx, configFile, string(buf))
		if err != nil {
			// config file not generated properly
			panic(err)
		}

		val, err := expr.Eval(ctx, rel.EmptyScope)
		if err != nil {
			panic(err)
		}
		t := val.(rel.Tuple)
		bundledConfig = bundleConfig{
			mainRoot: t.MustGet("main_root").String(),
			mainFile: t.MustGet("main_file").String(),
		}
	})
	return bundledConfig
}

func isRunningBundle(ctx context.Context) bool {
	return ctx.Value(runBundleMode) != nil
}

func withBundleConfig(ctx context.Context, b bundleConfig) context.Context {
	return context.WithValue(ctx, bundleConfKey, b)
}

func fromBundleConfig(ctx context.Context) bundleConfig {
	return ctx.Value(bundleConfKey).(bundleConfig)
}

func isBundling(ctx context.Context) bool {
	return ctx.Value(bundleFsKey) != nil
}

// SetupBundle adds necessary values to the context for bundling arrai scripts.
func SetupBundle(ctx context.Context, filePath string, source []byte) (_ context.Context, err error) {
	ctx = ctxfs.WithZipFs(ctx, bundleFsKey)

	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return ctx, err
	}

	root, err := findRootFromModule(ctx, filepath.Dir(filePath))
	if err != nil {
		if err != errModuleNotExist {
			return ctx, err
		}
		mainFile := path.Join(noModuleDir, filepath.Base(filePath))
		if err = ctxfs.ZipCreate(ctx, bundleFsKey, mainFile, source); err != nil {
			return ctx, err
		}

		ctx = withBundleConfig(ctx, bundleConfig{mainFile: mainFile, absRootPath: filepath.Dir(filePath)})
		return ctx, createConfig(ctx)
	}
	f, err := ctxfs.SourceFsFrom(ctx).Open(filepath.Join(root, ModuleRootSentinel))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return ctx, err
	}

	moduleRoot := rootModuleRE.FindStringSubmatch(string(buf))
	if moduleRoot == nil {
		//TODO: maybe treat it as no module?
		return ctx, errors.New("sentinel does not show module path")
	}

	err = ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(moduleDir, moduleRoot[1], ModuleRootSentinel), buf)
	if err != nil {
		return ctx, err
	}

	filePathRelativeToSentinel := toUnixPath(strings.TrimPrefix(filePath, root))
	mainPath := path.Join(moduleDir, moduleRoot[1], filePathRelativeToSentinel)

	if err = ctxfs.ZipCreate(ctx, bundleFsKey, mainPath, source); err != nil {
		return ctx, err
	}

	ctx = withBundleConfig(ctx, bundleConfig{moduleRoot[1], mainPath, root})
	return ctx, createConfig(ctx)
}

//FIXME: handle nested root
func addLocalRoot(ctx context.Context, rootPath string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	if exists, err := ctxfs.FileExists(ctx, bundleFsKey, rootPath); err != nil {
		return err
	} else if !exists {
		return nil
	}

	return ctxfs.ZipCreate(ctx, bundleFsKey, toUnixPath(rootPath), source)
}

func bundleLocalFile(ctx context.Context, filePath string) (err error) {
	if !isBundling(ctx) {
		return nil
	}
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return err
	}

	//FIXME: not very clean
	if filepath.Ext(filePath) == "" {
		filePath += ".arrai"
	}

	f, err := ctxfs.SourceFsFrom(ctx).Open(filePath)
	if err != nil {
		return err
	}

	source, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var dir string
	config := fromBundleConfig(ctx)
	if config.mainRoot != "" {
		dir = moduleDir
		filePath = path.Join(config.mainRoot, toUnixPath(strings.TrimPrefix(filePath, config.absRootPath)))
	} else {
		dir = noModuleDir
		filePath = toUnixPath(strings.TrimPrefix(filePath, config.absRootPath))
	}

	return ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(dir, filePath), source)
}

func bundleModule(ctx context.Context, relImportPath string, m *mod.Module) error {
	if !isBundling(ctx) {
		return nil
	}

	if filepath.Ext(relImportPath) == "" {
		relImportPath += ".arrai"
	}

	f, err := ctxfs.SourceFsFrom(ctx).Open(filepath.Join(m.Dir, relImportPath))
	if err != nil {
		return err
	}

	source, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(moduleDir, m.Name, relImportPath), source)
}

func bundleRemoteFile(ctx context.Context, url string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}

	//TODO: test with deep imports
	url = strings.TrimPrefix("http://", strings.TrimPrefix("https://", url))
	return ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(moduleDir, url), source)
}

func toUnixPath(p string) string {
	if runtime.GOOS == "windows" {
		if vol := filepath.VolumeName(p); strings.HasPrefix(p, vol) {
			p = strings.TrimPrefix(p, vol)
		}
		return strings.ReplaceAll(p, "\\", "/")
	}
	return p
}
