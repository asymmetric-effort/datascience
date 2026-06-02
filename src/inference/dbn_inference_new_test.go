//go:build unit

package inference

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

func buildSimpleDBN(t *testing.T) *DBNInference {
	t.Helper()

	// Single interface node X with 2 states.
	// Initial: P(X) = [0.6, 0.4]
	initF, err := factors.NewDiscreteFactor([]string{"X"}, []int{2}, []float64{0.6, 0.4})
	if err != nil {
		t.Fatal(err)
	}

	// Transition: P(X | X_prev)
	transF, err := factors.NewDiscreteFactor(
		[]string{"X", "X_prev"}, []int{2, 2},
		[]float64{0.7, 0.3, 0.4, 0.6})
	if err != nil {
		t.Fatal(err)
	}

	return NewDBNInference(
		[]*factors.DiscreteFactor{initF},
		[]*factors.DiscreteFactor{transF},
		[]string{"X"},
	)
}

func TestDBNBackwardInference(t *testing.T) {
	dbn := buildSimpleDBN(t)

	evidence := []map[string]int{
		{}, // t=0
		{}, // t=1
	}

	result, err := dbn.BackwardInference([]string{"X"}, evidence, 0)
	if err != nil {
		t.Fatalf("BackwardInference failed: %v", err)
	}

	// The result should be a valid probability distribution over X.
	vars := result.Variables()
	if len(vars) != 1 || vars[0] != "X" {
		t.Errorf("expected variable X, got %v", vars)
	}

	data := result.Values().Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("expected sum ~1.0, got %f", sum)
	}
}

func TestDBNBackwardInferenceErrors(t *testing.T) {
	dbn := buildSimpleDBN(t)

	_, err := dbn.BackwardInference(nil, []map[string]int{{}}, 0)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}

	_, err = dbn.BackwardInference([]string{"X"}, nil, 0)
	if err == nil {
		t.Error("expected error for empty evidence sequence")
	}

	_, err = dbn.BackwardInference([]string{"X"}, []map[string]int{{}}, -1)
	if err == nil {
		t.Error("expected error for negative target time step")
	}

	_, err = dbn.BackwardInference([]string{"X"}, []map[string]int{{}}, 5)
	if err == nil {
		t.Error("expected error for out-of-range target time step")
	}
}

func TestDBNQuery(t *testing.T) {
	dbn := buildSimpleDBN(t)

	evidence := []map[string]int{
		{}, // t=0
		{}, // t=1
		{}, // t=2
	}

	// Query at last time step should use forward inference.
	result, err := dbn.Query([]string{"X"}, evidence, -1)
	if err != nil {
		t.Fatalf("Query at last step failed: %v", err)
	}
	if result == nil {
		t.Fatal("Query returned nil")
	}

	// Query at time step 0 should use backward inference (smoothing).
	result, err = dbn.Query([]string{"X"}, evidence, 0)
	if err != nil {
		t.Fatalf("Query at step 0 failed: %v", err)
	}
	if result == nil {
		t.Fatal("Query returned nil for smoothing")
	}
}

func TestDBNQueryErrors(t *testing.T) {
	dbn := buildSimpleDBN(t)

	_, err := dbn.Query([]string{"X"}, nil, 0)
	if err == nil {
		t.Error("expected error for empty evidence sequence")
	}
}
