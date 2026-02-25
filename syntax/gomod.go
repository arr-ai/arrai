package syntax

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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

// retrieveModule downloads a Go module and returns its local directory.
// Since importPath may include a file path within the module (e.g.
// "github.com/org/repo/file.arrai"), it tries progressively shorter
// prefixes until it finds a valid module.
func retrieveModule(importPath, version string) (*goModule, error) {
	modPath := importPath
	var lastErr error
	for {
		parts := strings.Split(modPath, "/")
		if len(parts) < 3 {
			break
		}
		arg := modPath
		if version != "" {
			arg += "@" + version
		} else {
			arg += "@latest"
		}
		cmd := exec.Command("go", "mod", "download", "-json", arg) //nolint:gosec
		out, err := cmd.Output()
		if err == nil {
			var result struct {
				Path string `json:"Path"`
				Dir  string `json:"Dir"`
			}
			if err := json.Unmarshal(out, &result); err != nil {
				return nil, fmt.Errorf("go mod download %s: parsing output: %w", arg, err)
			}
			return &goModule{Name: result.Path, Dir: result.Dir}, nil
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
