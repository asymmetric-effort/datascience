//go:build unit

package tabgo

import (
	"testing"
)

func TestSeriesDtype(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		want   string
	}{
		{"empty", nil, "empty"},
		{"float64", []any{1.0, 2.0, 3.0}, "float64"},
		{"int", []any{1, 2, 3}, "int"},
		{"string", []any{"a", "b", "c"}, "string"},
		{"bool", []any{true, false, true}, "bool"},
		{"mixed int float", []any{1, 2.0, 3}, "float64"},
		{"mixed", []any{1, "a", true}, "mixed"},
		{"all nil", []any{nil, nil}, "empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeries("test", tt.values)
			got := s.Dtype()
			if got != tt.want {
				t.Errorf("Dtype() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSeriesShape(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3})
	shape := s.Shape()
	if shape != [1]int{3} {
		t.Errorf("Shape() = %v, want [3]", shape)
	}
}

func TestSeriesEmpty(t *testing.T) {
	s := NewSeries("x", nil)
	if !s.Empty() {
		t.Error("Empty() should return true for nil values")
	}
	s2 := NewSeries("x", []any{1})
	if s2.Empty() {
		t.Error("Empty() should return false for non-empty series")
	}
}

func TestSeriesIndex(t *testing.T) {
	s := NewSeries("x", []any{10, 20, 30})
	idx := s.Index()
	if len(idx) != 3 {
		t.Fatalf("Index() length = %d, want 3", len(idx))
	}
	for i, v := range idx {
		if v != i {
			t.Errorf("Index()[%d] = %d, want %d", i, v, i)
		}
	}
}

func TestSeriesHead(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3, 4, 5})
	h := s.Head(3)
	if h.Len() != 3 {
		t.Fatalf("Head(3).Len() = %d, want 3", h.Len())
	}
	vals := h.Values()
	if vals[0] != 1 || vals[2] != 3 {
		t.Errorf("Head values = %v, want [1 2 3]", vals)
	}
	// Head larger than len
	h2 := s.Head(10)
	if h2.Len() != 5 {
		t.Errorf("Head(10).Len() = %d, want 5", h2.Len())
	}
}

func TestSeriesTail(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3, 4, 5})
	tl := s.Tail(3)
	if tl.Len() != 3 {
		t.Fatalf("Tail(3).Len() = %d, want 3", tl.Len())
	}
	vals := tl.Values()
	if vals[0] != 3 || vals[2] != 5 {
		t.Errorf("Tail values = %v, want [3 4 5]", vals)
	}
}

func TestSeriesLoc(t *testing.T) {
	s := NewSeries("x", []any{10, 20, 30, 40, 50})
	loc := s.Loc([]int{1, 3})
	if loc.Len() != 2 {
		t.Fatalf("Loc length = %d, want 2", loc.Len())
	}
	vals := loc.Values()
	if vals[0] != 20 || vals[1] != 40 {
		t.Errorf("Loc values = %v, want [20 40]", vals)
	}
}

func TestSeriesIloc(t *testing.T) {
	s := NewSeries("x", []any{10, 20, 30})
	iloc := s.Iloc([]int{0, 2})
	vals := iloc.Values()
	if vals[0] != 10 || vals[1] != 30 {
		t.Errorf("Iloc values = %v, want [10 30]", vals)
	}
}

func TestSeriesWhere(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3, 4, 5})
	cond := []bool{true, false, true, false, true}
	result := s.Where(cond, 0)
	vals := result.Values()
	if vals[0] != 1 || vals[1] != 0 || vals[2] != 3 || vals[3] != 0 || vals[4] != 5 {
		t.Errorf("Where values = %v, want [1 0 3 0 5]", vals)
	}
}

func TestSeriesMask(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3, 4, 5})
	cond := []bool{true, false, true, false, true}
	result := s.Mask(cond, -1)
	vals := result.Values()
	if vals[0] != -1 || vals[1] != 2 || vals[2] != -1 || vals[3] != 4 || vals[4] != -1 {
		t.Errorf("Mask values = %v, want [-1 2 -1 4 -1]", vals)
	}
}

func TestSeriesAstype(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3})
	f := s.Astype("float64")
	vals := f.Values()
	if v, ok := vals[0].(float64); !ok || v != 1.0 {
		t.Errorf("Astype float64: got %v (%T)", vals[0], vals[0])
	}

	str := s.Astype("string")
	svals := str.Values()
	if svals[0] != "1" {
		t.Errorf("Astype string: got %v", svals[0])
	}
}

func TestSeriesCumsum(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	cs := s.Cumsum()
	vals := cs.Values()
	expected := []float64{1, 3, 6, 10}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cumsum[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestSeriesCumprod(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	cp := s.Cumprod()
	vals := cp.Values()
	expected := []float64{1, 2, 6, 24}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cumprod[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestSeriesCummax(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.0, 5.0})
	cm := s.Cummax()
	vals := cm.Values()
	expected := []float64{3, 3, 4, 4, 5}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cummax[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestSeriesCummin(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.0, 5.0})
	cm := s.Cummin()
	vals := cm.Values()
	expected := []float64{3, 1, 1, 1, 1}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cummin[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestSeriesDiff(t *testing.T) {
	s := NewSeries("x", []any{1.0, 3.0, 6.0, 10.0})
	d := s.Diff(1)
	vals := d.Values()
	if vals[0] != nil {
		t.Errorf("Diff[0] = %v, want nil", vals[0])
	}
	if toFloat64(vals[1]) != 2.0 {
		t.Errorf("Diff[1] = %v, want 2", vals[1])
	}
	if toFloat64(vals[2]) != 3.0 {
		t.Errorf("Diff[2] = %v, want 3", vals[2])
	}
}

func TestSeriesPctChange(t *testing.T) {
	s := NewSeries("x", []any{10.0, 11.0, 12.1})
	pc := s.PctChange(1)
	vals := pc.Values()
	if vals[0] != nil {
		t.Errorf("PctChange[0] = %v, want nil", vals[0])
	}
	if v := toFloat64(vals[1]); v < 0.099 || v > 0.101 {
		t.Errorf("PctChange[1] = %v, want ~0.1", v)
	}
}

func TestSeriesNlargest(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.0, 5.0, 9.0})
	nl := s.Nlargest(3)
	vals := nl.Values()
	if len(vals) != 3 {
		t.Fatalf("Nlargest(3) length = %d, want 3", len(vals))
	}
	if toFloat64(vals[0]) != 9.0 || toFloat64(vals[1]) != 5.0 || toFloat64(vals[2]) != 4.0 {
		t.Errorf("Nlargest values = %v, want [9 5 4]", vals)
	}
}

func TestSeriesNsmallest(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.0, 5.0, 9.0})
	ns := s.Nsmallest(3)
	vals := ns.Values()
	if len(vals) != 3 {
		t.Fatalf("Nsmallest(3) length = %d, want 3", len(vals))
	}
	if toFloat64(vals[0]) != 1.0 || toFloat64(vals[1]) != 1.0 || toFloat64(vals[2]) != 3.0 {
		t.Errorf("Nsmallest values = %v, want [1 1 3]", vals)
	}
}
