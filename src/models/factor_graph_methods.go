package models

import (
	"fmt"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// AddEdge adds an edge between a variable node and a factor node (identified
// by its index in the factor list). This verifies that the variable is in
// the factor's scope and that the factor exists.
func (fg *FactorGraph) AddEdge(variable string, factorIndex int) error {
	if _, exists := fg.variables[variable]; !exists {
		return fmt.Errorf("models: variable %q does not exist", variable)
	}
	if factorIndex < 0 || factorIndex >= len(fg.factorList) {
		return fmt.Errorf("models: factor index %d out of range [0, %d)", factorIndex, len(fg.factorList))
	}

	f := fg.factorList[factorIndex]
	found := false
	for _, v := range f.Variables() {
		if v == variable {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("models: variable %q is not in factor %d scope", variable, factorIndex)
	}

	// The edge is already implicit via varToFactors, but we verify and
	// ensure the mapping exists.
	for _, existing := range fg.varToFactors[variable] {
		if existing == f {
			return nil // Already connected.
		}
	}
	fg.varToFactors[variable] = append(fg.varToFactors[variable], f)
	return nil
}

// RemoveFactors removes all factors from the graph and clears the
// variable-to-factor mappings.
func (fg *FactorGraph) RemoveFactors() {
	fg.factorList = nil
	fg.varToFactors = make(map[string][]*factors.DiscreteFactor)
}

// GetCardinality returns the cardinality of the given variable.
// Returns an error if the variable does not exist.
func (fg *FactorGraph) GetCardinality(variable string) (int, error) {
	card, exists := fg.variables[variable]
	if !exists {
		return 0, fmt.Errorf("models: variable %q does not exist", variable)
	}
	return card, nil
}

// GetFactorNodes returns a list of factor descriptions. Each entry is the
// list of variables in the factor's scope.
func (fg *FactorGraph) GetFactorNodes() [][]string {
	result := make([][]string, len(fg.factorList))
	for i, f := range fg.factorList {
		result[i] = f.Variables()
	}
	return result
}

// ToJunctionTree converts the factor graph to a junction tree. This creates
// a moral graph from the factor scopes (each factor's variables form a
// clique), triangulates it, finds maximal cliques, and builds the clique tree.
func (fg *FactorGraph) ToJunctionTree() (*JunctionTree, error) {
	if err := fg.CheckModel(); err != nil {
		return nil, fmt.Errorf("models: cannot build junction tree: %w", err)
	}

	// Build an undirected graph from the factor scopes.
	moral := graphgo.NewGraph()
	for v := range fg.variables {
		moral.AddNode(v)
	}
	for _, f := range fg.factorList {
		vars := f.Variables()
		for i := 0; i < len(vars); i++ {
			for j := i + 1; j < len(vars); j++ {
				moral.AddEdge(vars[i], vars[j])
			}
		}
	}

	// Find elimination order using min-degree heuristic.
	order := minDegreeOrder(moral)

	// Triangulate.
	triangulated := graphgo.Triangulate(moral, order)

	// Find maximal cliques.
	cliques := graphgo.MaxCliques(triangulated)
	if len(cliques) == 0 {
		return &JunctionTree{
			cliques:       nil,
			tree:          graphgo.NewGraph(),
			separators:    make(map[string][]string),
			cliqueFactors: make(map[int][]*factors.DiscreteFactor),
		}, nil
	}

	// Build junction tree.
	tree, separators := graphgo.BuildJunctionTree(cliques)

	// Assign factors to cliques.
	cliqueFactors := assignFactorsToCliques(fg.factorList, cliques)

	return &JunctionTree{
		cliques:       cliques,
		tree:          tree,
		separators:    separators,
		cliqueFactors: cliqueFactors,
	}, nil
}

// GetPartitionFunction computes the partition function (sum of the product
// of all factors over all possible assignments).
func (fg *FactorGraph) GetPartitionFunction() (float64, error) {
	if len(fg.factorList) == 0 {
		return 0, fmt.Errorf("models: no factors in graph")
	}

	product, err := factors.FactorProduct(fg.factorList...)
	if err != nil {
		return 0, fmt.Errorf("models: GetPartitionFunction: %w", err)
	}

	data := product.Values().Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum, nil
}

// Copy returns a deep copy of the FactorGraph.
func (fg *FactorGraph) Copy() *FactorGraph {
	newFG := NewFactorGraph()

	for v, card := range fg.variables {
		newFG.variables[v] = card
	}

	newFG.factorList = make([]*factors.DiscreteFactor, len(fg.factorList))
	for i, f := range fg.factorList {
		newFG.factorList[i] = f.Copy()
	}

	// Rebuild varToFactors from the copied factors.
	for _, f := range newFG.factorList {
		for _, v := range f.Variables() {
			newFG.varToFactors[v] = append(newFG.varToFactors[v], f)
		}
	}

	return newFG
}

// GetPointMassMessage creates a delta factor for a variable where only
// the specified state has value 1.0 and all others are 0.0.
func (fg *FactorGraph) GetPointMassMessage(variable string, state int) (*factors.DiscreteFactor, error) {
	card, exists := fg.variables[variable]
	if !exists {
		return nil, fmt.Errorf("models: variable %q does not exist", variable)
	}
	if state < 0 || state >= card {
		return nil, fmt.Errorf("models: state %d out of range [0, %d) for variable %q", state, card, variable)
	}

	values := make([]float64, card)
	values[state] = 1.0

	return factors.NewDiscreteFactor([]string{variable}, []int{card}, values)
}

// GetUniformMessage creates a uniform factor for a variable where all
// states have equal probability (1/cardinality).
func (fg *FactorGraph) GetUniformMessage(variable string) (*factors.DiscreteFactor, error) {
	card, exists := fg.variables[variable]
	if !exists {
		return nil, fmt.Errorf("models: variable %q does not exist", variable)
	}

	values := make([]float64, card)
	p := 1.0 / float64(card)
	for i := range values {
		values[i] = p
	}

	return factors.NewDiscreteFactor([]string{variable}, []int{card}, values)
}
