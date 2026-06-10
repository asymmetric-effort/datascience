package utils

import (
	"fmt"
)

// GenerateStateNames produces default state names ["s0", "s1", ..., "s{n-1}"].
func GenerateStateNames(cardinality int) []string {
	names := make([]string, cardinality)
	for i := 0; i < cardinality; i++ {
		names[i] = fmt.Sprintf("s%d", i)
	}
	return names
}

// ValidateStateNames checks that names has the expected length and contains no duplicates.
func ValidateStateNames(names []string, cardinality int) error {
	if len(names) != cardinality {
		return fmt.Errorf("expected %d state names, got %d", cardinality, len(names))
	}
	seen := make(map[string]struct{}, cardinality)
	for _, n := range names {
		if _, exists := seen[n]; exists {
			return fmt.Errorf("duplicate state name: %q", n)
		}
		seen[n] = struct{}{}
	}
	return nil
}

// StateIndex returns the index of name within names, or an error if not found.
func StateIndex(names []string, name string) (int, error) {
	for i, n := range names {
		if n == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("state name %q not found", name)
}

// StateNameMap builds a mapping from variable name to state name to index.
// variables lists the variable names; stateNames maps each variable to its
// ordered slice of state names.
func StateNameMap(variables []string, stateNames map[string][]string) map[string]map[string]int {
	result := make(map[string]map[string]int, len(variables))
	for _, v := range variables {
		names := stateNames[v]
		m := make(map[string]int, len(names))
		for i, n := range names {
			m[n] = i
		}
		result[v] = m
	}
	return result
}
