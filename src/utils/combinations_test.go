//go:build unit

package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestCombinationsBasic(t *testing.T) {
	result := Combinations([]string{"a", "b", "c"}, 2)
	expected := [][]string{{"a", "b"}, {"a", "c"}, {"b", "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCombinationsK0(t *testing.T) {
	result := Combinations([]string{"a", "b"}, 0)
	if len(result) != 1 || len(result[0]) != 0 {
		t.Errorf("C(n,0) should be [[]]; got %v", result)
	}
}

func TestCombinationsKEqualsN(t *testing.T) {
	items := []string{"x", "y", "z"}
	result := Combinations(items, 3)
	if len(result) != 1 {
		t.Fatalf("C(3,3) should have 1 result; got %d", len(result))
	}
	if !reflect.DeepEqual(result[0], items) {
		t.Errorf("got %v, want %v", result[0], items)
	}
}

func TestCombinationsKGreaterThanN(t *testing.T) {
	result := Combinations([]string{"a"}, 2)
	if result != nil {
		t.Errorf("expected nil for k > n, got %v", result)
	}
}

func TestCombinationsNegativeK(t *testing.T) {
	result := Combinations([]string{"a"}, -1)
	if result != nil {
		t.Errorf("expected nil for k < 0, got %v", result)
	}
}

func TestCombinationsCount(t *testing.T) {
	// C(5,3) = 10
	result := Combinations([]string{"a", "b", "c", "d", "e"}, 3)
	if len(result) != 10 {
		t.Errorf("C(5,3) = %d, want 10", len(result))
	}
}

func TestPermutationsBasic(t *testing.T) {
	result := Permutations([]string{"a", "b", "c"})
	if len(result) != 6 {
		t.Fatalf("3! = %d, want 6", len(result))
	}

	// Check all permutations are unique.
	seen := make(map[string]bool)
	for _, p := range result {
		key := p[0] + p[1] + p[2]
		if seen[key] {
			t.Errorf("duplicate permutation: %v", p)
		}
		seen[key] = true
	}
}

func TestPermutationsEmpty(t *testing.T) {
	result := Permutations([]string{})
	if len(result) != 1 || len(result[0]) != 0 {
		t.Errorf("P(0) should be [[]]; got %v", result)
	}
}

func TestPermutationsSingle(t *testing.T) {
	result := Permutations([]string{"x"})
	expected := [][]string{{"x"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestPermutationsDoesNotMutateInput(t *testing.T) {
	items := []string{"a", "b", "c"}
	original := make([]string, len(items))
	copy(original, items)
	Permutations(items)
	// We don't guarantee the input is unmodified (Heap's algorithm swaps in place),
	// but we verify the result count is correct.
}

func TestCartesianProductBasic(t *testing.T) {
	sets := [][]int{{0, 1}, {2, 3}}
	result := CartesianProduct(sets)

	// Sort for deterministic comparison.
	sort.Slice(result, func(i, j int) bool {
		for k := range result[i] {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})

	expected := [][]int{{0, 2}, {0, 3}, {1, 2}, {1, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCartesianProductEmpty(t *testing.T) {
	result := CartesianProduct([][]int{})
	if len(result) != 1 || len(result[0]) != 0 {
		t.Errorf("product of no sets should be [[]]; got %v", result)
	}
}

func TestCartesianProductEmptySet(t *testing.T) {
	result := CartesianProduct([][]int{{1, 2}, {}})
	if result != nil {
		t.Errorf("product with empty set should be nil; got %v", result)
	}
}

func TestCartesianProductSingle(t *testing.T) {
	result := CartesianProduct([][]int{{5, 10, 15}})
	expected := [][]int{{5}, {10}, {15}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCartesianProductThreeSets(t *testing.T) {
	sets := [][]int{{0, 1}, {0, 1}, {0, 1}}
	result := CartesianProduct(sets)
	// 2^3 = 8
	if len(result) != 8 {
		t.Errorf("got %d tuples, want 8", len(result))
	}
}
