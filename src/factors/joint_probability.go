package factors

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const jpdTolerance = 1e-9

// JointProbabilityDistribution represents a joint probability distribution
// over discrete random variables. It wraps a DiscreteFactor and enforces
// the constraint that values must be non-negative and sum to 1.
type JointProbabilityDistribution struct {
	*DiscreteFactor
}

// NewJointProbabilityDistribution creates a new JointProbabilityDistribution.
// It validates that all values are non-negative and sum to 1.0 (within tolerance).
func NewJointProbabilityDistribution(variables []string, cardinality []int, values []float64) (*JointProbabilityDistribution, error) {
	df, err := NewDiscreteFactor(variables, cardinality, values)
	if err != nil {
		return nil, err
	}
	jpd := &JointProbabilityDistribution{DiscreteFactor: df}
	if err := jpd.Validate(); err != nil {
		return nil, err
	}
	return jpd, nil
}

// Validate checks that the distribution values are non-negative and sum to 1.0.
func (j *JointProbabilityDistribution) Validate() error {
	data := j.Values().Data()
	sum := 0.0
	for i, v := range data {
		if v < 0 {
			return fmt.Errorf("factors: negative probability %f at index %d", v, i)
		}
		sum += v
	}
	if math.Abs(sum-1.0) > jpdTolerance {
		return fmt.Errorf("factors: probabilities sum to %f, expected 1.0", sum)
	}
	return nil
}

// MarginalDistribution marginalizes out all variables not in the given list,
// returning a new JointProbabilityDistribution over the specified variables.
func (j *JointProbabilityDistribution) MarginalDistribution(variables []string) (*JointProbabilityDistribution, error) {
	if len(variables) == 0 {
		return nil, fmt.Errorf("factors: must specify at least one variable to keep")
	}

	// Determine which variables to marginalize out.
	keepSet := make(map[string]bool, len(variables))
	for _, v := range variables {
		if j.varIndex(v) == -1 {
			return nil, fmt.Errorf("factors: variable %q not in distribution", v)
		}
		keepSet[v] = true
	}

	var margVars []string
	for _, v := range j.variables {
		if !keepSet[v] {
			margVars = append(margVars, v)
		}
	}

	if len(margVars) == 0 {
		return j.Copy(), nil
	}

	df, err := j.DiscreteFactor.Marginalize(margVars)
	if err != nil {
		return nil, err
	}

	return &JointProbabilityDistribution{DiscreteFactor: df}, nil
}

// ConditionalDistribution computes P(variables | evidence) by reducing
// on the evidence, then normalizing the result.
func (j *JointProbabilityDistribution) ConditionalDistribution(variables []string, evidence map[string]int) (*DiscreteFactor, error) {
	if len(variables) == 0 {
		return nil, fmt.Errorf("factors: must specify at least one query variable")
	}
	if len(evidence) == 0 {
		return nil, fmt.Errorf("factors: must specify at least one evidence variable")
	}

	// Validate query variables exist and are not in evidence.
	for _, v := range variables {
		if j.varIndex(v) == -1 {
			return nil, fmt.Errorf("factors: query variable %q not in distribution", v)
		}
		if _, ok := evidence[v]; ok {
			return nil, fmt.Errorf("factors: variable %q cannot be both query and evidence", v)
		}
	}

	// First marginalize to keep only query variables + evidence variables.
	keepSet := make(map[string]bool)
	for _, v := range variables {
		keepSet[v] = true
	}
	for v := range evidence {
		keepSet[v] = true
	}
	var margVars []string
	for _, v := range j.variables {
		if !keepSet[v] {
			margVars = append(margVars, v)
		}
	}

	var working *DiscreteFactor
	if len(margVars) > 0 {
		var err error
		working, err = j.DiscreteFactor.Marginalize(margVars)
		if err != nil {
			return nil, err
		}
	} else {
		working = j.DiscreteFactor.Copy()
	}

	// Reduce on evidence.
	reduced, err := working.Reduce(evidence)
	if err != nil {
		return nil, err
	}

	// Normalize so the result sums to 1.
	reduced.Normalize()
	return reduced, nil
}

// CheckIndependence tests whether var1 is conditionally independent of var2
// given the variables in 'given', i.e. var1 _|_ var2 | given.
//
// It checks this by comparing P(var1, var2 | given) with
// P(var1 | given) * P(var2 | given) for all value combinations, using the
// supplied absolute tolerance atol.
func (j *JointProbabilityDistribution) CheckIndependence(var1, var2 string, given []string, atol float64) bool {
	if j.varIndex(var1) == -1 || j.varIndex(var2) == -1 {
		return false
	}
	for _, v := range given {
		if j.varIndex(v) == -1 {
			return false
		}
	}

	// Get cardinalities for the given variables.
	givenCards := make([]int, len(given))
	for i, v := range given {
		givenCards[i] = j.cardinality[j.varIndex(v)]
	}

	// Get cardinalities for var1 and var2.
	card1 := j.cardinality[j.varIndex(var1)]
	card2 := j.cardinality[j.varIndex(var2)]

	// Iterate over all assignments of given variables.
	givenSize := 1
	for _, c := range givenCards {
		givenSize *= c
	}

	for gFlat := 0; gFlat < givenSize; gFlat++ {
		evidence := make(map[string]int)
		rem := gFlat
		for i := len(given) - 1; i >= 0; i-- {
			evidence[given[i]] = rem % givenCards[i]
			rem /= givenCards[i]
		}

		// Get the joint factor over var1, var2 (conditioned on given if any).
		var jointFactor *DiscreteFactor
		if len(evidence) == 0 {
			// Unconditional case: marginalize to get P(var1, var2).
			mJoint, err := j.MarginalDistribution([]string{var1, var2})
			if err != nil {
				return false
			}
			jointFactor = mJoint.DiscreteFactor
		} else {
			var err error
			jointFactor, err = j.ConditionalDistribution([]string{var1, var2}, evidence)
			if err != nil {
				return false
			}
		}

		// Get P(var1 | given) and P(var2 | given) by marginalizing the joint.
		marg1, err := jointFactor.Marginalize([]string{var2})
		if err != nil {
			return false
		}
		marg2, err := jointFactor.Marginalize([]string{var1})
		if err != nil {
			return false
		}

		// Compare P(var1, var2 | given) with P(var1 | given) * P(var2 | given)
		for v1 := 0; v1 < card1; v1++ {
			for v2 := 0; v2 < card2; v2++ {
				jointAssign := map[string]int{var1: v1, var2: v2}
				pJoint := jointFactor.GetValue(jointAssign)

				p1 := marg1.GetValue(map[string]int{var1: v1})
				p2 := marg2.GetValue(map[string]int{var2: v2})

				if math.Abs(pJoint-p1*p2) > atol {
					return false
				}
			}
		}
	}

	return true
}

// Copy returns a deep copy of the JointProbabilityDistribution.
func (j *JointProbabilityDistribution) Copy() *JointProbabilityDistribution {
	return &JointProbabilityDistribution{
		DiscreteFactor: j.DiscreteFactor.Copy(),
	}
}

// GetIndependencies finds all pairwise conditional independencies in the
// distribution by testing P(X,Y|Z) = P(X|Z)*P(Y|Z) for all variable pairs
// and all possible conditioning sets Z (subsets of remaining variables).
func (j *JointProbabilityDistribution) GetIndependencies(atol float64) [][3][]string {
	vars := j.Variables()
	var result [][3][]string

	for i := 0; i < len(vars); i++ {
		for k := i + 1; k < len(vars); k++ {
			v1, v2 := vars[i], vars[k]

			var others []string
			for _, v := range vars {
				if v != v1 && v != v2 {
					others = append(others, v)
				}
			}

			if j.CheckIndependence(v1, v2, nil, atol) {
				result = append(result, [3][]string{{v1}, {v2}, {}})
			}

			nOthers := len(others)
			for mask := 1; mask < (1 << nOthers); mask++ {
				var condSet []string
				for b := 0; b < nOthers; b++ {
					if mask&(1<<b) != 0 {
						condSet = append(condSet, others[b])
					}
				}
				if j.CheckIndependence(v1, v2, condSet, atol) {
					cs := make([]string, len(condSet))
					copy(cs, condSet)
					sort.Strings(cs)
					result = append(result, [3][]string{{v1}, {v2}, cs})
				}
			}
		}
	}
	return result
}

// MinimalIMap finds a minimal I-map DAG for this distribution using a greedy
// approach. It returns the edge list of the DAG.
func (j *JointProbabilityDistribution) MinimalIMap(ordering []string, atol float64) [][2]string {
	var edges [][2]string

	for i, node := range ordering {
		preceding := make([]string, i)
		copy(preceding, ordering[:i])

		parents := make([]string, len(preceding))
		copy(parents, preceding)

		for k := 0; k < len(parents); {
			candidate := make([]string, 0, len(parents)-1)
			candidate = append(candidate, parents[:k]...)
			candidate = append(candidate, parents[k+1:]...)

			removed := parents[k]
			if j.CheckIndependence(node, removed, candidate, atol) {
				parents = candidate
			} else {
				k++
			}
		}

		for _, p := range parents {
			edges = append(edges, [2]string{p, node})
		}
	}

	return edges
}

// IsIMap checks if the given DAG (represented as edges) is an I-map of this
// distribution.
func (j *JointProbabilityDistribution) IsIMap(edges [][2]string, atol float64) bool {
	vars := j.Variables()
	parentMap := make(map[string][]string)
	for _, v := range vars {
		parentMap[v] = nil
	}
	for _, e := range edges {
		parentMap[e[1]] = append(parentMap[e[1]], e[0])
	}

	childMap := make(map[string]map[string]bool)
	for _, v := range vars {
		childMap[v] = make(map[string]bool)
	}
	for _, e := range edges {
		childMap[e[0]][e[1]] = true
	}

	for _, node := range vars {
		parents := parentMap[node]
		parentSet := make(map[string]bool)
		for _, p := range parents {
			parentSet[p] = true
		}

		for _, other := range vars {
			if other == node || parentSet[other] || childMap[node][other] {
				continue
			}
			if !j.CheckIndependence(node, other, parents, atol) {
				if !jpdHasDirectedPath(edges, node, other) {
					return false
				}
			}
		}
	}

	return true
}

// jpdHasDirectedPath checks if there is a directed path from src to dst.
func jpdHasDirectedPath(edges [][2]string, src, dst string) bool {
	children := make(map[string][]string)
	for _, e := range edges {
		children[e[0]] = append(children[e[0]], e[1])
	}
	visited := make(map[string]bool)
	queue := []string{src}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == dst {
			return true
		}
		if visited[curr] {
			continue
		}
		visited[curr] = true
		queue = append(queue, children[curr]...)
	}
	return false
}

// ToFactor converts the JointProbabilityDistribution to a DiscreteFactor.
func (j *JointProbabilityDistribution) ToFactor() *DiscreteFactor {
	return j.DiscreteFactor.Copy()
}

// PMap returns all conditional independencies as a string representation.
func (j *JointProbabilityDistribution) PMap(atol float64) string {
	indeps := j.GetIndependencies(atol)
	if len(indeps) == 0 {
		return "{}"
	}

	var parts []string
	for _, ind := range indeps {
		x := strings.Join(ind[0], ", ")
		y := strings.Join(ind[1], ", ")
		z := strings.Join(ind[2], ", ")
		if z == "" {
			parts = append(parts, fmt.Sprintf("%s _|_ %s", x, y))
		} else {
			parts = append(parts, fmt.Sprintf("%s _|_ %s | %s", x, y, z))
		}
	}
	return "{\n  " + strings.Join(parts, ",\n  ") + "\n}"
}
