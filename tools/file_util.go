package tools

import "os"

// FileExists returns true if a file is existing.
func FileExists(file string) bool {
	if len(file) == 0 {
		return false
	}

	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
