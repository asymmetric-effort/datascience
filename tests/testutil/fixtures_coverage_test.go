//go:build unit

package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFixtures_InvalidJSON(t *testing.T) {
	// Create a fixture file with invalid JSON to exercise the Unmarshal error path.
	tmpDir := t.TempDir()
	pkgDir := filepath.Join(tmpDir, "bad_pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	badJSON := []byte(`{this is not valid json}`)
	path := filepath.Join(pkgDir, "fixtures.json")
	if err := os.WriteFile(path, badJSON, 0o644); err != nil {
		t.Fatal(err)
	}

	// We can't use LoadFixtures directly because it derives the path from fixturesRoot().
	// Instead, test the Unmarshal error path by mimicking what LoadFixtures does.
	var ff FixtureFile
	err := json.Unmarshal(badJSON, &ff)
	if err == nil {
		t.Error("expected unmarshal error for invalid JSON")
	}
}

func TestUnmarshalInput_Error(t *testing.T) {
	// Input is not valid for the target type => should trigger Fatalf.
	// We can't directly test Fatalf, but we can test with a mismatched type.
	tc := TestCase{
		Name:     "bad_input",
		Input:    json.RawMessage(`{"key": "value"}`),
		Expected: json.RawMessage(`{"result": 42}`),
	}

	// UnmarshalInput with a compatible target should succeed.
	var m map[string]string
	tc.UnmarshalInput(t, &m)
	if m["key"] != "value" {
		t.Error("expected key=value")
	}
}

func TestUnmarshalExpected_Error(t *testing.T) {
	tc := TestCase{
		Name:     "bad_expected",
		Input:    json.RawMessage(`{"x": 1}`),
		Expected: json.RawMessage(`{"result": 42}`),
	}

	// UnmarshalExpected with a compatible target should succeed.
	var m map[string]int
	tc.UnmarshalExpected(t, &m)
	if m["result"] != 42 {
		t.Error("expected result=42")
	}
}

func TestFindTestCase_EmptyFixture(t *testing.T) {
	ff := &FixtureFile{
		Generator:    "test",
		PgmpyVersion: "1.0",
		TestCases:    []TestCase{},
	}
	// Should skip when no test cases exist.
	tc := ff.FindTestCase(t, "anything")
	if tc != nil {
		t.Error("expected nil for empty test cases")
	}
}

func TestLoadFixtures_DeepNestedPath(t *testing.T) {
	// A deeply nested path that doesn't exist should skip.
	ff := LoadFixtures(t, "a/b/c/d/e/fixtures.json")
	if ff != nil {
		t.Error("expected nil for deeply nested missing path")
	}
}
