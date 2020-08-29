package tools

import (
	"context"
	"os"

	"github.com/arr-ai/arrai/pkg/ctxfs"
)

// FileExists returns true if a file is existing. This is meant to be used with
// SourceFs related operations.
func FileExists(ctx context.Context, file string) (bool, error) {
	if len(file) == 0 {
		return false, nil
	}
	info, err := ctxfs.SourceFsFrom(ctx).Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}
