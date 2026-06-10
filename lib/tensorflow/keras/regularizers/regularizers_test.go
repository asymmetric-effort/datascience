package regularizers

import (
	"math"
	"testing"
)

func TestL1(t *testing.T) {
	r := NewL1(0.01)
	if r.Name() != "l1" {
		t.Errorf("Name() = %q", r.Name())
	}
	penalty := r.Penalty([]float64{1, -2, 3})
	// 0.01 * (1+2+3) = 0.06
	if math.Abs(penalty-0.06) > 1e-10 {
		t.Errorf("Penalty = %f, want 0.06", penalty)
	}
}

func TestL1Empty(t *testing.T) {
	r := NewL1(0.01)
	if r.Penalty([]float64{}) != 0 {
		t.Error("expected 0 for empty weights")
	}
}

func TestL2(t *testing.T) {
	r := NewL2(0.01)
	if r.Name() != "l2" {
		t.Errorf("Name() = %q", r.Name())
	}
	penalty := r.Penalty([]float64{1, 2, 3})
	// 0.01 * (1+4+9) = 0.14
	if math.Abs(penalty-0.14) > 1e-10 {
		t.Errorf("Penalty = %f, want 0.14", penalty)
	}
}

func TestL2Empty(t *testing.T) {
	r := NewL2(0.01)
	if r.Penalty([]float64{}) != 0 {
		t.Error("expected 0 for empty weights")
	}
}

func TestL1L2(t *testing.T) {
	r := NewL1L2(0.01, 0.02)
	if r.Name() != "l1_l2" {
		t.Errorf("Name() = %q", r.Name())
	}
	penalty := r.Penalty([]float64{1, -2})
	// L1: 0.01*(1+2)=0.03, L2: 0.02*(1+4)=0.10, total=0.13
	if math.Abs(penalty-0.13) > 1e-10 {
		t.Errorf("Penalty = %f, want 0.13", penalty)
	}
}

func TestOrthogonal(t *testing.T) {
	r := NewOrthogonalRegularizer(1.0, 2)
	if r.Name() != "orthogonal" {
		t.Errorf("Name() = %q", r.Name())
	}
	// Identity matrix should have zero penalty.
	penalty := r.Penalty([]float64{1, 0, 0, 1})
	if math.Abs(penalty) > 1e-10 {
		t.Errorf("Penalty for identity = %f, want 0", penalty)
	}
}

func TestOrthogonalNonIdentity(t *testing.T) {
	r := NewOrthogonalRegularizer(1.0, 2)
	penalty := r.Penalty([]float64{2, 0, 0, 2})
	// W^T W = [[4,0],[0,4]], I = [[1,0],[0,1]], diff = [[3,0],[0,3]], ||diff||^2 = 18
	if math.Abs(penalty-18.0) > 1e-10 {
		t.Errorf("Penalty = %f, want 18", penalty)
	}
}

func TestOrthogonalWrongSize(t *testing.T) {
	r := NewOrthogonalRegularizer(1.0, 3)
	// Wrong size weights should return 0.
	penalty := r.Penalty([]float64{1, 0, 0, 1})
	if penalty != 0 {
		t.Errorf("Penalty for wrong size = %f, want 0", penalty)
	}
}

func TestRegularizerInterface(t *testing.T) {
	var _ Regularizer = &L1{}
	var _ Regularizer = &L2{}
	var _ Regularizer = &L1L2{}
	var _ Regularizer = &OrthogonalRegularizer{}
}
