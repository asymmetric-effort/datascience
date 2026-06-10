//go:build unit

package utils

import (
	"testing"
)

func TestGenerateStateNames(t *testing.T) {
	tests := []struct {
		cardinality int
		expected    []string
	}{
		{0, []string{}},
		{1, []string{"s0"}},
		{3, []string{"s0", "s1", "s2"}},
		{5, []string{"s0", "s1", "s2", "s3", "s4"}},
	}
	for _, tc := range tests {
		names := GenerateStateNames(tc.cardinality)
		if len(names) != len(tc.expected) {
			t.Errorf("GenerateStateNames(%d): got len %d, want %d", tc.cardinality, len(names), len(tc.expected))
			continue
		}
		for i, n := range names {
			if n != tc.expected[i] {
				t.Errorf("GenerateStateNames(%d)[%d] = %q, want %q", tc.cardinality, i, n, tc.expected[i])
			}
		}
	}
}

func TestValidateStateNames(t *testing.T) {
	// valid
	if err := ValidateStateNames([]string{"a", "b", "c"}, 3); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wrong length
	if err := ValidateStateNames([]string{"a", "b"}, 3); err == nil {
		t.Error("expected error for wrong length, got nil")
	}

	// duplicates
	if err := ValidateStateNames([]string{"a", "b", "a"}, 3); err == nil {
		t.Error("expected error for duplicates, got nil")
	}

	// empty valid
	if err := ValidateStateNames([]string{}, 0); err != nil {
		t.Errorf("unexpected error for empty: %v", err)
	}
}

func TestStateIndex(t *testing.T) {
	names := []string{"low", "medium", "high"}

	idx, err := StateIndex(names, "medium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 {
		t.Errorf("got %d, want 1", idx)
	}

	idx, err = StateIndex(names, "low")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 0 {
		t.Errorf("got %d, want 0", idx)
	}

	_, err = StateIndex(names, "missing")
	if err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestStateNameMap(t *testing.T) {
	variables := []string{"X", "Y"}
	stateNames := map[string][]string{
		"X": {"low", "high"},
		"Y": {"a", "b", "c"},
	}

	result := StateNameMap(variables, stateNames)

	if len(result) != 2 {
		t.Fatalf("expected 2 variables, got %d", len(result))
	}

	// Check X
	xMap := result["X"]
	if xMap["low"] != 0 || xMap["high"] != 1 {
		t.Errorf("X map incorrect: %v", xMap)
	}

	// Check Y
	yMap := result["Y"]
	if yMap["a"] != 0 || yMap["b"] != 1 || yMap["c"] != 2 {
		t.Errorf("Y map incorrect: %v", yMap)
	}
}

func TestStateNameMapEmpty(t *testing.T) {
	result := StateNameMap(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
