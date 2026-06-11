//go:build unit

package safepath

import (
	"strings"
	"testing"
)

func TestValidateAcceptsValidPaths(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"model.bif", "model.bif"},
		{"output/model.bif", "output/model.bif"},
		{"/tmp/model.bif", "/tmp/model.bif"},
		{"./model.bif", "model.bif"},
		{"a/b/c/d.json", "a/b/c/d.json"},
	}
	for _, tc := range cases {
		got, err := Validate(tc.input)
		if err != nil {
			t.Errorf("Validate(%q) returned error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("Validate(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestValidateRejectsTraversal(t *testing.T) {
	cases := []string{
		"../secret.txt",
		"../../etc/passwd",
		"foo/../../bar",
		"foo/../../../etc/shadow",
		"../",
		"..",
	}
	for _, path := range cases {
		_, err := Validate(path)
		if err == nil {
			t.Errorf("Validate(%q) should have rejected directory traversal", path)
			continue
		}
		if !strings.Contains(err.Error(), "directory traversal") {
			t.Errorf("Validate(%q) error = %v, want 'directory traversal'", path, err)
		}
	}
}

func TestValidateRejectsNullByte(t *testing.T) {
	_, err := Validate("model\x00.bif")
	if err == nil {
		t.Fatal("Validate should reject null bytes")
	}
	if !strings.Contains(err.Error(), "null byte") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCleansDots(t *testing.T) {
	got, err := Validate("./foo/../foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "foo/bar" {
		t.Errorf("got %q, want %q", got, "foo/bar")
	}
}
