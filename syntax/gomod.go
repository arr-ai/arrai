package syntax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// goModule holds the result of module resolution (from go list or go mod download).
type goModule struct {
	Name string
	Dir  string
}

// requiredModule is one entry from `go list -m -json all` for a non-main module.
type requiredModule struct {
	Path    string
	Version string
	Dir     string // effective directory (honours replace)
}

// extractVersion splits "module@version" into ("module", "version").
// If there is no "@", version is empty.
func extractVersion(path string) (module, version string) {
	module, version, _ = strings.Cut(path, "@")
	return
}

// requiredModulesByRoot caches `go list -m -json all` results keyed by module root
// directory (empty string = process working directory).
var requiredModulesByRoot sync.Map // map[string]map[string]requiredModule

// resetRequiredModulesCache clears the memoized go.mod graphs. Tests only.
func resetRequiredModulesCache() {
	requiredModulesByRoot = sync.Map{}
}

// requiredModuleOf returns the longest-prefix module matching importPath from
// the go.mod graph rooted at moduleRoot. Dir is the effective module directory
// (including replace targets) when go list reported one.
func requiredModuleOf(moduleRoot, importPath string) (requiredModule, bool) {
	mods := loadRequiredModules(moduleRoot)
	var best requiredModule
	found := false
	for path, mod := range mods {
		if path == importPath || strings.HasPrefix(importPath, path+"/") {
			if !found || len(path) > len(best.Path) {
				best, found = mod, true
			}
		}
	}
	return best, found
}

// loadRequiredModules runs `go list -m -json all` in moduleRoot (or the process
// cwd when moduleRoot is empty) and returns modules keyed by path. Results are
// memoized per moduleRoot. Returns an empty map if there is no enclosing module
// or the module graph can't be resolved.
func loadRequiredModules(moduleRoot string) map[string]requiredModule {
	if cached, ok := requiredModulesByRoot.Load(moduleRoot); ok {
		return cached.(map[string]requiredModule) //nolint:forcetypeassert
	}

	mods := map[string]requiredModule{}
	restore, err := snapshotGoModFiles(moduleRoot)
	if err != nil {
		requiredModulesByRoot.Store(moduleRoot, mods)
		return mods
	}
	defer restore()

	// -mod=mod lets this self-heal a go.sum that's missing entries (e.g. a
	// module required but not needed to build the main module's packages),
	// matching how `go mod download` behaves. Under the default -mod=readonly,
	// `go list` would instead fail outright and we'd wrongly treat the module
	// as unpinned. snapshotGoModFiles/restore above undo any such edits once
	// we're done, so resolving an import never leaves go.mod/go.sum changed
	// in the caller's project.
	cmd := exec.Command("go", "list", "-mod=mod", "-m", "-json", "all") //nolint:gosec
	if moduleRoot != "" {
		cmd.Dir = moduleRoot
	}
	out, err := cmd.Output()
	if err != nil {
		requiredModulesByRoot.Store(moduleRoot, mods)
		return mods
	}
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var m struct {
			Path    string
			Version string
			Dir     string
			Main    bool
			Replace *struct {
				Path string
				Dir  string
			}
		}
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			requiredModulesByRoot.Store(moduleRoot, mods)
			return mods
		}
		if m.Main {
			continue
		}
		dir := m.Dir
		if m.Replace != nil && m.Replace.Dir != "" {
			dir = m.Replace.Dir
		}
		if m.Path != "" {
			mods[m.Path] = requiredModule{Path: m.Path, Version: m.Version, Dir: dir}
		}
	}
	requiredModulesByRoot.Store(moduleRoot, mods)
	return mods
}

// fileSnapshot captures whether a file existed and its content, so it can be
// put back exactly as found.
type fileSnapshot struct {
	path    string
	existed bool
	content []byte
}

func snapshotFile(path string) (fileSnapshot, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fileSnapshot{path: path}, nil
		}
		return fileSnapshot{}, err
	}
	return fileSnapshot{path: path, existed: true, content: content}, nil
}

// restore is called via defer with nothing to propagate a failure to; it's a
// best-effort revert of a self-heal `go list` may have made.
func (s fileSnapshot) restore() {
	if s.existed {
		_ = os.WriteFile(s.path, s.content, 0o600) //nolint:errcheck
		return
	}
	_ = os.Remove(s.path) //nolint:errcheck
}

// snapshotGoModFiles saves the current go.mod/go.sum in dir (the process cwd
// if dir is "") and returns a func that restores them exactly as found. `go
// list -mod=mod` can rewrite either file to self-heal missing entries; the
// caller's own project files shouldn't come out of an import resolution
// changed, so loadRequiredModules reverts them once it has what it needs.
func snapshotGoModFiles(dir string) (restore func(), err error) {
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	goMod, err := snapshotFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return nil, err
	}
	goSum, err := snapshotFile(filepath.Join(dir, "go.sum"))
	if err != nil {
		return nil, err
	}
	return func() {
		goMod.restore()
		goSum.restore()
	}, nil
}

// seedRequiredModules marks the memoized go.mod graph for moduleRoot as already
// loaded. Tests only.
func seedRequiredModules(moduleRoot string, mods map[string]requiredModule) {
	requiredModulesByRoot.Store(moduleRoot, mods)
}

// retrieveModule resolves importPath to a local module directory.
// moduleRoot is the importing project's go.mod directory (may be empty to use
// the process working directory). If that go.mod already requires a matching
// module, its pinned version and effective Dir (honouring replace) are used,
// and it is an error if that pinned version can't be resolved. Otherwise
// importPath prefixes are tried at the given version (or "latest").
func retrieveModule(importPath, version, moduleRoot string) (*goModule, error) {
	if version == "" {
		if pinned, ok := requiredModuleOf(moduleRoot, importPath); ok {
			if pinned.Dir != "" {
				return &goModule{Name: pinned.Path, Dir: pinned.Dir}, nil
			}
			if pinned.Version != "" {
				m, err := downloadModule(pinned.Path, pinned.Version)
				if err != nil {
					// go.mod pins a version for this module: report the failure
					// rather than silently resolving to a different (latest)
					// version below, which would defeat the pin.
					return nil, fmt.Errorf("go.mod requires %s@%s but it could not be downloaded: %w",
						pinned.Path, pinned.Version, err)
				}
				return m, nil
			}
		}
	}

	modPath := importPath
	var lastErr error
	for {
		parts := strings.Split(modPath, "/")
		if len(parts) < 3 {
			break
		}
		m, err := downloadModule(modPath, version)
		if err == nil {
			return m, nil
		}
		lastErr = err
		modPath = strings.Join(parts[:len(parts)-1], "/")
	}
	if lastErr != nil {
		if ee, ok := lastErr.(*exec.ExitError); ok {
			return nil, fmt.Errorf("go mod download %s: %s", importPath, ee.Stderr)
		}
		return nil, fmt.Errorf("go mod download %s: %w", importPath, lastErr)
	}
	return nil, fmt.Errorf("go mod download %s: module not found", importPath)
}

// downloadModule runs `go mod download -json` for modPath at version (or
// "latest" if version is empty) and returns the resulting module info.
func downloadModule(modPath, version string) (*goModule, error) {
	arg := modPath
	if version != "" {
		arg += "@" + version
	} else {
		arg += "@latest"
	}
	cmd := exec.Command("go", "mod", "download", "-json", arg) //nolint:gosec
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var result struct {
		Path string `json:"Path"`
		Dir  string `json:"Dir"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("go mod download %s: parsing output: %w", arg, err)
	}
	return &goModule{Name: result.Path, Dir: result.Dir}, nil
}
