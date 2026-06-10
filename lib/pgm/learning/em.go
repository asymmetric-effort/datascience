package learning

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/inference"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ExpectationMaximization implements the EM algorithm for parameter estimation
// in Bayesian networks with latent (hidden) variables.
type ExpectationMaximization struct {
	bn         *models.BayesianNetwork
	data       *tabgo.DataFrame
	latentVars []string
	maxIter    int
	tol        float64
	iterations int
	converged  bool
}

// cpdStats holds expected sufficient statistics for a single CPD during EM.
type cpdStats struct {
	variable     string
	variableCard int
	evidence     []string
	evidenceCard []int
	// counts is a flat array of size variableCard * numParentConfig,
	// indexed as counts[childState*numParentConfig + parentConfig].
	counts          []float64
	numParentConfig int
}

// NewEM creates a new ExpectationMaximization estimator.
//
// bn is the Bayesian network whose CPD parameters will be estimated.
// data is the observed dataset (columns correspond to observed variable names).
// latentVars lists variables that are never observed in the data.
// maxIter is the maximum number of EM iterations.
// tol is the convergence tolerance: if the maximum absolute change in any CPD
// parameter between iterations is less than tol, the algorithm stops.
func NewEM(bn *models.BayesianNetwork, data *tabgo.DataFrame, latentVars []string, maxIter int, tol float64) *ExpectationMaximization {
	lv := make([]string, len(latentVars))
	copy(lv, latentVars)
	return &ExpectationMaximization{
		bn:         bn,
		data:       data,
		latentVars: lv,
		maxIter:    maxIter,
		tol:        tol,
	}
}

// Iterations returns the number of EM iterations performed.
func (em *ExpectationMaximization) Iterations() int {
	return em.iterations
}

// Converged returns true if the algorithm converged before reaching maxIter.
func (em *ExpectationMaximization) Converged() bool {
	return em.converged
}

// GetParameters returns the estimated CPDs for all nodes in the network as a
// map from node name to TabularCPD. This matches pgmpy's
// ExpectationMaximization.get_parameters(). The CPDs must have been estimated
// via Estimate() before calling this method.
func (em *ExpectationMaximization) GetParameters() (map[string]*factors.TabularCPD, error) {
	if em.bn == nil {
		return nil, fmt.Errorf("learning: BayesianNetwork is nil")
	}

	nodes := em.bn.Nodes()
	result := make(map[string]*factors.TabularCPD, len(nodes))
	for _, node := range nodes {
		cpd := em.bn.GetCPD(node)
		if cpd == nil {
			return nil, fmt.Errorf("learning: no CPD found for node %q; call Estimate first", node)
		}
		result[node] = cpd
	}
	return result, nil
}

// Estimate runs the EM algorithm and updates the Bayesian network's CPDs.
func (em *ExpectationMaximization) Estimate() error {
	nodes := em.bn.Nodes()

	// Determine which variables are latent vs observed.
	latentSet := make(map[string]bool, len(em.latentVars))
	for _, v := range em.latentVars {
		latentSet[v] = true
	}

	// Build cardinality map from the BN's state definitions.
	cardMap := make(map[string]int, len(nodes))
	for _, node := range nodes {
		states := em.bn.GetStates(node)
		if states == nil {
			return fmt.Errorf("learning: node %q has no states defined", node)
		}
		cardMap[node] = len(states)
	}

	// Build a state-name-to-index map for all variables.
	stateIndex := make(map[string]map[string]int, len(nodes))
	for _, node := range nodes {
		states := em.bn.GetStates(node)
		m := make(map[string]int, len(states))
		for i, s := range states {
			m[s] = i
		}
		stateIndex[node] = m
	}

	// Pre-extract observed data as integer state indices, one slice per
	// observed variable.
	observedCols := make(map[string][]int)
	dfCols := em.data.Columns()
	dfColSet := make(map[string]bool, len(dfCols))
	for _, c := range dfCols {
		dfColSet[c] = true
	}
	for _, node := range nodes {
		if latentSet[node] {
			continue
		}
		if !dfColSet[node] {
			return fmt.Errorf("learning: observed variable %q not found in data", node)
		}
		vals := em.data.Column(node).Values()
		col := make([]int, len(vals))
		for i, v := range vals {
			s := fmt.Sprintf("%v", v)
			idx, ok := stateIndex[node][s]
			if !ok {
				return fmt.Errorf("learning: unknown state %q for variable %q", s, node)
			}
			col[i] = idx
		}
		observedCols[node] = col
	}

	nRows := em.data.Len()

	// Step 1: Initialize CPDs.
	if err := em.initializeCPDs(nodes, latentSet, cardMap, observedCols, nRows); err != nil {
		return err
	}

	// EM loop.
	for iter := 0; iter < em.maxIter; iter++ {
		em.iterations = iter + 1

		// Build expected sufficient statistics accumulators.
		statsMap := make(map[string]*cpdStats, len(nodes))
		for _, node := range nodes {
			parents := em.bn.Parents(node)
			evCard := make([]int, len(parents))
			numPC := 1
			for i, p := range parents {
				evCard[i] = cardMap[p]
				numPC *= cardMap[p]
			}
			statsMap[node] = &cpdStats{
				variable:        node,
				variableCard:    cardMap[node],
				evidence:        parents,
				evidenceCard:    evCard,
				counts:          make([]float64, cardMap[node]*numPC),
				numParentConfig: numPC,
			}
		}

		// E-step + accumulate statistics.
		for row := 0; row < nRows; row++ {
			evidence := make(map[string]int, len(observedCols))
			for v, col := range observedCols {
				evidence[v] = col[row]
			}

			latentPosterior, err := em.computeLatentPosterior(evidence)
			if err != nil {
				return fmt.Errorf("learning: E-step failed at row %d: %w", row, err)
			}

			em.accumulateStats(statsMap, latentPosterior, evidence, cardMap)
		}

		// M-step: build new CPDs from expected counts.
		maxDelta := 0.0
		for _, node := range nodes {
			st := statsMap[node]
			values := make([][]float64, st.variableCard)
			for cs := 0; cs < st.variableCard; cs++ {
				values[cs] = make([]float64, st.numParentConfig)
			}

			for pc := 0; pc < st.numParentConfig; pc++ {
				colSum := 0.0
				for cs := 0; cs < st.variableCard; cs++ {
					colSum += st.counts[cs*st.numParentConfig+pc]
				}
				if colSum == 0 {
					for cs := 0; cs < st.variableCard; cs++ {
						values[cs][pc] = 1.0 / float64(st.variableCard)
					}
				} else {
					for cs := 0; cs < st.variableCard; cs++ {
						values[cs][pc] = st.counts[cs*st.numParentConfig+pc] / colSum
					}
				}
			}

			newCPD, err := factors.NewTabularCPD(node, st.variableCard, values, st.evidence, st.evidenceCard)
			if err != nil {
				return fmt.Errorf("learning: M-step failed for %q: %w", node, err)
			}

			oldCPD := em.bn.GetCPD(node)
			if oldCPD != nil {
				oldData := oldCPD.ToFactor().Values().Data()
				newData := newCPD.ToFactor().Values().Data()
				for i := range oldData {
					d := math.Abs(oldData[i] - newData[i])
					if d > maxDelta {
						maxDelta = d
					}
				}
			}

			if err := em.bn.AddCPD(newCPD); err != nil {
				return fmt.Errorf("learning: failed to set CPD for %q: %w", node, err)
			}
		}

		// Check convergence.
		if maxDelta < em.tol {
			em.converged = true
			return nil
		}
	}

	return nil
}

// initializeCPDs sets initial CPDs on the BN. Observed-only nodes get MLE
// estimates from data; nodes that are latent or have latent parents get
// uniform distributions.
func (em *ExpectationMaximization) initializeCPDs(
	nodes []string,
	latentSet map[string]bool,
	cardMap map[string]int,
	observedCols map[string][]int,
	nRows int,
) error {
	for _, node := range nodes {
		parents := em.bn.Parents(node)
		evCard := make([]int, len(parents))
		numPC := 1
		for i, p := range parents {
			evCard[i] = cardMap[p]
			numPC *= cardMap[p]
		}

		// Check if this node or any parent is latent.
		hasLatent := latentSet[node]
		if !hasLatent {
			for _, p := range parents {
				if latentSet[p] {
					hasLatent = true
					break
				}
			}
		}

		values := make([][]float64, cardMap[node])
		for cs := 0; cs < cardMap[node]; cs++ {
			values[cs] = make([]float64, numPC)
		}

		if hasLatent {
			// Slightly asymmetric initialization to break symmetry.
			// Pure uniform initialization can trap EM at a symmetric
			// fixed point where the latent variable is uninformative.
			for cs := 0; cs < cardMap[node]; cs++ {
				for pc := 0; pc < numPC; pc++ {
					// Add a small deterministic perturbation based on
					// the child state and parent config indices.
					base := 1.0 / float64(cardMap[node])
					perturbation := 0.05 * float64(cs-cardMap[node]/2) / float64(1+pc+cs)
					values[cs][pc] = base + perturbation
				}
			}
			// Renormalize each column to ensure it sums to 1.
			for pc := 0; pc < numPC; pc++ {
				colSum := 0.0
				for cs := 0; cs < cardMap[node]; cs++ {
					if values[cs][pc] < 1e-10 {
						values[cs][pc] = 1e-10
					}
					colSum += values[cs][pc]
				}
				for cs := 0; cs < cardMap[node]; cs++ {
					values[cs][pc] /= colSum
				}
			}
		} else {
			// MLE from observed data.
			counts := make([]float64, cardMap[node]*numPC)
			for row := 0; row < nRows; row++ {
				childState := observedCols[node][row]
				pc := parentConfigFromMap(parents, observedCols, row, cardMap)
				counts[childState*numPC+pc]++
			}
			for pc := 0; pc < numPC; pc++ {
				colSum := 0.0
				for cs := 0; cs < cardMap[node]; cs++ {
					colSum += counts[cs*numPC+pc]
				}
				if colSum == 0 {
					for cs := 0; cs < cardMap[node]; cs++ {
						values[cs][pc] = 1.0 / float64(cardMap[node])
					}
				} else {
					for cs := 0; cs < cardMap[node]; cs++ {
						values[cs][pc] = counts[cs*numPC+pc] / colSum
					}
				}
			}
		}

		cpd, err := factors.NewTabularCPD(node, cardMap[node], values, parents, evCard)
		if err != nil {
			return fmt.Errorf("learning: initialization failed for %q: %w", node, err)
		}
		if err := em.bn.AddCPD(cpd); err != nil {
			return err
		}
	}
	return nil
}

// parentConfigFromMap computes the flat parent configuration index for a data
// row given observed column data. Uses row-major ordering matching TabularCPD.
func parentConfigFromMap(
	parents []string,
	observedCols map[string][]int,
	row int,
	cardMap map[string]int,
) int {
	pc := 0
	stride := 1
	for i := len(parents) - 1; i >= 0; i-- {
		pc += observedCols[parents[i]][row] * stride
		stride *= cardMap[parents[i]]
	}
	return pc
}

// parentConfigFromAssignment computes the flat parent configuration index
// from a variable assignment map.
func parentConfigFromAssignment(
	parents []string,
	assignment map[string]int,
	cardMap map[string]int,
) int {
	pc := 0
	stride := 1
	for i := len(parents) - 1; i >= 0; i-- {
		pc += assignment[parents[i]] * stride
		stride *= cardMap[parents[i]]
	}
	return pc
}

// computeLatentPosterior computes P(latent vars | evidence) using variable
// elimination. Returns a DiscreteFactor over the latent variables.
func (em *ExpectationMaximization) computeLatentPosterior(
	evidence map[string]int,
) (*factors.DiscreteFactor, error) {
	if len(em.latentVars) == 0 {
		return factors.NewDiscreteFactor(nil, nil, []float64{1.0})
	}

	markovFactors, err := em.bn.ToMarkovFactors()
	if err != nil {
		return nil, err
	}

	ve := inference.NewVariableElimination(markovFactors)
	result, err := ve.Query(em.latentVars, evidence)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// accumulateStats enumerates all latent variable assignments weighted by
// the posterior probability, and accumulates expected counts into statsMap.
func (em *ExpectationMaximization) accumulateStats(
	statsMap map[string]*cpdStats,
	latentPosterior *factors.DiscreteFactor,
	evidence map[string]int,
	cardMap map[string]int,
) {
	// Enumerate all latent variable assignments.
	latentCards := make([]int, len(em.latentVars))
	totalLatent := 1
	for i, v := range em.latentVars {
		latentCards[i] = cardMap[v]
		totalLatent *= cardMap[v]
	}

	for lFlat := 0; lFlat < totalLatent; lFlat++ {
		// Decode latent assignment.
		latentAssignment := make(map[string]int, len(em.latentVars))
		rem := lFlat
		for i := len(em.latentVars) - 1; i >= 0; i-- {
			latentAssignment[em.latentVars[i]] = rem % latentCards[i]
			rem /= latentCards[i]
		}

		// Get weight = P(latent assignment | evidence).
		weight := latentPosterior.GetValue(latentAssignment)

		// Build full assignment: observed + latent.
		fullAssignment := make(map[string]int, len(evidence)+len(latentAssignment))
		for k, v := range evidence {
			fullAssignment[k] = v
		}
		for k, v := range latentAssignment {
			fullAssignment[k] = v
		}

		// Accumulate weighted counts for each node.
		for _, st := range statsMap {
			childState := fullAssignment[st.variable]
			pc := parentConfigFromAssignment(st.evidence, fullAssignment, cardMap)
			st.counts[childState*st.numParentConfig+pc] += weight
		}
	}
}
