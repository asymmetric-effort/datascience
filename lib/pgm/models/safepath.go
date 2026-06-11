package models

import "github.com/asymmetric-effort/datascience/internal/safepath"

// safePath validates that a file path does not contain directory traversal
// sequences or other dangerous patterns.
func safePath(path string) (string, error) {
	return safepath.Validate(path)
}
