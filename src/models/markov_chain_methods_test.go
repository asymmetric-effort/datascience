//go:build unit

package models

import (
	"math"
	"testing"
)

func buildSimpleMarkovChain(t *testing.T) *MarkovChain {
	t.Helper()
	tm := [][]float64{
		{0.7, 0.3},
		{0.4, 0.6},
	}
	mc, err := NewMarkovChain(tm, []string{"sunny", "rainy"})
	if err != nil {
		t.Fatalf("NewMarkovChain: %v", err)
	}
	return mc
}

func TestMarkovChainSetStartState(t *testing.T) {
	mc := buildSimpleMarkovChain(t)

	state, err := mc.SetStartState(0)
	if err != nil {
		t.Fatalf("SetStartState(0): %v", err)
	}
	if state != 0 {
		t.Errorf("expected 0, got %d", state)
	}

	state, err = mc.SetStartState(1)
	if err != nil {
		t.Fatalf("SetStartState(1): %v", err)
	}
	if state != 1 {
		t.Errorf("expected 1, got %d", state)
	}
}

func TestMarkovChainSetStartStateInvalid(t *testing.T) {
	mc := buildSimpleMarkovChain(t)

	if _, err := mc.SetStartState(-1); err == nil {
		t.Error("expected error for negative state")
	}
	if _, err := mc.SetStartState(5); err == nil {
		t.Error("expected error for out-of-range state")
	}
}

func TestMarkovChainAddVariable(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	oldN := mc.NumStates()

	mc.AddVariable("cloudy")

	if mc.NumStates() != oldN+1 {
		t.Errorf("expected %d states, got %d", oldN+1, mc.NumStates())
	}

	names := mc.StateNames()
	if names[2] != "cloudy" {
		t.Errorf("expected last state name 'cloudy', got %q", names[2])
	}

	// New row should be a valid distribution.
	tm := mc.TransitionMatrix()
	sum := 0.0
	for _, v := range tm[2] {
		sum += v
	}
	if math.Abs(sum-1.0) > 1e-6 {
		t.Errorf("new row sums to %f, expected 1.0", sum)
	}
}

func TestMarkovChainAddVariablesFrom(t *testing.T) {
	mc1 := buildSimpleMarkovChain(t)
	tm2 := [][]float64{
		{0.5, 0.5},
		{0.5, 0.5},
	}
	mc2, _ := NewMarkovChain(tm2, []string{"cloudy", "rainy"})

	mc1.AddVariablesFrom(mc2)

	// "rainy" already exists, so only "cloudy" should be added.
	if mc1.NumStates() != 3 {
		t.Errorf("expected 3 states, got %d", mc1.NumStates())
	}
	names := mc1.StateNames()
	found := false
	for _, n := range names {
		if n == "cloudy" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'cloudy' in state names")
	}
}

func TestMarkovChainAddVariablesFromNil(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	mc.AddVariablesFrom(nil) // Should not panic.
	if mc.NumStates() != 2 {
		t.Errorf("expected 2 states, got %d", mc.NumStates())
	}
}

func TestMarkovChainAddTransitionModel(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	newTM := [][]float64{
		{0.5, 0.5},
		{0.3, 0.7},
	}

	if err := mc.AddTransitionModel(newTM); err != nil {
		t.Fatalf("AddTransitionModel: %v", err)
	}

	got := mc.TransitionMatrix()
	if got[0][0] != 0.5 || got[1][1] != 0.7 {
		t.Errorf("transition matrix not updated correctly: %v", got)
	}
}

func TestMarkovChainAddTransitionModelWrongSize(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	if err := mc.AddTransitionModel([][]float64{{1.0}}); err == nil {
		t.Error("expected error for wrong matrix size")
	}
}

func TestMarkovChainAddTransitionModelInvalidRow(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	if err := mc.AddTransitionModel([][]float64{
		{0.3, 0.3},
		{0.5, 0.5},
	}); err == nil {
		t.Error("expected error for row not summing to 1")
	}
}

func TestMarkovChainProbFromSample(t *testing.T) {
	mc := buildSimpleMarkovChain(t)

	// Generate a long sample and estimate probabilities.
	samples, err := mc.Sample(10000, 0, 42)
	if err != nil {
		t.Fatalf("Sample: %v", err)
	}

	estimated, err := mc.ProbFromSample(samples)
	if err != nil {
		t.Fatalf("ProbFromSample: %v", err)
	}

	// The estimated T[0][0] should be close to 0.7.
	if math.Abs(estimated[0][0]-0.7) > 0.05 {
		t.Errorf("estimated T[0][0] = %f, expected ~0.7", estimated[0][0])
	}
	if math.Abs(estimated[1][1]-0.6) > 0.05 {
		t.Errorf("estimated T[1][1] = %f, expected ~0.6", estimated[1][1])
	}
}

func TestMarkovChainProbFromSampleShort(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	_, err := mc.ProbFromSample([]int{0})
	if err == nil {
		t.Error("expected error for sequence too short")
	}
}

func TestMarkovChainProbFromSampleInvalidState(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	_, err := mc.ProbFromSample([]int{0, 5})
	if err == nil {
		t.Error("expected error for invalid state index")
	}
}

func TestMarkovChainGenerateSample(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	samples, err := mc.GenerateSample(50, 0, 42)
	if err != nil {
		t.Fatalf("GenerateSample: %v", err)
	}
	if len(samples) != 50 {
		t.Errorf("expected 50 samples, got %d", len(samples))
	}
	if samples[0] != 0 {
		t.Errorf("expected first sample 0, got %d", samples[0])
	}
}

func TestMarkovChainIsStationarity(t *testing.T) {
	mc := buildSimpleMarkovChain(t)

	pi, err := mc.StationaryDistribution()
	if err != nil {
		t.Fatalf("StationaryDistribution: %v", err)
	}

	isStationary, err := mc.IsStationarity(pi)
	if err != nil {
		t.Fatalf("IsStationarity: %v", err)
	}
	if !isStationary {
		t.Error("stationary distribution should be stationary")
	}
}

func TestMarkovChainIsStationarityFalse(t *testing.T) {
	mc := buildSimpleMarkovChain(t)

	nonStationary := []float64{0.9, 0.1}
	isStationary, err := mc.IsStationarity(nonStationary)
	if err != nil {
		t.Fatalf("IsStationarity: %v", err)
	}
	if isStationary {
		t.Error("non-stationary distribution should not be stationary")
	}
}

func TestMarkovChainIsStationarityWrongLength(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	_, err := mc.IsStationarity([]float64{1.0})
	if err == nil {
		t.Error("expected error for wrong distribution length")
	}
}

func TestMarkovChainRandomState(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	state, err := mc.RandomState(42)
	if err != nil {
		t.Fatalf("RandomState: %v", err)
	}
	if state < 0 || state >= mc.NumStates() {
		t.Errorf("state %d out of range", state)
	}
}

func TestMarkovChainCopy(t *testing.T) {
	mc := buildSimpleMarkovChain(t)
	cpy := mc.Copy()

	if cpy.NumStates() != mc.NumStates() {
		t.Errorf("copy has %d states, expected %d", cpy.NumStates(), mc.NumStates())
	}

	// Modify copy and ensure original is unaffected.
	cpy.AddVariable("new_state")
	if mc.NumStates() != 2 {
		t.Error("original was affected by copy modification")
	}

	// Verify transition matrix is a deep copy.
	cpyTM := cpy.TransitionMatrix()
	origTM := mc.TransitionMatrix()
	if cpyTM[0][0] != origTM[0][0] {
		t.Error("copy does not have same transition matrix values")
	}
}

func TestMarkovChainCopyNilNames(t *testing.T) {
	tm := [][]float64{{0.5, 0.5}, {0.5, 0.5}}
	mc, _ := NewMarkovChain(tm, nil)
	cpy := mc.Copy()
	if cpy.StateNames() != nil {
		t.Error("expected nil state names in copy")
	}
}
