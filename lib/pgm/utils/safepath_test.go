//go:build unit

package utils

import "testing"

func TestSafePathDelegates(t *testing.T) {
	// Verify SafePath delegates to internal/safepath.Validate correctly.
	got, err := SafePath("valid/path.bif")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "valid/path.bif" {
		t.Errorf("got %q, want %q", got, "valid/path.bif")
	}

	_, err = SafePath("../../etc/passwd")
	if err == nil {
		t.Fatal("expected error for traversal path")
	}
}
