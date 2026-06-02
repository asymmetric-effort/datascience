//go:build unit

package tabgo

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewMultiIndex(t *testing.T) {
	levels := [][]any{{"a", "b"}, {1, 2, 3}}
	codes := [][]int{{0, 0, 1, 1}, {0, 1, 2, 0}}
	names := []string{"letters", "numbers"}

	mi, err := NewMultiIndex(levels, codes, names)
	if err != nil {
		t.Fatalf("NewMultiIndex error: %v", err)
	}
	if mi.Len() != 4 {
		t.Fatalf("Len() = %d, want 4", mi.Len())
	}
	if mi.Nlevels() != 2 {
		t.Fatalf("Nlevels() = %d, want 2", mi.Nlevels())
	}
}

func TestNewMultiIndex_Validation(t *testing.T) {
	// No levels.
	_, err := NewMultiIndex(nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for empty levels")
	}

	// Mismatched levels/codes count.
	_, err = NewMultiIndex([][]any{{"a"}}, [][]int{{0}, {0}}, []string{"x"})
	if err == nil {
		t.Fatal("expected error for mismatched levels/codes count")
	}

	// Mismatched names count.
	_, err = NewMultiIndex([][]any{{"a"}}, [][]int{{0}}, []string{"x", "y"})
	if err == nil {
		t.Fatal("expected error for mismatched names count")
	}

	// Unequal code lengths.
	_, err = NewMultiIndex([][]any{{"a"}, {"b"}}, [][]int{{0}, {0, 0}}, []string{"x", "y"})
	if err == nil {
		t.Fatal("expected error for unequal code lengths")
	}

	// Code out of range.
	_, err = NewMultiIndex([][]any{{"a"}}, [][]int{{5}}, []string{"x"})
	if err == nil {
		t.Fatal("expected error for out-of-range code")
	}
}

func TestMultiIndexNames(t *testing.T) {
	mi, _ := NewMultiIndex([][]any{{"a"}}, [][]int{{0}}, []string{"lvl"})
	names := mi.Names()
	if len(names) != 1 || names[0] != "lvl" {
		t.Fatalf("Names() = %v, want [lvl]", names)
	}
	// Mutating returned slice should not affect original.
	names[0] = "changed"
	if mi.Names()[0] != "lvl" {
		t.Fatal("Names() did not return a copy")
	}
}

func TestMultiIndexLevels(t *testing.T) {
	mi, _ := NewMultiIndex([][]any{{"a", "b"}, {10, 20}}, [][]int{{0, 1}, {0, 1}}, []string{"x", "y"})
	levels := mi.Levels()
	if len(levels) != 2 {
		t.Fatalf("Levels() returned %d levels, want 2", len(levels))
	}
	if fmt.Sprintf("%v", levels[0]) != "[a b]" {
		t.Fatalf("Levels()[0] = %v, want [a b]", levels[0])
	}
	// Mutating returned slice should not affect original.
	levels[0][0] = "z"
	if mi.Levels()[0][0] == "z" {
		t.Fatal("Levels() did not return a copy")
	}
}

func TestMultiIndexGetLevel(t *testing.T) {
	mi, _ := NewMultiIndex([][]any{{"a", "b"}, {1, 2}}, [][]int{{0, 1}, {0, 1}}, []string{"letters", "numbers"})

	lv, err := mi.GetLevel("letters")
	if err != nil {
		t.Fatalf("GetLevel error: %v", err)
	}
	if fmt.Sprintf("%v", lv) != "[a b]" {
		t.Fatalf("GetLevel(letters) = %v, want [a b]", lv)
	}

	_, err = mi.GetLevel("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent level name")
	}
}

func TestMultiIndexGet(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1, 0}, {0, 1, 1}},
		[]string{"x", "y"},
	)
	got := mi.Get(0)
	if fmt.Sprintf("%v", got) != "[a 1]" {
		t.Fatalf("Get(0) = %v, want [a 1]", got)
	}
	got = mi.Get(2)
	if fmt.Sprintf("%v", got) != "[a 2]" {
		t.Fatalf("Get(2) = %v, want [a 2]", got)
	}
	// Out of range.
	if mi.Get(-1) != nil {
		t.Fatal("Get(-1) should return nil")
	}
	if mi.Get(100) != nil {
		t.Fatal("Get(100) should return nil")
	}
}

func TestMultiIndexSet(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)

	// Set to existing value.
	err := mi.Set(0, []any{"b", 2})
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}
	got := mi.Get(0)
	if fmt.Sprintf("%v", got) != "[b 2]" {
		t.Fatalf("after Set, Get(0) = %v, want [b 2]", got)
	}

	// Set to new value (not previously in levels).
	err = mi.Set(1, []any{"c", 3})
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}
	got = mi.Get(1)
	if fmt.Sprintf("%v", got) != "[c 3]" {
		t.Fatalf("after Set new, Get(1) = %v, want [c 3]", got)
	}

	// Error: wrong number of values.
	err = mi.Set(0, []any{"a"})
	if err == nil {
		t.Fatal("expected error for wrong number of values")
	}

	// Error: out of range.
	err = mi.Set(99, []any{"a", 1})
	if err == nil {
		t.Fatal("expected error for out-of-range row")
	}
}

func TestMultiIndexEquals(t *testing.T) {
	mi1, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	mi2, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	if !mi1.Equals(mi2) {
		t.Fatal("expected Equals to return true for identical indices")
	}

	// Different values.
	mi3, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{1, 0}, {0, 1}},
		[]string{"x", "y"},
	)
	if mi1.Equals(mi3) {
		t.Fatal("expected Equals to return false for different indices")
	}

	// Different names.
	mi4, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "z"},
	)
	if mi1.Equals(mi4) {
		t.Fatal("expected Equals to return false for different names")
	}

	// Nil comparison.
	if mi1.Equals(nil) {
		t.Fatal("expected Equals(nil) to return false")
	}
}

func TestMultiIndexCopy(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	cp := mi.Copy()

	if !mi.Equals(cp) {
		t.Fatal("Copy should be equal to original")
	}

	// Mutating copy should not affect original.
	_ = cp.Set(0, []any{"b", 2})
	if mi.Equals(cp) {
		t.Fatal("mutating copy should not affect original")
	}
}

func TestMultiIndexString(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	s := mi.String()
	if !strings.Contains(s, "MultiIndex") {
		t.Fatalf("String() should contain 'MultiIndex', got: %s", s)
	}
	if !strings.Contains(s, "names=") {
		t.Fatalf("String() should contain 'names=', got: %s", s)
	}
}

func TestMultiIndexStringEmpty(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a"}},
		[][]int{{}},
		[]string{"x"},
	)
	// Len is 0 since codes are empty.
	if mi.Len() != 0 {
		t.Fatalf("expected Len()=0, got %d", mi.Len())
	}
	s := mi.String()
	if s != "MultiIndex([])" {
		t.Fatalf("String() for empty = %q, want MultiIndex([])", s)
	}
}

func TestMultiIndexStringManyRows(t *testing.T) {
	// More than 10 rows should show truncation.
	levels := [][]any{make([]any, 0)}
	for i := 0; i < 15; i++ {
		levels[0] = append(levels[0], i)
	}
	codes := [][]int{make([]int, 15)}
	for i := 0; i < 15; i++ {
		codes[0][i] = i
	}
	mi, _ := NewMultiIndex(levels, codes, []string{"n"})
	s := mi.String()
	if !strings.Contains(s, "... (5 more)") {
		t.Fatalf("String() should contain truncation message, got: %s", s)
	}
}

func TestNewMultiIndexFromArrays(t *testing.T) {
	arrays := [][]any{
		{"a", "a", "b", "b"},
		{1, 2, 1, 2},
	}
	mi, err := NewMultiIndexFromArrays(arrays, []string{"letters", "numbers"})
	if err != nil {
		t.Fatalf("NewMultiIndexFromArrays error: %v", err)
	}
	if mi.Len() != 4 {
		t.Fatalf("Len() = %d, want 4", mi.Len())
	}
	got := mi.Get(0)
	if fmt.Sprintf("%v", got) != "[a 1]" {
		t.Fatalf("Get(0) = %v, want [a 1]", got)
	}
}

func TestNewMultiIndexFromArraysNilNames(t *testing.T) {
	arrays := [][]any{{"x", "y"}}
	mi, err := NewMultiIndexFromArrays(arrays, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if mi.Names()[0] != "level_0" {
		t.Fatalf("expected auto-generated name, got %q", mi.Names()[0])
	}
}

func TestMultiIndexDroplevel(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	dropped, err := mi.Droplevel(0)
	if err != nil {
		t.Fatalf("Droplevel error: %v", err)
	}
	if dropped.Nlevels() != 1 {
		t.Fatalf("after Droplevel, Nlevels() = %d, want 1", dropped.Nlevels())
	}

	// Cannot drop only remaining level.
	_, err = dropped.Droplevel(0)
	if err == nil {
		t.Fatal("expected error when dropping the only level")
	}
}

func TestMultiIndexSortlevel(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"c", "a", "b"}, {3, 1, 2}},
		[][]int{{0, 1, 2}, {0, 1, 2}},
		[]string{"x", "y"},
	)
	sorted, err := mi.Sortlevel(0)
	if err != nil {
		t.Fatalf("Sortlevel error: %v", err)
	}
	// After sorting by level 0: a, b, c
	got := sorted.Get(0)
	if fmt.Sprintf("%v", got[0]) != "a" {
		t.Fatalf("expected first element to be 'a', got %v", got[0])
	}
}

func TestMultiIndexGetLevelByIndex(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1}, {0, 1}},
		[]string{"x", "y"},
	)
	lv := mi.GetLevelByIndex(0)
	if fmt.Sprintf("%v", lv) != "[a b]" {
		t.Fatalf("GetLevelByIndex(0) = %v, want [a b]", lv)
	}
	if mi.GetLevelByIndex(-1) != nil {
		t.Fatal("GetLevelByIndex(-1) should return nil")
	}
	if mi.GetLevelByIndex(99) != nil {
		t.Fatal("GetLevelByIndex(99) should return nil")
	}
}

func TestMultiIndexGetLevelSeries(t *testing.T) {
	mi, _ := NewMultiIndex(
		[][]any{{"a", "b"}, {1, 2}},
		[][]int{{0, 1, 0}, {0, 1, 0}},
		[]string{"x", "y"},
	)
	s := mi.GetLevelSeries(0)
	if s.Len() != 3 {
		t.Fatalf("GetLevelSeries Len = %d, want 3", s.Len())
	}
	if s.Name() != "x" {
		t.Fatalf("GetLevelSeries Name = %q, want x", s.Name())
	}
	vals := s.Values()
	if fmt.Sprintf("%v", vals) != "[a b a]" {
		t.Fatalf("GetLevelSeries values = %v, want [a b a]", vals)
	}
}

func TestMultiIndexDeepCopyInput(t *testing.T) {
	levels := [][]any{{"a", "b"}}
	codes := [][]int{{0, 1}}
	names := []string{"x"}

	mi, _ := NewMultiIndex(levels, codes, names)

	// Mutate originals.
	levels[0][0] = "z"
	codes[0][0] = 1
	names[0] = "changed"

	if mi.Get(0)[0] == "z" {
		t.Fatal("NewMultiIndex did not deep copy levels")
	}
	if mi.Names()[0] == "changed" {
		t.Fatal("NewMultiIndex did not copy names")
	}
}
