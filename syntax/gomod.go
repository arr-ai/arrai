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

// goModule holds the result of `go mod download -json`.
type goModule struct {
	Name string
	Dir  string
}

// extractVersion splits "module@version" into ("module", "version").
// If there is no "@", version is empty.
func extractVersion(path string) (module, version string) {
	module, version, _ = strings.Cut(path, "@")
	return
}

var (
	requiredModuleVersionsOnce sync.Once
	requiredModuleVersions     map[string]string
)

// requiredModuleVersionOf returns the version already pinned for importPath (or
// an ancestor package path of it) in the enclosing project's go.mod/go.sum, as
// reported by `go list -m -json all`. It returns false if the enclosing
// project has no go.mod, or no such requirement exists.
func requiredModuleVersionOf(importPath string) (modPath, version string, ok bool) {
	requiredModuleVersionsOnce.Do(func() {
		requiredModuleVersions = loadRequiredModuleVersions()
	})
	for path, ver := range requiredModuleVersions {
		if (path == importPath || strings.HasPrefix(importPath, path+"/")) && len(path) > len(modPath) {
			modPath, version, ok = path, ver, true
		}
	}
	return
}

// loadRequiredModuleVersions runs `go list -m -json all` in the working
// directory and returns the versions already required by its go.mod, keyed
// by module path. It returns an empty map if there is no enclosing module or
// the module graph can't be resolved.
func loadRequiredModuleVersions() map[string]string {
	versions := map[string]string{}
	cmd := exec.Command("go", "list", "-m", "-json", "all") //nolint:gosec
	out, err := cmd.Output()
	if err != nil {
		return versions
	}
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var m struct {
			Path    string
			Version string
			Main    bool
		}
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return versions
		}
		if !m.Main && m.Version != "" {
			versions[m.Path] = m.Version
		}
	}
	return versions
}

// retrieveModule downloads a Go module and returns its local directory.
// If the enclosing project's go.mod already requires a module matching
// importPath, that pinned version is used. Otherwise, since importPath may
// include a file path within the module (e.g. "github.com/org/repo/file.arrai"),
// it tries progressively shorter prefixes at the given version (or "latest"
// when none is given) until it finds a valid module.
func retrieveModule(importPath, version string) (*goModule, error) {
	if version == "" {
		if modPath, pinned, ok := requiredModuleVersionOf(importPath); ok {
			if m, err := downloadModule(modPath, pinned); err == nil {
				return m, nil
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
