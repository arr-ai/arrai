package syntax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
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
	cmd := exec.Command("go", "list", "-m", "-json", "all") //nolint:gosec
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

// seedRequiredModules marks the memoized go.mod graph for moduleRoot as already
// loaded. Tests only.
func seedRequiredModules(moduleRoot string, mods map[string]requiredModule) {
	requiredModulesByRoot.Store(moduleRoot, mods)
}

// retrieveModule resolves importPath to a local module directory.
// moduleRoot is the importing project's go.mod directory (may be empty to use
// the process working directory). If that go.mod already requires a matching
// module, its pinned version and effective Dir (honouring replace) are used.
// Otherwise importPath prefixes are tried at the given version (or "latest").
func retrieveModule(importPath, version, moduleRoot string) (*goModule, error) {
	if version == "" {
		if pinned, ok := requiredModuleOf(moduleRoot, importPath); ok {
			if pinned.Dir != "" {
				return &goModule{Name: pinned.Path, Dir: pinned.Dir}, nil
			}
			if pinned.Version != "" {
				if m, err := downloadModule(pinned.Path, pinned.Version); err == nil {
					return m, nil
				}
			}
			// Fall through to latest-version resolution below.
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
