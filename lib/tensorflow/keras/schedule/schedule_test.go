package schedule

import (
	"math"
	"testing"
)

func TestExponentialDecay(t *testing.T) {
	s := NewExponentialDecay(0.1, 100, 0.96, false)
	if s.Name() != "exponential_decay" {
		t.Errorf("Name() = %q", s.Name())
	}
	lr0 := s.LearningRate(0)
	if math.Abs(lr0-0.1) > 1e-10 {
		t.Errorf("LR(0) = %f, want 0.1", lr0)
	}
	lr100 := s.LearningRate(100)
	expected := 0.1 * 0.96
	if math.Abs(lr100-expected) > 1e-10 {
		t.Errorf("LR(100) = %f, want %f", lr100, expected)
	}
}

func TestExponentialDecayStaircase(t *testing.T) {
	s := NewExponentialDecay(0.1, 100, 0.96, true)
	// Step 50: floor(50/100)=0, so LR = 0.1 * 0.96^0 = 0.1
	lr50 := s.LearningRate(50)
	if math.Abs(lr50-0.1) > 1e-10 {
		t.Errorf("LR(50) = %f, want 0.1", lr50)
	}
	// Step 100: floor(100/100)=1, so LR = 0.1 * 0.96
	lr100 := s.LearningRate(100)
	expected := 0.1 * 0.96
	if math.Abs(lr100-expected) > 1e-10 {
		t.Errorf("LR(100) = %f, want %f", lr100, expected)
	}
}

func TestCosineDecay(t *testing.T) {
	s := NewCosineDecay(0.1, 1000, 0.0)
	if s.Name() != "cosine_decay" {
		t.Errorf("Name() = %q", s.Name())
	}
	lr0 := s.LearningRate(0)
	if math.Abs(lr0-0.1) > 1e-10 {
		t.Errorf("LR(0) = %f, want 0.1", lr0)
	}
	// At half, cosine = 0.5*(1+cos(pi*0.5)) = 0.5*(1+0) = 0.5
	lr500 := s.LearningRate(500)
	if math.Abs(lr500-0.05) > 1e-10 {
		t.Errorf("LR(500) = %f, want 0.05", lr500)
	}
	// At end.
	lr1000 := s.LearningRate(1000)
	if math.Abs(lr1000) > 1e-10 {
		t.Errorf("LR(1000) = %f, want 0", lr1000)
	}
}

func TestCosineDecayAlpha(t *testing.T) {
	s := NewCosineDecay(0.1, 100, 0.01)
	lr := s.LearningRate(100)
	if math.Abs(lr-0.001) > 1e-10 {
		t.Errorf("LR(100) = %f, want 0.001", lr)
	}
}

func TestCosineDecayBeyondTotal(t *testing.T) {
	s := NewCosineDecay(0.1, 100, 0.01)
	lr := s.LearningRate(200)
	if math.Abs(lr-0.001) > 1e-10 {
		t.Errorf("LR(200) = %f, want 0.001", lr)
	}
}

func TestPiecewiseConstant(t *testing.T) {
	s := NewPiecewiseConstant([]int{100, 200}, []float64{0.1, 0.01, 0.001})
	if s.Name() != "piecewise_constant" {
		t.Errorf("Name() = %q", s.Name())
	}
	if s.LearningRate(0) != 0.1 {
		t.Errorf("LR(0) = %f, want 0.1", s.LearningRate(0))
	}
	if s.LearningRate(50) != 0.1 {
		t.Errorf("LR(50) = %f, want 0.1", s.LearningRate(50))
	}
	if s.LearningRate(100) != 0.01 {
		t.Errorf("LR(100) = %f, want 0.01", s.LearningRate(100))
	}
	if s.LearningRate(300) != 0.001 {
		t.Errorf("LR(300) = %f, want 0.001", s.LearningRate(300))
	}
}

func TestWarmupSchedule(t *testing.T) {
	decay := NewExponentialDecay(0.1, 100, 0.96, false)
	s := NewWarmupSchedule(10, 0.1, decay)
	if s.Name() != "warmup" {
		t.Errorf("Name() = %q", s.Name())
	}
	// During warmup: linear increase.
	lr0 := s.LearningRate(0) // (0+1)/10 * 0.1 = 0.01
	if math.Abs(lr0-0.01) > 1e-10 {
		t.Errorf("LR(0) = %f, want 0.01", lr0)
	}
	lr5 := s.LearningRate(4) // (4+1)/10 * 0.1 = 0.05
	if math.Abs(lr5-0.05) > 1e-10 {
		t.Errorf("LR(4) = %f, want 0.05", lr5)
	}
	lr9 := s.LearningRate(9) // (9+1)/10 * 0.1 = 0.1
	if math.Abs(lr9-0.1) > 1e-10 {
		t.Errorf("LR(9) = %f, want 0.1", lr9)
	}
	// After warmup: uses decay schedule (step offset by warmup steps).
	lr10 := s.LearningRate(10) // decay step 0 => 0.1
	if math.Abs(lr10-0.1) > 1e-10 {
		t.Errorf("LR(10) = %f, want 0.1", lr10)
	}
}

func TestScheduleInterface(t *testing.T) {
	var _ Schedule = &ExponentialDecay{}
	var _ Schedule = &CosineDecay{}
	var _ Schedule = &PiecewiseConstant{}
	var _ Schedule = &WarmupSchedule{}
}
