//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestRollingMean(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}, {4.0}, {5.0}},
	)
	result := df.Rolling(3).Mean()
	vals := result.Column("x").Values()
	if vals[0] != nil || vals[1] != nil {
		t.Errorf("Rolling(3).Mean(): first 2 should be nil, got %v %v", vals[0], vals[1])
	}
	if v := toFloat64(vals[2]); v != 2.0 {
		t.Errorf("Rolling(3).Mean()[2] = %v, want 2.0", v)
	}
	if v := toFloat64(vals[3]); v != 3.0 {
		t.Errorf("Rolling(3).Mean()[3] = %v, want 3.0", v)
	}
	if v := toFloat64(vals[4]); v != 4.0 {
		t.Errorf("Rolling(3).Mean()[4] = %v, want 4.0", v)
	}
}

func TestRollingSum(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}, {4.0}},
	)
	result := df.Rolling(2).Sum()
	vals := result.Column("x").Values()
	if vals[0] != nil {
		t.Errorf("Rolling(2).Sum()[0] = %v, want nil", vals[0])
	}
	if v := toFloat64(vals[1]); v != 3.0 {
		t.Errorf("Rolling(2).Sum()[1] = %v, want 3.0", v)
	}
}

func TestRollingMin(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {1.0}, {4.0}, {1.0}},
	)
	result := df.Rolling(2).Min()
	vals := result.Column("x").Values()
	if v := toFloat64(vals[1]); v != 1.0 {
		t.Errorf("Rolling(2).Min()[1] = %v, want 1.0", v)
	}
}

func TestRollingMax(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {1.0}, {4.0}, {1.0}},
	)
	result := df.Rolling(2).Max()
	vals := result.Column("x").Values()
	if v := toFloat64(vals[1]); v != 3.0 {
		t.Errorf("Rolling(2).Max()[1] = %v, want 3.0", v)
	}
	if v := toFloat64(vals[2]); v != 4.0 {
		t.Errorf("Rolling(2).Max()[2] = %v, want 4.0", v)
	}
}

func TestExpandingMean(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}, {4.0}},
	)
	result := df.Expanding().Mean()
	vals := result.Column("x").Values()
	if v := toFloat64(vals[0]); v != 1.0 {
		t.Errorf("Expanding().Mean()[0] = %v, want 1.0", v)
	}
	if v := toFloat64(vals[1]); v != 1.5 {
		t.Errorf("Expanding().Mean()[1] = %v, want 1.5", v)
	}
	if v := toFloat64(vals[2]); v != 2.0 {
		t.Errorf("Expanding().Mean()[2] = %v, want 2.0", v)
	}
	if v := toFloat64(vals[3]); v != 2.5 {
		t.Errorf("Expanding().Mean()[3] = %v, want 2.5", v)
	}
}

func TestExpandingSum(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := df.Expanding().Sum()
	vals := result.Column("x").Values()
	if v := toFloat64(vals[2]); v != 6.0 {
		t.Errorf("Expanding().Sum()[2] = %v, want 6.0", v)
	}
}

func TestEWMMean(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := df.EWM(3).Mean()
	vals := result.Column("x").Values()
	// First value should equal input
	if v := toFloat64(vals[0]); v != 1.0 {
		t.Errorf("EWM(3).Mean()[0] = %v, want 1.0", v)
	}
	// Subsequent values should be weighted averages
	if v := toFloat64(vals[1]); v <= 1.0 || v >= 2.0 {
		t.Errorf("EWM(3).Mean()[1] = %v, expected between 1 and 2", v)
	}
}

func TestEWMStd(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}, {4.0}},
	)
	result := df.EWM(3).Std()
	vals := result.Column("x").Values()
	// First value should be 0
	if v := toFloat64(vals[0]); v != 0.0 {
		t.Errorf("EWM(3).Std()[0] = %v, want 0.0", v)
	}
	// Subsequent values should be positive
	if v := toFloat64(vals[1]); v <= 0 {
		t.Errorf("EWM(3).Std()[1] = %v, want > 0", v)
	}
}

func TestEWMVar(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := df.EWM(3).Var()
	vals := result.Column("x").Values()
	if v := toFloat64(vals[0]); v != 0.0 {
		t.Errorf("EWM(3).Var()[0] = %v, want 0.0", v)
	}
}

func TestCummaxDataFrame(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {1.0}, {4.0}, {1.0}, {5.0}},
	)
	result := Cummax(df)
	vals := result.Column("x").Values()
	expected := []float64{3, 3, 4, 4, 5}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cummax[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestCumminDataFrame(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {1.0}, {4.0}, {0.5}, {5.0}},
	)
	result := Cummin(df)
	vals := result.Column("x").Values()
	expected := []float64{3, 1, 1, 0.5, 0.5}
	for i, e := range expected {
		if toFloat64(vals[i]) != e {
			t.Errorf("Cummin[%d] = %v, want %v", i, vals[i], e)
		}
	}
}

func TestRollingStd(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{2.0}, {4.0}, {4.0}, {4.0}, {5.0}},
	)
	result := df.Rolling(3).Std()
	vals := result.Column("x").Values()
	// window [2,4,4] has std ~= 1.1547
	if v := toFloat64(vals[2]); math.Abs(v-1.1547) > 0.01 {
		t.Errorf("Rolling(3).Std()[2] = %v, want ~1.1547", v)
	}
}
