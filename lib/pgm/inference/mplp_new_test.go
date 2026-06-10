//go:build unit

package inference

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

func buildMPLPFactors(t *testing.T) []*factors.DiscreteFactor {
	t.Helper()

	fA, err := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.6, 0.4})
	if err != nil {
		t.Fatal(err)
	}
	fB, err := factors.NewDiscreteFactor([]string{"B"}, []int{2}, []float64{0.3, 0.7})
	if err != nil {
		t.Fatal(err)
	}
	fAB, err := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2},
		[]float64{0.9, 0.1, 0.3, 0.7})
	if err != nil {
		t.Fatal(err)
	}
	fBC, err := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2},
		[]float64{0.8, 0.2, 0.4, 0.6})
	if err != nil {
		t.Fatal(err)
	}

	return []*factors.DiscreteFactor{fA, fB, fAB, fBC}
}

func TestMPLPFindTriangles(t *testing.T) {
	fs := buildMPLPFactors(t)
	m := NewMPLP(fs)
	triangles := m.FindTriangles()
	// A, B, C form a triangle since A-B (from fAB) and B-C (from fBC),
	// but A and C don't share a factor directly.
	t.Logf("found %d triangles", len(triangles))
}

func TestMPLPGetIntegralityGap(t *testing.T) {
	fs := buildMPLPFactors(t)
	m := NewMPLP(fs)
	gap, err := m.GetIntegralityGap([]string{"A", "B", "C"}, nil, 100, 1e-6)
	if err != nil {
		t.Fatalf("GetIntegralityGap failed: %v", err)
	}
	if gap < 0 {
		t.Errorf("expected non-negative gap, got %f", gap)
	}
	t.Logf("integrality gap: %f", gap)
}

func TestMPLPQuery(t *testing.T) {
	fs := buildMPLPFactors(t)
	m := NewMPLP(fs)
	result, err := m.Query([]string{"A"}, nil, 100, 1e-6)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if result == nil {
		t.Fatal("Query returned nil")
	}

	vars := result.Variables()
	if len(vars) != 1 || vars[0] != "A" {
		t.Errorf("expected [A], got %v", vars)
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

func TestMPLPQueryErrors(t *testing.T) {
	fs := buildMPLPFactors(t)
	m := NewMPLP(fs)

	_, err := m.Query(nil, nil, 100, 1e-6)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestMPLPMAPQuery(t *testing.T) {
	fs := buildMPLPFactors(t)
	m := NewMPLP(fs)
	assignment, err := m.MAPQuery([]string{"A", "B", "C"}, nil)
	if err != nil {
		t.Fatalf("MAPQuery failed: %v", err)
	}
	if len(assignment) != 3 {
		t.Errorf("expected 3 assignments, got %d", len(assignment))
	}
	for _, v := range []string{"A", "B", "C"} {
		if _, ok := assignment[v]; !ok {
			t.Errorf("missing assignment for %s", v)
		}
	}
}
