//go:build unit

package gpu

import (
	"math"
	"testing"
)

func TestAcceleratedVESimple(t *testing.T) {
	b := NewCPUBackend()

	// Simple network: A -> B
	// P(A) = [0.4, 0.6]
	// P(B|A) = [[0.2, 0.8], [0.5, 0.5]]
	// Query: P(B) = sum_A P(A)*P(B|A)
	//   P(B=0) = 0.4*0.2 + 0.6*0.5 = 0.08 + 0.30 = 0.38
	//   P(B=1) = 0.4*0.8 + 0.6*0.5 = 0.32 + 0.30 = 0.62
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.4, 0.6}},
		{Variables: []string{"A", "B"}, Shape: []int{2, 2}, Values: []float64{0.2, 0.8, 0.5, 0.5}},
	}

	result, err := AcceleratedVE(b, factors, []string{"B"}, []string{"A"}, nil)
	if err != nil {
		t.Fatalf("AcceleratedVE: %v", err)
	}

	// The result should be a normalized distribution.
	sum := b.Sum(result)
	if math.Abs(sum-1.0) > 1e-6 {
		t.Errorf("result should sum to 1, got %v", sum)
	}

	// Verify it has the right number of elements (2 for binary B).
	if len(result) < 2 {
		t.Fatalf("expected at least 2 elements, got %d", len(result))
	}
}

func TestAcceleratedVEWithEvidence(t *testing.T) {
	b := NewCPUBackend()

	// P(A) = [0.4, 0.6]
	// P(B|A) shape [2,2] = [[0.2, 0.8], [0.5, 0.5]]
	// Evidence: A=0
	// P(B|A=0) = [0.2, 0.8] normalized = [0.2, 0.8]
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.4, 0.6}},
		{Variables: []string{"A", "B"}, Shape: []int{2, 2}, Values: []float64{0.2, 0.8, 0.5, 0.5}},
	}

	result, err := AcceleratedVE(b, factors, []string{"B"}, []string{}, map[string]int{"A": 0})
	if err != nil {
		t.Fatalf("AcceleratedVE with evidence: %v", err)
	}

	sum := b.Sum(result)
	if math.Abs(sum-1.0) > 1e-6 {
		t.Errorf("result should sum to 1, got %v", sum)
	}
}

func TestAcceleratedVENoFactors(t *testing.T) {
	b := NewCPUBackend()
	_, err := AcceleratedVE(b, nil, []string{"A"}, nil, nil)
	if err == nil {
		t.Error("expected error for empty factors")
	}
}

func TestAcceleratedBPSimple(t *testing.T) {
	b := NewCPUBackend()

	// Two cliques sharing variable B: {A,B} and {B,C}
	cliques := [][]string{{"A", "B"}, {"B", "C"}}
	cliqueFactors := []FactorData{
		{Variables: []string{"A", "B"}, Shape: []int{2, 2}, Values: []float64{0.5, 0.8, 0.1, 0.3}},
		{Variables: []string{"B", "C"}, Shape: []int{2, 2}, Values: []float64{0.5, 0.5, 0.4, 0.6}},
	}

	err := AcceleratedBP(b, cliques, cliqueFactors)
	if err != nil {
		t.Fatalf("AcceleratedBP: %v", err)
	}

	// After BP, beliefs should be updated (non-zero).
	for i, cf := range cliqueFactors {
		allZero := true
		for _, v := range cf.Values {
			if v != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			t.Errorf("clique %d belief should not be all zeros", i)
		}
	}
}

func TestAcceleratedBPMismatch(t *testing.T) {
	b := NewCPUBackend()
	err := AcceleratedBP(b, [][]string{{"A"}}, []FactorData{})
	if err == nil {
		t.Error("expected error for mismatched lengths")
	}
}

func TestAcceleratedSampleBasic(t *testing.T) {
	b := NewCPUBackend()

	// Simple chain: A -> B
	// P(A) = [0.3, 0.7]
	// P(B|A) = [[0.9, 0.1], [0.2, 0.8]]
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.3, 0.7}},
		{Variables: []string{"A", "B"}, Shape: []int{2, 2}, Values: []float64{0.9, 0.1, 0.2, 0.8}},
	}

	samples := AcceleratedSample(b, factors, []string{"A", "B"}, 100)
	if len(samples) != 100 {
		t.Fatalf("expected 100 samples, got %d", len(samples))
	}

	for i, row := range samples {
		if len(row) != 2 {
			t.Fatalf("sample %d: expected 2 values, got %d", i, len(row))
		}
		if row[0] < 0 || row[0] > 1 {
			t.Errorf("sample %d: A=%d out of range", i, row[0])
		}
		if row[1] < 0 || row[1] > 1 {
			t.Errorf("sample %d: B=%d out of range", i, row[1])
		}
	}
}

func TestAcceleratedSampleSingleVar(t *testing.T) {
	b := NewCPUBackend()

	// Deterministic: P(A) = [0, 1] so A should always be 1.
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0, 1}},
	}

	samples := AcceleratedSample(b, factors, []string{"A"}, 50)
	for i, row := range samples {
		if row[0] != 1 {
			t.Errorf("sample %d: expected A=1, got %d", i, row[0])
		}
	}
}

func TestAcceleratedSampleEmpty(t *testing.T) {
	b := NewCPUBackend()
	samples := AcceleratedSample(b, nil, []string{"A"}, 0)
	if len(samples) != 0 {
		t.Errorf("expected 0 samples, got %d", len(samples))
	}
}

// --- Helper function tests ---

func TestIndexOf(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if indexOf(slice, "b") != 1 {
		t.Errorf("indexOf: expected 1")
	}
	if indexOf(slice, "d") != -1 {
		t.Errorf("indexOf: expected -1 for missing")
	}
}

func TestRemoveString(t *testing.T) {
	slice := []string{"a", "b", "c"}
	result := removeString(slice, 1)
	if len(result) != 2 || result[0] != "a" || result[1] != "c" {
		t.Errorf("removeString: expected [a,c], got %v", result)
	}
}

func TestSharedVars(t *testing.T) {
	a := []string{"A", "B", "C"}
	b := []string{"B", "C", "D"}
	shared := sharedVars(a, b)
	if len(shared) != 2 || shared[0] != "B" || shared[1] != "C" {
		t.Errorf("sharedVars: expected [B,C], got %v", shared)
	}
}

func TestSharedVarsNone(t *testing.T) {
	a := []string{"A"}
	b := []string{"B"}
	shared := sharedVars(a, b)
	if shared != nil {
		t.Errorf("sharedVars: expected nil, got %v", shared)
	}
}

func TestFactorDataStruct(t *testing.T) {
	f := FactorData{
		Variables: []string{"X", "Y"},
		Shape:     []int{2, 3},
		Values:    []float64{1, 2, 3, 4, 5, 6},
	}
	if len(f.Variables) != 2 {
		t.Errorf("expected 2 variables, got %d", len(f.Variables))
	}
	if len(f.Values) != 6 {
		t.Errorf("expected 6 values, got %d", len(f.Values))
	}
}

// TestAcceleratedVEEvidenceNotInFactor covers the case where evidence
// references a variable not present in a factor (axis < 0 → continue).
// Covers accelerated.go lines 43-44.
func TestAcceleratedVEEvidenceNotInFactor(t *testing.T) {
	b := NewCPUBackend()
	// Factor only mentions A, but evidence is for Z which is absent.
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.4, 0.6}},
	}
	result, err := AcceleratedVE(b, factors, []string{"A"}, []string{}, map[string]int{"Z": 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

// TestAcceleratedVERemainingFactors covers the else branch where a factor
// does not mention the elimination variable (remaining path).
// Covers accelerated.go lines 61-63.
func TestAcceleratedVERemainingFactors(t *testing.T) {
	b := NewCPUBackend()
	// Factor for A, factor for B (independent). Eliminate A.
	// The B factor goes into "remaining" (not involved with A).
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.4, 0.6}},
		{Variables: []string{"B"}, Shape: []int{2}, Values: []float64{0.3, 0.7}},
	}
	result, err := AcceleratedVE(b, factors, []string{"B"}, []string{"A"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := b.Sum(result)
	if math.Abs(sum-1.0) > 1e-6 {
		t.Errorf("result should sum to 1, got %v", sum)
	}
}

// TestAcceleratedVEElimVarNotInAnyFactor covers the case where an
// elimination variable is not present in any factor (len(involved)==0).
// Covers accelerated.go lines 65-66.
func TestAcceleratedVEElimVarNotInAnyFactor(t *testing.T) {
	b := NewCPUBackend()
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.4, 0.6}},
	}
	// Eliminate "Z" which doesn't exist in any factor.
	result, err := AcceleratedVE(b, factors, []string{"A"}, []string{"Z"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

// TestAcceleratedVESingleInvolvedFactor covers elimination where only one
// factor involves the eliminated variable (no multiplication needed).
func TestAcceleratedVESingleInvolvedFactor(t *testing.T) {
	b := NewCPUBackend()
	factors := []FactorData{
		{Variables: []string{"A", "B"}, Shape: []int{2, 2}, Values: []float64{0.2, 0.8, 0.5, 0.5}},
	}
	result, err := AcceleratedVE(b, factors, []string{"B"}, []string{"A"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

// TestAcceleratedBPThreeCliques covers the BP message passing path where
// a clique receives messages from multiple neighbors, triggering the
// condition nb[1]==from && nb[0]!=to. Covers accelerated.go lines 166-168.
func TestAcceleratedBPThreeCliques(t *testing.T) {
	b := NewCPUBackend()
	// Three cliques: {A,B}, {B,C}, {A,C} - each pair shares a variable.
	// When updating message from clique 0 to clique 1, clique 0 should
	// incorporate the incoming message from clique 2 (nb[0]=2, nb[1]=0, to=1).
	//
	// Use cardinality 1 for all variables so that tensor sizes stay at 1
	// element regardless of how many outer products are computed during
	// the 100 BP iterations.
	cliques := [][]string{{"A", "B"}, {"B", "C"}, {"A", "C"}}
	cliqueFactors := []FactorData{
		{Variables: []string{"A", "B"}, Shape: []int{1, 1}, Values: []float64{0.5}},
		{Variables: []string{"B", "C"}, Shape: []int{1, 1}, Values: []float64{0.4}},
		{Variables: []string{"A", "C"}, Shape: []int{1, 1}, Values: []float64{0.7}},
	}
	err := AcceleratedBP(b, cliques, cliqueFactors)
	if err != nil {
		t.Fatalf("AcceleratedBP: %v", err)
	}
	for i, cf := range cliqueFactors {
		allZero := true
		for _, v := range cf.Values {
			if v != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			t.Errorf("clique %d belief should not be all zeros", i)
		}
	}
}

// TestAcceleratedSampleMissingFactor covers the case where topoOrder
// contains a variable with no corresponding factor.
// Covers accelerated.go lines 245-248.
func TestAcceleratedSampleMissingFactor(t *testing.T) {
	b := NewCPUBackend()
	// Factor for A only, but topoOrder includes B which has no factor.
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.3, 0.7}},
	}
	samples := AcceleratedSample(b, factors, []string{"A", "B"}, 10)
	if len(samples) != 10 {
		t.Fatalf("expected 10 samples, got %d", len(samples))
	}
	for i, row := range samples {
		if len(row) != 2 {
			t.Fatalf("sample %d: expected 2 values, got %d", i, len(row))
		}
		// B should always be 0 (default for missing factor).
		if row[1] != 0 {
			t.Errorf("sample %d: expected B=0 for missing factor, got %d", i, row[1])
		}
	}
}

// TestAcceleratedSampleParentNotYetSampled covers the case where a
// factor references a parent variable that hasn't been sampled yet
// (not in sample map). Covers accelerated.go lines 259-260.
func TestAcceleratedSampleParentNotYetSampled(t *testing.T) {
	b := NewCPUBackend()
	// Factor for B depends on C (parent), but C comes after B in topoOrder.
	// When processing B, C hasn't been sampled yet → exists=false → continue.
	factors := []FactorData{
		{Variables: []string{"C", "B"}, Shape: []int{2, 2}, Values: []float64{0.9, 0.1, 0.2, 0.8}},
		{Variables: []string{"C"}, Shape: []int{2}, Values: []float64{0.5, 0.5}},
	}
	// B comes before C in topo order, so when sampling B, parent C is not yet sampled.
	samples := AcceleratedSample(b, factors, []string{"B", "C"}, 10)
	if len(samples) != 10 {
		t.Fatalf("expected 10 samples, got %d", len(samples))
	}
}

// TestAcceleratedSampleTwoParents covers sampling with a factor that has
// two parent variables, both already sampled, exercising multi-parent reduction.
func TestAcceleratedSampleTwoParents(t *testing.T) {
	b := NewCPUBackend()
	factors := []FactorData{
		{Variables: []string{"A"}, Shape: []int{2}, Values: []float64{0.5, 0.5}},
		{Variables: []string{"B"}, Shape: []int{2}, Values: []float64{0.3, 0.7}},
		// C depends on both A and B.
		{Variables: []string{"A", "B", "C"}, Shape: []int{2, 2, 2},
			Values: []float64{0.1, 0.9, 0.3, 0.7, 0.4, 0.6, 0.8, 0.2}},
	}
	samples := AcceleratedSample(b, factors, []string{"A", "B", "C"}, 10)
	if len(samples) != 10 {
		t.Fatalf("expected 10 samples, got %d", len(samples))
	}
	for i, row := range samples {
		if len(row) != 3 {
			t.Fatalf("sample %d: expected 3 values, got %d", i, len(row))
		}
	}
}

// TestMessageDistanceDifferentLengths covers messageDistance when
// input vectors have different lengths. Covers accelerated.go lines 365-367.
func TestMessageDistanceDifferentLengths(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1.0, 2.0}
	c := []float64{1.0, 2.0, 3.0}
	dist := messageDistance(b, a, c)
	if !math.IsInf(dist, 1) {
		t.Errorf("expected +Inf for mismatched lengths, got %v", dist)
	}
}

// TestSampleCategoricalFallback covers the case where the cumulative
// sum never exceeds u, returning the last index.
// Covers accelerated.go line 390.
func TestSampleCategoricalFallback(t *testing.T) {
	// Probabilities that sum to less than 1 due to floating point,
	// combined with a seed that produces u very close to 1.
	// We need u >= cumulative for all elements.
	// Since u = float64(h%1000000)/1000000.0, the max is 0.999999.
	// If probs sum to < 0.999999, the loop won't return early.
	//
	// Use very small probabilities that sum to much less than 1.
	probs := []float64{0.0, 0.0, 0.0}
	// Any seed will cause fallback since all probs are 0.
	result := sampleCategorical(probs, 0, 0)
	if result != len(probs)-1 {
		t.Errorf("expected last index %d, got %d", len(probs)-1, result)
	}
}

// TestCPUBackendFactorMaximize3D covers FactorMaximize on a 3D tensor,
// exercising the newStrides computation loop for multi-dimensional results.
// Covers cpu_backend.go lines 392-394 (newStrides loop with len >= 2).
func TestCPUBackendFactorMaximize3D(t *testing.T) {
	b := NewCPUBackend()
	// 3D tensor with shape [2,2,3], maximize over axis 1.
	// Result shape: [2,3], newStrides has 2 elements, loop body executes.
	// Values in row-major: [[[1,2,3],[4,5,6]], [[7,8,9],[10,11,12]]]
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	result, newShape := b.FactorMaximize(values, []int{2, 2, 3}, 1)
	expectedShape := []int{2, 3}
	if len(newShape) != len(expectedShape) {
		t.Fatalf("expected shape %v, got %v", expectedShape, newShape)
	}
	for i := range expectedShape {
		if newShape[i] != expectedShape[i] {
			t.Fatalf("expected shape %v, got %v", expectedShape, newShape)
		}
	}
	// Max over axis 1: max([1,4])=4, max([2,5])=5, max([3,6])=6,
	//                   max([7,10])=10, max([8,11])=11, max([9,12])=12
	expected := []float64{4, 5, 6, 10, 11, 12}
	if len(result) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("result[%d]: expected %v, got %v", i, expected[i], result[i])
		}
	}
}
