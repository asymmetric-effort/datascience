package optimizer

import (
	"math"
	"testing"
)

func TestSGDBasic(t *testing.T) {
	opt := NewSGD(0.1, 0.0)
	params := []float64{1.0, 2.0}
	grads := []float64{0.5, -0.5}
	opt.Step(0, params, grads)

	// params -= lr * grads => [1 - 0.05, 2 + 0.05] = [0.95, 2.05]
	if math.Abs(params[0]-0.95) > 1e-10 {
		t.Errorf("params[0] = %f, want 0.95", params[0])
	}
	if math.Abs(params[1]-2.05) > 1e-10 {
		t.Errorf("params[1] = %f, want 2.05", params[1])
	}
}

func TestSGDMomentum(t *testing.T) {
	opt := NewSGD(0.1, 0.9)
	params := []float64{1.0}
	grads := []float64{1.0}

	// Step 1: v = 0.9*0 + 1.0 = 1.0, param = 1 - 0.1*1.0 = 0.9
	opt.Step(0, params, grads)
	if math.Abs(params[0]-0.9) > 1e-10 {
		t.Errorf("step 1: params[0] = %f, want 0.9", params[0])
	}

	// Step 2: v = 0.9*1.0 + 1.0 = 1.9, param = 0.9 - 0.1*1.9 = 0.71
	opt.Step(0, params, grads)
	if math.Abs(params[0]-0.71) > 1e-10 {
		t.Errorf("step 2: params[0] = %f, want 0.71", params[0])
	}
}

func TestSGDLearningRate(t *testing.T) {
	opt := NewSGD(0.01, 0.0)
	if opt.LearningRate() != 0.01 {
		t.Errorf("LearningRate() = %f, want 0.01", opt.LearningRate())
	}
	opt.SetLearningRate(0.001)
	if opt.LearningRate() != 0.001 {
		t.Errorf("LearningRate() = %f, want 0.001", opt.LearningRate())
	}
}

func TestSGDMultipleParams(t *testing.T) {
	opt := NewSGD(0.1, 0.0)
	params1 := []float64{1.0}
	params2 := []float64{2.0}
	grads1 := []float64{1.0}
	grads2 := []float64{-1.0}

	opt.Step(0, params1, grads1)
	opt.Step(1, params2, grads2)

	if math.Abs(params1[0]-0.9) > 1e-10 {
		t.Errorf("params1[0] = %f, want 0.9", params1[0])
	}
	if math.Abs(params2[0]-2.1) > 1e-10 {
		t.Errorf("params2[0] = %f, want 2.1", params2[0])
	}
}

func TestAdamBasic(t *testing.T) {
	opt := NewAdam(0.001)
	params := []float64{1.0, 2.0}
	grads := []float64{0.1, -0.1}
	opt.Step(0, params, grads)

	// After one step, params should have moved.
	if params[0] >= 1.0 {
		t.Error("params[0] should have decreased")
	}
	if params[1] <= 2.0 {
		t.Error("params[1] should have increased")
	}
}

func TestAdamConvergence(t *testing.T) {
	// Minimize f(x) = x^2, gradient = 2x
	opt := NewAdam(0.1)
	params := []float64{5.0}
	for range 1000 {
		grads := []float64{2 * params[0]}
		opt.Step(0, params, grads)
	}
	if math.Abs(params[0]) > 0.01 {
		t.Errorf("Adam did not converge to 0: params = %f", params[0])
	}
}

func TestAdamWithParams(t *testing.T) {
	opt := NewAdamWithParams(0.01, 0.8, 0.99, 1e-7)
	if opt.LearningRate() != 0.01 {
		t.Errorf("LearningRate() = %f, want 0.01", opt.LearningRate())
	}
	params := []float64{1.0}
	grads := []float64{0.5}
	opt.Step(0, params, grads)
	if params[0] >= 1.0 {
		t.Error("params should have moved")
	}
}

func TestAdamLearningRate(t *testing.T) {
	opt := NewAdam(0.001)
	if opt.LearningRate() != 0.001 {
		t.Errorf("LearningRate() = %f, want 0.001", opt.LearningRate())
	}
	opt.SetLearningRate(0.01)
	if opt.LearningRate() != 0.01 {
		t.Errorf("LearningRate() = %f, want 0.01", opt.LearningRate())
	}
}

func TestRMSPropBasic(t *testing.T) {
	opt := NewRMSProp(0.01)
	params := []float64{1.0}
	grads := []float64{0.5}
	opt.Step(0, params, grads)
	if params[0] >= 1.0 {
		t.Error("params should have decreased")
	}
}

func TestRMSPropConvergence(t *testing.T) {
	opt := NewRMSProp(0.01)
	params := []float64{5.0}
	for range 2000 {
		grads := []float64{2 * params[0]}
		opt.Step(0, params, grads)
	}
	if math.Abs(params[0]) > 0.1 {
		t.Errorf("RMSProp did not converge near 0: params = %f", params[0])
	}
}

func TestRMSPropLearningRate(t *testing.T) {
	opt := NewRMSProp(0.01)
	if opt.LearningRate() != 0.01 {
		t.Errorf("LearningRate() = %f, want 0.01", opt.LearningRate())
	}
	opt.SetLearningRate(0.001)
	if opt.LearningRate() != 0.001 {
		t.Errorf("LearningRate() = %f, want 0.001", opt.LearningRate())
	}
}
