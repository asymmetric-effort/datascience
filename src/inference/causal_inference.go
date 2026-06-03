package inference

import (
	"fmt"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// CausalInference provides causal reasoning over a Bayesian network using
// the do-calculus. It supports interventional queries (do-operator), average
// treatment effect estimation, and backdoor adjustment.
type CausalInference struct {
	bn *models.BayesianNetwork
}

// NewCausalInference creates a new CausalInference engine from a validated
// BayesianNetwork. The network is deep-copied so the caller's original is
// not modified during interventions.
func NewCausalInference(bn *models.BayesianNetwork) (*CausalInference, error) {
	if bn == nil {
		return nil, fmt.Errorf("inference: BayesianNetwork must not be nil")
	}
	if err := bn.CheckModel(); err != nil {
		return nil, fmt.Errorf("inference: invalid BayesianNetwork: %w", err)
	}
	return &CausalInference{bn: bn.Copy()}, nil
}

// Query computes P(queryVars | do(doVars), evidence) using the truncated
// factorization (graph mutilation) approach:
//  1. Mutilate the graph: for each do-variable, remove all incoming edges.
//  2. Replace the do-variable's factor with a delta function.
//  3. Reduce remaining factors by evidence.
//  4. Use variable elimination on the mutilated model.
func (ci *CausalInference) Query(queryVars []string, doVars map[string]int, evidence map[string]int) (*factors.DiscreteFactor, error) {
	if len(queryVars) == 0 {
		return nil, fmt.Errorf("inference: queryVars must not be empty")
	}

	// Build the mutilated factor list directly from the CPDs, replacing
	// do-variable CPDs with delta factors (which also implicitly removes
	// parent dependencies, achieving graph mutilation at the factor level).
	doSet := make(map[string]bool, len(doVars))
	for v := range doVars {
		doSet[v] = true
	}

	var mutilatedFactors []*factors.DiscreteFactor
	for _, node := range ci.bn.Nodes() {
		cpd := ci.bn.GetCPD(node)
		if cpd == nil {
			return nil, fmt.Errorf("inference: no CPD for node %q", node)
		}

		if doSet[node] {
			// Replace this CPD with a delta factor over just the do-variable.
			doVal := doVars[node]
			card := cpd.VariableCard()
			if doVal < 0 || doVal >= card {
				return nil, fmt.Errorf("inference: do-value %d out of range for variable %q (card %d)", doVal, node, card)
			}
			deltaVals := make([]float64, card)
			deltaVals[doVal] = 1.0
			deltaFactor, err := factors.NewDiscreteFactor([]string{node}, []int{card}, deltaVals)
			if err != nil {
				return nil, fmt.Errorf("inference: failed to create delta factor for %q: %w", node, err)
			}
			mutilatedFactors = append(mutilatedFactors, deltaFactor)
		} else {
			mutilatedFactors = append(mutilatedFactors, cpd.ToFactor())
		}
	}

	ve := NewVariableElimination(mutilatedFactors)
	result, err := ve.Query(queryVars, evidence)
	if err != nil {
		return nil, fmt.Errorf("inference: variable elimination failed: %w", err)
	}

	return result, nil
}

// ATE computes the Average Treatment Effect of a binary treatment on a
// discrete outcome variable:
//
//	ATE = E[outcome | do(treatment=treatmentValues[1])] - E[outcome | do(treatment=treatmentValues[0])]
//
// For a discrete outcome with states 0, 1, ..., k-1, the expected value
// E[outcome | do(treatment=v)] = sum_i ( i * P(outcome=i | do(treatment=v)) ).
func (ci *CausalInference) ATE(treatment string, outcome string, treatmentValues [2]int) (float64, error) {
	expectations := [2]float64{}

	for idx, tVal := range treatmentValues {
		doVars := map[string]int{treatment: tVal}
		result, err := ci.Query([]string{outcome}, doVars, nil)
		if err != nil {
			return 0, fmt.Errorf("inference: ATE query failed for do(%s=%d): %w", treatment, tVal, err)
		}

		// Compute expected value: sum_i(i * P(outcome=i))
		vars := result.Variables()
		card := result.Cardinality()
		if len(vars) != 1 {
			return 0, fmt.Errorf("inference: ATE expected single-variable result, got %v", vars)
		}
		ev := 0.0
		assignment := make(map[string]int, 1)
		for i := 0; i < card[0]; i++ {
			assignment[vars[0]] = i
			ev += float64(i) * result.GetValue(assignment)
		}
		expectations[idx] = ev
	}

	return expectations[1] - expectations[0], nil
}

// BackdoorAdjustment estimates the ATE using the backdoor adjustment formula
// from observational data. The adjustment set must satisfy the backdoor
// criterion (use IsValidBackdoor to verify).
//
// The formula is:
//
//	P(outcome | do(treatment)) = sum_z P(outcome | treatment, z) * P(z)
//
// where z ranges over all configurations of the adjustment set variables.
// The ATE is then:
//
//	E[outcome | do(treatment=1)] - E[outcome | do(treatment=0)]
//
// using treatment values 0 and 1.
//
// All variables in the data must be integer-valued (discrete states).
func (ci *CausalInference) BackdoorAdjustment(treatment, outcome string, adjustmentSet []string, data *tabgo.DataFrame) (float64, error) {
	if data == nil {
		return 0, fmt.Errorf("inference: data must not be nil")
	}
	if !ci.IsValidBackdoor(treatment, outcome, adjustmentSet) {
		return 0, fmt.Errorf("inference: adjustment set %v does not satisfy backdoor criterion for (%s, %s)", adjustmentSet, treatment, outcome)
	}

	nRows := data.Len()
	if nRows == 0 {
		return 0, fmt.Errorf("inference: data has no rows")
	}

	// Read all relevant columns as int slices.
	treatmentData := data.Column(treatment).Int()
	outcomeData := data.Column(outcome).Int()

	adjData := make([][]int, len(adjustmentSet))
	for i, v := range adjustmentSet {
		adjData[i] = data.Column(v).Int()
	}

	// Enumerate unique adjustment set configurations and count them.
	configCount := make(map[string]int)
	configMap := make(map[string][]int)

	for row := 0; row < nRows; row++ {
		vals := make([]int, len(adjustmentSet))
		for j := range adjustmentSet {
			vals[j] = adjData[j][row]
		}
		key := fmt.Sprintf("%v", vals)
		configCount[key]++
		if _, exists := configMap[key]; !exists {
			configMap[key] = vals
		}
	}

	// For each treatment value (0 and 1), compute E[outcome | do(treatment=t)]
	// using the backdoor formula.
	expectations := [2]float64{}
	for tIdx, tVal := range [2]int{0, 1} {
		totalExpectation := 0.0

		for key, adjVals := range configMap {
			// P(Z = z) = count(z) / N
			pZ := float64(configCount[key]) / float64(nRows)

			// E[outcome | treatment=t, Z=z] from data
			sumOutcome := 0.0
			count := 0
			for row := 0; row < nRows; row++ {
				if treatmentData[row] != tVal {
					continue
				}
				match := true
				for j := range adjustmentSet {
					if adjData[j][row] != adjVals[j] {
						match = false
						break
					}
				}
				if match {
					sumOutcome += float64(outcomeData[row])
					count++
				}
			}

			if count > 0 {
				eOutcome := sumOutcome / float64(count)
				totalExpectation += eOutcome * pZ
			}
			// If count == 0, this treatment-adjustment configuration was never
			// observed; we skip it (contributes 0).
		}

		expectations[tIdx] = totalExpectation
	}

	return expectations[1] - expectations[0], nil
}

// IsValidBackdoor checks whether the given adjustment set satisfies the
// backdoor criterion for estimating the causal effect of treatment on outcome.
//
// The backdoor criterion requires:
//  1. No node in adjustmentSet is a descendant of treatment.
//  2. adjustmentSet d-separates treatment from outcome in the graph where
//     all edges out of treatment have been removed (the manipulated graph).
func (ci *CausalInference) IsValidBackdoor(treatment, outcome string, adjustmentSet []string) bool {
	// Build a DiGraph from the BN structure.
	g := bnToDigraph(ci.bn)

	// Check criterion 1: no adjustment variable is a descendant of treatment.
	descendants := allDescendants(g, treatment)
	for _, z := range adjustmentSet {
		if descendants[z] {
			return false
		}
	}

	// Build the manipulated graph: remove all edges out of treatment.
	mutilatedG := g.Copy()
	for _, child := range g.Successors(treatment) {
		_ = mutilatedG.RemoveEdge(treatment, child)
	}

	// Check criterion 2: treatment and outcome are d-separated given adjustmentSet
	// in the manipulated graph.
	xSet := map[string]bool{treatment: true}
	ySet := map[string]bool{outcome: true}
	zSet := make(map[string]bool, len(adjustmentSet))
	for _, z := range adjustmentSet {
		zSet[z] = true
	}

	return graphgo.DSeparation(mutilatedG, xSet, ySet, zSet)
}

// bnToDigraph reconstructs a graphgo.DiGraph from the BayesianNetwork's
// public API. This is needed because the BN's internal DAG is unexported.
func bnToDigraph(bn *models.BayesianNetwork) *graphgo.DiGraph {
	g := graphgo.NewDiGraph()
	for _, node := range bn.Nodes() {
		g.AddNode(node)
	}
	for _, edge := range bn.Edges() {
		g.AddEdge(edge[0], edge[1])
	}
	return g
}

// allDescendants returns the set of all descendants of the given node
// (not including the node itself) using BFS.
func allDescendants(g *graphgo.DiGraph, node string) map[string]bool {
	desc := make(map[string]bool)
	queue := g.Successors(node)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if desc[curr] {
			continue
		}
		desc[curr] = true
		queue = append(queue, g.Successors(curr)...)
	}
	return desc
}

// allAncestors returns the set of all ancestors of the given node
// (not including the node itself) using BFS.
func allAncestors(g *graphgo.DiGraph, node string) map[string]bool {
	return graphgo.Ancestors(g, node)
}

// GetAllBackdoorAdjustmentSets enumerates all valid backdoor adjustment sets
// for estimating the causal effect of treatment on outcome. Only feasible
// for small graphs.
func (ci *CausalInference) GetAllBackdoorAdjustmentSets(treatment, outcome string) [][]string {
	g := bnToDigraph(ci.bn)

	// Candidate variables: non-descendants of treatment, excluding treatment and outcome.
	descendants := allDescendants(g, treatment)
	nodes := ci.bn.Nodes()
	var candidates []string
	for _, n := range nodes {
		if n == treatment || n == outcome || descendants[n] {
			continue
		}
		candidates = append(candidates, n)
	}

	var results [][]string
	nCand := len(candidates)
	for mask := 0; mask < (1 << nCand); mask++ {
		var subset []string
		for i := 0; i < nCand; i++ {
			if mask&(1<<i) != 0 {
				subset = append(subset, candidates[i])
			}
		}
		if subset == nil {
			subset = []string{}
		}
		if ci.IsValidBackdoor(treatment, outcome, subset) {
			results = append(results, subset)
		}
	}
	return results
}

// IsValidFrontdoorAdjustmentSet checks whether the given set satisfies the
// front-door criterion for estimating the causal effect of treatment on outcome.
//
// The front-door criterion requires:
//  1. The set intercepts all directed paths from treatment to outcome.
//  2. No unblocked back-door path from treatment to any variable in the set.
//  3. All back-door paths from each set variable to outcome are blocked by treatment.
func (ci *CausalInference) IsValidFrontdoorAdjustmentSet(treatment, outcome string, frontdoorSet []string) bool {
	if len(frontdoorSet) == 0 {
		return false
	}
	g := bnToDigraph(ci.bn)

	fdSet := make(map[string]bool, len(frontdoorSet))
	for _, v := range frontdoorSet {
		fdSet[v] = true
	}

	// Condition 1: intercepts all directed paths from treatment to outcome.
	if !interceptsAllPaths(g, treatment, outcome, fdSet) {
		return false
	}

	// Condition 2: no unblocked back-door path from treatment to any fd variable.
	manipulated := g.Copy()
	for _, child := range g.Successors(treatment) {
		_ = manipulated.RemoveEdge(treatment, child)
	}
	xSet := map[string]bool{treatment: true}
	for _, m := range frontdoorSet {
		ySet := map[string]bool{m: true}
		if !graphgo.DSeparation(manipulated, xSet, ySet, map[string]bool{}) {
			return false
		}
	}

	// Condition 3: all back-door paths from each fd variable to outcome blocked by treatment.
	treatmentSet := map[string]bool{treatment: true}
	for _, m := range frontdoorSet {
		manipulatedM := g.Copy()
		for _, child := range g.Successors(m) {
			_ = manipulatedM.RemoveEdge(m, child)
		}
		mSet := map[string]bool{m: true}
		ySet := map[string]bool{outcome: true}
		if !graphgo.DSeparation(manipulatedM, mSet, ySet, treatmentSet) {
			return false
		}
	}

	return true
}

// interceptsAllPaths checks whether every directed path from src to dst passes
// through at least one node in the interceptSet.
func interceptsAllPaths(g *graphgo.DiGraph, src, dst string, interceptSet map[string]bool) bool {
	visited := make(map[string]bool)
	var dfs func(string) bool
	dfs = func(node string) bool {
		if node == dst {
			return true
		}
		visited[node] = true
		for _, child := range g.Successors(node) {
			if visited[child] {
				continue
			}
			if interceptSet[child] {
				continue
			}
			if dfs(child) {
				return true
			}
		}
		return false
	}
	visited[src] = true
	for _, child := range g.Successors(src) {
		if interceptSet[child] {
			continue
		}
		if dfs(child) {
			return false
		}
	}
	return true
}

// GetAllFrontdoorAdjustmentSets enumerates all valid front-door adjustment sets.
// Only feasible for small graphs.
func (ci *CausalInference) GetAllFrontdoorAdjustmentSets(treatment, outcome string) [][]string {
	g := bnToDigraph(ci.bn)

	// Candidates are mediators: nodes on directed paths from treatment to outcome.
	descT := allDescendants(g, treatment)
	ancO := allAncestors(g, outcome)
	var candidates []string
	for _, n := range ci.bn.Nodes() {
		if n == treatment || n == outcome {
			continue
		}
		if descT[n] && ancO[n] {
			candidates = append(candidates, n)
		}
	}

	var results [][]string
	nCand := len(candidates)
	for mask := 1; mask < (1 << nCand); mask++ {
		var subset []string
		for i := 0; i < nCand; i++ {
			if mask&(1<<i) != 0 {
				subset = append(subset, candidates[i])
			}
		}
		if ci.IsValidFrontdoorAdjustmentSet(treatment, outcome, subset) {
			results = append(results, subset)
		}
	}
	return results
}

// GetScalingIndicators returns variables that could serve as scaling indicators
// for instrumental variable analysis. These are exogenous variables (roots)
// that are not ancestors of the outcome through a backdoor path.
func (ci *CausalInference) GetScalingIndicators(treatment, outcome string) []string {
	g := bnToDigraph(ci.bn)
	var indicators []string
	for _, n := range ci.bn.Nodes() {
		if n == treatment || n == outcome {
			continue
		}
		if g.InDegree(n) == 0 {
			indicators = append(indicators, n)
		}
	}
	return indicators
}

// GetIVs finds instrumental variables for the causal effect of treatment on outcome.
// An IV must: (1) be associated with treatment, (2) affect outcome only through
// treatment, (3) have no common cause with outcome (d-separated from outcome
// given treatment in the manipulated graph).
func (ci *CausalInference) GetIVs(treatment, outcome string) []string {
	g := bnToDigraph(ci.bn)
	var ivs []string

	for _, n := range ci.bn.Nodes() {
		if n == treatment || n == outcome {
			continue
		}
		// Condition 1: n is associated with treatment (not d-separated marginally).
		nSet := map[string]bool{n: true}
		tSet := map[string]bool{treatment: true}
		if graphgo.DSeparation(g, nSet, tSet, map[string]bool{}) {
			continue
		}

		// Condition 2: n affects outcome only through treatment.
		// Remove treatment and check if n can still reach outcome.
		manipulated := g.Copy()
		manipulated.RemoveNode(treatment)
		oSet := map[string]bool{outcome: true}
		if !graphgo.DSeparation(manipulated, nSet, oSet, map[string]bool{}) {
			continue
		}

		// Condition 3: d-separated from outcome given treatment in manipulated graph.
		manipulated2 := g.Copy()
		for _, child := range g.Successors(treatment) {
			_ = manipulated2.RemoveEdge(treatment, child)
		}
		if !graphgo.DSeparation(manipulated2, nSet, oSet, tSet) {
			continue
		}

		ivs = append(ivs, n)
	}
	return ivs
}

// GetConditionalIVs finds conditional instrumental variables given a conditioning set.
func (ci *CausalInference) GetConditionalIVs(treatment, outcome string, conditioningSet []string) []string {
	g := bnToDigraph(ci.bn)
	condSet := make(map[string]bool, len(conditioningSet))
	for _, v := range conditioningSet {
		condSet[v] = true
	}

	var ivs []string
	for _, n := range ci.bn.Nodes() {
		if n == treatment || n == outcome || condSet[n] {
			continue
		}
		nSet := map[string]bool{n: true}
		tSet := map[string]bool{treatment: true}

		// Must be associated with treatment given conditioning set.
		if graphgo.DSeparation(g, nSet, tSet, condSet) {
			continue
		}

		// Must be d-separated from outcome given treatment + conditioning set.
		zSet := make(map[string]bool)
		for k, v := range condSet {
			zSet[k] = v
		}
		zSet[treatment] = true
		oSet := map[string]bool{outcome: true}

		manipulated := g.Copy()
		for _, child := range g.Successors(treatment) {
			_ = manipulated.RemoveEdge(treatment, child)
		}
		if !graphgo.DSeparation(manipulated, nSet, oSet, zSet) {
			continue
		}

		ivs = append(ivs, n)
	}
	return ivs
}

// GetTotalConditionalIVs returns all conditional IVs across all possible conditioning sets.
// Only feasible for small graphs.
func (ci *CausalInference) GetTotalConditionalIVs(treatment, outcome string) map[string][]string {
	result := make(map[string][]string)
	nodes := ci.bn.Nodes()
	var candidates []string
	for _, n := range nodes {
		if n != treatment && n != outcome {
			candidates = append(candidates, n)
		}
	}

	// For each candidate, try conditioning on subsets of others.
	for _, iv := range candidates {
		var others []string
		for _, c := range candidates {
			if c != iv {
				others = append(others, c)
			}
		}

		nOthers := len(others)
		for mask := 0; mask < (1 << nOthers); mask++ {
			var condSet []string
			for i := 0; i < nOthers; i++ {
				if mask&(1<<i) != 0 {
					condSet = append(condSet, others[i])
				}
			}
			civs := ci.GetConditionalIVs(treatment, outcome, condSet)
			for _, c := range civs {
				if c == iv {
					key := fmt.Sprintf("%v", condSet)
					result[iv] = append(result[iv], key)
					break
				}
			}
		}
	}
	return result
}

// canIdentifyByBackdoor returns true if there exists at least one valid
// backdoor adjustment set for the causal effect of treatment on outcome.
func (ci *CausalInference) canIdentifyByBackdoor(treatment, outcome string) bool {
	return len(ci.GetAllBackdoorAdjustmentSets(treatment, outcome)) > 0
}

// canIdentifyByFrontdoor returns true if there exists at least one valid
// front-door adjustment set for the causal effect of treatment on outcome.
func (ci *CausalInference) canIdentifyByFrontdoor(treatment, outcome string) bool {
	return len(ci.GetAllFrontdoorAdjustmentSets(treatment, outcome)) > 0
}

// canIdentifyByIV returns true if there exists at least one instrumental
// variable for the causal effect of treatment on outcome.
func (ci *CausalInference) canIdentifyByIV(treatment, outcome string) bool {
	return len(ci.GetIVs(treatment, outcome)) > 0
}

// IdentificationMethod returns which identification method applies for
// estimating the causal effect of treatment on outcome.
// Returns "backdoor", "frontdoor", "iv", or "none".
func (ci *CausalInference) IdentificationMethod(treatment, outcome string) string {
	if ci.canIdentifyByBackdoor(treatment, outcome) {
		return "backdoor"
	}
	if ci.canIdentifyByFrontdoor(treatment, outcome) {
		return "frontdoor"
	}
	if ci.canIdentifyByIV(treatment, outcome) {
		return "iv"
	}
	return "none"
}

// EstimateATE estimates the Average Treatment Effect from observational data
// using the appropriate identification method (backdoor or frontdoor adjustment).
// Uses binary treatment values [0, 1].
func (ci *CausalInference) EstimateATE(treatment, outcome string, data *tabgo.DataFrame) (float64, error) {
	if data == nil {
		return 0, fmt.Errorf("inference: data must not be nil")
	}

	// Try backdoor first.
	backdoorSets := ci.GetAllBackdoorAdjustmentSets(treatment, outcome)
	if len(backdoorSets) > 0 {
		// Use the first (smallest) valid adjustment set.
		var bestSet []string
		for _, s := range backdoorSets {
			if bestSet == nil || len(s) < len(bestSet) {
				bestSet = s
			}
		}
		return ci.BackdoorAdjustment(treatment, outcome, bestSet, data)
	}

	// Try model-based approach using do-calculus.
	return ci.ATE(treatment, outcome, [2]int{0, 1})
}

// GetProperBackdoorGraph returns the manipulated graph where all edges out of
// treatment have been removed. This is the graph used for backdoor criterion checks.
func (ci *CausalInference) GetProperBackdoorGraph(treatment string) *graphgo.DiGraph {
	g := bnToDigraph(ci.bn)
	manipulated := g.Copy()
	for _, child := range g.Successors(treatment) {
		_ = manipulated.RemoveEdge(treatment, child)
	}
	return manipulated
}

// IsValidAdjustmentSet checks whether the adjustment set satisfies the backdoor
// criterion. Delegates to the IsValidBackdoor method.
func (ci *CausalInference) IsValidAdjustmentSet(treatment, outcome string, adjustmentSet []string) bool {
	return ci.IsValidBackdoor(treatment, outcome, adjustmentSet)
}

// GetMinimalAdjustmentSet finds a minimal valid backdoor adjustment set by
// starting with parents of treatment and greedily removing variables.
func (ci *CausalInference) GetMinimalAdjustmentSet(treatment, outcome string) ([]string, error) {
	g := bnToDigraph(ci.bn)
	parents := g.Parents(treatment)

	if !ci.IsValidBackdoor(treatment, outcome, parents) {
		return nil, fmt.Errorf("inference: parents of %s do not form a valid adjustment set", treatment)
	}

	minimal := make([]string, len(parents))
	copy(minimal, parents)

	for i := 0; i < len(minimal); {
		candidate := make([]string, 0, len(minimal)-1)
		candidate = append(candidate, minimal[:i]...)
		candidate = append(candidate, minimal[i+1:]...)
		if ci.IsValidBackdoor(treatment, outcome, candidate) {
			minimal = candidate
		} else {
			i++
		}
	}

	return minimal, nil
}
