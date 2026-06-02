package models

import (
	"fmt"
	"math"
	"math/rand"
)

// SetStartState is a convenience that validates a start state index.
// It returns an error if the index is out of range; otherwise it returns
// the validated index. (MarkovChain is stateless regarding "current state",
// so this simply validates.)
func (mc *MarkovChain) SetStartState(state int) (int, error) {
	if state < 0 || state >= mc.NumStates() {
		return -1, fmt.Errorf("models: state %d out of range [0, %d)", state, mc.NumStates())
	}
	return state, nil
}

// AddVariable adds a new state to the Markov chain. The transition
// probabilities for the new state (row and column) are initialized to
// a uniform distribution. Existing rows are re-normalized.
func (mc *MarkovChain) AddVariable(name string) {
	n := mc.NumStates()
	newN := n + 1

	// Extend each existing row with 0 probability, then re-normalize.
	newMat := make([][]float64, newN)
	for i := 0; i < n; i++ {
		row := make([]float64, newN)
		copy(row, mc.transitionMatrix[i])
		row[n] = 0.0
		// Re-normalize.
		sum := 0.0
		for _, v := range row {
			sum += v
		}
		if sum > 0 {
			for j := range row {
				row[j] /= sum
			}
		}
		newMat[i] = row
	}
	// New state row: uniform.
	newRow := make([]float64, newN)
	for j := 0; j < newN; j++ {
		newRow[j] = 1.0 / float64(newN)
	}
	newMat[n] = newRow

	mc.transitionMatrix = newMat

	if mc.stateNames != nil {
		mc.stateNames = append(mc.stateNames, name)
	}
}

// AddVariablesFrom copies state names and transition probabilities from
// another MarkovChain, adding any states not already present. If the
// source has unnamed states, they are added with generated names.
func (mc *MarkovChain) AddVariablesFrom(other *MarkovChain) {
	if other == nil {
		return
	}
	otherNames := other.StateNames()
	if otherNames == nil {
		// Generate names for unnamed states.
		for i := 0; i < other.NumStates(); i++ {
			name := fmt.Sprintf("state_%d", mc.NumStates())
			mc.AddVariable(name)
		}
		return
	}

	existing := make(map[string]bool)
	if mc.stateNames != nil {
		for _, s := range mc.stateNames {
			existing[s] = true
		}
	}
	for _, name := range otherNames {
		if !existing[name] {
			mc.AddVariable(name)
		}
	}
}

// AddTransitionModel sets the transition matrix for the chain. The provided
// matrix must be square with the same size as the current number of states,
// and each row must sum to 1 (within tolerance).
func (mc *MarkovChain) AddTransitionModel(matrix [][]float64) error {
	n := mc.NumStates()
	if len(matrix) != n {
		return fmt.Errorf("models: matrix has %d rows, expected %d", len(matrix), n)
	}

	const tol = 1e-6
	for i, row := range matrix {
		if len(row) != n {
			return fmt.Errorf("models: matrix row %d has length %d, expected %d", i, len(row), n)
		}
		sum := 0.0
		for _, v := range row {
			if v < -tol {
				return fmt.Errorf("models: matrix has negative value %f at row %d", v, i)
			}
			sum += v
		}
		if math.Abs(sum-1.0) > tol {
			return fmt.Errorf("models: matrix row %d sums to %f, expected 1.0", i, sum)
		}
	}

	// Deep copy.
	mat := make([][]float64, n)
	for i, row := range matrix {
		mat[i] = make([]float64, n)
		copy(mat[i], row)
	}
	mc.transitionMatrix = mat
	return nil
}

// ProbFromSample estimates transition probabilities from a sequence of
// state indices. Returns a new transition matrix where T[i][j] is the
// fraction of times state i was followed by state j.
func (mc *MarkovChain) ProbFromSample(sequence []int) ([][]float64, error) {
	n := mc.NumStates()
	if len(sequence) < 2 {
		return nil, fmt.Errorf("models: sequence must have at least 2 elements")
	}

	counts := make([][]float64, n)
	rowTotals := make([]float64, n)
	for i := 0; i < n; i++ {
		counts[i] = make([]float64, n)
	}

	for k := 0; k < len(sequence)-1; k++ {
		from := sequence[k]
		to := sequence[k+1]
		if from < 0 || from >= n || to < 0 || to >= n {
			return nil, fmt.Errorf("models: state index out of range at position %d", k)
		}
		counts[from][to]++
		rowTotals[from]++
	}

	result := make([][]float64, n)
	for i := 0; i < n; i++ {
		result[i] = make([]float64, n)
		if rowTotals[i] > 0 {
			for j := 0; j < n; j++ {
				result[i][j] = counts[i][j] / rowTotals[i]
			}
		} else {
			// Uniform for states never visited.
			for j := 0; j < n; j++ {
				result[i][j] = 1.0 / float64(n)
			}
		}
	}

	return result, nil
}

// GenerateSample generates a slice of n sampled state indices starting from
// startState, using the given seed for reproducibility.
func (mc *MarkovChain) GenerateSample(n int, startState int, seed int64) ([]int, error) {
	return mc.Sample(n, startState, seed)
}

// IsStationarity tests whether the given distribution is stationary for
// this Markov chain (i.e., pi * T = pi within tolerance).
func (mc *MarkovChain) IsStationarity(dist []float64) (bool, error) {
	n := mc.NumStates()
	if len(dist) != n {
		return false, fmt.Errorf("models: distribution length %d does not match %d states", len(dist), n)
	}

	const tol = 1e-6

	// Compute pi * T.
	piT := make([]float64, n)
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			piT[j] += dist[i] * mc.transitionMatrix[i][j]
		}
	}

	for j := 0; j < n; j++ {
		if math.Abs(piT[j]-dist[j]) > tol {
			return false, nil
		}
	}
	return true, nil
}

// RandomState samples a state index from the stationary distribution.
func (mc *MarkovChain) RandomState(seed int64) (int, error) {
	pi, err := mc.StationaryDistribution()
	if err != nil {
		return -1, err
	}

	rng := rand.New(rand.NewSource(seed))
	u := rng.Float64()
	cumSum := 0.0
	for i, p := range pi {
		cumSum += p
		if u < cumSum {
			return i, nil
		}
	}
	return len(pi) - 1, nil
}

// Copy returns a deep copy of the MarkovChain.
func (mc *MarkovChain) Copy() *MarkovChain {
	n := mc.NumStates()
	mat := make([][]float64, n)
	for i, row := range mc.transitionMatrix {
		mat[i] = make([]float64, n)
		copy(mat[i], row)
	}

	var names []string
	if mc.stateNames != nil {
		names = make([]string, len(mc.stateNames))
		copy(names, mc.stateNames)
	}

	return &MarkovChain{
		transitionMatrix: mat,
		stateNames:       names,
	}
}
