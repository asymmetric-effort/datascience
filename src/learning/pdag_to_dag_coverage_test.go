//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
)

// ---------------------------------------------------------------------------
// pdagToDAG coverage tests — cycle detection during edge orientation
// ---------------------------------------------------------------------------

func TestPDAGToDAG_FullyDirected(t *testing.T) {
	// All edges are already directed. No undirected edges to orient.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddDirectedEdge("A", "B")
	pdag.AddDirectedEdge("B", "C")

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}

	es := make(map[[2]string]bool)
	for _, e := range edges {
		es[e] = true
	}
	if !es[[2]string{"A", "B"}] {
		t.Error("expected A->B")
	}
	if !es[[2]string{"B", "C"}] {
		t.Error("expected B->C")
	}
}

func TestPDAGToDAG_SingleUndirectedEdge(t *testing.T) {
	// One undirected edge, trivially orientable.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("X", "Y")
	pdag.AddUndirectedEdge("X", "Y")

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
}

func TestPDAGToDAG_UndirectedChain(t *testing.T) {
	// A -- B -- C (all undirected). Should orient without cycles.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddUndirectedEdge("A", "B")
	pdag.AddUndirectedEdge("B", "C")

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}

	// Verify it's a valid DAG (checked internally by BN.AddEdge).
	nodes := bn.Nodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}
}

func TestPDAGToDAG_UndirectedForceReversal(t *testing.T) {
	// Create a situation where the first orientation attempt creates a cycle,
	// forcing the reversal path.
	//
	// Directed: A -> B
	// Undirected: B -- A (this is impossible since A->B directed already exists)
	//
	// Better approach: Directed B -> A, Undirected A -- C, Directed C -> B.
	// When orienting A--C, try A->C first. A->C plus C->B plus B->A gives a cycle!
	// So it must orient C->A instead.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddDirectedEdge("B", "A")   // B -> A
	pdag.AddDirectedEdge("C", "B")   // C -> B
	pdag.AddUndirectedEdge("A", "C") // A -- C

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}

	es := make(map[[2]string]bool)
	for _, e := range edges {
		es[e] = true
	}

	// A->C would create cycle A->C->B->A. So must be C->A.
	if es[[2]string{"A", "C"}] {
		t.Error("A->C would create a cycle; expected C->A")
	}
	if !es[[2]string{"C", "A"}] {
		t.Error("expected C->A after reversal")
	}
}

func TestPDAGToDAG_MixedDirectedAndUndirected(t *testing.T) {
	// Directed: X -> Z
	// Undirected: X -- Y, Y -- Z
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("X", "Y", "Z")
	pdag.AddDirectedEdge("X", "Z")
	pdag.AddUndirectedEdge("X", "Y")
	pdag.AddUndirectedEdge("Y", "Z")

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
}

func TestPDAGToDAG_EmptyPDAG(t *testing.T) {
	pdag := graphgo.NewPDAG()
	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}
	if len(bn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(bn.Nodes()))
	}
}

func TestPDAGToDAG_SingleNode(t *testing.T) {
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A")
	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}
	if len(bn.Nodes()) != 1 {
		t.Errorf("expected 1 node, got %d", len(bn.Nodes()))
	}
}

func TestPDAGToDAG_MultipleUndirectedWithCycleRisk(t *testing.T) {
	// Diamond: A -- B, A -- C, B -- D, C -- D
	// with directed edge D -> A.
	// Orienting A->B, A->C, then B->D creates B->D but C->D also needed.
	// But D -> A already exists, so A->B->D->A cycle? Only if B->D is chosen.
	// Actually: D->A is directed, A--B undirected.
	// Try A->B: no cycle (D->A->B path, but no back path B->...->D yet).
	// Try A->C: no cycle.
	// Try B->D: cycle D->A->B->D! So must reverse to D->B.
	// Try C->D: cycle D->A->C->D! So must reverse to D->C.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C", "D")
	pdag.AddDirectedEdge("D", "A")
	pdag.AddUndirectedEdge("A", "B")
	pdag.AddUndirectedEdge("A", "C")
	pdag.AddUndirectedEdge("B", "D")
	pdag.AddUndirectedEdge("C", "D")

	bn, err := pdagToDAG(pdag)
	if err != nil {
		t.Fatalf("pdagToDAG failed: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 5 {
		t.Fatalf("expected 5 edges, got %d: %v", len(edges), edges)
	}

	// Verify no cycle by checking it's a valid DAG.
	es := make(map[[2]string]bool)
	for _, e := range edges {
		es[e] = true
	}
	// D->A should be preserved.
	if !es[[2]string{"D", "A"}] {
		t.Error("expected D->A to be preserved")
	}
}
