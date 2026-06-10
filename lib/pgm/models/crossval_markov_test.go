//go:build unit

package models_test

import (
	"math"
	"sort"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/tests/testutil"
)

func TestCrossval_MarkovNetworkTriangle(t *testing.T) {
	ff := testutil.LoadFixtures(t, "markov_network/fixtures.json")
	tc := ff.FindTestCase(t, "markov_network_triangle")

	var input struct {
		Edges   [][]string `json:"edges"`
		Factors []struct {
			Variables   []string  `json:"variables"`
			Cardinality []int     `json:"cardinality"`
			Values      []float64 `json:"values"`
		} `json:"factors"`
	}
	tc.UnmarshalInput(t, &input)

	var expected struct {
		Nodes             []string   `json:"nodes"`
		Edges             [][]string `json:"edges"`
		NumNodes          int        `json:"num_nodes"`
		NumEdges          int        `json:"num_edges"`
		IsValid           bool       `json:"is_valid"`
		PartitionFunction float64    `json:"partition_function"`
		NumFactors        int        `json:"num_factors"`
	}
	tc.UnmarshalExpected(t, &expected)

	mn := models.NewMarkovNetwork()

	// Collect and add all unique nodes.
	nodeSet := make(map[string]bool)
	for _, edge := range input.Edges {
		nodeSet[edge[0]] = true
		nodeSet[edge[1]] = true
	}
	for node := range nodeSet {
		if err := mn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q) failed: %v", node, err)
		}
	}

	// Add edges.
	for _, edge := range input.Edges {
		if err := mn.AddEdge(edge[0], edge[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q) failed: %v", edge[0], edge[1], err)
		}
	}

	// Add factors.
	for i, fd := range input.Factors {
		f, err := factors.NewDiscreteFactor(fd.Variables, fd.Cardinality, fd.Values)
		if err != nil {
			t.Fatalf("NewDiscreteFactor (factor %d) failed: %v", i, err)
		}
		if err := mn.AddFactor(f); err != nil {
			t.Fatalf("AddFactor (factor %d) failed: %v", i, err)
		}
	}

	// Verify nodes.
	gotNodes := mn.Nodes()
	sort.Strings(gotNodes)
	expectedNodes := make([]string, len(expected.Nodes))
	copy(expectedNodes, expected.Nodes)
	sort.Strings(expectedNodes)

	if len(gotNodes) != len(expectedNodes) {
		t.Fatalf("nodes: expected %v, got %v", expectedNodes, gotNodes)
	}
	for i := range expectedNodes {
		if gotNodes[i] != expectedNodes[i] {
			t.Errorf("nodes[%d]: expected %q, got %q", i, expectedNodes[i], gotNodes[i])
		}
	}

	if len(gotNodes) != expected.NumNodes {
		t.Errorf("num_nodes: expected %d, got %d", expected.NumNodes, len(gotNodes))
	}

	// Verify edges.
	gotEdges := mn.Edges()
	if len(gotEdges) != expected.NumEdges {
		t.Errorf("num_edges: expected %d, got %d", expected.NumEdges, len(gotEdges))
	}

	// Verify check_model.
	err := mn.CheckModel()
	if expected.IsValid {
		if err != nil {
			t.Errorf("CheckModel: expected valid, got error: %v", err)
		}
	} else {
		if err == nil {
			t.Errorf("CheckModel: expected invalid, got nil error")
		}
	}

	// Verify factor count.
	gotFactors := mn.GetFactors()
	if len(gotFactors) != expected.NumFactors {
		t.Errorf("num_factors: expected %d, got %d", expected.NumFactors, len(gotFactors))
	}

	// Verify partition function.
	Z, err := mn.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction failed: %v", err)
	}
	if math.Abs(Z-expected.PartitionFunction) > 1e-6 {
		t.Errorf("partition_function: expected %f, got %f", expected.PartitionFunction, Z)
	}
}

func TestCrossval_MarkovNetworkChain(t *testing.T) {
	ff := testutil.LoadFixtures(t, "markov_network/fixtures.json")
	tc := ff.FindTestCase(t, "markov_network_chain")

	var input struct {
		Edges   [][]string `json:"edges"`
		Factors []struct {
			Variables   []string  `json:"variables"`
			Cardinality []int     `json:"cardinality"`
			Values      []float64 `json:"values"`
		} `json:"factors"`
	}
	tc.UnmarshalInput(t, &input)

	var expected struct {
		Nodes             []string   `json:"nodes"`
		Edges             [][]string `json:"edges"`
		NumNodes          int        `json:"num_nodes"`
		NumEdges          int        `json:"num_edges"`
		IsValid           bool       `json:"is_valid"`
		PartitionFunction float64    `json:"partition_function"`
		NumFactors        int        `json:"num_factors"`
	}
	tc.UnmarshalExpected(t, &expected)

	mn := models.NewMarkovNetwork()

	// Add nodes.
	nodeSet := make(map[string]bool)
	for _, edge := range input.Edges {
		nodeSet[edge[0]] = true
		nodeSet[edge[1]] = true
	}
	for node := range nodeSet {
		if err := mn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q) failed: %v", node, err)
		}
	}

	// Add edges.
	for _, edge := range input.Edges {
		if err := mn.AddEdge(edge[0], edge[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q) failed: %v", edge[0], edge[1], err)
		}
	}

	// Add factors.
	for i, fd := range input.Factors {
		f, err := factors.NewDiscreteFactor(fd.Variables, fd.Cardinality, fd.Values)
		if err != nil {
			t.Fatalf("NewDiscreteFactor (factor %d) failed: %v", i, err)
		}
		if err := mn.AddFactor(f); err != nil {
			t.Fatalf("AddFactor (factor %d) failed: %v", i, err)
		}
	}

	// Verify structure.
	if len(mn.Nodes()) != expected.NumNodes {
		t.Errorf("num_nodes: expected %d, got %d", expected.NumNodes, len(mn.Nodes()))
	}
	if len(mn.Edges()) != expected.NumEdges {
		t.Errorf("num_edges: expected %d, got %d", expected.NumEdges, len(mn.Edges()))
	}
	if len(mn.GetFactors()) != expected.NumFactors {
		t.Errorf("num_factors: expected %d, got %d", expected.NumFactors, len(mn.GetFactors()))
	}

	// Verify validity.
	err := mn.CheckModel()
	if expected.IsValid {
		if err != nil {
			t.Errorf("CheckModel: expected valid, got error: %v", err)
		}
	} else {
		if err == nil {
			t.Errorf("CheckModel: expected invalid, got nil error")
		}
	}

	// Verify partition function.
	Z, err := mn.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction failed: %v", err)
	}
	if math.Abs(Z-expected.PartitionFunction) > 1e-6 {
		t.Errorf("partition_function: expected %f, got %f", expected.PartitionFunction, Z)
	}
}
