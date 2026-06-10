package models

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// AddNode is an alias for AddCluster. It adds a cluster with the given
// variables and no factors, returning the cluster index.
func (cg *ClusterGraph) AddNode(variables []string) int {
	return cg.AddCluster(variables, nil)
}

// AddFactors adds factors to the cluster at the given index.
func (cg *ClusterGraph) AddFactors(clusterIndex int, newFactors []*factors.DiscreteFactor) error {
	if clusterIndex < 0 || clusterIndex >= len(cg.clusters) {
		return fmt.Errorf("models: cluster index %d out of range [0, %d)", clusterIndex, len(cg.clusters))
	}
	fs := make([]*factors.DiscreteFactor, len(newFactors))
	copy(fs, newFactors)
	cg.clusters[clusterIndex].Factors = append(cg.clusters[clusterIndex].Factors, fs...)
	return nil
}

// GetFactors returns all factors across all clusters.
func (cg *ClusterGraph) GetFactors() []*factors.DiscreteFactor {
	var result []*factors.DiscreteFactor
	for _, c := range cg.clusters {
		for _, f := range c.Factors {
			result = append(result, f)
		}
	}
	return result
}

// RemoveFactors removes all factors from all clusters.
func (cg *ClusterGraph) RemoveFactors() {
	for i := range cg.clusters {
		cg.clusters[i].Factors = nil
	}
}

// CliqueBeliefs runs a simple loopy belief propagation on the cluster graph
// and returns the belief (product of factors) for each cluster as a map
// from cluster index to a DiscreteFactor. For each cluster, the belief is
// the product of all its factors, normalized.
func (cg *ClusterGraph) CliqueBeliefs() (map[int]*factors.DiscreteFactor, error) {
	beliefs := make(map[int]*factors.DiscreteFactor)

	for i, c := range cg.clusters {
		if len(c.Factors) == 0 {
			continue
		}

		product, err := factors.FactorProduct(c.Factors...)
		if err != nil {
			return nil, fmt.Errorf("models: CliqueBeliefs cluster %d: %w", i, err)
		}
		product.Normalize()
		beliefs[i] = product
	}

	return beliefs, nil
}

// GetCardinality collects the cardinality of each variable from the
// factors in the cluster graph. Returns a map from variable name to
// cardinality.
func (cg *ClusterGraph) GetCardinality() map[string]int {
	result := make(map[string]int)
	for _, c := range cg.clusters {
		for _, f := range c.Factors {
			vars := f.Variables()
			card := f.Cardinality()
			for j, v := range vars {
				result[v] = card[j]
			}
		}
	}
	return result
}

// GetPartitionFunction computes the partition function (the sum over all
// assignments of the product of all factors in the cluster graph).
func (cg *ClusterGraph) GetPartitionFunction() (float64, error) {
	var allFactors []*factors.DiscreteFactor
	for _, c := range cg.clusters {
		allFactors = append(allFactors, c.Factors...)
	}
	if len(allFactors) == 0 {
		return 0, fmt.Errorf("models: no factors in cluster graph")
	}

	product, err := factors.FactorProduct(allFactors...)
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

// Copy returns a deep copy of the ClusterGraph.
func (cg *ClusterGraph) Copy() *ClusterGraph {
	newCG := &ClusterGraph{
		clusters: make([]Cluster, len(cg.clusters)),
		edges:    make([]ClusterEdge, len(cg.edges)),
	}

	for i, c := range cg.clusters {
		vars := make([]string, len(c.Variables))
		copy(vars, c.Variables)
		fs := make([]*factors.DiscreteFactor, len(c.Factors))
		for j, f := range c.Factors {
			fs[j] = f.Copy()
		}
		newCG.clusters[i] = Cluster{Variables: vars, Factors: fs}
	}

	for i, e := range cg.edges {
		ss := make([]string, len(e.SepSet))
		copy(ss, e.SepSet)
		newCG.edges[i] = ClusterEdge{
			Cluster1: e.Cluster1,
			Cluster2: e.Cluster2,
			SepSet:   ss,
		}
	}

	return newCG
}
