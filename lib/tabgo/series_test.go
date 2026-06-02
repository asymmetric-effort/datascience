//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestNewSeries(t *testing.T) {
	data := []any{1, 2, 3}
	s := NewSeries("x", data)
	if s.Name() != "x" {
		t.Fatalf("Name() = %q, want %q", s.Name(), "x")
	}
	if s.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", s.Len())
	}
	// mutating original should not affect series
	data[0] = 99
	if s.Values()[0] == 99 {
		t.Fatal("NewSeries did not copy input data")
	}
}

func TestSeriesValues(t *testing.T) {
	s := NewSeries("v", []any{"a", "b"})
	vals := s.Values()
	vals[0] = "z"
	if s.Values()[0] == "z" {
		t.Fatal("Values() did not return a copy")
	}
}

func TestSeriesFloat64(t *testing.T) {
	tests := []struct {
		name string
		data []any
		want []float64
	}{
		{"float64", []any{1.5, 2.5}, []float64{1.5, 2.5}},
		{"int", []any{1, 2}, []float64{1, 2}},
		{"string", []any{"3.14", "2.71"}, []float64{3.14, 2.71}},
		{"mixed", []any{int64(10), float32(1.5)}, []float64{10, 1.5}},
		{"invalid string", []any{"abc"}, []float64{0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSeries("f", tc.data)
			got := s.Float64()
			if len(got) != len(tc.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tc.want))
			}
			for i := range got {
				diff := got[i] - tc.want[i]
				if diff < -0.01 || diff > 0.01 {
					t.Errorf("[%d] = %f, want %f", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestSeriesInt(t *testing.T) {
	tests := []struct {
		name string
		data []any
		want []int
	}{
		{"int", []any{1, 2, 3}, []int{1, 2, 3}},
		{"float64", []any{1.9, 2.1}, []int{1, 2}},
		{"string", []any{"42", "7"}, []int{42, 7}},
		{"invalid", []any{"abc"}, []int{0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSeries("i", tc.data)
			got := s.Int()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValueCounts(t *testing.T) {
	s := NewSeries("c", []any{"a", "b", "a", "c", "b", "a"})
	vc := s.ValueCounts()
	if vc["a"] != 3 || vc["b"] != 2 || vc["c"] != 1 {
		t.Fatalf("ValueCounts = %v", vc)
	}
}

func TestUnique(t *testing.T) {
	s := NewSeries("u", []any{3, 1, 2, 1, 3})
	u := s.Unique()
	if len(u) != 3 {
		t.Fatalf("Unique len = %d, want 3", len(u))
	}
	// order of first appearance: 3, 1, 2
	want := []any{3, 1, 2}
	if !reflect.DeepEqual(u, want) {
		t.Fatalf("Unique = %v, want %v", u, want)
	}
}

func TestNUnique(t *testing.T) {
	s := NewSeries("n", []any{"x", "y", "x"})
	if s.NUnique() != 2 {
		t.Fatalf("NUnique = %d, want 2", s.NUnique())
	}
}

func TestEmptySeries(t *testing.T) {
	s := NewSeries("empty", nil)
	if s.Len() != 0 {
		t.Fatalf("Len = %d, want 0", s.Len())
	}
	if len(s.Float64()) != 0 {
		t.Fatal("Float64 on empty should return empty slice")
	}
	if len(s.Unique()) != 0 {
		t.Fatal("Unique on empty should return empty slice")
	}
	if s.NUnique() != 0 {
		t.Fatal("NUnique on empty should be 0")
	}
}
