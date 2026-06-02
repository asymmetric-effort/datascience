package models

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/src/base"
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// MarkovNetwork represents a Markov random field (MRF) — an undirected
// graphical model where each factor (potential) is defined over a clique
// of the graph.
type MarkovNetwork struct {
	graph        *base.UndirectedGraph
	factorList   []*factors.DiscreteFactor
	varToFactors map[string][]*factors.DiscreteFactor
}

// NewMarkovNetwork creates a new empty MarkovNetwork.
func NewMarkovNetwork() *MarkovNetwork {
	return &MarkovNetwork{
		graph:        base.NewUndirectedGraph(),
		factorList:   nil,
		varToFactors: make(map[string][]*factors.DiscreteFactor),
	}
}

// AddNode adds a node (variable) to the network.
func (mn *MarkovNetwork) AddNode(node string) error {
	return mn.graph.AddNode(node)
}

// AddEdge adds an undirected edge between u and v. Both nodes must exist.
func (mn *MarkovNetwork) AddEdge(u, v string) error {
	return mn.graph.AddEdge(u, v)
}

// Nodes returns a sorted list of all nodes in the network.
func (mn *MarkovNetwork) Nodes() []string {
	return mn.graph.Nodes()
}

// Edges returns all undirected edges in canonical form (A < B), sorted
// lexicographically.
func (mn *MarkovNetwork) Edges() [][2]string {
	raw := mn.graph.Edges()
	result := make([][2]string, len(raw))
	for i, e := range raw {
		result[i] = [2]string{e.A, e.B}
	}
	return result
}

// Neighbors returns the sorted neighbors of a node in the undirected graph.
func (mn *MarkovNetwork) Neighbors(node string) []string {
	return mn.graph.Neighbors(node)
}

// AddFactor adds a factor to the network. All variables in the factor's scope
// must be nodes in the graph. Returns an error if any variable is missing or
// if the factor is nil.
func (mn *MarkovNetwork) AddFactor(f *factors.DiscreteFactor) error {
	if f == nil {
		return fmt.Errorf("models: factor must not be nil")
	}
	vars := f.Variables()
	for _, v := range vars {
		if !mn.graph.HasNode(v) {
			return fmt.Errorf("models: factor references unknown node %q", v)
		}
	}
	mn.factorList = append(mn.factorList, f)
	for _, v := range vars {
		mn.varToFactors[v] = append(mn.varToFactors[v], f)
	}
	return nil
}

// RemoveFactor removes all factors whose variable set exactly matches the
// given variables (order-independent).
func (mn *MarkovNetwork) RemoveFactor(variables []string) {
	target := make([]string, len(variables))
	copy(target, variables)
	sort.Strings(target)
	targetKey := strings.Join(target, "\x00")

	var kept []*factors.DiscreteFactor
	for _, f := range mn.factorList {
		fvars := f.Variables()
		sorted := make([]string, len(fvars))
		copy(sorted, fvars)
		sort.Strings(sorted)
		if strings.Join(sorted, "\x00") == targetKey {
			continue // remove this factor
		}
		kept = append(kept, f)
	}
	mn.factorList = kept

	// Rebuild varToFactors index.
	mn.varToFactors = make(map[string][]*factors.DiscreteFactor)
	for _, f := range mn.factorList {
		for _, v := range f.Variables() {
			mn.varToFactors[v] = append(mn.varToFactors[v], f)
		}
	}
}

// GetFactors returns all factors, in insertion order.
func (mn *MarkovNetwork) GetFactors() []*factors.DiscreteFactor {
	result := make([]*factors.DiscreteFactor, len(mn.factorList))
	copy(result, mn.factorList)
	return result
}

// GetFactorsOf returns all factors that include the given variable in their
// scope. Returns nil if the variable has no factors.
func (mn *MarkovNetwork) GetFactorsOf(variable string) []*factors.DiscreteFactor {
	fs := mn.varToFactors[variable]
	if len(fs) == 0 {
		return nil
	}
	result := make([]*factors.DiscreteFactor, len(fs))
	copy(result, fs)
	return result
}

// CheckModel validates the Markov network:
//  1. Every factor's scope variables must be connected by edges: for every
//     pair of variables in a factor, the edge must exist in the graph.
//  2. Every node must be covered by at least one factor.
func (mn *MarkovNetwork) CheckModel() error {
	if len(mn.factorList) == 0 {
		return fmt.Errorf("models: Markov network has no factors")
	}

	// Check that every pair of variables within each factor has an edge
	// (or the factor is unary).
	for i, f := range mn.factorList {
		vars := f.Variables()
		for j := 0; j < len(vars); j++ {
			if !mn.graph.HasNode(vars[j]) {
				return fmt.Errorf("models: factor %d references unknown node %q", i, vars[j])
			}
			for k := j + 1; k < len(vars); k++ {
				if !mn.graph.HasEdge(vars[j], vars[k]) {
					return fmt.Errorf("models: factor %d has variables %q and %q but no edge exists between them",
						i, vars[j], vars[k])
				}
			}
		}
	}

	// Check that every node is covered by at least one factor.
	nodes := mn.graph.Nodes()
	for _, node := range nodes {
		if len(mn.varToFactors[node]) == 0 {
			return fmt.Errorf("models: node %q is not covered by any factor", node)
		}
	}

	return nil
}

// GetPartitionFunction computes Z, the partition function, by summing the
// product of all factors over all joint assignments. This is only feasible
// for small models.
func (mn *MarkovNetwork) GetPartitionFunction() (float64, error) {
	if len(mn.factorList) == 0 {
		return 0, fmt.Errorf("models: no factors in the network")
	}

	// Compute the full joint factor by multiplying all factors together.
	product, err := factors.FactorProduct(mn.factorList...)
	if err != nil {
		return 0, fmt.Errorf("models: partition function: %w", err)
	}

	// Sum all values in the product factor.
	data := product.Values().Data()
	z := 0.0
	for _, v := range data {
		z += v
	}
	return z, nil
}

// MarkovBlanket returns the Markov blanket of a node, which in an undirected
// model is simply the set of neighbors.
func (mn *MarkovNetwork) MarkovBlanket(node string) []string {
	return mn.graph.Neighbors(node)
}

// ToJunctionTree constructs a junction tree from the Markov network by:
//  1. Building a graphgo.Graph from the undirected graph
//  2. Finding a min-degree elimination order
//  3. Triangulating the graph
//  4. Finding maximal cliques
//  5. Building the junction tree
//  6. Assigning factors to cliques
func (mn *MarkovNetwork) ToJunctionTree() (*JunctionTree, error) {
	if err := mn.CheckModel(); err != nil {
		return nil, fmt.Errorf("models: cannot build junction tree: %w", err)
	}

	// Build a graphgo.Graph from the UndirectedGraph's public API.
	g := graphgo.NewGraph()
	for _, node := range mn.graph.Nodes() {
		g.AddNode(node)
	}
	for _, e := range mn.graph.Edges() {
		g.AddEdge(e.A, e.B)
	}

	// Find elimination order using min-degree heuristic.
	order := minDegreeOrder(g)

	// Triangulate.
	triangulated := graphgo.Triangulate(g, order)

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
	cliqueFactors := assignFactorsToCliques(mn.factorList, cliques)

	return &JunctionTree{
		cliques:       cliques,
		tree:          tree,
		separators:    separators,
		cliqueFactors: cliqueFactors,
	}, nil
}

// Copy returns a deep copy of the MarkovNetwork.
func (mn *MarkovNetwork) Copy() *MarkovNetwork {
	newFactors := make([]*factors.DiscreteFactor, len(mn.factorList))
	for i, f := range mn.factorList {
		newFactors[i] = f.Copy()
	}

	newMN := &MarkovNetwork{
		graph:        mn.graph.Copy(),
		factorList:   newFactors,
		varToFactors: make(map[string][]*factors.DiscreteFactor),
	}

	// Rebuild varToFactors from the copied factors.
	for _, f := range newMN.factorList {
		for _, v := range f.Variables() {
			newMN.varToFactors[v] = append(newMN.varToFactors[v], f)
		}
	}

	return newMN
}
