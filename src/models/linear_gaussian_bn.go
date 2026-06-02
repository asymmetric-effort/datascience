package models

import (
	"fmt"
	"sort"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// LinearGaussianBayesianNetwork is a Bayesian network where every node has a
// LinearGaussianCPD instead of a TabularCPD. It embeds *BayesianNetwork for
// graph structure (nodes, edges) and stores LG CPDs in a separate map.
type LinearGaussianBayesianNetwork struct {
	*BayesianNetwork
	lgCPDs map[string]*factors.LinearGaussianCPD
}

// NewLinearGaussianBayesianNetwork creates a new empty LinearGaussianBayesianNetwork.
func NewLinearGaussianBayesianNetwork() *LinearGaussianBayesianNetwork {
	return &LinearGaussianBayesianNetwork{
		BayesianNetwork: NewBayesianNetwork(),
		lgCPDs:          make(map[string]*factors.LinearGaussianCPD),
	}
}

// AddLinearGaussianCPD stores a LinearGaussianCPD for its variable.
// It validates that the variable exists in the graph and that the CPD's
// evidence matches the node's parents in the DAG.
func (lgbn *LinearGaussianBayesianNetwork) AddLinearGaussianCPD(cpd *factors.LinearGaussianCPD) error {
	if cpd == nil {
		return fmt.Errorf("models: cpd must not be nil")
	}

	v := cpd.Variable()
	if !lgbn.dag.HasNode(v) {
		return fmt.Errorf("models: variable %q is not a node in the network", v)
	}

	// Verify evidence matches parents.
	parents := lgbn.dag.Parents(v) // sorted
	evidence := cpd.Evidence()
	sortedEvidence := make([]string, len(evidence))
	copy(sortedEvidence, evidence)
	sort.Strings(sortedEvidence)

	if len(parents) != len(sortedEvidence) {
		return fmt.Errorf("models: LG CPD for %q has evidence %v but node has parents %v",
			v, evidence, parents)
	}
	for i := range parents {
		if parents[i] != sortedEvidence[i] {
			return fmt.Errorf("models: LG CPD for %q has evidence %v but node has parents %v",
				v, evidence, parents)
		}
	}

	lgbn.lgCPDs[v] = cpd
	return nil
}

// GetLinearGaussianCPD returns the LinearGaussianCPD for the given variable,
// or nil if none is set.
func (lgbn *LinearGaussianBayesianNetwork) GetLinearGaussianCPD(variable string) *factors.LinearGaussianCPD {
	return lgbn.lgCPDs[variable]
}

// CheckModel validates the LinearGaussianBayesianNetwork:
//  1. Every node has a LinearGaussianCPD.
//  2. Each CPD's evidence matches the node's parents in the DAG.
//  3. Each CPD passes Validate().
func (lgbn *LinearGaussianBayesianNetwork) CheckModel() error {
	nodes := lgbn.dag.Nodes()

	for _, node := range nodes {
		cpd, ok := lgbn.lgCPDs[node]
		if !ok {
			return fmt.Errorf("models: node %q has no LinearGaussianCPD", node)
		}

		// Check evidence matches parents.
		parents := lgbn.dag.Parents(node) // sorted
		evidence := cpd.Evidence()
		sortedEvidence := make([]string, len(evidence))
		copy(sortedEvidence, evidence)
		sort.Strings(sortedEvidence)

		if len(parents) != len(sortedEvidence) {
			return fmt.Errorf("models: LG CPD for %q has evidence %v but node has parents %v",
				node, evidence, parents)
		}
		for i := range parents {
			if parents[i] != sortedEvidence[i] {
				return fmt.Errorf("models: LG CPD for %q has evidence %v but node has parents %v",
					node, evidence, parents)
			}
		}

		if err := cpd.Validate(); err != nil {
			return fmt.Errorf("models: LG CPD for %q failed validation: %w", node, err)
		}
	}
	return nil
}

// Copy returns a deep copy of the LinearGaussianBayesianNetwork.
func (lgbn *LinearGaussianBayesianNetwork) Copy() *LinearGaussianBayesianNetwork {
	newLGCPDs := make(map[string]*factors.LinearGaussianCPD, len(lgbn.lgCPDs))
	for k, v := range lgbn.lgCPDs {
		newLGCPDs[k] = v.Copy()
	}
	return &LinearGaussianBayesianNetwork{
		BayesianNetwork: lgbn.BayesianNetwork.Copy(),
		lgCPDs:          newLGCPDs,
	}
}
