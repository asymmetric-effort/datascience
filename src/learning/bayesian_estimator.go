package learning

import (
	"fmt"
	"sort"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// PriorType specifies the type of prior used in Bayesian parameter estimation.
type PriorType int

const (
	// BDeu is the Bayesian Dirichlet equivalent uniform prior.
	// Pseudo-count = equivalentSampleSize / (num_parent_configs * num_states).
	BDeu PriorType = iota
	// K2 uses a uniform pseudo-count of 1 for every cell.
	K2
	// UniformPrior uses a pseudo-count of 1/num_states for every cell.
	UniformPrior
)

// BayesianEstimator learns CPD parameters for a BayesianNetwork from data
// using Bayesian estimation with configurable priors.
type BayesianEstimator struct {
	bn                   *models.BayesianNetwork
	data                 *tabgo.DataFrame
	priorType            PriorType
	equivalentSampleSize float64
}

// NewBayesianEstimator creates a new BayesianEstimator.
//
// bn is the Bayesian network whose structure (nodes, edges, states) must
// already be defined. data is a DataFrame whose columns correspond to the
// network's node names and whose values are state indices (int). priorType
// selects the prior family and equivalentSampleSize is the ESS used by BDeu
// (ignored for K2 and UniformPrior).
func NewBayesianEstimator(bn *models.BayesianNetwork, data *tabgo.DataFrame, priorType PriorType, equivalentSampleSize float64) *BayesianEstimator {
	return &BayesianEstimator{
		bn:                   bn,
		data:                 data,
		priorType:            priorType,
		equivalentSampleSize: equivalentSampleSize,
	}
}

// Estimate fits CPDs for every node in the network using Bayesian estimation
// and adds them to the network.
func (be *BayesianEstimator) Estimate() error {
	nodes := be.bn.Nodes()
	for _, node := range nodes {
		cpd, err := be.estimateNode(node)
		if err != nil {
			return fmt.Errorf("learning: bayesian estimation failed for node %q: %w", node, err)
		}
		if err := be.bn.AddCPD(cpd); err != nil {
			return fmt.Errorf("learning: failed to add CPD for node %q: %w", node, err)
		}
	}
	return nil
}

// GetParameters returns the CPD for the given node from the network, or an
// error if no CPD has been estimated yet.
func (be *BayesianEstimator) GetParameters(node string) (*factors.TabularCPD, error) {
	cpd := be.bn.GetCPD(node)
	if cpd == nil {
		return nil, fmt.Errorf("learning: no CPD found for node %q", node)
	}
	return cpd, nil
}

// estimateNode computes the Bayesian-estimated CPD for a single node.
func (be *BayesianEstimator) estimateNode(node string) (*factors.TabularCPD, error) {
	states := be.bn.GetStates(node)
	if len(states) == 0 {
		return nil, fmt.Errorf("no states defined for node %q", node)
	}
	numStates := len(states)

	parents := be.bn.Parents(node)
	sort.Strings(parents)

	// Collect parent cardinalities and build state-index maps.
	parentCards := make([]int, len(parents))
	parentStateMaps := make([]map[any]int, len(parents))
	for i, p := range parents {
		ps := be.bn.GetStates(p)
		if len(ps) == 0 {
			return nil, fmt.Errorf("no states defined for parent %q of node %q", p, node)
		}
		parentCards[i] = len(ps)
		m := make(map[any]int, len(ps))
		for si, s := range ps {
			m[s] = si
		}
		parentStateMaps[i] = m
	}

	numParentConfigs := 1
	for _, c := range parentCards {
		numParentConfigs *= c
	}

	// Build state-index map for the node itself.
	nodeStateMap := make(map[any]int, numStates)
	for si, s := range states {
		nodeStateMap[s] = si
	}

	// Count occurrences: counts[childState][parentConfig].
	counts := make([][]float64, numStates)
	for i := range counts {
		counts[i] = make([]float64, numParentConfigs)
	}

	nodeVals := be.data.Column(node).Values()
	parentVals := make([][]any, len(parents))
	for i, p := range parents {
		parentVals[i] = be.data.Column(p).Values()
	}

	nRows := be.data.Len()
	for r := 0; r < nRows; r++ {
		childIdx, ok := nodeStateMap[nodeVals[r]]
		if !ok {
			continue // skip unknown states
		}

		parentConfig := 0
		skip := false
		for pi := 0; pi < len(parents); pi++ {
			idx, ok2 := parentStateMaps[pi][parentVals[pi][r]]
			if !ok2 {
				skip = true
				break
			}
			parentConfig = parentConfig*parentCards[pi] + idx
		}
		if skip {
			continue
		}

		counts[childIdx][parentConfig]++
	}

	// Add pseudo-counts based on prior type.
	pseudoCount := be.pseudoCount(numStates, numParentConfigs)
	for i := 0; i < numStates; i++ {
		for j := 0; j < numParentConfigs; j++ {
			counts[i][j] += pseudoCount
		}
	}

	// Normalize each column (parent config) so it sums to 1.
	for j := 0; j < numParentConfigs; j++ {
		colSum := 0.0
		for i := 0; i < numStates; i++ {
			colSum += counts[i][j]
		}
		if colSum == 0 {
			// Uniform fallback if column is all zeros (should not happen with priors).
			for i := 0; i < numStates; i++ {
				counts[i][j] = 1.0 / float64(numStates)
			}
		} else {
			for i := 0; i < numStates; i++ {
				counts[i][j] /= colSum
			}
		}
	}

	// Build the TabularCPD.
	evidence := make([]string, len(parents))
	copy(evidence, parents)
	evidenceCard := make([]int, len(parentCards))
	copy(evidenceCard, parentCards)

	cpd, err := factors.NewTabularCPD(node, numStates, counts, evidence, evidenceCard)
	if err != nil {
		return nil, fmt.Errorf("failed to create CPD for node %q: %w", node, err)
	}
	return cpd, nil
}

// pseudoCount returns the pseudo-count to add to each cell based on the prior type.
func (be *BayesianEstimator) pseudoCount(numStates, numParentConfigs int) float64 {
	switch be.priorType {
	case BDeu:
		return be.equivalentSampleSize / float64(numParentConfigs*numStates)
	case K2:
		return 1.0
	case UniformPrior:
		return 1.0 / float64(numStates)
	default:
		return 1.0
	}
}
