//go:build unit

package numgo

import (
	"math"
	"testing"
)

// --- product() tests ---

func TestProductBasic(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want int
	}{
		{"single", []int{5}, 5},
		{"two", []int{3, 4}, 12},
		{"three", []int{2, 3, 5}, 30},
		{"one element is 1", []int{1, 7, 1}, 7},
		{"all ones", []int{1, 1, 1, 1}, 1},
		{"empty slice", []int{}, 1},
		{"large valid", []int{1000, 1000, 100}, 100_000_000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := product(tc.vals)
			if got != tc.want {
				t.Errorf("product(%v) = %d, want %d", tc.vals, got, tc.want)
			}
		})
	}
}

func TestProductZeroDimension(t *testing.T) {
	cases := []struct {
		name string
		vals []int
	}{
		{"single zero", []int{0}},
		{"zero first", []int{0, 5, 3}},
		{"zero middle", []int{4, 0, 3}},
		{"zero last", []int{4, 5, 0}},
		{"all zeros", []int{0, 0, 0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := product(tc.vals)
			if got != 0 {
				t.Errorf("product(%v) = %d, want 0", tc.vals, got)
			}
		})
	}
}

func TestProductPanicsOnNegative(t *testing.T) {
	cases := []struct {
		name string
		vals []int
	}{
		{"single negative", []int{-1}},
		{"negative first", []int{-3, 4, 5}},
		{"negative middle", []int{3, -4, 5}},
		{"negative last", []int{3, 4, -5}},
		{"large negative", []int{-math.MaxInt}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("product(%v) did not panic on negative dimension", tc.vals)
				}
			}()
			product(tc.vals)
		})
	}
}

func TestProductPanicsOnOverflow(t *testing.T) {
	cases := []struct {
		name string
		vals []int
	}{
		{"two large", []int{math.MaxInt, 2}},
		{"sqrt overflow", []int{1 << 33, 1 << 33}},
		{"three-way overflow", []int{1 << 22, 1 << 22, 1 << 22}},
		{"gradual overflow", []int{1_000_000_000, 1_000_000_000, 1_000_000_000}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("product(%v) did not panic on overflow", tc.vals)
				}
			}()
			product(tc.vals)
		})
	}
}

// --- productUnsafe() tests ---

func TestProductUnsafeBasic(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want int
	}{
		{"single", []int{5}, 5},
		{"two", []int{3, 4}, 12},
		{"three", []int{2, 3, 5}, 30},
		{"empty slice", []int{}, 1},
		{"with zero", []int{3, 0, 5}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := productUnsafe(tc.vals)
			if got != tc.want {
				t.Errorf("productUnsafe(%v) = %d, want %d", tc.vals, got, tc.want)
			}
		})
	}
}

func TestProductUnsafeMatchesProductForValidInputs(t *testing.T) {
	// For all valid (positive, non-overflowing) inputs, both must agree.
	cases := [][]int{
		{1},
		{2, 3},
		{4, 5, 6},
		{1, 1, 1, 1, 1},
		{100, 200, 50},
		{7},
		{10, 10, 10, 10},
	}
	for _, vals := range cases {
		safe := product(vals)
		unsafe := productUnsafe(vals)
		if safe != unsafe {
			t.Errorf("product(%v) = %d, productUnsafe(%v) = %d — mismatch", vals, safe, vals, unsafe)
		}
	}
}

func TestProductUnsafeAllowsNegative(t *testing.T) {
	// productUnsafe does not validate — negative values just multiply through.
	got := productUnsafe([]int{3, -2})
	if got != -6 {
		t.Errorf("productUnsafe([3, -2]) = %d, want -6", got)
	}
}

func TestProductUnsafeNoOverflowCheck(t *testing.T) {
	// productUnsafe silently wraps on overflow — verify it doesn't panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("productUnsafe panicked on overflow, but should not: %v", r)
		}
	}()
	_ = productUnsafe([]int{math.MaxInt, 2})
}
