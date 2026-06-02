//go:build unit

package models

import (
	"testing"
)

func TestNaiveBayesAddEdge_Valid(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// Edge from class to feature already exists, so this should fail with "already exists".
	err = nb.AddEdge("C", "X")
	if err == nil {
		t.Error("expected error for duplicate edge, got nil")
	}
}

func TestNaiveBayesAddEdge_NonClassParent(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// Edge from non-class node should fail.
	err = nb.AddEdge("X", "Y")
	if err == nil {
		t.Error("expected error for non-class parent")
	}
}

func TestNaiveBayesAddEdge_NonFeatureChild(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// Edge to non-feature should fail.
	err = nb.AddEdge("C", "Z")
	if err == nil {
		t.Error("expected error for non-feature child")
	}
}

func TestNaiveBayesAddEdgesFrom_Valid(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y", "Z"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// Edges already exist so this should return an error.
	err = nb.AddEdgesFrom("C", []string{"X", "Y"})
	if err == nil {
		t.Error("expected error for duplicate edges")
	}
}

func TestNaiveBayesAddEdgesFrom_InvalidParent(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	err = nb.AddEdgesFrom("X", []string{"Y"})
	if err == nil {
		t.Error("expected error for non-class parent in AddEdgesFrom")
	}
}

func TestNaiveBayesActiveTrailNodes_NoObserved(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y", "Z"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// From C with no observations: all features should be reachable.
	trail, err := nb.ActiveTrailNodes("C", nil)
	if err != nil {
		t.Fatalf("ActiveTrailNodes: %v", err)
	}
	if len(trail) != 3 {
		t.Errorf("expected 3 active trail nodes from C, got %d: %v", len(trail), trail)
	}
}

func TestNaiveBayesActiveTrailNodes_FromFeature(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y", "Z"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// From X with no observations: C and through C, Y and Z.
	trail, err := nb.ActiveTrailNodes("X", nil)
	if err != nil {
		t.Fatalf("ActiveTrailNodes: %v", err)
	}
	if len(trail) != 3 {
		t.Errorf("expected 3 active trail nodes from X (C, Y, Z), got %d: %v", len(trail), trail)
	}
}

func TestNaiveBayesActiveTrailNodes_ClassObserved(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y", "Z"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// From X with C observed: C blocks the path, X is isolated from Y, Z.
	observed := map[string]bool{"C": true}
	trail, err := nb.ActiveTrailNodes("X", observed)
	if err != nil {
		t.Fatalf("ActiveTrailNodes: %v", err)
	}
	if len(trail) != 0 {
		t.Errorf("expected 0 active trail nodes from X with C observed, got %d: %v", len(trail), trail)
	}
}

func TestNaiveBayesActiveTrailNodes_InvalidNode(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	_, err = nb.ActiveTrailNodes("Q", nil)
	if err == nil {
		t.Error("expected error for nonexistent node")
	}
}

func TestNaiveBayesLocalIndependencies(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y", "Z"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	ind := nb.LocalIndependencies()
	if ind == nil {
		t.Fatal("LocalIndependencies returned nil")
	}

	// With 3 features, each feature is independent of the other 2 given C.
	// That gives us 3 assertions.
	assertions := ind.GetAssertions()
	if len(assertions) != 3 {
		t.Errorf("expected 3 independence assertions, got %d", len(assertions))
	}
}

func TestNaiveBayesLocalIndependencies_TwoFeatures(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	ind := nb.LocalIndependencies()
	assertions := ind.GetAssertions()

	// 2 features: (X _|_ Y | C) and (Y _|_ X | C) are equivalent, so
	// after deduplication we get 1 assertion.
	if len(assertions) != 1 {
		t.Errorf("expected 1 independence assertion (deduplicated), got %d", len(assertions))
	}

	// Check that the given set is always the class variable.
	for _, a := range assertions {
		given := a.Given()
		if len(given) != 1 || given[0] != "C" {
			t.Errorf("expected given=[C], got %v", given)
		}
	}
}

func TestNaiveBayesActiveTrailNodes_FeatureObserved(t *testing.T) {
	nb, err := NewNaiveBayes("C", []string{"X", "Y"})
	if err != nil {
		t.Fatalf("NewNaiveBayes: %v", err)
	}

	// From C with X observed: Y should still be reachable.
	observed := map[string]bool{"X": true}
	trail, err := nb.ActiveTrailNodes("C", observed)
	if err != nil {
		t.Fatalf("ActiveTrailNodes: %v", err)
	}
	found := false
	for _, n := range trail {
		if n == "Y" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Y in active trail from C with X observed, got %v", trail)
	}
}
