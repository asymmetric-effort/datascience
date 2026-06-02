//go:build unit

package models

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

func buildSimpleClusterGraph(t *testing.T) *ClusterGraph {
	t.Helper()
	cg := NewClusterGraph()

	f1, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	f2, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 3}, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6})

	cg.AddCluster([]string{"A", "B"}, []*factors.DiscreteFactor{f1})
	cg.AddCluster([]string{"B", "C"}, []*factors.DiscreteFactor{f2})
	_ = cg.AddEdge(0, 1, []string{"B"})

	return cg
}

func TestClusterGraphAddNode(t *testing.T) {
	cg := NewClusterGraph()
	idx := cg.AddNode([]string{"X", "Y"})
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}

	clusters := cg.Clusters()
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	// AddNode uses nil factors.
	if len(clusters[0].Factors) != 0 {
		t.Errorf("expected 0 factors, got %d", len(clusters[0].Factors))
	}
}

func TestClusterGraphAddNodeMultiple(t *testing.T) {
	cg := NewClusterGraph()
	idx0 := cg.AddNode([]string{"A"})
	idx1 := cg.AddNode([]string{"B"})
	if idx0 != 0 || idx1 != 1 {
		t.Errorf("expected indices 0, 1; got %d, %d", idx0, idx1)
	}
}

func TestClusterGraphAddFactors(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A", "B"}, nil)

	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})

	if err := cg.AddFactors(0, []*factors.DiscreteFactor{f}); err != nil {
		t.Fatalf("AddFactors: %v", err)
	}

	clusters := cg.Clusters()
	if len(clusters[0].Factors) != 1 {
		t.Errorf("expected 1 factor, got %d", len(clusters[0].Factors))
	}
}

func TestClusterGraphAddFactorsMultiple(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A", "B"}, nil)

	f1, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	f2, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.6, 0.4})

	_ = cg.AddFactors(0, []*factors.DiscreteFactor{f1})
	_ = cg.AddFactors(0, []*factors.DiscreteFactor{f2})

	clusters := cg.Clusters()
	if len(clusters[0].Factors) != 2 {
		t.Errorf("expected 2 factors, got %d", len(clusters[0].Factors))
	}
}

func TestClusterGraphAddFactorsInvalidIndex(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A"}, nil)

	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	if err := cg.AddFactors(5, []*factors.DiscreteFactor{f}); err == nil {
		t.Error("expected error for out-of-range index")
	}
	if err := cg.AddFactors(-1, []*factors.DiscreteFactor{f}); err == nil {
		t.Error("expected error for negative index")
	}
}

func TestClusterGraphGetFactors(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	allFactors := cg.GetFactors()
	if len(allFactors) != 2 {
		t.Errorf("expected 2 factors, got %d", len(allFactors))
	}
}

func TestClusterGraphGetFactorsEmpty(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A"}, nil)
	allFactors := cg.GetFactors()
	if len(allFactors) != 0 {
		t.Errorf("expected 0 factors, got %d", len(allFactors))
	}
}

func TestClusterGraphRemoveFactors(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	cg.RemoveFactors()

	for i, c := range cg.Clusters() {
		if len(c.Factors) != 0 {
			t.Errorf("cluster %d still has %d factors", i, len(c.Factors))
		}
	}
}

func TestClusterGraphCliqueBeliefs(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	beliefs, err := cg.CliqueBeliefs()
	if err != nil {
		t.Fatalf("CliqueBeliefs: %v", err)
	}

	if len(beliefs) != 2 {
		t.Errorf("expected 2 beliefs, got %d", len(beliefs))
	}

	// Each belief should be normalized (sum to 1).
	for idx, belief := range beliefs {
		data := belief.Values().Data()
		sum := 0.0
		for _, v := range data {
			sum += v
		}
		if math.Abs(sum-1.0) > 1e-6 {
			t.Errorf("belief %d sums to %f, expected 1.0", idx, sum)
		}
	}
}

func TestClusterGraphCliqueBeliefsNoFactors(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A"}, nil)

	beliefs, err := cg.CliqueBeliefs()
	if err != nil {
		t.Fatalf("CliqueBeliefs: %v", err)
	}
	// Cluster with no factors should not have a belief.
	if len(beliefs) != 0 {
		t.Errorf("expected 0 beliefs, got %d", len(beliefs))
	}
}

func TestClusterGraphGetCardinality(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	card := cg.GetCardinality()

	if card["A"] != 2 {
		t.Errorf("expected A cardinality 2, got %d", card["A"])
	}
	if card["B"] != 2 {
		t.Errorf("expected B cardinality 2, got %d", card["B"])
	}
	if card["C"] != 3 {
		t.Errorf("expected C cardinality 3, got %d", card["C"])
	}
}

func TestClusterGraphGetCardinalityEmpty(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"A"}, nil)
	card := cg.GetCardinality()
	if len(card) != 0 {
		t.Errorf("expected empty cardinality map, got %v", card)
	}
}

func TestClusterGraphGetPartitionFunction(t *testing.T) {
	cg := NewClusterGraph()
	f, _ := factors.NewDiscreteFactor([]string{"X"}, []int{2}, []float64{0.3, 0.7})
	cg.AddCluster([]string{"X"}, []*factors.DiscreteFactor{f})

	z, err := cg.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if math.Abs(z-1.0) > 1e-6 {
		t.Errorf("expected partition function 1.0, got %f", z)
	}
}

func TestClusterGraphGetPartitionFunctionNoFactors(t *testing.T) {
	cg := NewClusterGraph()
	cg.AddCluster([]string{"X"}, nil)
	_, err := cg.GetPartitionFunction()
	if err == nil {
		t.Error("expected error for no factors")
	}
}

func TestClusterGraphCopy(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	cpy := cg.Copy()

	if len(cpy.Clusters()) != len(cg.Clusters()) {
		t.Errorf("copy has %d clusters, expected %d", len(cpy.Clusters()), len(cg.Clusters()))
	}
	if len(cpy.Edges()) != len(cg.Edges()) {
		t.Errorf("copy has %d edges, expected %d", len(cpy.Edges()), len(cg.Edges()))
	}

	// Verify copy is independent.
	cpy.RemoveFactors()
	origFactors := cg.GetFactors()
	if len(origFactors) != 2 {
		t.Error("original was affected by copy modification")
	}
}

func TestClusterGraphCopyCheckModel(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	cpy := cg.Copy()
	if err := cpy.CheckModel(); err != nil {
		t.Fatalf("copied cluster graph CheckModel: %v", err)
	}
}

func TestClusterGraphCopyDeepFactors(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	cpy := cg.Copy()

	// Verify factors are deep-copied.
	origClusters := cg.Clusters()
	cpyClusters := cpy.Clusters()

	if origClusters[0].Factors[0] == cpyClusters[0].Factors[0] {
		t.Error("factors should be deep-copied, not shared")
	}
}

func TestClusterGraphCopyEdges(t *testing.T) {
	cg := buildSimpleClusterGraph(t)
	cpy := cg.Copy()

	cpyEdges := cpy.Edges()
	if len(cpyEdges) != 1 {
		t.Fatalf("expected 1 edge in copy, got %d", len(cpyEdges))
	}
	if cpyEdges[0].SepSet[0] != "B" {
		t.Errorf("expected sep set [B], got %v", cpyEdges[0].SepSet)
	}
}
