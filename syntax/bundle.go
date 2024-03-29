package syntax

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"

	"github.com/anz-bank/pkg/mod"
	"github.com/spf13/afero"
	"github.com/spf13/afero/zipfs"
)

const windows = "windows"

type bundleConfig struct {
	mainRoot, mainFile string

	// absRootPath is only meant for bundling, to trim absolute filepaths.
	// It is not stored to the config file or used during running bundled
	// scripts.
	absRootPath string
}

// module that is currently being imported by a script
type moduleData struct {
	// moduleName is the name of the go module e.h. github.com/arr-ai/arrai
	moduleName string

	// modulePath is the location of the module. Modules are cached by go mod
	// e.g. /go/pkg/mod/github.com/arr-ai/arrai@v0.200.0/
	modulePath string
}

type bundleKey int

const (
	// ModuleDir contains all the module files in a bundled scripts.
	ModuleDir = "/module"

	// NoModuleDir contains all the module files without root in a bundled scripts.
	NoModuleDir = "/unnamed"

	// BundleConfig contains configurations to run a bundled scripts.
	BundleConfig = "/config.arrai"
	arraiExt     = ".arrai"

	// random name for scripts without modules. This value is used if a script with
	// no module is compiled into binary.
	unnamedModule = "unnamed.com/unnamed/unnamed"

	bundleFsKey bundleKey = iota
	bundleConfKey
	runBundleMode
	currentModule
)

var (
	rootModuleRE           = regexp.MustCompile("^module ([^\n]+)\n")
	errSentinelHasNoModule = errors.New("sentinel does not show module path")
)

func (b bundleConfig) String() string {
	return fmt.Sprintf("(main_root: %q, main_file: %q)", b.mainRoot, b.mainFile)
}

func createConfig(ctx context.Context) error {
	return ctxfs.ZipCreate(
		ctx, bundleFsKey,
		BundleConfig, []byte(fromBundleConfig(ctx).String()),
	)
}

// WithBundleRun adds necessary values to the context to allow running bundled arrai scripts.
func WithBundleRun(ctx context.Context, buf []byte) (context.Context, error) {
	r := bytes.NewReader(buf)
	z, err := zip.NewReader(r, r.Size())
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctxfs.SourceFsOnto(ctx, zipfs.New(z)), runBundleMode, true), nil
}

// GetMainBundleSource gets the path for the main arrai script in the bundled arrai scripts.
func GetMainBundleSource(ctx context.Context) (context.Context, []byte, string) {
	if !isRunningBundle(ctx) {
		return ctx, nil, ""
	}
	ctx = withBundledConfig(ctx)
	mainFile := bundleToValidPath(ctx, fromBundleConfig(ctx).mainFile)
	fs := ctxfs.SourceFsFrom(ctx)
	buf, err := afero.ReadFile(fs, mainFile)
	if err != nil {
		panic(fmt.Errorf("not bundled properly, main file not accessible: %s", err))
	}
	return ctx, buf, mainFile
}

// GetModuleFromBundle fetches the module path of the bundle from the bundle's buffer.
func GetModuleFromBundle(ctx context.Context, buf []byte) (context.Context, string, error) {
	ctx, err := WithBundleRun(ctx, buf)
	if err != nil {
		return ctx, "", err
	}
	ctx = withBundledConfig(ctx)
	return ctx, fromBundleConfig(ctx).mainRoot, nil
}

// OutputArraiz writes the zip binary to the provided writer.
func OutputArraiz(ctx context.Context, w io.Writer) error {
	if !isBundling(ctx) {
		return errors.New("cannot output bundled arrai because it is not bundling")
	}

	return ctxfs.OutputZip(ctx, bundleFsKey, w)
}

// withBundledConfig is used add bundled scripts configuration to the context.
func withBundledConfig(ctx context.Context) context.Context {
	if !isRunningBundle(ctx) {
		//FIXME: return error?
		return ctx
	}

	buf, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), BundleConfig)
	if err != nil {
		// config file not generated
		panic(err)
	}

	expr, err := Compile(ctx, BundleConfig, string(buf))
	if err != nil {
		// config file not generated properly
		panic(err)
	}

	val, err := expr.Eval(ctx, rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	t := val.(rel.Tuple)
	root := t.MustGet("main_root").String()
	if root == "{}" {
		root = unnamedModule
	}
	root = ctxfs.ToUnixPath(root)

	return context.WithValue(ctx, bundleConfKey, bundleConfig{
		mainRoot: root,
		mainFile: t.MustGet("main_file").String(),
	})
}

func isRunningBundle(ctx context.Context) bool {
	return ctx.Value(runBundleMode) != nil
}

func withBundleConfig(ctx context.Context, b bundleConfig) context.Context {
	return context.WithValue(ctx, bundleConfKey, b)
}

// fromBundleConfig is meant to be used for fetching configurations during bundling.
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
		mainFile := path.Join(NoModuleDir, filepath.Base(filePath))
		if err = ctxfs.ZipCreate(ctx, bundleFsKey, mainFile, source); err != nil {
			return ctx, err
		}

		ctx = withBundleConfig(ctx, bundleConfig{mainFile: mainFile, absRootPath: filepath.Dir(filePath)})
		return ctx, createConfig(ctx)
	}

	buf, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), filepath.Join(root, ModuleRootSentinel))
	if err != nil {
		return ctx, err
	}

	moduleRoot := rootModuleRE.FindStringSubmatch(string(buf))
	if moduleRoot == nil {
		//TODO: maybe treat it as no module?
		return ctx, errSentinelHasNoModule
	}

	err = ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(ModuleDir, moduleRoot[1], ModuleRootSentinel), buf)
	if err != nil {
		return ctx, err
	}

	filePathRelativeToSentinel := ctxfs.ToUnixPath(strings.TrimPrefix(filePath, root))
	mainPath := path.Join(ModuleDir, moduleRoot[1], filePathRelativeToSentinel)

	if err = ctxfs.ZipCreate(ctx, bundleFsKey, mainPath, source); err != nil {
		return ctx, err
	}

	ctx = withBundleConfig(ctx, bundleConfig{moduleRoot[1], mainPath, root})
	return ctx, createConfig(ctx)
}

func addModuleSentinel(ctx context.Context, rootPath string) (err error) {
	if !isBundling(ctx) {
		return nil
	}

	rootPath, err = filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	rootPath = filepath.Join(rootPath, ModuleRootSentinel)

	buf, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), rootPath)
	if err != nil {
		return err
	}
	var sentinelLocation string
	if isImportModule(ctx) {
		sentinelLocation, err = createModulePath(ctx, rootPath)
		if err != nil {
			return err
		}
	} else {
		rootPath = strings.TrimPrefix(rootPath, fromBundleConfig(ctx).absRootPath)
		sentinelLocation = path.Join(fromBundleConfig(ctx).mainRoot, rootPath)
	}

	pathInBundle := path.Join(ModuleDir, sentinelLocation)
	if exists, err := ctxfs.FileExists(ctx, bundleFsKey, pathInBundle); err != nil {
		return err
	} else if exists {
		return nil
	}

	return ctxfs.ZipCreate(ctx, bundleFsKey, pathInBundle, buf)
}

func createModulePath(ctx context.Context, filePath string) (modulePath string, err error) {
	if !filepath.IsAbs(filePath) {
		return "", fmt.Errorf("filepath is not absolute: %s", filePath)
	}
	filePath = ctxfs.ToUnixPath(filePath)
	currModule := getCurrentModule(ctx)

	relPath, err := filepath.Rel(currModule.modulePath, filePath)
	if err != nil {
		return "", err
	}
	return path.Join(currModule.moduleName, ctxfs.ToUnixPath(relPath)), nil
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
		filePath += arraiExt
	}

	source, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), filePath)
	if err != nil {
		return err
	}

	var dir string
	if isImportModule(ctx) {
		dir = ModuleDir
		filePath, err = createModulePath(ctx, filePath)
		if err != nil {
			return err
		}
	} else {
		config := fromBundleConfig(ctx)
		if config.mainRoot != "" {
			dir = ModuleDir
			filePath = path.Join(config.mainRoot, ctxfs.ToUnixPath(strings.TrimPrefix(filePath, config.absRootPath)))
		} else {
			dir = NoModuleDir
			filePath = ctxfs.ToUnixPath(strings.TrimPrefix(filePath, config.absRootPath))
		}
	}

	return ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(dir, filePath), source)
}

func bundleModule(ctx context.Context, relImportPath string, m *mod.Module) (context.Context, error) {
	if !isBundling(ctx) {
		return ctx, nil
	}

	if filepath.Ext(relImportPath) == "" {
		relImportPath += arraiExt
	}

	source, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), filepath.Join(m.Dir, relImportPath))
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, currentModule, moduleData{ctxfs.ToUnixPath(m.Name), ctxfs.ToUnixPath(m.Dir)})
	return ctx, ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(ModuleDir, m.Name, relImportPath), source)
}

func isImportModule(ctx context.Context) bool {
	return ctx.Value(currentModule) != nil
}

func getCurrentModule(ctx context.Context) moduleData {
	return ctx.Value(currentModule).(moduleData)
}

func bundleRemoteFile(ctx context.Context, url string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}

	//TODO: test with deep imports
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")
	return ctxfs.ZipCreate(ctx, bundleFsKey, path.Join(ModuleDir, url), source)
}

// since config file uses UNIX path, they need to be converted to windows path
// (without volume) to work with afero zipfs.
func bundleToValidPath(ctx context.Context, p string) string {
	if runtime.GOOS == windows && isRunningBundle(ctx) {
		p = strings.TrimPrefix(p, filepath.VolumeName(p))
		return strings.ReplaceAll(p, "/", "\\")
	}
	return p
}
