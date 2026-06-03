package inference

import (
	"fmt"
	"math/rand"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// ApproxInference estimates marginal distributions via likelihood-weighted
// sampling. It operates on a set of discrete factors (as obtained from a
// Bayesian network's ToMarkovFactors or constructed directly), similar to
// VariableElimination, but trades exactness for scalability.
type ApproxInference struct {
	factors []*factors.DiscreteFactor
	rng     *rand.Rand
}

// NewApproxInference creates a new ApproxInference engine from the given
// factor list. Each factor is deep-copied so the caller's originals are not
// modified during inference. The seed controls the random number generator
// for reproducibility.
func NewApproxInference(factorList []*factors.DiscreteFactor, seed int64) *ApproxInference {
	copied := make([]*factors.DiscreteFactor, len(factorList))
	for i, f := range factorList {
		copied[i] = f.Copy()
	}
	return &ApproxInference{
		factors: copied,
		rng:     rand.New(rand.NewSource(seed)),
	}
}

// Query approximates P(queryVars | evidence) using likelihood-weighted
// sampling.
//
// Steps:
//  1. Collect all variables and their cardinalities from all factors.
//  2. Reduce factors by evidence to obtain conditional factors.
//  3. For nSamples iterations, draw a uniform sample over all non-evidence
//     variables, compute the weight as the product of all (reduced) factor
//     values at that assignment, and accumulate the weighted counts for the
//     query variable assignments.
//  4. Normalize the accumulated counts to get an approximate marginal.
//  5. Return the result as a DiscreteFactor.
func (ai *ApproxInference) Query(queryVars []string, evidence map[string]int, nSamples int) (*factors.DiscreteFactor, error) {
	if len(queryVars) == 0 {
		return nil, fmt.Errorf("approx_inference: queryVars must not be empty")
	}
	if nSamples <= 0 {
		return nil, fmt.Errorf("approx_inference: nSamples must be positive")
	}

	// Step 1: Collect all variables and cardinalities.
	cardMap := make(map[string]int)
	for _, f := range ai.factors {
		vars := f.Variables()
		card := f.Cardinality()
		for i, v := range vars {
			if _, ok := cardMap[v]; !ok {
				cardMap[v] = card[i]
			}
		}
	}

	// Validate query variables exist.
	for _, v := range queryVars {
		if _, ok := cardMap[v]; !ok {
			return nil, fmt.Errorf("approx_inference: query variable %q not found in any factor", v)
		}
	}

	// Validate evidence variables exist and values are in range.
	for v, val := range evidence {
		c, ok := cardMap[v]
		if !ok {
			return nil, fmt.Errorf("approx_inference: evidence variable %q not found in any factor", v)
		}
		if val < 0 || val >= c {
			return nil, fmt.Errorf("approx_inference: evidence value %d out of range for variable %q (card %d)", val, v, c)
		}
	}

	// Step 2: Reduce factors by evidence.
	workingFactors, err := reduceAll(ai.factors, evidence)
	if err != nil {
		return nil, fmt.Errorf("approx_inference: evidence reduction failed: %w", err)
	}

	// Build the list of free (non-evidence) variables to sample over.
	var freeVars []string
	var freeCard []int
	for _, f := range ai.factors {
		for _, v := range f.Variables() {
			if _, isEvidence := evidence[v]; isEvidence {
				continue
			}
			// Only add if not already seen.
			found := false
			for _, fv := range freeVars {
				if fv == v {
					found = true
					break
				}
			}
			if !found {
				freeVars = append(freeVars, v)
				freeCard = append(freeCard, cardMap[v])
			}
		}
	}

	// Build the query variable cardinalities and index positions.
	queryCard := make([]int, len(queryVars))
	for i, v := range queryVars {
		queryCard[i] = cardMap[v]
	}

	// Compute total size of the query result table.
	querySize := 1
	for _, c := range queryCard {
		querySize *= c
	}
	counts := make([]float64, querySize)

	// Step 3: Likelihood-weighted sampling.
	assignment := make(map[string]int, len(freeVars))
	for s := 0; s < nSamples; s++ {
		// Draw a uniform random assignment for all free variables.
		for i, v := range freeVars {
			assignment[v] = ai.rng.Intn(freeCard[i])
		}

		// Compute the weight: product of all reduced factor values at this
		// assignment.
		weight := 1.0
		for _, f := range workingFactors {
			fVars := f.Variables()
			fAssign := make(map[string]int, len(fVars))
			for _, fv := range fVars {
				fAssign[fv] = assignment[fv]
			}
			weight *= f.GetValue(fAssign)
		}

		if weight == 0 {
			continue
		}

		// Compute the flat index for the query variables.
		queryFlat := 0
		stride := 1
		for i := len(queryVars) - 1; i >= 0; i-- {
			queryFlat += assignment[queryVars[i]] * stride
			stride *= queryCard[i]
		}
		counts[queryFlat] += weight
	}

	// Step 4: Normalize.
	sum := 0.0
	for _, c := range counts {
		sum += c
	}
	if sum == 0 {
		return nil, fmt.Errorf("approx_inference: all samples had zero weight; try increasing nSamples")
	}
	for i := range counts {
		counts[i] /= sum
	}

	// Step 5: Return as DiscreteFactor.
	return factors.NewDiscreteFactor(queryVars, queryCard, counts)
}

// GetDistribution approximates the full joint distribution over all variables
// using likelihood-weighted sampling with the given number of samples.
func (ai *ApproxInference) GetDistribution(nSamples int) (*factors.DiscreteFactor, error) {
	if nSamples <= 0 {
		return nil, fmt.Errorf("approx_inference: nSamples must be positive")
	}

	// Collect all variables and cardinalities.
	var allVars []string
	cardMap := make(map[string]int)
	for _, f := range ai.factors {
		vars := f.Variables()
		card := f.Cardinality()
		for i, v := range vars {
			if _, ok := cardMap[v]; !ok {
				cardMap[v] = card[i]
				allVars = append(allVars, v)
			}
		}
	}

	allCard := make([]int, len(allVars))
	for i, v := range allVars {
		allCard[i] = cardMap[v]
	}

	// Compute total size.
	totalSize := 1
	for _, c := range allCard {
		totalSize *= c
	}
	counts := make([]float64, totalSize)

	// Sample.
	assignment := make(map[string]int, len(allVars))
	for s := 0; s < nSamples; s++ {
		for i, v := range allVars {
			assignment[v] = ai.rng.Intn(allCard[i])
		}

		weight := 1.0
		for _, f := range ai.factors {
			fVars := f.Variables()
			fAssign := make(map[string]int, len(fVars))
			for _, fv := range fVars {
				fAssign[fv] = assignment[fv]
			}
			weight *= f.GetValue(fAssign)
		}

		if weight == 0 {
			continue
		}

		flat := 0
		stride := 1
		for i := len(allVars) - 1; i >= 0; i-- {
			flat += assignment[allVars[i]] * stride
			stride *= allCard[i]
		}
		counts[flat] += weight
	}

	// Normalize.
	sum := 0.0
	for _, c := range counts {
		sum += c
	}
	if sum == 0 {
		return nil, fmt.Errorf("approx_inference: all samples had zero weight")
	}
	for i := range counts {
		counts[i] /= sum
	}

	return factors.NewDiscreteFactor(allVars, allCard, counts)
}

// QueryRejection approximates P(queryVars | evidence) using rejection
// sampling. It draws samples from the prior (product of all factors) and
// rejects those that do not match the observed evidence.
//
// This is simpler but less efficient than likelihood-weighted sampling,
// especially when evidence is unlikely.
func (ai *ApproxInference) QueryRejection(queryVars []string, evidence map[string]int, nSamples int) (*factors.DiscreteFactor, error) {
	if len(queryVars) == 0 {
		return nil, fmt.Errorf("approx_inference: queryVars must not be empty")
	}
	if nSamples <= 0 {
		return nil, fmt.Errorf("approx_inference: nSamples must be positive")
	}

	cardMap := make(map[string]int)
	var allVars []string
	for _, f := range ai.factors {
		vars := f.Variables()
		card := f.Cardinality()
		for i, v := range vars {
			if _, ok := cardMap[v]; !ok {
				cardMap[v] = card[i]
				allVars = append(allVars, v)
			}
		}
	}

	for _, v := range queryVars {
		if _, ok := cardMap[v]; !ok {
			return nil, fmt.Errorf("approx_inference: query variable %q not found in any factor", v)
		}
	}
	for v, val := range evidence {
		c, ok := cardMap[v]
		if !ok {
			return nil, fmt.Errorf("approx_inference: evidence variable %q not found in any factor", v)
		}
		if val < 0 || val >= c {
			return nil, fmt.Errorf("approx_inference: evidence value %d out of range for variable %q (card %d)", val, v, c)
		}
	}

	queryCard := make([]int, len(queryVars))
	for i, v := range queryVars {
		queryCard[i] = cardMap[v]
	}
	querySize := 1
	for _, c := range queryCard {
		querySize *= c
	}
	counts := make([]float64, querySize)

	allCard := make([]int, len(allVars))
	for i, v := range allVars {
		allCard[i] = cardMap[v]
	}

	assignment := make(map[string]int, len(allVars))
	accepted := 0

	for s := 0; s < nSamples; s++ {
		// Draw uniform random assignment.
		for i, v := range allVars {
			assignment[v] = ai.rng.Intn(allCard[i])
		}

		// Compute weight from factor product.
		weight := 1.0
		for _, f := range ai.factors {
			fVars := f.Variables()
			fAssign := make(map[string]int, len(fVars))
			for _, fv := range fVars {
				fAssign[fv] = assignment[fv]
			}
			weight *= f.GetValue(fAssign)
		}
		if weight == 0 {
			continue
		}

		// Reject if evidence doesn't match.
		match := true
		for v, val := range evidence {
			if assignment[v] != val {
				match = false
				break
			}
		}
		if !match {
			continue
		}

		accepted++
		queryFlat := 0
		stride := 1
		for i := len(queryVars) - 1; i >= 0; i-- {
			queryFlat += assignment[queryVars[i]] * stride
			stride *= queryCard[i]
		}
		counts[queryFlat] += weight
	}

	if accepted == 0 {
		return nil, fmt.Errorf("approx_inference: no samples matched evidence; try increasing nSamples")
	}

	sum := 0.0
	for _, c := range counts {
		sum += c
	}
	if sum == 0 {
		return nil, fmt.Errorf("approx_inference: all accepted samples had zero weight")
	}
	for i := range counts {
		counts[i] /= sum
	}

	return factors.NewDiscreteFactor(queryVars, queryCard, counts)
}

// QueryGibbs approximates P(queryVars | evidence) using Gibbs sampling
// (Markov Chain Monte Carlo). It initializes all non-evidence variables
// randomly, then iteratively resamples each variable from its full
// conditional distribution. After a burn-in period, samples are collected.
func (ai *ApproxInference) QueryGibbs(queryVars []string, evidence map[string]int, nSamples int, burnIn int) (*factors.DiscreteFactor, error) {
	if len(queryVars) == 0 {
		return nil, fmt.Errorf("approx_inference: queryVars must not be empty")
	}
	if nSamples <= 0 {
		return nil, fmt.Errorf("approx_inference: nSamples must be positive")
	}
	if burnIn < 0 {
		burnIn = 0
	}

	cardMap := make(map[string]int)
	var allVars []string
	for _, f := range ai.factors {
		vars := f.Variables()
		card := f.Cardinality()
		for i, v := range vars {
			if _, ok := cardMap[v]; !ok {
				cardMap[v] = card[i]
				allVars = append(allVars, v)
			}
		}
	}

	for _, v := range queryVars {
		if _, ok := cardMap[v]; !ok {
			return nil, fmt.Errorf("approx_inference: query variable %q not found in any factor", v)
		}
	}
	for v, val := range evidence {
		c, ok := cardMap[v]
		if !ok {
			return nil, fmt.Errorf("approx_inference: evidence variable %q not found in any factor", v)
		}
		if val < 0 || val >= c {
			return nil, fmt.Errorf("approx_inference: evidence value %d out of range for variable %q (card %d)", val, v, c)
		}
	}

	// Determine free (non-evidence) variables.
	evidenceSet := make(map[string]bool, len(evidence))
	for v := range evidence {
		evidenceSet[v] = true
	}
	var freeVars []string
	for _, v := range allVars {
		if !evidenceSet[v] {
			freeVars = append(freeVars, v)
		}
	}

	// Precompute: for each variable, which factors involve it?
	varFactors := make(map[string][]*factors.DiscreteFactor)
	for _, f := range ai.factors {
		for _, v := range f.Variables() {
			varFactors[v] = append(varFactors[v], f)
		}
	}

	// Initialize assignment.
	assignment := make(map[string]int, len(allVars))
	for v, val := range evidence {
		assignment[v] = val
	}
	for _, v := range freeVars {
		assignment[v] = ai.rng.Intn(cardMap[v])
	}

	// Query result accumulators.
	queryCard := make([]int, len(queryVars))
	for i, v := range queryVars {
		queryCard[i] = cardMap[v]
	}
	querySize := 1
	for _, c := range queryCard {
		querySize *= c
	}
	counts := make([]float64, querySize)

	totalIter := burnIn + nSamples
	for iter := 0; iter < totalIter; iter++ {
		// Resample each free variable from its full conditional.
		for _, v := range freeVars {
			card := cardMap[v]
			probs := make([]float64, card)
			for s := 0; s < card; s++ {
				assignment[v] = s
				p := 1.0
				for _, f := range varFactors[v] {
					fVars := f.Variables()
					fAssign := make(map[string]int, len(fVars))
					for _, fv := range fVars {
						fAssign[fv] = assignment[fv]
					}
					p *= f.GetValue(fAssign)
				}
				probs[s] = p
			}

			// Normalize and sample.
			sum := 0.0
			for _, p := range probs {
				sum += p
			}
			if sum == 0 {
				assignment[v] = ai.rng.Intn(card)
				continue
			}
			r := ai.rng.Float64() * sum
			cumulative := 0.0
			chosen := card - 1
			for s := 0; s < card; s++ {
				cumulative += probs[s]
				if r <= cumulative {
					chosen = s
					break
				}
			}
			assignment[v] = chosen
		}

		// Collect sample after burn-in.
		if iter >= burnIn {
			queryFlat := 0
			stride := 1
			for i := len(queryVars) - 1; i >= 0; i-- {
				queryFlat += assignment[queryVars[i]] * stride
				stride *= queryCard[i]
			}
			counts[queryFlat]++
		}
	}

	// Normalize.
	sum := 0.0
	for _, c := range counts {
		sum += c
	}
	if sum == 0 {
		return nil, fmt.Errorf("approx_inference: Gibbs sampling produced no valid samples")
	}
	for i := range counts {
		counts[i] /= sum
	}

	return factors.NewDiscreteFactor(queryVars, queryCard, counts)
}

// GetDistributionWithEvidence approximates the distribution over queryVars
// given evidence, using the specified number of samples. This extends
// GetDistribution to support evidence conditioning.
func (ai *ApproxInference) GetDistributionWithEvidence(queryVars []string, evidence map[string]int, nSamples int) (*factors.DiscreteFactor, error) {
	return ai.Query(queryVars, evidence, nSamples)
}

// MAPQuery approximates the MAP (Maximum A Posteriori) assignment for
// queryVars given evidence, using sampling to find the assignment with
// the highest approximate probability.
func (ai *ApproxInference) MAPQuery(queryVars []string, evidence map[string]int, nSamples int) (map[string]int, error) {
	if len(queryVars) == 0 {
		return nil, fmt.Errorf("approx_inference: queryVars must not be empty")
	}
	if nSamples <= 0 {
		return nil, fmt.Errorf("approx_inference: nSamples must be positive")
	}

	result, err := ai.Query(queryVars, evidence, nSamples)
	if err != nil {
		return nil, err
	}

	// Find the assignment with maximum value.
	vars := result.Variables()
	card := result.Cardinality()
	totalSize := 1
	for _, c := range card {
		totalSize *= c
	}

	bestVal := -1.0
	bestAssignment := make(map[string]int, len(vars))

	for flat := 0; flat < totalSize; flat++ {
		assignment := make(map[string]int, len(vars))
		rem := flat
		for i := len(vars) - 1; i >= 0; i-- {
			assignment[vars[i]] = rem % card[i]
			rem /= card[i]
		}
		val := result.GetValue(assignment)
		if val > bestVal {
			bestVal = val
			bestAssignment = assignment
		}
	}

	return bestAssignment, nil
}
