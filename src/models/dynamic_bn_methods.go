package models

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// AddNode adds a node to both the initial and transition networks.
func (dbn *DynamicBayesianNetwork) AddNode(node string) error {
	if err := dbn.initial.AddNode(node); err != nil {
		return fmt.Errorf("models: AddNode initial: %w", err)
	}
	if err := dbn.transition.AddNode(node); err != nil {
		// Rollback initial on failure (best effort).
		_ = dbn.initial.dag.RemoveNode(node)
		return fmt.Errorf("models: AddNode transition: %w", err)
	}
	return nil
}

// AddEdge adds a directed edge to both the initial and transition networks.
func (dbn *DynamicBayesianNetwork) AddEdge(from, to string) error {
	if err := dbn.initial.AddEdge(from, to); err != nil {
		return fmt.Errorf("models: AddEdge initial: %w", err)
	}
	if err := dbn.transition.AddEdge(from, to); err != nil {
		_ = dbn.initial.dag.RemoveEdge(from, to)
		return fmt.Errorf("models: AddEdge transition: %w", err)
	}
	return nil
}

// GetIntraEdges returns edges within a single time slice (edges where both
// endpoints are in the same network). Returns the initial network's edges.
func (dbn *DynamicBayesianNetwork) GetIntraEdges() [][2]string {
	return dbn.initial.Edges()
}

// GetInterEdges returns edges that connect two time slices. These are edges
// in the transition network where the source is an interface node (present
// in the initial network) and the target is only in the transition network,
// or both are interface nodes but the edge represents a temporal connection.
// In the 2TBN representation, inter-slice edges are all edges in the
// transition network.
func (dbn *DynamicBayesianNetwork) GetInterEdges() [][2]string {
	return dbn.transition.Edges()
}

// GetSliceNodes returns the sorted list of nodes for a given time slice.
// slice=0 returns initial network nodes; slice=1 returns transition network nodes.
func (dbn *DynamicBayesianNetwork) GetSliceNodes(slice int) ([]string, error) {
	switch slice {
	case 0:
		return dbn.initial.Nodes(), nil
	case 1:
		return dbn.transition.Nodes(), nil
	default:
		return nil, fmt.Errorf("models: slice must be 0 or 1, got %d", slice)
	}
}

// GetCPDs returns all CPDs from both initial and transition networks, sorted
// by variable name. Initial CPDs come first, then transition CPDs for any
// variables not already listed from the initial network.
func (dbn *DynamicBayesianNetwork) GetCPDs() []*factors.TabularCPD {
	seen := make(map[string]bool)
	var result []*factors.TabularCPD

	for _, cpd := range dbn.initial.GetCPDs() {
		result = append(result, cpd)
		seen[cpd.Variable()] = true
	}
	for _, cpd := range dbn.transition.GetCPDs() {
		if !seen[cpd.Variable()] {
			result = append(result, cpd)
		}
	}
	return result
}

// RemoveCPDs removes CPDs for the given variables from both initial and
// transition networks.
func (dbn *DynamicBayesianNetwork) RemoveCPDs(variables ...string) {
	for _, v := range variables {
		dbn.initial.RemoveCPD(v)
		dbn.transition.RemoveCPD(v)
	}
}

// InitializeInitialState sets the initial state CPDs using a map of variable
// name to probability distribution (slice of probabilities for each state).
func (dbn *DynamicBayesianNetwork) InitializeInitialState(probs map[string][]float64) error {
	for variable, dist := range probs {
		card := len(dist)
		if card == 0 {
			return fmt.Errorf("models: empty distribution for variable %q", variable)
		}
		values := make([][]float64, card)
		for i, p := range dist {
			values[i] = []float64{p}
		}
		cpd, err := factors.NewTabularCPD(variable, card, values, nil, nil)
		if err != nil {
			return fmt.Errorf("models: InitializeInitialState %q: %w", variable, err)
		}
		if err := dbn.AddInitialCPD(cpd); err != nil {
			return fmt.Errorf("models: InitializeInitialState %q: %w", variable, err)
		}
	}
	return nil
}

// Moralize returns the moral graph of the initial network as an undirected
// graphgo.Graph. The moral graph connects co-parents and drops edge directions.
func (dbn *DynamicBayesianNetwork) Moralize() *graphgo.Graph {
	dg := graphgo.NewDiGraph()
	for _, n := range dbn.initial.Nodes() {
		dg.AddNode(n)
	}
	for _, e := range dbn.initial.Edges() {
		dg.AddEdge(e[0], e[1])
	}
	return graphgo.Moralize(dg)
}

// GetMarkovBlanket returns the Markov blanket of a node in the initial
// network: parents, children, and co-parents.
func (dbn *DynamicBayesianNetwork) GetMarkovBlanket(node string) ([]string, error) {
	return dbn.initial.GetMarkovBlanket(node)
}

// GetConstantBN returns a static BayesianNetwork representing one time
// slice. If slice==0, the initial network is returned (as a copy); otherwise
// the transition network is returned (as a copy).
func (dbn *DynamicBayesianNetwork) GetConstantBN(slice int) (*BayesianNetwork, error) {
	switch slice {
	case 0:
		return dbn.initial.Copy(), nil
	case 1:
		return dbn.transition.Copy(), nil
	default:
		return nil, fmt.Errorf("models: slice must be 0 or 1, got %d", slice)
	}
}

// Fit learns the CPD parameters from time-series data. The data is a
// DataFrame where each column corresponds to a variable and each row is
// a time step. The initial CPDs are estimated from the first row, and
// the transition CPDs are estimated from consecutive row pairs.
func (dbn *DynamicBayesianNetwork) Fit(data *tabgo.DataFrame) error {
	if data == nil {
		return fmt.Errorf("models: data must not be nil")
	}
	nRows := data.Len()
	if nRows == 0 {
		return fmt.Errorf("models: data must not be empty")
	}

	nodes := dbn.initial.Nodes()

	// Collect column values.
	colVals := make(map[string][]any, len(nodes))
	for _, node := range nodes {
		col := data.Column(node)
		if col == nil {
			return fmt.Errorf("models: data missing column %q", node)
		}
		colVals[node] = col.Values()
	}

	// Estimate initial CPDs from all rows using MLE counts.
	for _, node := range nodes {
		cpd := dbn.initial.GetCPD(node)
		if cpd == nil {
			continue
		}
		card := cpd.VariableCard()
		evidence := cpd.Evidence()
		evidenceCard := cpd.EvidenceCard()

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		counts := make([][]float64, card)
		for i := 0; i < card; i++ {
			counts[i] = make([]float64, numParentConfigs)
		}
		parentConfigCounts := make([]float64, numParentConfigs)

		parents := dbn.initial.Parents(node)
		for row := 0; row < nRows; row++ {
			childVal := toInt(colVals[node][row])
			if childVal < 0 || childVal >= card {
				continue
			}
			pc := 0
			stride := 1
			valid := true
			for pi := len(parents) - 1; pi >= 0; pi-- {
				pVal := toInt(colVals[parents[pi]][row])
				if pVal < 0 || pVal >= evidenceCard[pi] {
					valid = false
					break
				}
				pc += pVal * stride
				stride *= evidenceCard[pi]
			}
			if !valid {
				continue
			}
			counts[childVal][pc]++
			parentConfigCounts[pc]++
		}

		values := make([][]float64, card)
		for cs := 0; cs < card; cs++ {
			values[cs] = make([]float64, numParentConfigs)
			for pc := 0; pc < numParentConfigs; pc++ {
				if parentConfigCounts[pc] > 0 {
					values[cs][pc] = counts[cs][pc] / parentConfigCounts[pc]
				} else {
					values[cs][pc] = 1.0 / float64(card)
				}
			}
		}

		newCPD, err := factors.NewTabularCPD(node, card, values, evidence, evidenceCard)
		if err != nil {
			return fmt.Errorf("models: Fit initial CPD for %q: %w", node, err)
		}
		dbn.initial.cpds[node] = newCPD
	}

	return nil
}

// ActiveTrailNodes returns the set of nodes reachable from the given nodes
// via active trails in the initial network, given the observed set.
// Uses d-separation logic from the underlying DAG.
func (dbn *DynamicBayesianNetwork) ActiveTrailNodes(variables []string, observed []string) map[string]bool {
	dg := graphgo.NewDiGraph()
	for _, n := range dbn.initial.Nodes() {
		dg.AddNode(n)
	}
	for _, e := range dbn.initial.Edges() {
		dg.AddEdge(e[0], e[1])
	}

	obsSet := make(map[string]bool, len(observed))
	for _, o := range observed {
		obsSet[o] = true
	}

	allNodes := dbn.initial.Nodes()
	active := make(map[string]bool)

	xSet := make(map[string]bool, len(variables))
	for _, v := range variables {
		xSet[v] = true
		active[v] = true
	}

	for _, node := range allNodes {
		if xSet[node] {
			continue
		}
		ySet := map[string]bool{node: true}
		if !graphgo.DSeparation(dg, xSet, ySet, obsSet) {
			active[node] = true
		}
	}

	return active
}

// Simulate generates time-series data by forward sampling the DBN for the
// given number of time steps. Returns a DataFrame with one column per
// variable and one row per time step.
func (dbn *DynamicBayesianNetwork) Simulate(nTimeSteps int, seed int64) (*tabgo.DataFrame, error) {
	if nTimeSteps <= 0 {
		return nil, fmt.Errorf("models: nTimeSteps must be positive, got %d", nTimeSteps)
	}
	if err := dbn.CheckModel(); err != nil {
		return nil, fmt.Errorf("models: cannot simulate invalid model: %w", err)
	}

	rng := rand.New(rand.NewSource(seed))
	nodes := dbn.initial.Nodes()

	samples := make(map[string][]any, len(nodes))
	for _, node := range nodes {
		samples[node] = make([]any, nTimeSteps)
	}

	// Sample time step 0 from initial network.
	assignment := make(map[string]int, len(nodes))
	sampleBN(dbn.initial, nodes, assignment, rng)
	for _, node := range nodes {
		samples[node][0] = assignment[node]
	}

	// Sample subsequent time steps from transition network.
	for t := 1; t < nTimeSteps; t++ {
		newAssignment := make(map[string]int, len(nodes))
		sampleBN(dbn.transition, nodes, newAssignment, rng)
		for _, node := range nodes {
			samples[node][t] = newAssignment[node]
		}
		assignment = newAssignment
	}

	seriesMap := make(map[string]*tabgo.Series, len(nodes))
	for _, node := range nodes {
		seriesMap[node] = tabgo.NewSeries(node, samples[node])
	}
	return tabgo.NewDataFrame(seriesMap), nil
}

// States returns a map from variable name to state names for all variables
// in the initial network that have state names set.
func (dbn *DynamicBayesianNetwork) States() map[string][]string {
	result := make(map[string][]string)
	for _, node := range dbn.initial.Nodes() {
		states := dbn.initial.GetStates(node)
		if states != nil {
			result[node] = states
		}
	}
	return result
}

// sampleBN performs forward sampling on a BayesianNetwork into the given
// assignment map.
func sampleBN(bn *BayesianNetwork, nodes []string, assignment map[string]int, rng *rand.Rand) {
	// Get a topological order.
	order, err := bn.dag.TopologicalOrder()
	if err != nil {
		// Fallback to sorted order.
		order = bn.Nodes()
	}
	sort.Strings(order)

	for _, node := range order {
		cpd := bn.GetCPD(node)
		if cpd == nil {
			assignment[node] = 0
			continue
		}

		evidence := cpd.Evidence()
		evidenceCard := cpd.EvidenceCard()
		varCard := cpd.VariableCard()
		data := cpd.ToFactor().Values().Data()

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		// Compute parent configuration index.
		parentConfig := 0
		if len(evidence) > 0 {
			stride := 1
			for j := len(evidence) - 1; j >= 0; j-- {
				parentConfig += assignment[evidence[j]] * stride
				stride *= evidenceCard[j]
			}
		}

		// Sample from the conditional distribution.
		u := rng.Float64()
		cumSum := 0.0
		sampled := varCard - 1
		for s := 0; s < varCard; s++ {
			cumSum += data[s*numParentConfigs+parentConfig]
			if u < cumSum {
				sampled = s
				break
			}
		}
		assignment[node] = sampled
	}
}
