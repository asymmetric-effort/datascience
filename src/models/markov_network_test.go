//go:build unit

package models

import (
	"math"
	"sort"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// buildTriangleMRF constructs a simple 3-node triangle Markov network:
//
//	A -- B
//	|  /
//	C
//
// Each variable is binary (cardinality 2). Three pairwise factors are defined.
func buildTriangleMRF(t *testing.T) *MarkovNetwork {
	t.Helper()

	mn := NewMarkovNetwork()
	for _, node := range []string{"A", "B", "C"} {
		if err := mn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q): %v", node, err)
		}
	}
	for _, edge := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}} {
		if err := mn.AddEdge(edge[0], edge[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", edge[0], edge[1], err)
		}
	}

	// Factor phi_AB: scope {A, B}, cardinality {2, 2}
	// Values: phi(0,0)=1, phi(0,1)=2, phi(1,0)=3, phi(1,1)=4
	fAB, err := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("NewDiscreteFactor(AB): %v", err)
	}
	// Factor phi_AC: scope {A, C}, cardinality {2, 2}
	fAC, err := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{2, 1, 1, 2})
	if err != nil {
		t.Fatalf("NewDiscreteFactor(AC): %v", err)
	}
	// Factor phi_BC: scope {B, C}, cardinality {2, 2}
	fBC, err := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})
	if err != nil {
		t.Fatalf("NewDiscreteFactor(BC): %v", err)
	}

	for _, f := range []*factors.DiscreteFactor{fAB, fAC, fBC} {
		if err := mn.AddFactor(f); err != nil {
			t.Fatalf("AddFactor: %v", err)
		}
	}

	return mn
}

func TestNewMarkovNetwork(t *testing.T) {
	mn := NewMarkovNetwork()
	if len(mn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(mn.Nodes()))
	}
	if len(mn.GetFactors()) != 0 {
		t.Errorf("expected 0 factors, got %d", len(mn.GetFactors()))
	}
}

func TestMarkovNetwork_AddNodeAndEdge(t *testing.T) {
	mn := NewMarkovNetwork()
	if err := mn.AddNode("X"); err != nil {
		t.Fatal(err)
	}
	if err := mn.AddNode("Y"); err != nil {
		t.Fatal(err)
	}

	// Duplicate node.
	if err := mn.AddNode("X"); err == nil {
		t.Error("expected error adding duplicate node")
	}

	if err := mn.AddEdge("X", "Y"); err != nil {
		t.Fatal(err)
	}

	// Duplicate edge.
	if err := mn.AddEdge("X", "Y"); err == nil {
		t.Error("expected error adding duplicate edge")
	}

	// Edge with missing node.
	if err := mn.AddEdge("X", "Z"); err == nil {
		t.Error("expected error adding edge to missing node")
	}

	nodes := mn.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}

	edges := mn.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
}

func TestMarkovNetwork_NodesEdgesNeighbors(t *testing.T) {
	mn := buildTriangleMRF(t)

	nodes := mn.Nodes()
	expected := []string{"A", "B", "C"}
	if len(nodes) != len(expected) {
		t.Fatalf("expected %d nodes, got %d", len(expected), len(nodes))
	}
	for i, n := range nodes {
		if n != expected[i] {
			t.Errorf("node[%d] = %q, want %q", i, n, expected[i])
		}
	}

	edges := mn.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}

	neighborsA := mn.Neighbors("A")
	sort.Strings(neighborsA)
	if len(neighborsA) != 2 || neighborsA[0] != "B" || neighborsA[1] != "C" {
		t.Errorf("Neighbors(A) = %v, want [B C]", neighborsA)
	}
}

func TestMarkovNetwork_AddFactor(t *testing.T) {
	mn := NewMarkovNetwork()
	_ = mn.AddNode("A")
	_ = mn.AddNode("B")
	_ = mn.AddEdge("A", "B")

	// nil factor.
	if err := mn.AddFactor(nil); err == nil {
		t.Error("expected error adding nil factor")
	}

	// Factor referencing unknown node.
	fBad, _ := factors.NewDiscreteFactor([]string{"A", "Z"}, []int{2, 2}, []float64{1, 2, 3, 4})
	if err := mn.AddFactor(fBad); err == nil {
		t.Error("expected error adding factor with unknown node")
	}

	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	if err := mn.AddFactor(f); err != nil {
		t.Fatalf("AddFactor: %v", err)
	}
	if len(mn.GetFactors()) != 1 {
		t.Errorf("expected 1 factor, got %d", len(mn.GetFactors()))
	}
}

func TestMarkovNetwork_RemoveFactor(t *testing.T) {
	mn := buildTriangleMRF(t)
	if len(mn.GetFactors()) != 3 {
		t.Fatalf("expected 3 factors, got %d", len(mn.GetFactors()))
	}

	mn.RemoveFactor([]string{"A", "B"})
	if len(mn.GetFactors()) != 2 {
		t.Errorf("expected 2 factors after removal, got %d", len(mn.GetFactors()))
	}

	// Removing non-existent factor does nothing.
	mn.RemoveFactor([]string{"X", "Y"})
	if len(mn.GetFactors()) != 2 {
		t.Errorf("expected 2 factors, got %d", len(mn.GetFactors()))
	}
}

func TestMarkovNetwork_GetFactorsOf(t *testing.T) {
	mn := buildTriangleMRF(t)

	fs := mn.GetFactorsOf("A")
	if len(fs) != 2 {
		t.Errorf("expected 2 factors for A, got %d", len(fs))
	}

	fs = mn.GetFactorsOf("B")
	if len(fs) != 2 {
		t.Errorf("expected 2 factors for B, got %d", len(fs))
	}

	fs = mn.GetFactorsOf("Z")
	if fs != nil {
		t.Errorf("expected nil for unknown variable, got %v", fs)
	}
}

func TestMarkovNetwork_CheckModel(t *testing.T) {
	mn := buildTriangleMRF(t)
	if err := mn.CheckModel(); err != nil {
		t.Errorf("CheckModel failed on valid model: %v", err)
	}
}

func TestMarkovNetwork_CheckModel_NoFactors(t *testing.T) {
	mn := NewMarkovNetwork()
	_ = mn.AddNode("A")
	if err := mn.CheckModel(); err == nil {
		t.Error("expected error for model with no factors")
	}
}

func TestMarkovNetwork_CheckModel_UncoveredNode(t *testing.T) {
	mn := NewMarkovNetwork()
	_ = mn.AddNode("A")
	_ = mn.AddNode("B")
	_ = mn.AddNode("C")
	_ = mn.AddEdge("A", "B")
	_ = mn.AddEdge("A", "C")
	_ = mn.AddEdge("B", "C")

	// Only add factor for A-B, leaving C uncovered.
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	_ = mn.AddFactor(f)

	if err := mn.CheckModel(); err == nil {
		t.Error("expected error for uncovered node C")
	}
}

func TestMarkovNetwork_CheckModel_MissingEdge(t *testing.T) {
	mn := NewMarkovNetwork()
	_ = mn.AddNode("A")
	_ = mn.AddNode("B")
	_ = mn.AddNode("C")
	_ = mn.AddEdge("A", "B")
	// No edge between A and C.

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	_ = mn.AddFactor(fAB)

	// Factor referencing A and C without an edge between them.
	// We need to add node C to the factor but there's no edge A-C.
	// Since AddFactor checks nodes exist but not edges, we need to add C edge for factor.
	// Actually, AddFactor only checks nodes exist. CheckModel checks edges.
	fAC, _ := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})
	_ = mn.AddFactor(fAC)

	// Also cover C with a unary factor.
	fC, _ := factors.NewDiscreteFactor([]string{"C"}, []int{2}, []float64{1, 1})
	_ = mn.AddFactor(fC)

	if err := mn.CheckModel(); err == nil {
		t.Error("expected error for factor with missing edge between variables")
	}
}

func TestMarkovNetwork_GetPartitionFunction(t *testing.T) {
	mn := buildTriangleMRF(t)

	// Manually compute Z = sum over all (a,b,c) in {0,1}^3 of phi_AB(a,b) * phi_AC(a,c) * phi_BC(b,c)
	//
	// phi_AB: (0,0)=1, (0,1)=2, (1,0)=3, (1,1)=4
	// phi_AC: (0,0)=2, (0,1)=1, (1,0)=1, (1,1)=2
	// phi_BC: (0,0)=1, (0,1)=1, (1,0)=1, (1,1)=1
	//
	// (0,0,0): 1*2*1 = 2
	// (0,0,1): 1*1*1 = 1
	// (0,1,0): 2*2*1 = 4
	// (0,1,1): 2*1*1 = 2
	// (1,0,0): 3*1*1 = 3
	// (1,0,1): 3*2*1 = 6
	// (1,1,0): 4*1*1 = 4
	// (1,1,1): 4*2*1 = 8
	// Z = 2+1+4+2+3+6+4+8 = 30

	z, err := mn.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if math.Abs(z-30.0) > 1e-10 {
		t.Errorf("partition function = %f, want 30.0", z)
	}
}

func TestMarkovNetwork_GetPartitionFunction_NoFactors(t *testing.T) {
	mn := NewMarkovNetwork()
	_, err := mn.GetPartitionFunction()
	if err == nil {
		t.Error("expected error for partition function with no factors")
	}
}

func TestMarkovNetwork_MarkovBlanket(t *testing.T) {
	mn := buildTriangleMRF(t)

	blanket := mn.MarkovBlanket("A")
	sort.Strings(blanket)
	if len(blanket) != 2 || blanket[0] != "B" || blanket[1] != "C" {
		t.Errorf("MarkovBlanket(A) = %v, want [B C]", blanket)
	}

	blanket = mn.MarkovBlanket("B")
	sort.Strings(blanket)
	if len(blanket) != 2 || blanket[0] != "A" || blanket[1] != "C" {
		t.Errorf("MarkovBlanket(B) = %v, want [A C]", blanket)
	}
}

func TestMarkovNetwork_ToJunctionTree(t *testing.T) {
	mn := buildTriangleMRF(t)

	jt, err := mn.ToJunctionTree()
	if err != nil {
		t.Fatalf("ToJunctionTree: %v", err)
	}

	cliques := jt.Cliques()
	if len(cliques) == 0 {
		t.Fatal("expected at least 1 clique")
	}

	// A triangle graph is already chordal, so the single maximal clique is {A, B, C}.
	if len(cliques) != 1 {
		t.Fatalf("expected 1 clique for triangle, got %d: %v", len(cliques), cliques)
	}
	c := cliques[0]
	sort.Strings(c)
	if len(c) != 3 || c[0] != "A" || c[1] != "B" || c[2] != "C" {
		t.Errorf("clique = %v, want [A B C]", c)
	}

	// All 3 factors should be assigned to the single clique.
	cfs := jt.GetCliqueFactors([]string{"A", "B", "C"})
	if len(cfs) != 3 {
		t.Errorf("expected 3 factors in clique, got %d", len(cfs))
	}

	// Verify the running intersection property.
	if err := jt.CheckModel(); err != nil {
		t.Errorf("junction tree CheckModel failed: %v", err)
	}
}

func TestMarkovNetwork_ToJunctionTree_Chain(t *testing.T) {
	// Build a chain: A -- B -- C (no A-C edge).
	mn := NewMarkovNetwork()
	for _, node := range []string{"A", "B", "C"} {
		_ = mn.AddNode(node)
	}
	_ = mn.AddEdge("A", "B")
	_ = mn.AddEdge("B", "C")

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})
	_ = mn.AddFactor(fAB)
	_ = mn.AddFactor(fBC)

	jt, err := mn.ToJunctionTree()
	if err != nil {
		t.Fatalf("ToJunctionTree: %v", err)
	}

	cliques := jt.Cliques()
	// A chain A-B-C should yield two cliques: {A,B} and {B,C}.
	if len(cliques) != 2 {
		t.Fatalf("expected 2 cliques for chain, got %d: %v", len(cliques), cliques)
	}

	if err := jt.CheckModel(); err != nil {
		t.Errorf("junction tree CheckModel failed: %v", err)
	}

	// Separator should contain B.
	seps := jt.SeparatorSets()
	found := false
	for _, sep := range seps {
		if len(sep) == 1 && sep[0] == "B" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected separator {B}, got %v", seps)
	}
}

func TestMarkovNetwork_Copy(t *testing.T) {
	mn := buildTriangleMRF(t)
	cp := mn.Copy()

	// Verify the copy has the same structure.
	if len(cp.Nodes()) != len(mn.Nodes()) {
		t.Errorf("copy nodes %d != original %d", len(cp.Nodes()), len(mn.Nodes()))
	}
	if len(cp.Edges()) != len(mn.Edges()) {
		t.Errorf("copy edges %d != original %d", len(cp.Edges()), len(mn.Edges()))
	}
	if len(cp.GetFactors()) != len(mn.GetFactors()) {
		t.Errorf("copy factors %d != original %d", len(cp.GetFactors()), len(mn.GetFactors()))
	}

	// Verify independence: add a node to the copy, original should be unaffected.
	_ = cp.AddNode("D")
	if len(mn.Nodes()) != 3 {
		t.Error("modifying copy affected original")
	}
}

// --- DiscreteMarkovNetwork tests ---

func TestNewDiscreteMarkovNetwork(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	if len(dmn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(dmn.Nodes()))
	}
}

func TestDiscreteMarkovNetwork_CheckModel(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	for _, node := range []string{"A", "B", "C"} {
		_ = dmn.AddNode(node)
	}
	for _, edge := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}} {
		_ = dmn.AddEdge(edge[0], edge[1])
	}

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	fAC, _ := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{2, 1, 1, 2})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})

	for _, f := range []*factors.DiscreteFactor{fAB, fAC, fBC} {
		if err := dmn.AddFactor(f); err != nil {
			t.Fatalf("AddFactor: %v", err)
		}
	}

	if err := dmn.CheckModel(); err != nil {
		t.Errorf("CheckModel failed on valid discrete model: %v", err)
	}
}

func TestDiscreteMarkovNetwork_AddFactor_NilFactor(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	if err := dmn.AddFactor(nil); err == nil {
		t.Error("expected error adding nil factor")
	}
}

func TestDiscreteMarkovNetwork_CheckModel_NegativeValues(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	_ = dmn.AddNode("A")
	_ = dmn.AddNode("B")
	_ = dmn.AddEdge("A", "B")

	// Factor with negative value: use MarkovNetwork.AddFactor to bypass
	// the discrete check on AddFactor, then CheckModel should catch it.
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, -1, 2, 3})
	// Use the embedded MarkovNetwork's AddFactor to bypass discrete checks.
	if err := dmn.MarkovNetwork.AddFactor(f); err != nil {
		t.Fatalf("MarkovNetwork.AddFactor: %v", err)
	}

	if err := dmn.CheckModel(); err == nil {
		t.Error("expected error for factor with negative values")
	}
}

func TestDiscreteMarkovNetwork_Copy(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	_ = dmn.AddNode("A")
	_ = dmn.AddNode("B")
	_ = dmn.AddEdge("A", "B")

	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	_ = dmn.AddFactor(f)

	cp := dmn.Copy()
	if len(cp.Nodes()) != 2 {
		t.Errorf("copy nodes = %d, want 2", len(cp.Nodes()))
	}
	if len(cp.GetFactors()) != 1 {
		t.Errorf("copy factors = %d, want 1", len(cp.GetFactors()))
	}

	// Verify independence.
	_ = cp.AddNode("C")
	if len(dmn.Nodes()) != 2 {
		t.Error("modifying copy affected original")
	}
}

func TestDiscreteMarkovNetwork_PartitionFunction(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	for _, node := range []string{"A", "B", "C"} {
		_ = dmn.AddNode(node)
	}
	for _, edge := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}} {
		_ = dmn.AddEdge(edge[0], edge[1])
	}

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	fAC, _ := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{2, 1, 1, 2})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})

	for _, f := range []*factors.DiscreteFactor{fAB, fAC, fBC} {
		_ = dmn.AddFactor(f)
	}

	z, err := dmn.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if math.Abs(z-30.0) > 1e-10 {
		t.Errorf("partition function = %f, want 30.0", z)
	}
}

func TestDiscreteMarkovNetwork_JunctionTree(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	for _, node := range []string{"A", "B", "C"} {
		_ = dmn.AddNode(node)
	}
	for _, edge := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "C"}} {
		_ = dmn.AddEdge(edge[0], edge[1])
	}

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	fAC, _ := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{2, 1, 1, 2})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{1, 1, 1, 1})

	for _, f := range []*factors.DiscreteFactor{fAB, fAC, fBC} {
		_ = dmn.AddFactor(f)
	}

	jt, err := dmn.ToJunctionTree()
	if err != nil {
		t.Fatalf("ToJunctionTree: %v", err)
	}

	if err := jt.CheckModel(); err != nil {
		t.Errorf("junction tree CheckModel failed: %v", err)
	}
}

func TestMarkovNetwork_UnaryFactor(t *testing.T) {
	// Test a model with a unary factor (single-variable potential).
	mn := NewMarkovNetwork()
	_ = mn.AddNode("A")
	_ = mn.AddNode("B")
	_ = mn.AddEdge("A", "B")

	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	fA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 1.5})
	fB, _ := factors.NewDiscreteFactor([]string{"B"}, []int{2}, []float64{1, 1})

	_ = mn.AddFactor(fAB)
	_ = mn.AddFactor(fA)
	_ = mn.AddFactor(fB)

	if err := mn.CheckModel(); err != nil {
		t.Errorf("CheckModel failed with unary factors: %v", err)
	}

	// Z = sum over (a,b) of phi_AB(a,b) * phi_A(a) * phi_B(b)
	// (0,0): 1*0.5*1 = 0.5
	// (0,1): 2*0.5*1 = 1.0
	// (1,0): 3*1.5*1 = 4.5
	// (1,1): 4*1.5*1 = 6.0
	// Z = 12.0
	z, err := mn.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if math.Abs(z-12.0) > 1e-10 {
		t.Errorf("partition function = %f, want 12.0", z)
	}
}
