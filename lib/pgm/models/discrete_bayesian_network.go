package models

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// DiscreteBayesianNetwork is a Bayesian network where all variables are
// discrete. It embeds *BayesianNetwork and adds discrete-specific validation,
// fitting (via injected estimator functions), and forward sampling.
type DiscreteBayesianNetwork struct {
	*BayesianNetwork
}

// NewDiscreteBayesianNetwork creates a new empty DiscreteBayesianNetwork.
func NewDiscreteBayesianNetwork() *DiscreteBayesianNetwork {
	return &DiscreteBayesianNetwork{
		BayesianNetwork: NewBayesianNetwork(),
	}
}

// AddCPD stores a CPD for its variable. In addition to the base BayesianNetwork
// checks, it validates that all cardinalities are positive and finite.
func (dbn *DiscreteBayesianNetwork) AddCPD(cpd *factors.TabularCPD) error {
	if cpd == nil {
		return fmt.Errorf("models: cpd must not be nil")
	}

	// Validate variable cardinality is positive and finite.
	vc := cpd.VariableCard()
	if vc <= 0 {
		return fmt.Errorf("models: variable cardinality must be positive, got %d", vc)
	}

	// Validate all evidence cardinalities are positive and finite.
	for i, ec := range cpd.EvidenceCard() {
		if ec <= 0 {
			return fmt.Errorf("models: evidence cardinality at index %d must be positive, got %d", i, ec)
		}
	}

	// Validate that the CPD values contain no NaN or Inf.
	f := cpd.ToFactor()
	for _, v := range f.Values().Data() {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("models: CPD for %q contains NaN or Inf values", cpd.Variable())
		}
	}

	return dbn.BayesianNetwork.AddCPD(cpd)
}

// CheckModel validates the discrete Bayesian network. It calls the base
// BayesianNetwork.CheckModel() and then performs additional discrete-specific
// checks:
//   - All state cardinalities must be positive integers (already enforced by
//     TabularCPD construction, but re-verified here).
//   - If state names are set for a variable, the number of state names must
//     match the variable's cardinality in its CPD.
//   - State names referenced across CPDs must be consistent: if a parent
//     variable has state names set, the evidence cardinality in child CPDs
//     must match the number of states.
func (dbn *DiscreteBayesianNetwork) CheckModel() error {
	if err := dbn.BayesianNetwork.CheckModel(); err != nil {
		return err
	}

	nodes := dbn.Nodes()
	for _, node := range nodes {
		cpd := dbn.GetCPD(node)
		if cpd == nil {
			continue // base CheckModel already catches this
		}

		// Verify variable cardinality is positive.
		if cpd.VariableCard() <= 0 {
			return fmt.Errorf("models: variable %q has non-positive cardinality %d", node, cpd.VariableCard())
		}

		// Verify evidence cardinalities are positive.
		for i, ec := range cpd.EvidenceCard() {
			if ec <= 0 {
				return fmt.Errorf("models: evidence cardinality at index %d for variable %q is non-positive: %d", i, node, ec)
			}
		}

		// If state names are set for this variable, check consistency with CPD cardinality.
		stateNames := dbn.GetStates(node)
		if stateNames != nil && len(stateNames) != cpd.VariableCard() {
			return fmt.Errorf("models: variable %q has %d state names but cardinality %d",
				node, len(stateNames), cpd.VariableCard())
		}

		// Check that evidence cardinalities match parent state names if set.
		evidence := cpd.Evidence()
		evidenceCard := cpd.EvidenceCard()
		for i, parent := range evidence {
			parentStates := dbn.GetStates(parent)
			if parentStates != nil && len(parentStates) != evidenceCard[i] {
				return fmt.Errorf("models: parent %q of %q has %d state names but evidence cardinality is %d",
					parent, node, len(parentStates), evidenceCard[i])
			}
		}

		// Validate CPD values contain no NaN or Inf.
		f := cpd.ToFactor()
		for _, v := range f.Values().Data() {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return fmt.Errorf("models: CPD for %q contains NaN or Inf values", node)
			}
		}
	}

	return nil
}

// FitWith fits the model parameters from data using a caller-supplied
// estimation function. This avoids circular imports between the models and
// learning packages.
//
// The estimateFn receives the underlying BayesianNetwork (so it can set CPDs)
// and the data, and should return an error on failure.
func (dbn *DiscreteBayesianNetwork) FitWith(estimateFn func(*BayesianNetwork, *tabgo.DataFrame) error, data *tabgo.DataFrame) error {
	if estimateFn == nil {
		return fmt.Errorf("models: estimateFn must not be nil")
	}
	if data == nil {
		return fmt.Errorf("models: data must not be nil")
	}
	return estimateFn(dbn.BayesianNetwork, data)
}

// Simulate generates n data points by forward sampling the network.
// It returns a DataFrame with one column per node, in topological order.
// The seed parameter controls the random number generator (use 0 for
// non-deterministic behavior).
func (dbn *DiscreteBayesianNetwork) Simulate(n int, seed int64) (*tabgo.DataFrame, error) {
	if n <= 0 {
		return nil, fmt.Errorf("models: n must be positive, got %d", n)
	}

	if err := dbn.CheckModel(); err != nil {
		return nil, fmt.Errorf("models: cannot simulate invalid model: %w", err)
	}

	// Get topological order for forward sampling.
	// We compute a deterministic topological order by doing a BFS with
	// sorted tie-breaking so that results are reproducible with the same seed.
	topoOrder, err := dbn.deterministicTopologicalOrder()
	if err != nil {
		return nil, fmt.Errorf("models: failed to get topological order: %w", err)
	}

	rng := rand.New(rand.NewSource(seed))

	// Pre-compute CPD data for each node to avoid repeated copies.
	type nodeInfo struct {
		evidence         []string
		evidenceCard     []int
		varCard          int
		data             []float64
		numParentConfigs int
	}
	nodeInfos := make(map[string]*nodeInfo, len(topoOrder))
	for _, node := range topoOrder {
		cpd := dbn.GetCPD(node)
		ec := cpd.EvidenceCard()
		npc := 1
		for _, c := range ec {
			npc *= c
		}
		f := cpd.ToFactor()
		nodeInfos[node] = &nodeInfo{
			evidence:         cpd.Evidence(),
			evidenceCard:     ec,
			varCard:          cpd.VariableCard(),
			data:             f.Values().Data(),
			numParentConfigs: npc,
		}
	}

	// Pre-allocate sample storage: variable -> slice of sampled state indices.
	samples := make(map[string][]any, len(topoOrder))
	for _, node := range topoOrder {
		samples[node] = make([]any, n)
	}

	for i := 0; i < n; i++ {
		// Sample each node in topological order.
		assignment := make(map[string]int, len(topoOrder))

		for _, node := range topoOrder {
			ni := nodeInfos[node]

			// Compute parent configuration index (column in CPD table).
			parentConfig := 0
			if len(ni.evidence) > 0 {
				stride := 1
				for j := len(ni.evidence) - 1; j >= 0; j-- {
					parentConfig += assignment[ni.evidence[j]] * stride
					stride *= ni.evidenceCard[j]
				}
			}

			// Build cumulative distribution for sampling.
			cumSum := 0.0
			u := rng.Float64()
			sampled := ni.varCard - 1 // default to last state
			for s := 0; s < ni.varCard; s++ {
				cumSum += ni.data[s*ni.numParentConfigs+parentConfig]
				if u < cumSum {
					sampled = s
					break
				}
			}

			assignment[node] = sampled
			samples[node][i] = sampled
		}
	}

	// Build DataFrame columns sorted alphabetically for determinism.
	colMap := make(map[string]*tabgo.Series, len(topoOrder))
	for _, node := range topoOrder {
		colMap[node] = tabgo.NewSeries(node, samples[node])
	}

	// Use sorted column names.
	sortedNames := make([]string, len(topoOrder))
	copy(sortedNames, topoOrder)
	sort.Strings(sortedNames)

	cols := make(map[string]*tabgo.Series, len(sortedNames))
	for _, name := range sortedNames {
		cols[name] = colMap[name]
	}

	return tabgo.NewDataFrame(cols), nil
}

// deterministicTopologicalOrder returns a topological ordering where ties
// (nodes with equal in-degree at each step) are broken alphabetically.
// This ensures reproducible sampling given the same seed.
func (dbn *DiscreteBayesianNetwork) deterministicTopologicalOrder() ([]string, error) {
	nodes := dbn.Nodes() // sorted
	inDegree := make(map[string]int, len(nodes))
	for _, n := range nodes {
		inDegree[n] = len(dbn.Parents(n))
	}

	// Seed with zero in-degree nodes (already sorted since nodes is sorted).
	var queue []string
	for _, n := range nodes {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}

	var order []string
	for len(queue) > 0 {
		// Sort queue to ensure deterministic tie-breaking.
		sort.Strings(queue)
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, child := range dbn.Children(node) { // Children returns sorted
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}

	if len(order) != len(nodes) {
		return nil, fmt.Errorf("models: graph contains a cycle")
	}
	return order, nil
}

// Copy returns a deep copy of the DiscreteBayesianNetwork.
func (dbn *DiscreteBayesianNetwork) Copy() *DiscreteBayesianNetwork {
	return &DiscreteBayesianNetwork{
		BayesianNetwork: dbn.BayesianNetwork.Copy(),
	}
}
