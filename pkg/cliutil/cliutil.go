package cliutil

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/arr-ai/arrai/tools"
)

// FileExists checks if a file in path exists and returns the CLI-appropriate error.
func FileExists(ctx context.Context, path string) error {
	if exists, err := tools.FileExists(ctx, path); err != nil {
		return err
	} else if !exists {
		if !strings.Contains(path, string([]rune{os.PathSeparator})) {
			return fmt.Errorf(`"%s": not a command and not found as a file in the current directory`, path)
		}
		return fmt.Errorf(`"%s": file not found`, path)
	}
	return nil
}
