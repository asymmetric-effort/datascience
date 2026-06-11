// Package safepath provides file path validation to prevent directory
// traversal attacks. It is shared across all libraries in this module.
package safepath

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Validate checks that a file path does not contain directory traversal
// sequences or other dangerous patterns. It returns the cleaned path or an
// error if the path is unsafe.
func Validate(path string) (string, error) {
	cleaned := filepath.Clean(path)

	// Reject paths containing ".." after cleaning.
	for _, part := range strings.Split(cleaned, string(filepath.Separator)) {
		if part == ".." {
			return "", fmt.Errorf("path contains directory traversal: %q", path)
		}
	}

	// Reject null bytes (can truncate paths in some syscalls).
	if strings.ContainsRune(path, 0) {
		return "", fmt.Errorf("path contains null byte: %q", path)
	}

	return cleaned, nil
}
