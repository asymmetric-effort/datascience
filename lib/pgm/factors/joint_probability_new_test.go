//go:build unit

package factors

import (
	"strings"
	"testing"
)

func buildTestJPD(t *testing.T) *JointProbabilityDistribution {
	t.Helper()
	// P(A, B) where A and B are independent.
	// P(A=0)=0.6, P(A=1)=0.4, P(B=0)=0.5, P(B=1)=0.5
	// Joint: P(A=0,B=0)=0.3, P(A=0,B=1)=0.3, P(A=1,B=0)=0.2, P(A=1,B=1)=0.2
	jpd, err := NewJointProbabilityDistribution(
		[]string{"A", "B"}, []int{2, 2},
		[]float64{0.3, 0.3, 0.2, 0.2},
	)
	if err != nil {
		t.Fatal(err)
	}
	return jpd
}

func buildDependentJPD(t *testing.T) *JointProbabilityDistribution {
	t.Helper()
	// P(A, B) where A and B are NOT independent.
	jpd, err := NewJointProbabilityDistribution(
		[]string{"A", "B"}, []int{2, 2},
		[]float64{0.4, 0.1, 0.1, 0.4},
	)
	if err != nil {
		t.Fatal(err)
	}
	return jpd
}

func TestGetIndependenciesIndependent(t *testing.T) {
	jpd := buildTestJPD(t)
	indeps := jpd.GetIndependencies(0.01)

	// A and B are independent, so we should find A _|_ B.
	found := false
	for _, ind := range indeps {
		if (ind[0][0] == "A" && ind[1][0] == "B") || (ind[0][0] == "B" && ind[1][0] == "A") {
			if len(ind[2]) == 0 {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected to find A _|_ B independence")
	}
}

func TestGetIndependenciesDependent(t *testing.T) {
	jpd := buildDependentJPD(t)
	indeps := jpd.GetIndependencies(0.01)

	// A and B are dependent, so marginal independence should NOT be found.
	for _, ind := range indeps {
		if len(ind[2]) == 0 {
			t.Error("A and B are dependent, should not find marginal independence")
		}
	}
}

func TestMinimalIMap(t *testing.T) {
	jpd := buildTestJPD(t)
	edges := jpd.MinimalIMap([]string{"A", "B"}, 0.01)

	// A and B are independent, so minimal I-map should have no edges.
	if len(edges) != 0 {
		t.Errorf("expected 0 edges for independent vars, got %d: %v", len(edges), edges)
	}
}

func TestMinimalIMapDependent(t *testing.T) {
	jpd := buildDependentJPD(t)
	edges := jpd.MinimalIMap([]string{"A", "B"}, 0.01)

	// A and B are dependent, so there should be an edge.
	if len(edges) != 1 {
		t.Errorf("expected 1 edge for dependent vars, got %d: %v", len(edges), edges)
	}
}

func TestIsIMap(t *testing.T) {
	jpd := buildTestJPD(t)

	// No edges: valid I-map for independent A and B.
	if !jpd.IsIMap(nil, 0.01) {
		t.Error("empty graph should be an I-map for independent variables")
	}

	// With an edge A->B: also a valid I-map (but not minimal).
	if !jpd.IsIMap([][2]string{{"A", "B"}}, 0.01) {
		t.Error("A->B should be a valid I-map (supergraph of minimal)")
	}
}

func TestJPDToFactor(t *testing.T) {
	jpd := buildTestJPD(t)
	f := jpd.ToFactor()
	if f == nil {
		t.Fatal("ToFactor returned nil")
	}

	vars := f.Variables()
	if len(vars) != 2 {
		t.Errorf("expected 2 variables, got %d", len(vars))
	}

	// Values should match.
	data := f.Values().Data()
	expected := []float64{0.3, 0.3, 0.2, 0.2}
	for i, exp := range expected {
		if data[i] != exp {
			t.Errorf("value[%d]: expected %f, got %f", i, exp, data[i])
		}
	}
}

func TestPMap(t *testing.T) {
	jpd := buildTestJPD(t)
	pmap := jpd.PMap(0.01)
	if pmap == "{}" {
		t.Error("expected non-empty PMap for independent variables")
	}
	if !strings.Contains(pmap, "_|_") {
		t.Error("expected _|_ in PMap output")
	}
}

func TestPMapDependent(t *testing.T) {
	jpd := buildDependentJPD(t)
	pmap := jpd.PMap(0.01)
	// No marginal independencies expected.
	if strings.Contains(pmap, "A _|_ B") && !strings.Contains(pmap, "|") {
		t.Error("should not find marginal independence for dependent variables")
	}
}
