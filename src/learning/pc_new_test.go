//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
)

func TestOrientColliders(t *testing.T) {
	// Build a PDAG skeleton: A - C - B (all undirected), A and B not adjacent.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddUndirectedEdge("A", "C")
	pdag.AddUndirectedEdge("B", "C")

	// Separating set for (A, B) does not contain C -> orient as A -> C <- B.
	sepSets := map[[2]string][]string{
		{"A", "B"}: {}, // C is not in the sep set.
	}

	OrientColliders(pdag, sepSets)

	// Check that A -> C and B -> C are now directed.
	if !pdag.HasDirectedEdge("A", "C") {
		t.Error("expected directed edge A -> C")
	}
	if !pdag.HasDirectedEdge("B", "C") {
		t.Error("expected directed edge B -> C")
	}
}

func TestOrientCollidersNoCollider(t *testing.T) {
	// Build: A - C - B, but sep set of (A,B) contains C -> no collider.
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddUndirectedEdge("A", "C")
	pdag.AddUndirectedEdge("B", "C")

	sepSets := map[[2]string][]string{
		{"A", "B"}: {"C"}, // C IS in the sep set.
	}

	OrientColliders(pdag, sepSets)

	// Edges should remain undirected.
	if pdag.HasDirectedEdge("A", "C") {
		t.Error("edge A -> C should not be directed")
	}
	if pdag.HasDirectedEdge("B", "C") {
		t.Error("edge B -> C should not be directed")
	}
	if !pdag.HasUndirectedEdge("A", "C") {
		t.Error("expected undirected edge A - C")
	}
}

func TestOrientCollidersAdjacentNotCollider(t *testing.T) {
	// A - C - B and A - B (shielded triple, not a collider).
	pdag := graphgo.NewPDAG()
	pdag.AddNodes("A", "B", "C")
	pdag.AddUndirectedEdge("A", "C")
	pdag.AddUndirectedEdge("B", "C")
	pdag.AddUndirectedEdge("A", "B")

	sepSets := map[[2]string][]string{}

	OrientColliders(pdag, sepSets)

	// All edges should remain undirected (A and B are adjacent).
	if pdag.HasDirectedEdge("A", "C") || pdag.HasDirectedEdge("B", "C") {
		t.Error("shielded triple should not produce a collider orientation")
	}
}
