//go:build unit

package tabgo

import (
	"testing"
)

func TestCatCategories(t *testing.T) {
	s := NewSeries("x", []any{"a", "b", "a", "c", "b", nil})
	cats := s.Cat().Categories()
	if len(cats) != 3 {
		t.Fatalf("Categories() length = %d, want 3", len(cats))
	}
	// Sorted: a, b, c
	if cats[0] != "a" || cats[1] != "b" || cats[2] != "c" {
		t.Errorf("Categories() = %v, want [a b c]", cats)
	}
}

func TestCatCodes(t *testing.T) {
	s := NewSeries("x", []any{"b", "a", "c", "a", nil})
	codes := s.Cat().Codes()
	vals := codes.Values()
	// Categories sorted: a=0, b=1, c=2
	if vals[0] != 1 || vals[1] != 0 || vals[2] != 2 || vals[3] != 0 || vals[4] != -1 {
		t.Errorf("Codes() = %v, want [1 0 2 0 -1]", vals)
	}
}

func TestCatRenameCategories(t *testing.T) {
	s := NewSeries("x", []any{"a", "b", "c"})
	renamed := s.Cat().RenameCategories(map[string]any{"a": "alpha", "b": "beta"})
	vals := renamed.Values()
	if vals[0] != "alpha" || vals[1] != "beta" || vals[2] != "c" {
		t.Errorf("RenameCategories() = %v", vals)
	}
}

func TestCatRemoveCategories(t *testing.T) {
	s := NewSeries("x", []any{"a", "b", "c", "a"})
	result := s.Cat().RemoveCategories([]any{"b"})
	vals := result.Values()
	if vals[0] != "a" || vals[1] != nil || vals[2] != "c" || vals[3] != "a" {
		t.Errorf("RemoveCategories() = %v", vals)
	}
}

func TestCatOrdered(t *testing.T) {
	s := NewSeries("x", []any{"a", "b"})
	if s.Cat().Ordered() {
		t.Error("Ordered() should return false")
	}
}
