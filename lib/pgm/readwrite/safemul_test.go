//go:build unit

package readwrite

import (
	"strings"
	"testing"
)

func TestSafeParentConfigsValid(t *testing.T) {
	cases := []struct {
		name  string
		cards []int
		want  int
	}{
		{"empty", []int{}, 1},
		{"single", []int{3}, 3},
		{"two", []int{3, 4}, 12},
		{"three", []int{2, 3, 5}, 30},
		{"all ones", []int{1, 1, 1}, 1},
		{"typical BN", []int{2, 3, 2}, 12},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := safeParentConfigs(tc.cards)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("safeParentConfigs(%v) = %d, want %d", tc.cards, got, tc.want)
			}
		})
	}
}

func TestSafeParentConfigsRejectsNegativeCardinality(t *testing.T) {
	_, err := safeParentConfigs([]int{3, -1, 2})
	if err == nil {
		t.Fatal("expected error for negative cardinality")
	}
	if !strings.Contains(err.Error(), "invalid cardinality") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSafeParentConfigsRejectsZeroCardinality(t *testing.T) {
	_, err := safeParentConfigs([]int{3, 0, 2})
	if err == nil {
		t.Fatal("expected error for zero cardinality")
	}
	if !strings.Contains(err.Error(), "invalid cardinality") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSafeParentConfigsRejectsOverflow(t *testing.T) {
	// Use maxCPDElements itself as a cardinality — first multiply succeeds
	// (1 * maxCPDElements = maxCPDElements, at the limit), then the second
	// multiply (maxCPDElements * 2) triggers the overflow/limit check.
	// But since maxCPDElements * 2 > maxCPDElements, it hits the limit first.
	// To test actual int overflow, we'd need values near math.MaxInt which
	// also exceed the limit. The limit check catches the dangerous case
	// before overflow can occur — which is the correct behavior.
	_, err := safeParentConfigs([]int{maxCPDElements, 2})
	if err == nil {
		t.Fatal("expected error for exceeding limit")
	}
}

func TestSafeParentConfigsRejectsExceedingLimit(t *testing.T) {
	// 3^40 ≈ 12 trillion — within int range but exceeds maxCPDElements
	cards := make([]int, 40)
	for i := range cards {
		cards[i] = 3
	}
	_, err := safeParentConfigs(cards)
	if err == nil {
		t.Fatal("expected error for exceeding limit")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSafeParentConfigsAtBoundary(t *testing.T) {
	// Just under the limit should succeed
	_, err := safeParentConfigs([]int{10000, 10000})
	if err != nil {
		t.Fatalf("100M should be within limit: %v", err)
	}

	// Just over the limit should fail
	_, err = safeParentConfigs([]int{10001, 10000})
	if err == nil {
		t.Fatal("100.01M should exceed limit")
	}
}
