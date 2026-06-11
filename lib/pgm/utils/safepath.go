package utils

import "github.com/asymmetric-effort/datascience/internal/safepath"

// SafePath validates that a file path does not contain directory traversal
// sequences or other dangerous patterns.
func SafePath(path string) (string, error) {
	return safepath.Validate(path)
}
