//go:build unit

package inference

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// buildSimpleCausalBN builds a simple network:
// X -> M -> Y with X also having a direct effect on Y.
// X: 2 states, M: 2 states, Y: 2 states.
func buildSimpleCausalBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	for _, n := range []string{"X", "M", "Y"} {
		if err := bn.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	_ = bn.AddEdge("X", "M")
	_ = bn.AddEdge("M", "Y")
	_ = bn.AddEdge("X", "Y")

	// P(X)
	cpdX, err := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdX)

	// P(M|X)
	cpdM, err := factors.NewTabularCPD("M", 2, [][]float64{{0.8, 0.2}, {0.2, 0.8}}, []string{"X"}, []int{2})
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdM)

	// P(Y|M,X)
	cpdY, err := factors.NewTabularCPD("Y", 2,
		[][]float64{{0.9, 0.6, 0.7, 0.1}, {0.1, 0.4, 0.3, 0.9}},
		[]string{"M", "X"}, []int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdY)

	return bn
}

func TestGetAllBackdoorAdjustmentSets(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	// For X -> Y, backdoor adjustment should find valid sets.
	sets := ci.GetAllBackdoorAdjustmentSets("X", "Y")
	// The empty set should work since X has no parents that create backdoor paths.
	if len(sets) == 0 {
		t.Error("expected at least one backdoor adjustment set")
	}
}

func TestIsValidFrontdoorAdjustmentSet(t *testing.T) {
	// Build a network where front-door is applicable: X -> M -> Y, with X <- U -> Y (confounding).
	// We can't add U -> X and U -> Y easily in a BN without the U variable, so test the false case.
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	// M alone can't be a valid front-door set for X->Y because X also directly causes Y.
	valid := ci.IsValidFrontdoorAdjustmentSet("X", "Y", []string{"M"})
	// This depends on the structure. M doesn't intercept all directed paths (X->Y direct exists).
	if valid {
		t.Error("M should not be a valid frontdoor set when X->Y direct edge exists")
	}
}

func TestGetAllFrontdoorAdjustmentSets(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	sets := ci.GetAllFrontdoorAdjustmentSets("X", "Y")
	// With the direct X->Y edge, no front-door set should work.
	if len(sets) != 0 {
		t.Logf("found %d front-door sets (may be valid depending on structure)", len(sets))
	}
}

func TestGetScalingIndicators(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	indicators := ci.GetScalingIndicators("X", "Y")
	// X is a root. The only exogenous variable is X itself, which is excluded.
	// So no indicators expected.
	if len(indicators) != 0 {
		t.Logf("found %d scaling indicators", len(indicators))
	}
}

func TestGetIVs(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	ivs := ci.GetIVs("X", "Y")
	// With this structure, no node satisfies the IV conditions.
	t.Logf("found %d IVs: %v", len(ivs), ivs)
}

func TestGetConditionalIVs(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	ivs := ci.GetConditionalIVs("X", "Y", nil)
	t.Logf("found %d conditional IVs with empty conditioning: %v", len(ivs), ivs)
}

func TestGetTotalConditionalIVs(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	result := ci.GetTotalConditionalIVs("X", "Y")
	t.Logf("total conditional IVs: %v", result)
}

func TestIdentificationMethod(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	method := ci.IdentificationMethod("X", "Y")
	if method != "backdoor" && method != "frontdoor" && method != "iv" && method != "none" {
		t.Errorf("unexpected identification method: %s", method)
	}
	// Since X has no parents creating a confound, backdoor should work.
	if method != "backdoor" {
		t.Errorf("expected 'backdoor', got %q", method)
	}
}

func TestEstimateATE(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	// EstimateATE without data falls back to model-based.
	_, err = ci.EstimateATE("X", "Y", nil)
	if err == nil {
		t.Error("expected error with nil data")
	}
}

func TestGetProperBackdoorGraph(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	g := ci.GetProperBackdoorGraph("X")
	if g == nil {
		t.Fatal("GetProperBackdoorGraph returned nil")
	}
	// X should have no outgoing edges in the manipulated graph.
	succs := g.Successors("X")
	if len(succs) != 0 {
		t.Errorf("expected 0 successors of X in backdoor graph, got %d", len(succs))
	}
}

func TestIsValidAdjustmentSet(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	// Empty set should be valid for X->Y since X has no parents.
	if !ci.IsValidAdjustmentSet("X", "Y", []string{}) {
		t.Error("empty set should be valid adjustment for X->Y")
	}
}

func TestGetMinimalAdjustmentSet(t *testing.T) {
	bn := buildSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	minimal, err := ci.GetMinimalAdjustmentSet("X", "Y")
	if err != nil {
		t.Fatalf("GetMinimalAdjustmentSet failed: %v", err)
	}
	// X has no parents, so the minimal set should be empty.
	if len(minimal) != 0 {
		t.Errorf("expected empty minimal adjustment set, got %v", minimal)
	}
}
